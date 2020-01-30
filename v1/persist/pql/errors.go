package pql

import (
	"errors"
	"fmt"
)

var (
	EOF                = errors.New("EOF")
	ErrInvalidEscape   = errors.New("Invalid escape sequence")
	ErrInvalidIdent    = errors.New("Invalid identifier")
	ErrInvalidQName    = errors.New("Invalid qualified name")
	ErrUnexpectedToken = errors.New("Unexpected token")
	ErrUnexpectedEOF   = errors.New("Unexpected end-of-file")
)

type Error struct {
	error
	span Span
}

func (e Error) Error() string {
	return fmt.Sprintf("%v %s", e.error.Error(), e.span.Describe())
}

func newErr(c error, s Span) *Error {
	return &Error{
		error: c,
		span:  s,
	}
}
