package pql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecAST(t *testing.T) {
	tests := []struct {
		Node    Node
		Context Context
		Expect  string
		Error   error
	}{
		{
			literalNode{
				node: newNode("This is the text, ok.", 0, len("This is the text, ok.")),
				text: "This is the text, ok.",
			},
			Context{},
			"This is the text, ok.",
			nil,
		},
		{
			exprLiteralNode{
				node:   newNode("{p.a}", 1, 3),
				prefix: "p",
				name:   "a",
			},
			Context{},
			"p.a",
			nil,
		},
		{
			exprMatchNode{
				node:   newNode("{p.*}", 1, 3),
				prefix: "p",
			},
			Context{},
			"",
			nil,
		},
		{
			exprMatchNode{
				node:   newNode("{p.*}", 1, 3),
				prefix: "p",
			},
			Context{
				Columns: []string{"a", "b", "c"},
			},
			"p.a, p.b, p.c",
			nil,
		},
		{
			exprListNode{
				node: newNode("{p.a, p.b}", 1, 8),
				sub: []Node{
					exprLiteralNode{
						node:   newNode("{p.a, p.b}", 1, 3),
						prefix: "p",
						name:   "a",
					},
					exprLiteralNode{
						node:   newNode("{p.a, p.b}", 6, 3),
						prefix: "p",
						name:   "b",
					},
				},
			},
			Context{},
			"p.a, p.b",
			nil,
		},
		{
			exprListNode{
				node: newNode("{p.a, x.*}", 1, 8),
				sub: []Node{
					exprLiteralNode{
						node:   newNode("{p.a, x.*}", 1, 3),
						prefix: "p",
						name:   "a",
					},
					exprMatchNode{
						node:   newNode("{p.a, x.*}", 6, 3),
						prefix: "x",
					},
				},
			},
			Context{
				Columns: []string{"a", "b", "c"},
			},
			"p.a, x.a, x.b, x.c",
			nil,
		},
	}
	for _, e := range tests {
		w := &strings.Builder{}
		err := e.Node.Exec(w, e.Context)
		if e.Error != nil {
			fmt.Println(">>>", err)
			assert.Equal(t, e.Error, err)
		} else if assert.Nil(t, err, fmt.Sprint(err)) {
			r := w.String()
			fmt.Println(">>>", r)
			assert.Equal(t, e.Expect, r)
		}
	}
}
