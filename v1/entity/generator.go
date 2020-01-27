package entity

import (
	"sort"
	"strconv"
	"strings"
)

type Generator struct {
	mapper *FieldMapper
	sorted bool
}

func NewGenerator(m *FieldMapper) *Generator {
	return &Generator{m, false}
}

func (g *Generator) Insert(table string, entity interface{}) (string, []interface{}) {
	cols := g.mapper.Columns(entity)
	if g.sorted {
		sort.Sort(cols)
	}

	b := &strings.Builder{}
	b.WriteString("INSERT INTO ")
	b.WriteString(table)
	b.WriteString(" (")

	for i, e := range cols.Cols {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e)
	}

	b.WriteString(") VALUES (")

	for i, _ := range cols.Vals {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("$")
		b.WriteString(strconv.FormatInt(int64(i+1), 10))
	}

	b.WriteString(")")

	return b.String(), cols.Vals
}

func (g *Generator) Update(table string, names []string, entity interface{}) (string, []interface{}) {
	cols := g.mapper.Columns(entity)
	if g.sorted {
		sort.Sort(cols)
	}

	var incl map[string]struct{}
	if len(names) > 0 {
		incl = make(map[string]struct{})
		for _, e := range names {
			incl[e] = struct{}{}
		}
	}

	b := &strings.Builder{}
	b.WriteString("UPDATE ")
	b.WriteString(table)
	b.WriteString(" SET ")

	var args []interface{}
	if incl != nil {
		args = make([]interface{}, 0, len(incl))
	}

	for i, e := range cols.Cols {
		var x int
		if incl != nil {
			if _, ok := incl[e]; !ok {
				continue
			}
			x = len(args)
		} else {
			x = i
		}
		if x > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e)
		b.WriteString(" = $")
		b.WriteString(strconv.FormatInt(int64(x+1), 10))
		if args != nil {
			args = append(args, cols.Vals[i])
		}
	}

	b.WriteString(" ")

	if args == nil {
		args = cols.Vals
	}
	return b.String(), args
}
