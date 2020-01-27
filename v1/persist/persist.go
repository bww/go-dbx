package persist

import (
	"errors"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
)

var (
	ErrInvalidField = errors.New("Invalid field")
)

type Persister interface {
	Context(...dbx.Context) dbx.Context
	Store(string, interface{}, []string, dbx.Context) error
}

type persister struct {
	cxt dbx.Context // default context
	fm  *entity.FieldMapper
	gen *entity.Generator
}

func New(cxt dbx.Context) Persister {
	return &persister{cxt}
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

	keys := p.fm.Keys(entity)
	for _, e := range keys.Vals {
		if !e.IsValid() {
			return ErrInvalidField
		}
		if e.IsZero() {
			insert = true
			break
		}
	}

	var sql string
	var args []interface{}
	if insert {
		sql, args = p.gen.Insert(table, entity)
	} else {
		sql, args = p.gen.Update(table, entity, cols)
	}

	cxt = p.Context(cxt)
	_, err := cxt.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}
