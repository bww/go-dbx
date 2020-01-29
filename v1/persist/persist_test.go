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

	ea := &testEntity{
		DontUseThisExportedEntity: &DontUseThisExportedEntity{
			B: "This is the value of B",
		},
		C: 999,
	}

	err = pst.Store(testTable, ea, nil, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Len(t, ea.A, 32)
	}

	var ca, cb string
	var cc int
	err = db.QueryRow(`SELECT a, b, c FROM `+testTable+` WHERE a = $1`, ea.A).Scan(&ca, &cb, &cc)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, ea.A, ca)
		assert.Equal(t, ea.B, cb)
		assert.Equal(t, ea.C, cc)
	}

	var eb testEntity
	err = db.QueryRowx(`SELECT a, b, c FROM `+testTable+` WHERE a = $1`, ea.A).StructScan(&eb)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, ea.A, eb.A)
		assert.Equal(t, ea.B, eb.B)
		assert.Equal(t, ea.C, eb.C)
	}

	var ec testEntity
	err = pst.Fetch(testTable, &ec, ea.A, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, ea.A, ec.A)
		assert.Equal(t, ea.B, ec.B)
		assert.Equal(t, ea.C, ec.C)
	}

}
