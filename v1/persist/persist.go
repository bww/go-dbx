package persist

import (
	"reflect"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
	"github.com/bww/go-dbx/v1/persist/ident"
	"github.com/bww/go-dbx/v1/persist/pql"
)

type Persister interface {
	With(dbx.Context) Persister
	Store(string, interface{}, []string) error
	Fetch(string, interface{}, interface{}) error
	Count(string, ...interface{}) (int, error)
	Select(interface{}, string, ...interface{}) error
	Delete(string, interface{}) error
}

type persister struct {
	dbx.Context
	fm  *entity.FieldMapper
	gen *entity.Generator
	ids ident.Generator
}

func New(cxt dbx.Context, fm *entity.FieldMapper, ids ident.Generator) Persister {
	return &persister{
		Context: cxt,
		fm:      fm,
		gen:     entity.NewGenerator(fm),
		ids:     ids,
	}
}

func (p *persister) With(cxt dbx.Context) Persister {
	return &persister{
		Context: cxt,
		fm:      p.fm,
		gen:     p.gen,
		ids:     p.ids,
	}
}

func (p *persister) Fetch(table string, ent, id interface{}) error {
	keys, _ := p.fm.Columns(ent)

	if len(keys.Cols) != 1 {
		return dbx.ErrInvalidKeyCount
	}

	sql, args := p.gen.Select(table, ent, &entity.Columns{
		Cols: keys.Cols,
		Vals: []interface{}{id},
	})

	err := p.Context.QueryRowx(sql, args...).StructScan(ent)
	if err != nil {
		return err
	}

	return nil
}

func (p *persister) Count(query string, args ...interface{}) (int, error) {
	var n int

	err := p.Context.QueryRow(query, args...).Scan(&n)
	if err != nil {
		return -1, err
	}

	return n, nil
}

func (p *persister) Select(ent interface{}, query string, args ...interface{}) error {
	val := reflect.ValueOf(ent)
	ind := reflect.Indirect(val)

	var many bool
	switch ind.Kind() {
	case reflect.Slice:
		many = true
	default:
		many = false
	}

	var typ reflect.Type
	if many {
		typ = ind.Type().Elem()
	} else {
		typ = ind.Type()
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	cols := p.fm.ColumnsForType(typ)

	prg, err := pql.Parse(query)
	if err != nil {
		return err
	}
	sql, err := prg.Text(pql.Context{Columns: cols})
	if err != nil {
		return err
	}

	if many {
		return p.selectMany(ent, val, cols, sql, args)
	} else {
		return p.selectOne(ent, val, cols, sql, args)
	}
}

func (p *persister) selectOne(ent interface{}, val reflect.Value, cols []string, sql string, args []interface{}) error {
	return p.Context.QueryRowx(sql, args...).StructScan(ent)
}

func (p *persister) selectMany(ent interface{}, val reflect.Value, cols []string, sql string, args []interface{}) error {
	if val.Kind() != reflect.Ptr {
		return dbx.ErrNotAPointer
	}

	rows, err := p.Context.Queryx(sql, args...)
	if err != nil {
		return err
	}

	eval := reflect.Indirect(val)
	etype := eval.Type().Elem()
	ctype := etype

	if etype.Kind() == reflect.Ptr {
		ctype = etype.Elem()
	}

	for rows.Next() {
		elem := reflect.New(ctype)
		err := rows.StructScan(elem.Interface())
		if err != nil {
			return err
		}
		if etype.Kind() != reflect.Ptr {
			elem = reflect.Indirect(elem)
		}
		eval = reflect.Append(eval, elem)
	}

	reflect.Indirect(val).Set(eval)
	return nil
}

func (p *persister) Store(table string, ent interface{}, cols []string) error {
	var insert bool

	keys, err := p.fm.Keys(ent)
	if err != nil {
		return err
	}

	for _, e := range keys.Vals {
		if !e.IsValid() {
			return dbx.ErrInvalidField
		}
		if e.IsZero() {
			insert = true
			break
		}
	}

	if insert {
		if len(keys.Vals) != 1 {
			return dbx.ErrInvalidKeyCount
		}
		keys.Vals[0].Set(p.ids()) // generate primary key
	}

	var sql string
	var args []interface{}
	if insert {
		sql, args = p.gen.Insert(table, ent)
	} else {
		sql, args = p.gen.Update(table, ent, cols)
	}

	_, err = p.Context.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p *persister) Delete(table string, ent interface{}) error {
	keys, _ := p.fm.Columns(ent)
	if len(keys.Cols) != 1 {
		return dbx.ErrInvalidKeyCount
	}

	sql, args := p.gen.Delete(table, ent, keys)
	_, err := p.Context.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}
