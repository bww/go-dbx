package persist

import (
	"reflect"

	"github.com/bww/go-util/rand"
	"github.com/bww/go-util/ulid"
	"github.com/bww/go-util/uuid"
)

type IdentFunc func() reflect.Value

func ULID() reflect.Value {
	return reflect.ValueOf(ulid.New())
}

func UUID() reflect.Value {
	return reflect.ValueOf(uuid.New())
}

func Random() reflect.Value {
	return reflect.ValueOf(rand.RandomString(32))
}
