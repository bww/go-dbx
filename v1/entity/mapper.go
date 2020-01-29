package entity

import (
	"reflect"
	"sync"

	"github.com/bww/go-dbx/v1"
	"github.com/jmoiron/sqlx/reflectx"
)

const Tag = "db"

var (
	initOnce      sync.Once
	defaultMapper *FieldMapper
)

func DefaultFieldMapper() *FieldMapper {
	initOnce.Do(func() {
		defaultMapper = NewFieldMapper()
	})
	return defaultMapper
}

type FieldMapper struct {
	*reflectx.Mapper
}

func NewFieldMapper() *FieldMapper {
	return &FieldMapper{
		Mapper: reflectx.NewMapper(Tag),
	}
}

func (m *FieldMapper) Keys(entity interface{}) (*Values, error) {
	var kcols []string
	var kvals []reflect.Value

	e := reflect.ValueOf(entity)
	x := m.TypeMap(e.Type())

	if e.Kind() != reflect.Ptr {
		return nil, dbx.ErrNotAPointer
	}

	for k, f := range x.Names {
		if f.Options != nil {
			if _, ok := f.Options["pk"]; ok {
				kcols = append(kcols, k)
				kvals = append(kvals, reflectx.FieldByIndexes(e, f.Index))
			}
		}
	}

	return &Values{kcols, kvals}, nil
}

func (m *FieldMapper) Columns(entity interface{}) (*Columns, *Columns) {
	var vcols, kcols []string
	var vvals, kvals []interface{}

	e := reflect.ValueOf(entity)
	x := m.TypeMap(e.Type())

	for k, f := range x.Names {
		var x interface{}

		v := reflectx.FieldByIndexes(e, f.Index)
		if v.IsValid() && v.CanInterface() {
			x = v.Interface()
		} else {
			x = f.Zero.Interface()
		}

		vcols = append(vcols, k)
		vvals = append(vvals, x)

		if f.Options != nil {
			if _, ok := f.Options["pk"]; ok {
				kcols = append(kcols, k)
				kvals = append(kvals, x)
			}
		}
	}

	return &Columns{kcols, kvals}, &Columns{vcols, vvals}
}
