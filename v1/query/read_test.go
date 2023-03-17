package query

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigOptions(t *testing.T) {
	now, later := time.Now(), time.Now().Add(time.Hour)
	tests := []struct {
		ReadOptions []ReadOption
		Expect      ReadConfig
	}{
		{
			[]ReadOption{},
			ReadConfig{},
		},
		{
			[]ReadOption{
				WithLimit(Range{2, 20}),
			},
			ReadConfig{
				Limit: Range{2, 20},
			},
		},
		{
			[]ReadOption{
				WithLimit(Range{2, 20}),
				WithOrder(Ascending),
			},
			ReadConfig{
				Limit: Range{2, 20},
				Order: Ascending,
			},
		},
		{
			[]ReadOption{
				WithTimeframe(NewTimeframe(now, later)),
			},
			ReadConfig{
				Timeframe: NewTimeframe(now, later),
			},
		},
		{
			[]ReadOption{
				WithTimeframe(NewTimeframe(now, later)),
				WithOrder(Descending),
			},
			ReadConfig{
				Timeframe: NewTimeframe(now, later),
				Order:     Descending,
			},
		},
		{
			[]ReadOption{
				UseReadConfig(ReadConfig{Limit: Range{100, 101}}),
			},
			ReadConfig{
				Limit: Range{100, 101},
			},
		},
	}
	for _, e := range tests {
		f := New(e.ReadOptions)
		assert.Equal(t, e.Expect, f)
	}
}
