package dbx

import (
	"database/sql"
	"log"

	"github.com/bww/go-util/v1/text"
	"github.com/jmoiron/sqlx"
)

// An SQL context. This defines a unified type that encompasses the basic
// methods of sqlx.DB and sqlx.Tx so they can be used interchangably.
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

func (d *DB) wrapTx(tx Tx) *wrappedTx {
	return newTx(tx, d.log, d.debug)
}

// A transactional SQL context. This defines a unified type that
// encompasses the basic methods of sqlx.Tx and other theoretical
// transaction implementsions so that they can be used interchangably.
type Tx interface {
	Context
	Commit() error
	Rollback() error
}

// Determine if a context is implemented by a transaction or not
func IsTx(cxt Context) bool {
	_, ok := cxt.(Tx)
	return ok
}

// A wrapped transaction. This is primarily useful for wrapping a
// transaction to manage logging and debugging parameters.
type wrappedTx struct {
	Tx
	log   *log.Logger
	debug bool
}

func newTx(tx Tx, l *log.Logger, d bool) *wrappedTx {
	if l == nil {
		l = defaultLogger
	}
	return &wrappedTx{Tx: tx, log: l, debug: d}
}

func (c *wrappedTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	if c.debug {
		c.log.Printf("dbx/exec: (%T) [%s] %v\n", c.Tx, text.CollapseSpaces(query), args)
	}
	return c.Tx.Exec(query, args...)
}

func (c *wrappedTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if c.debug {
		c.log.Printf("dbx/query/n: (%T) [%s] %v\n", c.Tx, text.CollapseSpaces(query), args)
	}
	return c.Tx.Query(query, args...)
}

func (c *wrappedTx) QueryRow(query string, args ...interface{}) *sql.Row {
	if c.debug {
		c.log.Printf("dbx/query/1: (%T) [%s] %v\n", c.Tx, text.CollapseSpaces(query), args)
	}
	return c.Tx.QueryRow(query, args...)
}

func (c *wrappedTx) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	if c.debug {
		c.log.Printf("dbx/query/n: (%T) [%s] %v\n", c.Tx, text.CollapseSpaces(query), args)
	}
	return c.Tx.Queryx(query, args...)
}

func (c *wrappedTx) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	if c.debug {
		c.log.Printf("dbx/query/1: (%T) [%s] %v\n", c.Tx, text.CollapseSpaces(query), args)
	}
	return c.Tx.QueryRowx(query, args...)
}

func (c *wrappedTx) Commit() error {
	if c.debug {
		c.log.Printf("dbx/commit: (%T)\n", c.Tx)
	}
	return c.Tx.Commit()
}

func (c *wrappedTx) Rollback() error {
	if c.debug {
		c.log.Printf("dbx/rollback: (%T)\n", c.Tx)
	}
	return c.Tx.Rollback()
}
