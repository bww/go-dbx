package ident

import (
	"reflect"

	"github.com/bww/go-util/rand"
	"github.com/bww/go-util/ulid"
	"github.com/bww/go-util/uuid"
)

type Generator func() reflect.Value

func ULID() reflect.Value {
	return reflect.ValueOf(ulid.New())
}

func UUID() reflect.Value {
	return reflect.ValueOf(uuid.New())
}

func AlphaNumeric(n int) Generator {
	return func() reflect.Value {
		return reflect.ValueOf(rand.RandomString(n))
	}
}
