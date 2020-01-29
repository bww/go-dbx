package pql

import (
	"errors"
)

var (
	EOF              = errors.New("EOF")
	ErrInvalidEscape = errors.New("Invalid escape sequence")
)

type Error struct {
	error
	span Span
}

func newErr(c error, s Span) *Error {
	return &Error{
		error: c,
		span:  s,
	}
}
