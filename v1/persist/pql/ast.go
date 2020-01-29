package pql

import (
	"fmt"
	"io"
	"strings"
)

type Context struct {
	Columns []string
}

type Node interface {
	Exec(io.Writer, Context) error
	Span() Span
}

type node struct {
	span Span
}

func newNode(t string, o, l int) node {
	return node{span: NewSpan(t, o, l)}
}

func (n node) Exec(w io.Writer, cxt Context) error {
	return nil // noop
}

func (n node) Span() Span {
	return n.span
}

type Program struct {
	sub []Node
}

func (p Program) Text(cxt Context) (string, error) {
	w := &strings.Builder{}
	err := p.Exec(w, cxt)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

func (p Program) Exec(w io.Writer, cxt Context) error {
	if len(p.sub) < 1 {
		return nil
	}
	for _, e := range p.sub {
		err := e.Exec(w, cxt)
		if err != nil {
			return err
		}
	}
	return nil
}

type literalNode struct {
	node
	text string
}

func (n literalNode) Exec(w io.Writer, cxt Context) error {
	_, err := w.Write([]byte(n.text))
	return err
}

type exprListNode struct {
	node
	sub []Node
}

func (n exprListNode) Exec(w io.Writer, cxt Context) error {
	if len(n.sub) < 1 {
		return nil
	}
	for i, e := range n.sub {
		if i > 0 {
			_, err := w.Write([]byte(", "))
			if err != nil {
				return err
			}
		}
		err := e.Exec(w, cxt)
		if err != nil {
			return err
		}
	}
	return nil
}

type exprLiteralNode struct {
	node
	prefix string
	name   string
}

func (n exprLiteralNode) Exec(w io.Writer, cxt Context) error {
	var t string
	if n.prefix != "" {
		t = fmt.Sprintf("%s.%s", n.prefix, n.name)
	} else {
		t = n.name
	}
	_, err := w.Write([]byte(t))
	if err != nil {
		return err
	}
	return nil
}

type exprMatchNode struct {
	node
	prefix string
}

func (n exprMatchNode) Exec(w io.Writer, cxt Context) error {
	for i, e := range cxt.Columns {
		if i > 0 {
			_, err := w.Write([]byte(", "))
			if err != nil {
				return err
			}
		}
		var t string
		if n.prefix != "" {
			t = fmt.Sprintf("%s.%s", n.prefix, e)
		} else {
			t = e
		}
		_, err := w.Write([]byte(t))
		if err != nil {
			return err
		}
	}
	return nil
}
