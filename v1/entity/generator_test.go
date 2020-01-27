package entity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratorInsert(t *testing.T) {
	tests := []struct {
		Entity interface{}
		Table  string
		SQL    string
		Args   []interface{}
	}{
		{
			testEntity{embedEntity{"BBB"}, "AAA"},
			"some_table",
			"INSERT INTO some_table (y, z) VALUES ($1, $2)",
			[]interface{}{"BBB", "AAA"},
		},
	}
	gen := &Generator{NewFieldMapper("db"), true}
	for _, e := range tests {
		sql, args := gen.Insert(e.Table, e.Entity)
		fmt.Println("-->", sql)
		assert.Equal(t, e.SQL, sql)
		assert.Equal(t, e.Args, args)
	}
}

func TestGeneratorUpdate(t *testing.T) {
	tests := []struct {
		Entity  interface{}
		Table   string
		Columns []string
		SQL     string
		Args    []interface{}
	}{
		{
			testEntity{embedEntity{"BBB"}, "AAA"},
			"some_table",
			[]string{"z", "y"},
			"UPDATE some_table SET y = $1, z = $2 ",
			[]interface{}{"BBB", "AAA"},
		},
		{
			testEntity{embedEntity{"BBB"}, "AAA"},
			"some_table",
			nil,
			"UPDATE some_table SET y = $1, z = $2 ",
			[]interface{}{"BBB", "AAA"},
		},
		{
			testEntity{embedEntity{"BBB"}, "AAA"},
			"some_table",
			[]string{"z"},
			"UPDATE some_table SET z = $1 ",
			[]interface{}{"AAA"},
		},
	}
	gen := &Generator{NewFieldMapper("db"), true}
	for _, e := range tests {
		sql, args := gen.Update(e.Table, e.Columns, e.Entity)
		fmt.Println("-->", sql)
		assert.Equal(t, e.SQL, sql)
		assert.Equal(t, e.Args, args)
	}
}
