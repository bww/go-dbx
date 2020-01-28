package dbx

import (
	"errors"
)

var (
	ErrInvalidField    = errors.New("Invalid field")
	ErrInvalidKeyCount = errors.New("Invalid primary key count")
	ErrNotAPointer     = errors.New("Not a pointer")
)
