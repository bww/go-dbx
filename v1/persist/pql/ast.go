package pql

import (
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

func (p literalNode) Exec(w io.Writer, cxt Context) error {
	_, err := w.Write([]byte(p.text))
	return err
}
