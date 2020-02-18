package entity

import (
	"reflect"
	"runtime"
	"sync"

	"github.com/bww/go-dbx/v1"
	"github.com/jmoiron/sqlx/reflectx"
)

const Tag = "db"
const (
	optionPrimaryKey = "pk"
	optionOmitEmpty  = "omitempty"
)

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
		Mapper: reflectx.NewMapperFunc(Tag, ignoreField),
	}
}

func (m *FieldMapper) KeysForType(typ reflect.Type) []string {
	var cols []string
	x := m.TypeMap(typ)
	for k, f := range x.Names {
		if isExplicitMapping(f) {
			if f.Options != nil {
				if _, ok := f.Options[optionPrimaryKey]; ok {
					cols = append(cols, k)
				}
			}
		}
	}
	return cols
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
		if isExplicitMapping(f) {
			if f.Options != nil {
				if _, ok := f.Options[optionPrimaryKey]; ok {
					kcols = append(kcols, k)
					kvals = append(kvals, reflectx.FieldByIndexes(e, f.Index))
				}
			}
		}
	}

	return &Values{kcols, kvals}, nil
}

func (m *FieldMapper) ColumnsForType(typ reflect.Type) []string {
	var cols []string
	x := m.TypeMap(typ)
	for k, f := range x.Names {
		if isExplicitMapping(f) {
			cols = append(cols, k)
		}
	}
	return cols
}

func (m *FieldMapper) Columns(entity interface{}) (*Columns, *Columns) {
	var vcols, kcols []string
	var vvals, kvals []interface{}

	e := reflect.ValueOf(entity)
	x := m.TypeMap(e.Type())

	for k, f := range x.Names {
		if isExplicitMapping(f) {
			var x interface{}

			v := reflectx.FieldByIndexes(e, f.Index)
			if v.IsValid() && v.CanInterface() {
				x = v.Interface()
			} else {
				x = f.Zero.Interface()
			}

			vcols = append(vcols, k)
			if f.Options == nil {
				vvals = append(vvals, x)
			} else if _, ok := f.Options[optionOmitEmpty]; ok && isEmptyValue(x, v) {
				vvals = append(vvals, nil)
			} else {
				vvals = append(vvals, x)
			}

			if f.Options != nil {
				if _, ok := f.Options[optionPrimaryKey]; ok {
					kcols = append(kcols, k)
					kvals = append(kvals, x)
				}
			}
		}
	}

	return &Columns{kcols, kvals}, &Columns{vcols, vvals}
}

// TraversalsByName returns a slice of int slices which represent the struct
// traversals for each mapped name and a slice of bools which indicates whether
// each such traversal is omit-empty.  Panics if t is not a struct or Indirectable
// to a struct.  Returns empty int slice for each name not found.
func (m *FieldMapper) TraversalsByName(t reflect.Type, names []string) ([][]int, []bool) {
	r := make([][]int, 0, len(names))
	o := make([]bool, 0, len(names))
	m.TraversalsByNameFunc(t, names, func(_ int, i []int, e bool) error {
		if i == nil {
			r = append(r, []int{})
		} else {
			r = append(r, i)
		}
		o = append(o, e)
		return nil
	})
	return r, o
}

// TraversalsByNameFunc traverses the mapped names and calls fn with the index of
// each name and the struct traversal represented by that name. Panics if t is not
// a struct or Indirectable to a struct. Returns the first error returned by fn or nil.
func (m *FieldMapper) TraversalsByNameFunc(t reflect.Type, names []string, fn func(int, []int, bool) error) error {
	t = reflectx.Deref(t)
	mustBe(t, reflect.Struct)
	tm := m.TypeMap(t)
	for i, name := range names {
		fi, ok := tm.Names[name]
		if !ok {
			if err := fn(i, nil, false); err != nil {
				return err
			}
		} else {
			var oe bool
			if fi.Options != nil {
				_, oe = fi.Options[optionOmitEmpty]
			}
			if err := fn(i, fi.Index, oe); err != nil {
				return err
			}
		}
	}
	return nil
}

type kinder interface {
	Kind() reflect.Kind
}

func mustBe(v kinder, expected reflect.Kind) {
	if k := v.Kind(); k != expected {
		panic(&reflect.ValueError{Method: methodName(), Kind: k})
	}
}

func methodName() string {
	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	if f != nil {
		return f.Name()
	} else {
		return "unknown method"
	}
}

const ignoreName = "__ignore_field__"

func ignoreField(n string) string {
	return ignoreName
}

func isExplicitMapping(f *reflectx.FieldInfo) bool {
	return f.Name != ignoreName
}
