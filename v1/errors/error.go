package errors

import (
	"fmt"

	"github.com/bww/go-util/v1/text"
)

type SQLError struct {
	cause error
	stmt  string
}

func NewWithSQL(err error, sql string) *SQLError {
	return &SQLError{
		cause: err,
		stmt:  sql,
	}
}

func (e *SQLError) Error() string {
	return e.cause.Error()
}

func (e *SQLError) String() string {
	return fmt.Sprintf("%v <%s>", e.Error(), text.CollapseSpaces(e.stmt))
}

func (e *SQLError) Unwrap() error {
	return e.cause
}

// Statement produces the SQL statement that produced the error, if this is known.
func (e *SQLError) Statement() string {
	return e.stmt
}
