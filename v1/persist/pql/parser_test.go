package pql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLiteral(t *testing.T) {
	tests := []struct {
		Text   string
		Expect Node
		Error  error
	}{
		{
			"",
			literalNode{
				node: node{
					span: NewSpan("", 0, 0),
				},
				text: "",
			},
			nil,
		},
		{
			"{___",
			literalNode{
				node: node{
					span: NewSpan("{___", 0, 0),
				},
				text: "",
			},
			nil,
		},
		{
			"This is the text, ok.",
			literalNode{
				node: node{
					span: NewSpan("This is the text, ok.", 0, len("This is the text, ok.")),
				},
				text: "This is the text, ok.",
			},
			nil,
		},
		{
			"This is the { text, ok.",
			literalNode{
				node: node{
					span: NewSpan("This is the { text, ok.", 0, 12),
				},
				text: "This is the ",
			},
			nil,
		},
		{
			`This is the \{ text, ok.`,
			literalNode{
				node: node{
					span: NewSpan(`This is the \{ text, ok.`, 0, len(`This is the \{ text, ok.`)),
				},
				text: `This is the { text, ok.`,
			},
			nil,
		},
		{
			`This is the \\ text, ok.`,
			literalNode{
				node: node{
					span: NewSpan(`This is the \\ text, ok.`, 0, len(`This is the \\ text, ok.`)),
				},
				text: `This is the \ text, ok.`,
			},
			nil,
		},
		{
			`This is the \_ text, ok.`,
			literalNode{},
			newErr(ErrInvalidEscape, NewSpan(`This is the \_ text, ok.`, 13, 1)),
		},
	}
	for _, e := range tests {
		n, err := parseLiteral(NewScanner(e.Text))
		if e.Error != nil {
			fmt.Println(">>>", err)
			assert.Equal(t, e.Error, err, e.Text)
		} else if assert.Nil(t, err, fmt.Sprint(err)) {
			fmt.Printf(">>> [%s]\n", n.Span().Excerpt())
			assert.Equal(t, e.Expect, n)
		}
	}
}
