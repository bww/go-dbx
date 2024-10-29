package errors

import (
	"errors"
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

// Wrap produces a new [github.com/bww/go-dbx/errors.Error] only if the provided
// error is not already one.
func Wrap(err error, sql string) error {
	var dbxerr *Error
	if errors.As(err, &dbxerr) {
		return err
	} else {
		return New(err, sql)
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v <%s>", e.cause, text.CollapseSpaces(e.stmt))
}

func (e *Error) Unwrap() error {
	return e.cause
}

// Statement produces the SQL statement that produced the error, if this is
// known.
func (e *Error) Statement() string {
	return e.stmt
}
