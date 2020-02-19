package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterOptions(t *testing.T) {
	tests := []struct {
		Options []Option
		Expect  Filter
	}{
		{
			[]Option{},
			Filter{},
		},
		{
			[]Option{
				WithRange(Range{2, 20}),
			},
			Filter{
				Limit: Range{2, 20},
			},
		},
	}
	for _, e := range tests {
		f := New(e.Options...)
		assert.Equal(t, e.Expect, f)
	}
}
