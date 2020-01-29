package pql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	tests := []struct {
		Text      string
		Next      rune
		SkipWhite string
	}{
		{
			"Hello.", 'H', "ello.",
		},
		{
			"H    ello.", 'H', "ello.",
		},
	}
	for _, e := range tests {
		s := NewScanner(e.Text)
		assert.Equal(t, e.Next, s.Next())
		s.SkipWhite()
		assert.Equal(t, e.SkipWhite, s.Text())
	}
}
