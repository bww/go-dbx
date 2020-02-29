package filter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFilterOptions(t *testing.T) {
	now, later := time.Now(), time.Now().Add(time.Hour)
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
				WithLimit(Range{2, 20}),
			},
			Filter{
				Limit: Range{2, 20},
			},
		},
		{
			[]Option{
				WithLimit(Range{2, 20}),
				WithOrder(Ascending),
			},
			Filter{
				Limit: Range{2, 20},
				Order: Ascending,
			},
		},
		{
			[]Option{
				WithTimeframe(NewTimeframe(now, later)),
			},
			Filter{
				Timeframe: NewTimeframe(now, later),
			},
		},
		{
			[]Option{
				WithTimeframe(NewTimeframe(now, later)),
				WithOrder(Descending),
			},
			Filter{
				Timeframe: NewTimeframe(now, later),
				Order:     Descending,
			},
		},
	}
	for _, e := range tests {
		f := New(e.Options...)
		assert.Equal(t, e.Expect, f)
	}
}
