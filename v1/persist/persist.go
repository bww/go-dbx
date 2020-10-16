package persist

import (
	"reflect"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
	"github.com/bww/go-dbx/v1/persist/ident"
	"github.com/bww/go-dbx/v1/persist/option"
	"github.com/bww/go-dbx/v1/persist/pql"
	"github.com/bww/go-dbx/v1/persist/registry"

	dbsql "database/sql"
)

type FetchRelatedPersister interface {
	FetchRelated(Persister, interface{}) error
}

type StoreRelatedPersister interface {
	StoreRelated(Persister, interface{}) error
}

type StoreReferencesPersister interface {
	StoreReferences(Persister, interface{}) error
}

type DeleteRelatedPersister interface {
	DeleteRelated(Persister, interface{}) error
}

type DeleteReferencesPersister interface {
	DeleteReferences(Persister, interface{}) error
}

type Persister interface {
	dbx.Context
	WithContext(dbx.Context) Persister
	WithOptions(...option.Option) Persister
	Config() option.Config
	Param(name string) (interface{}, bool)
	Store(string, interface{}, []string) error
	Fetch(string, interface{}, interface{}) error
	Count(string, ...interface{}) (int, error)
	Select(interface{}, string, ...interface{}) error
	Delete(string, interface{}) error
	DeleteWithID(string, reflect.Type, interface{}) error
}

type persister struct {
	dbx.Context
	fm   *entity.FieldMapper
	gen  *entity.Generator
	reg  *registry.Registry
	ids  ident.Generator
	conf option.Config
}

func New(cxt dbx.Context, fm *entity.FieldMapper, reg *registry.Registry, ids ident.Generator, opts ...option.Option) Persister {
	return &persister{
		Context: cxt,
		fm:      fm,
		gen:     entity.NewGenerator(fm),
		reg:     reg,
		ids:     ids,
		conf: option.NewConfig(option.Config{
			FetchRelated:  true,
			StoreRelated:  true,
			DeleteRelated: true,
		}, opts),
	}
}

func (p *persister) WithContext(cxt dbx.Context) Persister {
	return &persister{
		Context: cxt,
		fm:      p.fm,
		gen:     p.gen,
		reg:     p.reg,
		ids:     p.ids,
		conf:    p.conf,
	}
}

func (p *persister) WithOptions(opts ...option.Option) Persister {
	return &persister{
		Context: p.Context,
		fm:      p.fm,
		gen:     p.gen,
		reg:     p.reg,
		ids:     p.ids,
		conf:    option.NewConfig(p.conf, opts),
	}
}

func (p *persister) Config() option.Config {
	return p.conf
}

func (p *persister) Param(name string) interface{} {
	if m := p.conf.Params; m != nil {
		return m[name]
	}
	return nil
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

	raw := p.Context.QueryRowx(sql, args...)
	row := newRow(raw, p.fm)
	err := row.ScanStruct(ent)
	if err == dbsql.ErrNoRows {
		return dbx.ErrNotFound
	} else if err != nil {
		return err
	}

	if p.conf.FetchRelated {
		if pst, ok := p.reg.GetFor(ent); ok {
			if c, ok := pst.(FetchRelatedPersister); ok {
				err = c.FetchRelated(p, ent)
				if err != nil {
					return err
				}
			}
		}
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

	raw := p.Context.QueryRowx(sql, args...)
	row := newRow(raw, p.fm)
	err := row.ScanStruct(ent)
	if err == dbsql.ErrNoRows {
		return dbx.ErrNotFound
	} else if err != nil {
		return err
	}

	if p.conf.FetchRelated {
		if pst, ok := p.reg.Get(val.Type()); ok {
			if c, ok := pst.(FetchRelatedPersister); ok {
				err = c.FetchRelated(p, ent)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (p *persister) selectMany(ent interface{}, val reflect.Value, cols []string, sql string, args []interface{}) error {
	if val.Kind() != reflect.Ptr {
		return dbx.ErrNotAPointer
	}

	raws, err := p.Context.Queryx(sql, args...)
	if err != nil {
		return err
	}

	rows := newRows(raws, p.fm)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	eval := reflect.Indirect(val)
	etype := eval.Type().Elem()
	ctype := etype

	if etype.Kind() == reflect.Ptr {
		ctype = etype.Elem()
	}

	var rel FetchRelatedPersister
	if p.conf.FetchRelated {
		if pst, ok := p.reg.Get(ctype); ok {
			if c, ok := pst.(FetchRelatedPersister); ok {
				rel = c
			}
		}
	}

	for rows.Next() {
		elem := reflect.New(ctype)
		eint := elem.Interface()
		err := rows.ScanStruct(eint)
		if err != nil {
			return err
		}
		if rel != nil {
			err = rel.FetchRelated(p, eint)
			if err != nil {
				return err
			}
		}
		if etype.Kind() != reflect.Ptr {
			elem = reflect.Indirect(elem)
		}
		eval = reflect.Append(eval, elem)
	}

	err, rows = rows.Close(), nil
	if err != nil {
		return err
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

	if p.conf.StoreRelated {
		if pst, ok := p.reg.GetFor(ent); ok {
			if c, ok := pst.(StoreRelatedPersister); ok {
				err := c.StoreRelated(p, ent)
				if err != nil {
					return err
				}
			}
			if c, ok := pst.(StoreReferencesPersister); ok {
				err := c.StoreReferences(p, ent)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (p *persister) Delete(table string, ent interface{}) error {
	keys, _ := p.fm.Columns(ent)
	if len(keys.Cols) != 1 {
		return dbx.ErrInvalidKeyCount
	}

	if p.conf.DeleteRelated {
		if pst, ok := p.reg.GetFor(ent); ok {
			if c, ok := pst.(DeleteReferencesPersister); ok {
				err := c.DeleteReferences(p, ent)
				if err != nil {
					return err
				}
			}
			if c, ok := pst.(DeleteRelatedPersister); ok {
				err := c.DeleteRelated(p, ent)
				if err != nil {
					return err
				}
			}
		}
	}

	sql, args := p.gen.Delete(table, keys)
	_, err := p.Context.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p *persister) DeleteWithID(table string, typ reflect.Type, id interface{}) error {
	keys := p.fm.KeysForType(typ)
	if len(keys) != 1 {
		return dbx.ErrInvalidKeyCount
	}

	sql, args := p.gen.Delete(table, &entity.Columns{
		Cols: keys,
		Vals: []interface{}{id},
	})
	_, err := p.Context.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}
