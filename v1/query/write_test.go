package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteConfigOptions(t *testing.T) {
	tests := []struct {
		WriteOptions []WriteOption
		Expect       WriteConfig
	}{
		{
			[]WriteOption{},
			WriteConfig{},
		},
		{
			[]WriteOption{
				WithCascade(),
			},
			WriteConfig{
				Cascade: true,
			},
		},
		{
			[]WriteOption{
				WithCascade(),
			},
			WriteConfig{
				Cascade: true,
			},
		},
		{
			[]WriteOption{
				WithWriteParams(Params{"A": "1", "B": "2"}),
			},
			WriteConfig{
				Params: Params{"A": "1", "B": "2"},
			},
		},
		{
			[]WriteOption{
				UseWriteConfig(WriteConfig{
					Cascade: true,
					Params:  Params{"A": "1", "B": "2"},
				}),
			},
			WriteConfig{
				Cascade: true,
				Params:  Params{"A": "1", "B": "2"},
			},
		},
	}
	for _, e := range tests {
		c := WriteConfig{}.WithOptions(e.WriteOptions)
		assert.Equal(t, e.Expect, c)
	}
}
