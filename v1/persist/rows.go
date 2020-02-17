package persist

import (
	"database/sql"
	"reflect"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

type Row struct {
	*sqlx.Row
	mapper *entity.FieldMapper
}

func newRow(r *sqlx.Row, m *entity.FieldMapper) *Row {
	return &Row{Row: r, mapper: m}
}

func (r *Row) ScanStruct(dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return dbx.ErrNotAPointer
	}

	v = v.Elem()
	columns, err := r.Columns()
	if err != nil {
		return err
	}

	m := r.mapper
	fields, omits := m.TraversalsByName(v.Type(), columns)
	temps := make([]bool, len(columns))
	values := make([]interface{}, len(columns))

	// initialize for scanning
	err = fieldsByTraversal(v, fields, omits, temps, values, true)
	if err != nil {
		return err
	}

	// scan out values, potentically including indirect placeholders for omittable fields
	err = r.Scan(values...)
	if err != nil {
		return err
	}

	// copy over valid omittable values to their fields
	err = finalizeFields(v, fields, temps, values)
	if err != nil {
		return err
	}

	return r.Err()
}

type Rows struct {
	*sql.Rows
	mapper *entity.FieldMapper
	fields [][]int
	omits  []bool
	temps  []bool
	values []interface{}
}

func newRows(r *sql.Rows, m *entity.FieldMapper) *Rows {
	return &Rows{Rows: r, mapper: m}
}

func (r *Rows) ScanStruct(dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return dbx.ErrNotAPointer
	}

	v = v.Elem()
	if r.fields == nil {
		columns, err := r.Columns()
		if err != nil {
			return err
		}
		m := r.mapper
		r.fields, r.omits = m.TraversalsByName(v.Type(), columns)
		r.temps = make([]bool, len(columns))
		r.values = make([]interface{}, len(columns))
	}

	// initialize for scanning
	err := fieldsByTraversal(v, r.fields, r.omits, r.temps, r.values, true)
	if err != nil {
		return err
	}

	// scan out values, potentically including indirect placeholders for omittable fields
	err = r.Scan(r.values...)
	if err != nil {
		return err
	}

	// copy over valid omittable values to their fields
	err = finalizeFields(v, r.fields, r.temps, r.values)
	if err != nil {
		return err
	}

	return r.Err()
}

func fieldsByTraversal(v reflect.Value, fields [][]int, omits, temps []bool, values []interface{}, ptrs bool) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return dbx.ErrNotAStruct
	}
	for i, field := range fields {
		if len(field) == 0 {
			values[i], temps[i] = new(interface{}), false
		} else {
			f := reflectx.FieldByIndexes(v, field)
			if omits[i] && ptrs && f.Kind() != reflect.Ptr {
				values[i], temps[i] = reflect.New(reflect.PtrTo(f.Type())).Interface(), true
			} else if ptrs {
				values[i], temps[i] = f.Addr().Interface(), false
			} else {
				values[i], temps[i] = f.Interface(), false
			}
		}
	}
	return nil
}

func finalizeFields(v reflect.Value, fields [][]int, temps []bool, values []interface{}) error {
	for i, e := range temps {
		if e {
			t := reflect.Indirect(reflect.ValueOf(values[i]))
			f := reflectx.FieldByIndexes(v, fields[i])
			if !t.IsNil() {
				f.Set(reflect.Indirect(t))
			} else { // explicitly set the zero value if the value is missing
				f.Set(reflect.Zero(f.Type()))
			}
		}
	}
	return nil
}
