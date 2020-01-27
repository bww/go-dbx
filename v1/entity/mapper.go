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

func (m *FieldMapper) Columns(entity interface{}) ([]string, []interface{}) {
	var cols []string
	var vals []interface{}

	e := reflect.ValueOf(entity)
	x := m.TypeMap(e.Type())

	for k, f := range x.Names {
		v := reflectx.FieldByIndexes(e, f.Index)
		if v.IsValid() && !v.IsZero() {
			vals = append(vals, v.Interface())
		} else {
			vals = append(vals, f.Zero.Interface())
		}
		cols = append(cols, k)
	}

	return cols, vals
}

var defaultMapper = NewFieldMapper(defaultTag)

func Columns(entity interface{}) ([]string, []interface{}) {
	return defaultMapper.Columns(entity)
}
