package persist

import (
	"fmt"
	"os"
	"testing"

	"github.com/bww/go-dbx/v1/entity"
	"github.com/bww/go-dbx/v1/persist/ident"
	"github.com/bww/go-dbx/v1/test"
	"github.com/bww/go-util/env"
	"github.com/bww/go-util/urls"
	"github.com/stretchr/testify/assert"
)

const testTable = "dbx_v1_persist_test"

type DontUseThisExportedEntity struct {
	B string `db:"b"`
}

type testEntity struct {
	*DontUseThisExportedEntity
	A string `db:"a,pk"`
	C int    `db:"c"`
}

func TestMain(m *testing.M) {
	test.Init(testTable, test.WithMigrations(urls.File(env.Etc("migrations"))))
	os.Exit(m.Run())
}

func TestPersist(t *testing.T) {
	db := test.DB()
	pst := New(db, entity.NewFieldMapper(), ident.AlphaNumeric(32))
	var err error

	e1 := &testEntity{
		DontUseThisExportedEntity: &DontUseThisExportedEntity{
			B: "This is the value of B",
		},
		C: 888,
	}

	err = pst.Store(testTable, e1, nil, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Len(t, e1.A, 32)
	}

	e2 := &testEntity{
		DontUseThisExportedEntity: &DontUseThisExportedEntity{
			B: "Never is this the value of B",
		},
		C: 999,
	}

	err = pst.Store(testTable, e2, nil, nil)
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
	err = db.QueryRowx(`SELECT a, b, c FROM `+testTable+` WHERE a = $1`, e1.A).StructScan(&eb)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, e1.A, eb.A)
		assert.Equal(t, e1.B, eb.B)
		assert.Equal(t, e1.C, eb.C)
	}

	var ec testEntity
	err = pst.Fetch(testTable, &ec, e1.A, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, e1.A, ec.A)
		assert.Equal(t, e1.B, ec.B)
		assert.Equal(t, e1.C, ec.C)
	}

	var ed []*testEntity
	err = pst.Select(testTable, &ed, nil, `SELECT {*} FROM `+testTable+` ORDER BY c`)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		if assert.Len(t, ed, 2) {
			assert.Equal(t, e1, ed[0])
			assert.Equal(t, e2, ed[1])
		}
	}

}
