package persist

import (
	"fmt"
	"reflect"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
	"github.com/jmoiron/sqlx"
)

type Row struct {
	*sqlx.Row
	mapper *entity.FieldMapper
}

func newRow(r *sqlx.Row, m *entity.FieldMapper) *Row {
	return &Row{Row: r, mapper: m}
}

func (r *Row) ScanStruct(dest interface{}) error {
	var err error
	var scanned bool
	defer func() {
		if !scanned {
			fmt.Printf("dbx: Row was not closed due to a parameter error; you are not using DBX correctly: %v\n", err)
			r.Row.Scan() // scan to force sqlx to close its rows if we haven't done so already
		}
	}()

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

	// Scan out values, potentically including indirect placeholders for omittable fields;
	// NOTE THAT Scan() MUST BE INVOKED AS CALLING SCAN IS WHAT CLOSES THE UNDERLYING ROWS
	// THAT ARE USED BY sqlx.
	err, scanned = r.Scan(values...), true
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
	*sqlx.Rows
	mapper *entity.FieldMapper
	fields [][]int
	omits  []bool
	temps  []bool
	values []interface{}
}

func newRows(r *sqlx.Rows, m *entity.FieldMapper) *Rows {
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
			f := entity.FieldByIndexes(v, field)
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
			f := entity.FieldByIndexes(v, fields[i])
			if !t.IsNil() {
				f.Set(reflect.Indirect(t))
			} else { // explicitly set the zero value if the value is missing
				f.Set(reflect.Zero(f.Type()))
			}
		}
	}
	return nil
}
