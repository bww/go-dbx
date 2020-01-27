package persist

import (
	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
)

type Persister interface {
	Context(...dbx.Context) dbx.Context
}

type persister struct {
	cxt dbx.Context // default context
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

func (p *persister) Store(entity interface{}) error {
	return nil
}
