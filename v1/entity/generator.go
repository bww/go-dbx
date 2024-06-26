package entity

import (
	"sort"
	"strconv"
	"strings"
)

type Generator struct {
	fm     *FieldMapper
	sorted bool
}

func NewGenerator(m *FieldMapper) *Generator {
	return &Generator{m, false}
}

func (g *Generator) Select(table string, entity interface{}, keys *Columns) (string, []interface{}) {
	_, cols := g.fm.Columns(entity)
	if g.sorted {
		sort.Sort(cols)
		sort.Sort(keys)
	}

	var n, x int
	args := make([]interface{}, 0, len(keys.Vals))

	b := &strings.Builder{}
	b.WriteString("SELECT ")

	n = 0
	for i, e := range cols.Cols {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e)
	}

	b.WriteString(" FROM ")
	b.WriteString(table)
	b.WriteString(" WHERE ")

	n = 0
	for i, e := range keys.Cols {
		if n > 0 {
			b.WriteString(" AND ")
		}
		b.WriteString(e)
		b.WriteString(" = $")
		b.WriteString(strconv.FormatInt(int64(x+1), 10))
		args = append(args, keys.Vals[i])
		x++
		n++
	}

	return b.String(), args
}

func (g *Generator) Insert(table string, entity interface{}) (string, []interface{}) {
	_, cols := g.fm.Columns(entity)
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

func (g *Generator) Upsert(table string, entity interface{}, names []string) (string, []interface{}) {
	keys, cols := g.fm.Columns(entity)
	if g.sorted {
		sort.Sort(keys)
		sort.Sort(cols)
	}

	kset := make(map[string]struct{})
	for _, k := range keys.Cols {
		kset[k] = struct{}{}
	}

	ucols := make(map[string]int)
	for i, e := range cols.Cols {
		if _, ok := kset[e]; !ok {
			ucols[e] = i
		}
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

	b.WriteString(") ON CONFLICT (")

	for i, e := range keys.Cols {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e)
	}

	b.WriteString(") DO UPDATE SET ")

	var n int
	for e, i := range ucols {
		if n > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e)
		b.WriteString(" = $")
		b.WriteString(strconv.FormatInt(int64(i+1), 10))
		n++
	}

	return b.String(), cols.Vals
}

func (g *Generator) Update(table string, entity interface{}, names []string) (string, []interface{}) {
	keys, cols := g.fm.Columns(entity)
	if g.sorted {
		sort.Sort(keys)
		sort.Sort(cols)
	}

	var incl map[string]struct{}
	if len(names) > 0 {
		incl = make(map[string]struct{})
		for _, e := range names {
			incl[e] = struct{}{}
		}
	}

	var n, x int
	b := &strings.Builder{}
	b.WriteString("UPDATE ")
	b.WriteString(table)
	b.WriteString(" SET ")

	args := make([]interface{}, 0, len(incl))

	n = 0
	for i, e := range cols.Cols {
		if incl != nil {
			if _, ok := incl[e]; !ok {
				continue
			}
		}
		if n > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e)
		b.WriteString(" = $")
		b.WriteString(strconv.FormatInt(int64(x+1), 10))
		args = append(args, cols.Vals[i])
		n++
		x++
	}

	b.WriteString(" WHERE ")

	n = 0
	for i, e := range keys.Cols {
		if n > 0 {
			b.WriteString(" AND ")
		}
		b.WriteString(e)
		b.WriteString(" = $")
		b.WriteString(strconv.FormatInt(int64(x+1), 10))
		args = append(args, keys.Vals[i])
		x++
		n++
	}

	return b.String(), args
}

func (g *Generator) Delete(table string, keys *Columns) (string, []interface{}) {
	if g.sorted {
		sort.Sort(keys)
	}

	var n, x int
	args := make([]interface{}, 0, len(keys.Vals))

	b := &strings.Builder{}
	b.WriteString("DELETE FROM ")
	b.WriteString(table)
	b.WriteString(" WHERE ")

	n = 0
	for i, e := range keys.Cols {
		if n > 0 {
			b.WriteString(" AND ")
		}
		b.WriteString(e)
		b.WriteString(" = $")
		b.WriteString(strconv.FormatInt(int64(x+1), 10))
		args = append(args, keys.Vals[i])
		x++
		n++
	}

	return b.String(), args
}
