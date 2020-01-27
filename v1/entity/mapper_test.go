package entity

import (
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
		c, v := m.Columns(e.Entity)
		c, v = sortColumnsAndValues(c, v)
		assert.Equal(t, e.Columns, c)
		assert.Equal(t, e.Values, v)
	}
}
