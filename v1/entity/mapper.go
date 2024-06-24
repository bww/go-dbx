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
	optionOmitPQL    = "omitpql" // don't expand in PQL expressions
	optionOmitEmpty  = "omitempty"
)

var (
	initOnce      sync.Once
	defaultMapper *FieldMapper
)

// Synthetic is implemented by entities that produce additional "synthetic"
// columns to be stored alongside columns derived from fields of the entity.
//
// Synthetic columns are read-only, they will not be fetched back from the
// database. They are generally intended to be used in cases where derived
// values are computed from the state of an entity when it it stored, perhaps
// for indexing purposes.
type Synthetic interface {
	AdditionalColumns() *Columns
}

var typeOfSynthetic = reflect.TypeOf((*Synthetic)(nil)).Elem()

type FieldFilter func(*reflectx.FieldInfo) bool

func ExcludeFromPQL(f *reflectx.FieldInfo) bool {
	_, ok := f.Options[optionOmitPQL]
	return !ok // exclude (false) when present
}

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
					kvals = append(kvals, FieldByIndexes(e, f.Index))
				}
			}
		}
	}

	return &Values{kcols, kvals}, nil
}

func (m *FieldMapper) ColumnsForType(typ reflect.Type, filters ...FieldFilter) []string {
	var cols []string
	x := m.TypeMap(typ)
outer:
	for k, f := range x.Names {
		for _, filter := range filters {
			if !filter(f) {
				continue outer
			}
		}
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

			v, ok := FieldByIndexesRO(e, f.Index)
			if !ok { // the field could not be fully traversed; an intermediate value is nil
				vcols, vvals = append(vcols, k), append(vvals, nil)
				continue
			}
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
					kcols, kvals = append(kcols, k), append(kvals, x)
				}
			}
		}
	}

	rkcols := &Columns{kcols, kvals}
	rvcols := &Columns{vcols, vvals}

	if syn, ok := entity.(Synthetic); ok {
		if acols := syn.AdditionalColumns(); acols != nil {
			rvcols.Cols = append(rvcols.Cols, acols.Cols...)
			rvcols.Vals = append(rvcols.Vals, acols.Vals...)
		}
	}

	return rkcols, rvcols
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
	err := mustBe(t, reflect.Struct)
	if err != nil {
		return err
	}
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

func FieldByIndexes(v reflect.Value, indexes []int) reflect.Value {
	for _, x := range indexes {
		v = reflect.Indirect(v).Field(x)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			v.Set(reflect.New(reflectx.Deref(v.Type())))
		}
		if v.Kind() == reflect.Map && v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
		}
	}
	return v
}

func FieldByIndexesRO(v reflect.Value, indexes []int) (reflect.Value, bool) {
	for _, i := range indexes {
		v = reflect.Indirect(v)
		if v.IsValid() {
			v = v.Field(i)
		} else {
			return v, false
		}
	}
	return v, true
}

type kinder interface {
	Kind() reflect.Kind
}

func mustBe(v kinder, expected reflect.Kind) error {
	if k := v.Kind(); k != expected {
		return &reflect.ValueError{Method: methodName(), Kind: k}
	}
	return nil
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
