package pql

import (
	"fmt"
	"strings"
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
				node: newNode("{___", 0, 0),
				text: "",
			},
			nil,
		},
		{
			"This is the text, ok.",
			literalNode{
				node: newNode("This is the text, ok.", 0, len("This is the text, ok.")),
				text: "This is the text, ok.",
			},
			nil,
		},
		{
			"This is the { text, ok.",
			literalNode{
				node: newNode("This is the { text, ok.", 0, 12),
				text: "This is the ",
			},
			nil,
		},
		{
			`This is the \{ text, ok.`,
			literalNode{
				node: newNode(`This is the \{ text, ok.`, 0, len(`This is the \{ text, ok.`)),
				text: `This is the { text, ok.`,
			},
			nil,
		},
		{
			`This is the \\ text, ok.`,
			literalNode{
				node: newNode(`This is the \\ text, ok.`, 0, len(`This is the \\ text, ok.`)),
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
			fmt.Println("-->", err)
			assert.Equal(t, e.Error, err, e.Text)
		} else if assert.Nil(t, err, fmt.Sprint(err)) {
			fmt.Printf("--> [%s]\n", n.Span().Excerpt())
			assert.Equal(t, e.Expect, n)
		}
	}
}

func TestParseExpr(t *testing.T) {
	tests := []struct {
		Text    string
		Expect  Node
		Error   error
		Context Context
		Output  string
	}{
		{
			`{p}`,
			exprListNode{
				node: newNode(`{p}`, 1, 1),
				sub: []Node{
					exprLiteralNode{
						node:   newNode(`{p}`, 1, 1),
						prefix: "",
						name:   "p",
					},
				},
			},
			nil,
			Context{},
			"p",
		},
		{
			`{p, p.*}`,
			exprListNode{
				node: newNode(`{p, p.*}`, 1, 6),
				sub: []Node{
					exprLiteralNode{
						node:   newNode(`{p, p.*}`, 1, 1),
						prefix: "",
						name:   "p",
					},
					exprMatchNode{
						node:   newNode(`{p, p.*}`, 4, 3),
						prefix: "p",
					},
				},
			},
			nil,
			Context{
				Columns: []string{"A", "B"},
			},
			"p, p.A, p.B",
		},
		{
			`{ p  ,  p . * }`,
			exprListNode{
				node: newNode(`{ p  ,  p . * }`, 1, 13),
				sub: []Node{
					exprLiteralNode{
						node:   newNode(`{ p  ,  p . * }`, 2, 3),
						prefix: "",
						name:   "p",
					},
					exprMatchNode{
						node:   newNode(`{ p  ,  p . * }`, 8, 6),
						prefix: "p",
					},
				},
			},
			nil,
			Context{
				Columns: []string{"A", "B"},
			},
			"p, p.A, p.B",
		},
	}
	for _, e := range tests {
		fmt.Println(">>>", e.Text)
		n, err := parseMeta(NewScanner(e.Text))
		if e.Error != nil {
			fmt.Println("-->", err)
			assert.Equal(t, e.Error, err, e.Text)
		} else if assert.Nil(t, err, fmt.Sprint(err)) {
			fmt.Println("-->", n.Span().Describe())
			assert.Equal(t, e.Expect, n)
			w := &strings.Builder{}
			err := n.Exec(w, e.Context)
			if assert.Nil(t, err, fmt.Sprint(err)) {
				r := w.String()
				fmt.Println("-->", r)
				assert.Equal(t, e.Output, r)
			}
		}
	}
}
