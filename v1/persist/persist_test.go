package persist

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
	"github.com/bww/go-dbx/v1/persist/ident"
	"github.com/bww/go-dbx/v1/persist/registry"
	"github.com/bww/go-dbx/v1/test"
	"github.com/bww/go-util/v1/env"
	"github.com/bww/go-util/v1/ulid"
	"github.com/bww/go-util/v1/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.Init(testDB, test.WithMigrations(env.Etc("migrations")))
	os.Exit(m.Run())
}

const (
	testDB      = "dbx_v1_persist_test"
	firstTable  = "first_entity"
	secondTable = "second_entity"
	thirdTable  = "third_entity"
	fourthTable = "fourth_entity"
)

type DontUseThisTestEntity struct {
	B string `db:"b"`
}

type firstEntity struct {
	*DontUseThisTestEntity
	A string `db:"a,pk"`
	C int    `db:"c"`
	E int    `db:"e,omitempty"`
	D []*secondEntity
}

type secondEntity struct {
	X string `db:"x,pk"`
	Z int    `db:"z"`
}

type firstPersister struct{}

func (p *firstPersister) FetchRelated(pst Persister, ent interface{}) error {
	z := ent.(*firstEntity)
	q := `
    SELECT {a.*} FROM ` + secondTable + ` AS a
    INNER JOIN first_entity_r_second_entity AS r ON r.x = a.x
    WHERE r.a = $1
    ORDER BY a.z`

	var another []*secondEntity
	err := pst.Select(&another, q, z.A)
	if err != nil {
		return err
	}

	z.D = another
	return nil
}

func (p *firstPersister) StoreRelated(pst Persister, ent interface{}) error {
	z := ent.(*firstEntity)
	for _, e := range z.D {
		err := pst.Store("second_entity", e, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *firstPersister) StoreReferences(pst Persister, ent interface{}) error {
	z := ent.(*firstEntity)
	for _, e := range z.D {
		_, err := pst.Exec(`INSERT INTO first_entity_r_second_entity (a, x) VALUES ($1, $2)`, z.A, e.X)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *firstPersister) DeleteRelated(pst Persister, ent interface{}) error {
	z := ent.(*firstEntity)
	for _, e := range z.D {
		err := pst.Delete("second_entity", e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *firstPersister) DeleteReferences(pst Persister, ent interface{}) error {
	z := ent.(*firstEntity)
	_, err := pst.Exec(`DELETE FROM first_entity_r_second_entity WHERE a = $1`, z.A)
	if err != nil {
		return err
	}
	return nil
}

func TestPersist(t *testing.T) {
	var err error

	db := test.DB()
	reg := registry.New()
	pst := New(db, entity.NewFieldMapper(), reg, ident.AlphaNumeric(32)).WithOptions(Cascade(true))

	reg.Set(reflect.ValueOf((*firstEntity)(nil)).Type(), &firstPersister{})

	e1 := &firstEntity{
		DontUseThisTestEntity: &DontUseThisTestEntity{
			B: "This is the value of B",
		},
		C: 888,
	}

	err = pst.Store(firstTable, e1, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Len(t, e1.A, 32)
	}

	e2 := &firstEntity{
		DontUseThisTestEntity: &DontUseThisTestEntity{
			B: "Never is this the value of B",
		},
		C: 999,
		E: 111,
		D: []*secondEntity{
			{Z: 111},
			{Z: 222},
			{Z: 333},
		},
	}

	err = pst.Store(firstTable, e2, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Len(t, e2.A, 32)
	}

	var ca, cb string
	var cc int
	err = db.QueryRow(`SELECT a, b, c FROM `+firstTable+` WHERE a = $1`, e1.A).Scan(&ca, &cb, &cc)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, e1.A, ca)
		assert.Equal(t, e1.B, cb)
		assert.Equal(t, e1.C, cc)
	}

	var eb firstEntity
	err = db.QueryRowx(`SELECT a, b, c, e FROM `+firstTable+` WHERE a = $1`, e2.A).StructScan(&eb)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, e2.A, eb.A)
		assert.Equal(t, e2.B, eb.B)
		assert.Equal(t, e2.C, eb.C)
		assert.Equal(t, e2.E, eb.E)
	}

	var ec firstEntity
	err = pst.Fetch(firstTable, &ec, e1.A)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, e1.A, ec.A)
		assert.Equal(t, e1.B, ec.B)
		assert.Equal(t, e1.C, ec.C)
		assert.Equal(t, e1.E, ec.E)
	}

	var ed firstEntity
	err = pst.Fetch(firstTable, &ed, e2.A)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, e2.A, ed.A)
		assert.Equal(t, e2.B, ed.B)
		assert.Equal(t, e2.C, ed.C)
		assert.Equal(t, e2.E, ed.E)
	}

	err = pst.Fetch(firstTable, &ec, "THIS IS NOT A VALID IDENT, BRAH")
	if assert.NotNil(t, err, "Expected an error") {
		assert.Equal(t, dbx.ErrNotFound, err)
	}

	err = pst.Select(&ec, `SELECT {*} FROM `+firstTable+` WHERE a = 'THIS IS NOT A VALID IDENT, BRAH'`)
	if assert.NotNil(t, err, "Expected an error") {
		assert.Equal(t, dbx.ErrNotFound, err)
	}

	count, err := pst.Count(`SELECT COUNT(*) FROM ` + firstTable)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, 2, count)
	}

	var ef []*firstEntity
	err = pst.Select(&ef, `SELECT {*} FROM `+firstTable+` ORDER BY c`)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		if assert.Len(t, ef, 2) {
			assert.Equal(t, e1, ef[0])
			assert.Equal(t, e2, ef[1])
		}
	}

	err = pst.DeleteWithID(firstTable, reflect.ValueOf((*firstEntity)(nil)).Type(), e1.A)
	assert.Nil(t, err, fmt.Sprint(err))

	count, err = pst.Count(`SELECT COUNT(*) FROM ` + firstTable)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, 1, count)
	}

	err = pst.Delete(firstTable, e2)
	assert.Nil(t, err, fmt.Sprint(err))

	count, err = pst.Count(`SELECT COUNT(*) FROM ` + firstTable)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, 0, count)
	}

}

func TestPersistOptions(t *testing.T) {
	db := test.DB()

	var err error
	var (
		reg      = registry.New()
		pst      = New(db, entity.NewFieldMapper(), reg, ident.AlphaNumeric(32)).WithOptions(Cascade(true))
		nofetch  = pst.WithOptions(FetchRelated(false))
		nostore  = pst.WithOptions(StoreRelated(false))
		nodelete = pst.WithOptions(DeleteRelated(false))
	)

	reg.Set(reflect.ValueOf((*firstEntity)(nil)).Type(), &firstPersister{})

	e1 := &firstEntity{
		DontUseThisTestEntity: &DontUseThisTestEntity{
			B: "This is the value of B",
		},
		C: 888,
		D: []*secondEntity{
			&secondEntity{
				Z: 111,
			},
			&secondEntity{
				Z: 222,
			},
		},
	}

	err = nostore.Store(firstTable, e1, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		eR := &firstEntity{}
		err = nostore.Fetch(firstTable, eR, e1.A)
		if assert.Nil(t, err, fmt.Sprint(err)) {
			assert.Equal(t, &firstEntity{
				DontUseThisTestEntity: e1.DontUseThisTestEntity,
				A:                     e1.A,
				C:                     e1.C,
			}, eR)
		}
	}

	err = pst.Store(firstTable, e1, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		eR := &firstEntity{}
		err = nofetch.Fetch(firstTable, eR, e1.A)
		if assert.Nil(t, err, fmt.Sprint(err)) {
			assert.Equal(t, &firstEntity{
				DontUseThisTestEntity: e1.DontUseThisTestEntity,
				A:                     e1.A,
				C:                     e1.C,
			}, eR)
		}
	}

	e2 := &firstEntity{
		DontUseThisTestEntity: &DontUseThisTestEntity{
			B: "You ain't never seen a B like this",
		},
		C: 456,
		D: []*secondEntity{
			&secondEntity{
				Z: 333,
			},
			&secondEntity{
				Z: 444,
			},
		},
	}

	err = pst.Store(firstTable, e2, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		eR := &firstEntity{}
		err = pst.Fetch(firstTable, eR, e2.A)
		if assert.Nil(t, err, fmt.Sprint(err)) {
			assert.Equal(t, e2, eR)
		}
	}

	err = nodelete.Delete(firstTable, e2)
	assert.NotNil(t, err, "Expected consistency error")
}

type thirdEntity struct {
	Z string    `db:"z,pk"`
	A string    `db:"a,omitempty"`
	B int       `db:"b,omitempty"`
	C uint64    `db:"c,omitempty"`
	D bool      `db:"d,omitempty"`
	E float64   `db:"e,omitempty"`
	F []byte    `db:"f,omitempty"`
	G time.Time `db:"g,omitempty"`
	H uuid.UUID `db:"h,omitempty"`
	I ulid.ULID `db:"i,omitempty"`
}

func TestPersistOmitEmpty(t *testing.T) {
	db := test.DB()
	pst := New(db, entity.NewFieldMapper(), registry.New(), ident.AlphaNumeric(32))
	var err error

	e1 := &thirdEntity{}
	err = pst.Store(thirdTable, e1, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Len(t, e1.Z, 32)
	}

	var c1 thirdEntity
	err = pst.Fetch(thirdTable, &c1, e1.Z)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, "", c1.A)
		assert.Equal(t, int(0), c1.B)
		assert.Equal(t, uint64(0), c1.C)
		assert.Equal(t, false, c1.D)
		assert.Equal(t, float64(0), c1.E)
		assert.Equal(t, []byte(nil), c1.F)
		assert.Equal(t, time.Time{}, c1.G)
		assert.Equal(t, uuid.Zero, c1.H)
		assert.Equal(t, ulid.Zero, c1.I)
	}

	x1 := struct {
		Z                         string
		A, B, C, D, E, F, G, H, I interface{}
	}{}
	r1 := []interface{}{
		&x1.Z,
		&x1.A,
		&x1.B,
		&x1.C,
		&x1.D,
		&x1.E,
		&x1.F,
		&x1.G,
		&x1.H,
		&x1.I,
	}
	err = db.QueryRow(`SELECT z, a, b, c, d, e, f, g, h, i FROM `+thirdTable+` WHERE z = $1`, e1.Z).Scan(r1...)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, e1.Z, x1.Z)
		assert.Equal(t, nil, x1.A)
		assert.Equal(t, nil, x1.B)
		assert.Equal(t, nil, x1.C)
		assert.Equal(t, nil, x1.D)
		assert.Equal(t, nil, x1.E)
		assert.Equal(t, nil, x1.F)
		assert.Equal(t, nil, x1.G)
		assert.Equal(t, nil, x1.H)
		assert.Equal(t, nil, x1.I)
	}

	t1 := time.Now().UTC().Truncate(time.Millisecond)
	u1 := uuid.New()
	l1 := ulid.New()
	e2 := &thirdEntity{
		A: "String here.",
		B: 999,
		C: 888,
		D: true,
		E: 77.77,
		F: []byte("And here"),
		G: t1,
		H: u1,
		I: l1,
	}
	err = pst.Store(thirdTable, e2, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Len(t, e2.Z, 32)
	}

	var c2 thirdEntity
	err = pst.Fetch(thirdTable, &c2, e2.Z)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, "String here.", c2.A)
		assert.Equal(t, int(999), c2.B)
		assert.Equal(t, uint64(888), c2.C)
		assert.Equal(t, true, c2.D)
		assert.Equal(t, float64(77.77), c2.E)
		assert.Equal(t, []byte("And here"), c2.F)
		assert.Equal(t, t1, c2.G)
		assert.Equal(t, u1, c2.H)
		assert.Equal(t, l1, c2.I)
	}

}

func TestPersistFetchAndStoreInATightLoop(t *testing.T) {
	db := test.DB()
	pst := New(db, entity.NewFieldMapper(), registry.New(), ident.AlphaNumeric(32))
	var err error

	for i := 0; i < 10; i++ {
		var e2 secondEntity
		err = pst.Select(&e2, `SELECT {*} FROM `+fourthTable+` WHERE x = $1`, fmt.Sprint(i))
		if assert.NotNil(t, err, "Expected an error") {
			assert.Equal(t, dbx.ErrNotFound, err)
		}
		e1 := &secondEntity{
			X: fmt.Sprint(i),
			Z: i,
		}
		_, err = pst.Exec(`INSERT INTO `+fourthTable+` (x, z) VALUES ($1, $2) ON CONFLICT (x) DO UPDATE SET z = $2`, e1.X, e1.Z)
		assert.Nil(t, err, fmt.Sprint(err))
	}

}

func TestInvalidParamInSelectOneDoesntLeakConns(t *testing.T) {
	db := test.DB()
	pst := New(db, entity.NewFieldMapper(), registry.New(), ident.AlphaNumeric(32))
	var err error

	for i := 0; i < 10; i++ {
		var invalid int // invalid destination type
		err = pst.Select(&invalid, `SELECT {*} FROM `+fourthTable+` WHERE x = $1`, fmt.Sprint(i))
		assert.NotNil(t, err, "Expected an error")
	}

	// if we don't hang forever, we have succeeded
}
