package persist

import (
	// "fmt"
	"os"
	"testing"

	// "github.com/bww/go-dbx/v1/entity"
	"github.com/bww/go-dbx/v1/test"
	"github.com/bww/go-util/env"
	"github.com/bww/go-util/urls"
	// "github.com/stretchr/testify/assert"
)

type embedEntity struct {
	B string `db:"y"`
}

type testEntity struct {
	embedEntity
	A string `db:"z,pk"`
}

func TestMain(m *testing.M) {
	test.Init("dbx_v1_persist_test", test.WithMigrations(urls.File(env.Etc("db"))))
	os.Exit(m.Run())
}

func TestPersist(t *testing.T) {
	// tests := []struct {
	// 	Entity interface{}
	// 	Table  string
	// 	SQL    string
	// 	Args   []interface{}
	// }{
	// 	{
	// 		testEntity{embedEntity{"BBB"}, "AAA"},
	// 		"some_table",
	// 		"INSERT INTO some_table (y, z) VALUES ($1, $2)",
	// 		[]interface{}{"BBB", "AAA"},
	// 	},
	// }
	// gen := &Generator{NewFieldMapper("db"), true}
	// for _, e := range tests {
	// 	sql, args := gen.Insert(e.Table, e.Entity)
	// 	fmt.Println("-->", sql)
	// 	assert.Equal(t, e.SQL, sql)
	// 	assert.Equal(t, e.Args, args)
	// }
}
