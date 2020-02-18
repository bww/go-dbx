package persist

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
	"github.com/bww/go-dbx/v1/persist/ident"
	"github.com/bww/go-dbx/v1/persist/registry"
	"github.com/bww/go-dbx/v1/test"
	"github.com/bww/go-util/env"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.Init(testDB, test.WithMigrations(env.Etc("migrations")))
	os.Exit(m.Run())
}

const (
	testDB         = "dbx_v1_persist_test"
	testTable      = "test_entity"
	anotherTable   = "another_entity"
	omitemptyTable = "omitempty_entity"
)

type DontUseThisExportedEntity struct {
	B string `db:"b"`
}

type anotherEntity struct {
	X string `db:"x,pk"`
	Z int    `db:"z"`
}

type testEntity struct {
	*DontUseThisExportedEntity
	A string `db:"a,pk"`
	C int    `db:"c"`
	E int    `db:"e,omitempty"`
	D []*anotherEntity
}

type testPersister struct{}

func (p *testPersister) FetchRelated(pst Persister, ent interface{}) error {
	z := ent.(*testEntity)
	q := `
    SELECT {a.*} FROM ` + anotherTable + ` AS a
    INNER JOIN test_entity_r_another_entity AS r ON r.x = a.x
    WHERE r.a = $1
    ORDER BY a.z`

	var another []*anotherEntity
	err := pst.Select(&another, q, z.A)
	if err != nil {
		return err
	}

	z.D = another
	return nil
}

func (p *testPersister) StoreRelated(pst Persister, ent interface{}) error {
	z := ent.(*testEntity)
	for _, e := range z.D {
		err := pst.Store("another_entity", e, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *testPersister) StoreReferences(pst Persister, ent interface{}) error {
	z := ent.(*testEntity)
	for _, e := range z.D {
		_, err := pst.Exec(`INSERT INTO test_entity_r_another_entity (a, x) VALUES ($1, $2)`, z.A, e.X)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *testPersister) DeleteRelated(pst Persister, ent interface{}) error {
	z := ent.(*testEntity)
	for _, e := range z.D {
		err := pst.Delete("another_entity", e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *testPersister) DeleteReferences(pst Persister, ent interface{}) error {
	z := ent.(*testEntity)
	_, err := pst.Exec(`DELETE FROM test_entity_r_another_entity WHERE a = $1`, z.A)
	if err != nil {
		return err
	}
	return nil
}

func TestPersist(t *testing.T) {
	db := test.DB()
	reg := registry.New()
	pst := New(db, entity.NewFieldMapper(), reg, ident.AlphaNumeric(32))
	var err error

	reg.Set(reflect.ValueOf((*testEntity)(nil)).Type(), &testPersister{})

	e1 := &testEntity{
		DontUseThisExportedEntity: &DontUseThisExportedEntity{
			B: "This is the value of B",
		},
		C: 888,
	}

	err = pst.Store(testTable, e1, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Len(t, e1.A, 32)
	}

	e2 := &testEntity{
		DontUseThisExportedEntity: &DontUseThisExportedEntity{
			B: "Never is this the value of B",
		},
		C: 999,
		E: 111,
		D: []*anotherEntity{
			{Z: 111},
			{Z: 222},
			{Z: 333},
		},
	}

	err = pst.Store(testTable, e2, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Len(t, e2.A, 32)
	}

	var ca, cb string
	var cc int
	err = db.QueryRow(`SELECT a, b, c FROM `+testTable+` WHERE a = $1`, e1.A).Scan(&ca, &cb, &cc)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, e1.A, ca)
		assert.Equal(t, e1.B, cb)
		assert.Equal(t, e1.C, cc)
	}

	var eb testEntity
	err = db.QueryRowx(`SELECT a, b, c, e FROM `+testTable+` WHERE a = $1`, e2.A).StructScan(&eb)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, e2.A, eb.A)
		assert.Equal(t, e2.B, eb.B)
		assert.Equal(t, e2.C, eb.C)
		assert.Equal(t, e2.E, eb.E)
	}

	var ec testEntity
	err = pst.Fetch(testTable, &ec, e1.A)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, e1.A, ec.A)
		assert.Equal(t, e1.B, ec.B)
		assert.Equal(t, e1.C, ec.C)
		assert.Equal(t, e1.E, ec.E)
	}

	var ed testEntity
	err = pst.Fetch(testTable, &ed, e2.A)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, e2.A, ed.A)
		assert.Equal(t, e2.B, ed.B)
		assert.Equal(t, e2.C, ed.C)
		assert.Equal(t, e2.E, ed.E)
	}

	err = pst.Fetch(testTable, &ec, "THIS IS NOT A VALID IDENT, BRAH")
	if assert.NotNil(t, err, "Expected an error") {
		assert.Equal(t, dbx.ErrNotFound, err)
	}

	count, err := pst.Count(`SELECT COUNT(*) FROM ` + testTable)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, 2, count)
	}

	var ef []*testEntity
	err = pst.Select(&ef, `SELECT {*} FROM `+testTable+` ORDER BY c`)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		if assert.Len(t, ef, 2) {
			assert.Equal(t, e1, ef[0])
			assert.Equal(t, e2, ef[1])
		}
	}

	err = pst.DeleteWithID(testTable, reflect.ValueOf((*testEntity)(nil)).Type(), e1.A)
	assert.Nil(t, err, fmt.Sprint(err))

	count, err = pst.Count(`SELECT COUNT(*) FROM ` + testTable)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, 1, count)
	}

	err = pst.Delete(testTable, e2)
	assert.Nil(t, err, fmt.Sprint(err))

	count, err = pst.Count(`SELECT COUNT(*) FROM ` + testTable)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, 0, count)
	}

}

type emptyEntity struct {
	Z string  `db:"z,pk"`
	A string  `db:"a,omitempty"`
	B int     `db:"b,omitempty"`
	C uint64  `db:"c,omitempty"`
	D bool    `db:"d,omitempty"`
	E float64 `db:"e,omitempty"`
	F []byte  `db:"f,omitempty"`
}

func TestPersistOmitEmpty(t *testing.T) {
	pst := New(test.DB(), entity.NewFieldMapper(), registry.New(), ident.AlphaNumeric(32))
	var err error

	e1 := &emptyEntity{}
	err = pst.Store(omitemptyTable, e1, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Len(t, e1.Z, 32)
	}

	var c1 emptyEntity
	err = pst.Fetch(omitemptyTable, &c1, e1.Z)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, "", c1.A)
		assert.Equal(t, int(0), c1.B)
		assert.Equal(t, uint64(0), c1.C)
		assert.Equal(t, false, c1.D)
		assert.Equal(t, float64(0), c1.E)
		assert.Equal(t, []byte(nil), c1.F)
	}

	e2 := &emptyEntity{
		A: "String here.",
		B: 999,
		C: 888,
		D: true,
		E: 77.77,
		F: []byte("And here"),
	}
	err = pst.Store(omitemptyTable, e2, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Len(t, e2.Z, 32)
	}

	var c2 emptyEntity
	err = pst.Fetch(omitemptyTable, &c2, e2.Z)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, "String here.", c2.A)
		assert.Equal(t, int(999), c2.B)
		assert.Equal(t, uint64(888), c2.C)
		assert.Equal(t, true, c2.D)
		assert.Equal(t, float64(77.77), c2.E)
		assert.Equal(t, []byte("And here"), c2.F)
	}

}
