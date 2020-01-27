package dbx

import (
	"database/sql"
	"log"

	"github.com/bww/go-util/text"
)

// An SQL context. This defines a unified type that encompasses the basic
// methods of sql.DB and sql.Tx so they can be used interchangably
type Context interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if d.debug {
		d.log.Printf("dbx/exec: (%T) [%s] %v\n", d, text.CollapseSpaces(query), args)
	}
	return d.DB.Exec(query, args...)
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if d.debug {
		d.log.Printf("dbx/query/n: (%T) [%s] %v\n", d, text.CollapseSpaces(query), args)
	}
	return d.DB.Query(query, args...)
}

func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	if d.debug {
		d.log.Printf("dbx/query/1: (%T) [%s] %v\n", d, text.CollapseSpaces(query), args)
	}
	return d.DB.QueryRow(query, args...)
}

// A debug context which wraps another context and logs out statements
type DebugContext struct {
	Context
	log *log.Logger
}

func NewDebugContext(cxt Context, l *log.Logger) DebugContext {
	return DebugContext{Context: cxt, log: l}
}

func (d DebugContext) Exec(query string, args ...interface{}) (sql.Result, error) {
	d.log.Printf("dbx/exec: (%T) [%s] %v\n", d.Context, text.CollapseSpaces(query), args)
	return d.Context.Exec(query, args...)
}

func (d DebugContext) Query(query string, args ...interface{}) (*sql.Rows, error) {
	d.log.Printf("dbx/query/n: (%T) [%s] %v\n", d.Context, text.CollapseSpaces(query), args)
	return d.Context.Query(query, args...)
}

func (d DebugContext) QueryRow(query string, args ...interface{}) *sql.Row {
	d.log.Printf("dbx/query/1: (%T) [%s] %v\n", d.Context, text.CollapseSpaces(query), args)
	return d.Context.QueryRow(query, args...)
}
