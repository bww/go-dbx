package pql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPQL(t *testing.T) {

	s := "This here is the literal text"
	l := literalNode{
		node: node{
			span: NewSpan(s, 0, len(s)),
		},
		text: s,
	}

	p := Program{
		sub: []Node{l},
	}
	r, err := p.Text(Context{})
	if assert.Nil(t, err, fmt.Sprint(err)) {
		fmt.Println(">>>", r)
	}

}
