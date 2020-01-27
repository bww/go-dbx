package entity

import (
	"reflect"

	"github.com/jmoiron/sqlx/reflectx"
)

const defaultTag = "db"

type FieldMapper struct {
	*reflectx.Mapper
}

func NewFieldMapper(tag string) *FieldMapper {
	return &FieldMapper{
		Mapper: reflectx.NewMapper(tag),
	}
}

func (m *FieldMapper) Columns(entity interface{}) (*Columns, *Columns) {
	var vcols, kcols []string
	var vvals, kvals []interface{}

	e := reflect.ValueOf(entity)
	x := m.TypeMap(e.Type())

	for k, f := range x.Names {
		var x interface{}

		v := reflectx.FieldByIndexes(e, f.Index)
		if v.IsValid() && !v.IsZero() {
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
