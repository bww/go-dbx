package dbx

import (
	"database/sql"
	"log"

	"github.com/bww/go-util/v1/text"
	"github.com/jmoiron/sqlx"
)

// An SQL context. This defines a unified type that encompasses the basic
// methods of sql.DB and sql.Tx so they can be used interchangably.
type Context interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
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

func (d *DB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	if d.debug {
		d.log.Printf("dbx/query/n: (%T) [%s] %v\n", d, text.CollapseSpaces(query), args)
	}
	return d.DB.Queryx(query, args...)
}

func (d *DB) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	if d.debug {
		d.log.Printf("dbx/query/1: (%T) [%s] %v\n", d, text.CollapseSpaces(query), args)
	}
	return d.DB.QueryRowx(query, args...)
}

func (d *DB) Wrap(cxt Context) Context {
	return NewContext(cxt, d.log, d.debug)
}

// A wrapped context. This is primarily useful for wrapping transactions to manage
// logging and debugging parameters.
type wrappedContext struct {
	Context
	log   *log.Logger
	debug bool
}

func NewContext(cxt Context, l *log.Logger, d bool) Context {
	if l == nil {
		l = defaultLogger
	}
	return &wrappedContext{Context: cxt, log: l, debug: d}
}

func (c *wrappedContext) Exec(query string, args ...interface{}) (sql.Result, error) {
	if c.debug {
		c.log.Printf("dbx/exec: (%T) [%s] %v\n", c.Context, text.CollapseSpaces(query), args)
	}
	return c.Context.Exec(query, args...)
}

func (c *wrappedContext) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if c.debug {
		c.log.Printf("dbx/query/n: (%T) [%s] %v\n", c.Context, text.CollapseSpaces(query), args)
	}
	return c.Context.Query(query, args...)
}

func (c *wrappedContext) QueryRow(query string, args ...interface{}) *sql.Row {
	if c.debug {
		c.log.Printf("dbx/query/1: (%T) [%s] %v\n", c.Context, text.CollapseSpaces(query), args)
	}
	return c.Context.QueryRow(query, args...)
}

func (c *wrappedContext) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	if c.debug {
		c.log.Printf("dbx/query/n: (%T) [%s] %v\n", c.Context, text.CollapseSpaces(query), args)
	}
	return c.Context.Queryx(query, args...)
}

func (c *wrappedContext) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	if c.debug {
		c.log.Printf("dbx/query/1: (%T) [%s] %v\n", c.Context, text.CollapseSpaces(query), args)
	}
	return c.Context.QueryRowx(query, args...)
}
