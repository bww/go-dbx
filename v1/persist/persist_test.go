package persist

import (
	"fmt"
	"os"
	"testing"

	"github.com/bww/go-dbx/v1/entity"
	"github.com/bww/go-dbx/v1/test"
	"github.com/bww/go-util/env"
	"github.com/bww/go-util/urls"
	"github.com/stretchr/testify/assert"
)

const testTable = "dbx_v1_persist_test"

type embedEntity struct {
	B string `db:"b"`
}

type testEntity struct {
	embedEntity
	A string `db:"a,pk"`
}

func TestMain(m *testing.M) {
	test.Init(testTable, test.WithMigrations(urls.File(env.Etc("migrations"))))
	os.Exit(m.Run())
}

func TestPersist(t *testing.T) {
	var err error

	gen := entity.NewGenerator(entity.NewFieldMapper("db"))
	pst := New(test.DB(), Random)

	_ = gen

	ea := &testEntity{
		embedEntity: embedEntity{
			B: "This is the value of B",
		},
	}

	err = pst.Store(testTable, ea, nil, nil)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Len(t, ea.A, 32)
	}

}
