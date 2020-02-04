package dbx

import (
	"errors"
)

var (
	ErrNotFound        = errors.New("Not found")
	ErrImmutable       = errors.New("Immutable")
	ErrInvalidField    = errors.New("Invalid field")
	ErrInvalidKeyCount = errors.New("Invalid primary key count")
	ErrNotAPointer     = errors.New("Not a pointer")
)
