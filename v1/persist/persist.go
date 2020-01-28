package persist

import (
	"errors"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
)

var (
	ErrInvalidField    = errors.New("Invalid field")
	ErrInvalidKeyCount = errors.New("Invalid primary key count")
)

type Persister interface {
	Context(...dbx.Context) dbx.Context
	Store(string, interface{}, []string, dbx.Context) error
}

type persister struct {
	cxt dbx.Context // default context
	fm  *entity.FieldMapper
	gen *entity.Generator
	ids IdentFunc
}

func New(cxt dbx.Context, ids IdentFunc) Persister {
	fm := entity.NewFieldMapper(entity.Tag)
	gen := entity.NewGenerator(fm)
	return &persister{cxt: cxt, fm: fm, gen: gen, ids: ids}
}

func (p *persister) Context(cxts ...dbx.Context) dbx.Context {
	for _, e := range cxts {
		if e != nil {
			return e
		}
	}
	return p.cxt
}

func (p *persister) Store(table string, entity interface{}, cols []string, cxt dbx.Context) error {
	var insert bool

	keys, err := p.fm.Keys(entity)
	if err != nil {
		return err
	}

	for _, e := range keys.Vals {
		if !e.IsValid() {
			return ErrInvalidField
		}
		if e.IsZero() {
			insert = true
			break
		}
	}

	if insert {
		if len(keys.Vals) != 1 {
			return ErrInvalidKeyCount
		}
		keys.Vals[0].Set(p.ids()) // generate primary key
	}

	var sql string
	var args []interface{}
	if insert {
		sql, args = p.gen.Insert(table, entity)
	} else {
		sql, args = p.gen.Update(table, entity, cols)
	}

	cxt = p.Context(cxt)
	_, err = cxt.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}
