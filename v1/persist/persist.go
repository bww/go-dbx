package persist

import (
	"fmt"
	"reflect"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
	"github.com/bww/go-dbx/v1/persist/ident"
	"github.com/bww/go-dbx/v1/persist/pql"
)

type Persister interface {
	Context(...dbx.Context) dbx.Context
	Store(string, interface{}, []string, dbx.Context) error
	Fetch(string, interface{}, interface{}, dbx.Context) error
	Select(string, interface{}, dbx.Context, string, ...interface{}) error
}

type persister struct {
	cxt dbx.Context // default context
	fm  *entity.FieldMapper
	gen *entity.Generator
	ids ident.Generator
}

func New(cxt dbx.Context, fm *entity.FieldMapper, ids ident.Generator) Persister {
	return &persister{
		cxt: cxt,
		fm:  fm,
		gen: entity.NewGenerator(fm),
		ids: ids,
	}
}

func (p *persister) Context(cxts ...dbx.Context) dbx.Context {
	for _, e := range cxts {
		if e != nil {
			return e
		}
	}
	return p.cxt
}

func (p *persister) Fetch(table string, ent interface{}, id interface{}, cxt dbx.Context) error {
	keys, _ := p.fm.Columns(ent)

	if len(keys.Cols) != 1 {
		return dbx.ErrInvalidKeyCount
	}

	sql, args := p.gen.Select(table, ent, &entity.Columns{
		Cols: keys.Cols,
		Vals: []interface{}{id},
	})

	err := p.Context(cxt).QueryRowx(sql, args...).StructScan(ent)
	if err != nil {
		return err
	}

	return nil
}

func (p *persister) Select(table string, ent interface{}, cxt dbx.Context, query string, args ...interface{}) error {
	_, cols := p.fm.Columns(ent)

	prg, err := pql.Parse(query)
	if err != nil {
		return err
	}

	sql, err := prg.Text(pql.Context{Columns: cols.Cols})
	if err != nil {
		return err
	}

	val := reflect.Indirect(reflect.ValueOf(ent))
	switch val.Kind() {
	case reflect.Slice:
		return p.selectMany(table, ent, val, cols, cxt, sql, args)
	default:
		return p.selectOne(table, ent, val, cols, cxt, sql, args)
	}
}

func (p *persister) selectOne(table string, ent interface{}, val reflect.Value, cols *entity.Columns, cxt dbx.Context, sql string, args []interface{}) error {
	return p.Context(cxt).QueryRowx(sql, args...).StructScan(ent)
}

func (p *persister) selectMany(table string, ent interface{}, val reflect.Value, cols *entity.Columns, cxt dbx.Context, sql string, args []interface{}) error {
	rows, err := p.Context(cxt).Queryx(sql, args...)
	if err != nil {
		return err
	}

	fmt.Printf(">>> %+v\n", val)
	etype := val.Type().Elem()
	for rows.Next() {
		elem := reflect.New(etype)
		fmt.Printf(">>> %+v\n", elem)

		dest := elem
		if dest.Kind() != reflect.Ptr {
			dest = dest.Addr()
		}

		err := rows.StructScan(dest.Interface())
		if err != nil {
			return err
		}

		val = reflect.Append(val, elem)
	}

	return nil
}

func (p *persister) Store(table string, ent interface{}, cols []string, cxt dbx.Context) error {
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

	_, err = p.Context(cxt).Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}
