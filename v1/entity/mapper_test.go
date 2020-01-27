package entity

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldMapper(t *testing.T) {
	m := NewFieldMapper("db")
	tests := []struct {
		Entity  interface{}
		Columns []string
		Values  []interface{}
	}{
		{
			testEntity{embedEntity{"BBB"}, "AAA"},
			[]string{"y", "z"},
			[]interface{}{"BBB", "AAA"},
		},
	}
	for _, e := range tests {
		c := m.Columns(e.Entity)
		sort.Sort(c)
		assert.Equal(t, e.Columns, c.Cols)
		assert.Equal(t, e.Values, c.Vals)
	}
}
