package pql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScannerNext(t *testing.T) {
	tests := []struct {
		Text string
		Next []rune
	}{
		{
			"Hello.", []rune{'H', 'e', 'l', 'l', 'o', '.'},
		},
	}
	for _, e := range tests {
		s := NewScanner(e.Text)
		fmt.Println(">>>", e.Text)
		for i, x := range e.Next {
			if i%2 == 1 {
				assert.Equal(t, x, s.Peek())
			}
			assert.Equal(t, x, s.Next())
			s.Backup()
			assert.Equal(t, x, s.Peek())
			assert.Equal(t, x, s.Next())
		}
	}
}

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
