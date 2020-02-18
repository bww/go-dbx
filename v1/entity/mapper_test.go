package entity

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldMapper(t *testing.T) {
	m := NewFieldMapper()
	tests := []struct {
		Entity  interface{}
		Keys    []string
		Columns []string
		Values  []interface{}
	}{
		{
			testEntity{embedEntity{"BBB"}, "AAA", 999, 0},
			[]string{"z"},
			[]string{"e", "y", "z"},
			[]interface{}{nil, "BBB", "AAA"},
		},
		{
			testEntity{embedEntity{"BBB"}, "AAA", 999, 111},
			[]string{"z"},
			[]string{"e", "y", "z"},
			[]interface{}{111, "BBB", "AAA"},
		},
	}
	for _, e := range tests {
		k, c := m.Columns(e.Entity)
		sort.Sort(k)
		sort.Sort(c)
		assert.Equal(t, e.Keys, k.Cols)
		assert.Equal(t, e.Columns, c.Cols)
		assert.Equal(t, e.Values, c.Vals)
	}
}
