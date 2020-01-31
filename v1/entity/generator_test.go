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
	gen := &Generator{NewFieldMapper(), true}
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
			"UPDATE some_table SET y = $1, z = $2 WHERE z = $3",
			[]interface{}{"BBB", "AAA", "AAA"},
		},
		{
			testEntity{embedEntity{"BBB"}, "AAA"},
			"some_table",
			nil,
			"UPDATE some_table SET y = $1, z = $2 WHERE z = $3",
			[]interface{}{"BBB", "AAA", "AAA"},
		},
		{
			testEntity{embedEntity{"BBB"}, "AAA"},
			"some_table",
			[]string{"z"},
			"UPDATE some_table SET z = $1 WHERE z = $2",
			[]interface{}{"AAA", "AAA"},
		},
		{
			multiPKEntity{embedEntity{"BBB"}, "AAA", "CCC"},
			"some_table",
			[]string{"z"},
			"UPDATE some_table SET z = $1 WHERE x = $2 AND z = $3",
			[]interface{}{"AAA", "CCC", "AAA"},
		},
	}
	gen := &Generator{NewFieldMapper(), true}
	for _, e := range tests {
		sql, args := gen.Update(e.Table, e.Entity, e.Columns)
		fmt.Println("-->", sql)
		assert.Equal(t, e.SQL, sql)
		assert.Equal(t, e.Args, args)
	}
}

func TestGeneratorSelect(t *testing.T) {
	tests := []struct {
		Entity interface{}
		Table  string
		Keys   *Columns
		SQL    string
		Args   []interface{}
	}{
		{
			testEntity{},
			"some_table",
			&Columns{
				Cols: []string{"z"},
				Vals: []interface{}{"AAA"},
			},
			"SELECT y, z FROM some_table WHERE z = $1",
			[]interface{}{"AAA"},
		},
		{
			multiPKEntity{},
			"some_table",
			&Columns{
				Cols: []string{"z", "x"},
				Vals: []interface{}{"AAA", "CCC"},
			},
			"SELECT x, y, z FROM some_table WHERE x = $1 AND z = $2",
			[]interface{}{"CCC", "AAA"},
		},
	}
	gen := &Generator{NewFieldMapper(), true}
	for _, e := range tests {
		sql, args := gen.Select(e.Table, e.Entity, e.Keys)
		fmt.Println("-->", sql)
		assert.Equal(t, e.SQL, sql)
		assert.Equal(t, e.Args, args)
	}
}

func TestGeneratorDelete(t *testing.T) {
	tests := []struct {
		Entity interface{}
		Table  string
		Keys   *Columns
		SQL    string
		Args   []interface{}
	}{
		{
			testEntity{},
			"some_table",
			&Columns{
				Cols: []string{"z"},
				Vals: []interface{}{"AAA"},
			},
			"DELETE FROM some_table WHERE z = $1",
			[]interface{}{"AAA"},
		},
		{
			multiPKEntity{},
			"some_table",
			&Columns{
				Cols: []string{"z", "x"},
				Vals: []interface{}{"AAA", "CCC"},
			},
			"DELETE FROM some_table WHERE x = $1 AND z = $2",
			[]interface{}{"CCC", "AAA"},
		},
	}
	gen := &Generator{NewFieldMapper(), true}
	for _, e := range tests {
		sql, args := gen.Delete(e.Table, e.Entity, e.Keys)
		fmt.Println("-->", sql)
		assert.Equal(t, e.SQL, sql)
		assert.Equal(t, e.Args, args)
	}
}
