package errors

import (
	"fmt"

	"github.com/bww/go-util/v1/text"
)

type Error struct {
	cause error
	stmt  string
}

func New(err error, sql string) *Error {
	return &Error{
		cause: err,
		stmt:  sql,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v <%s>", e.cause, text.CollapseSpaces(e.stmt))
}

func (e *Error) Unwrap() error {
	return e.cause
}

// Statement produces the SQL statement that produced the error, if this is known.
func (e *Error) Statement() string {
	return e.stmt
}
