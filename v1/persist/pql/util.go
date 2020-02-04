package pql

import (
	"fmt"
)

func stringer(v interface{}) string {
	switch c := v.(type) {
	case nil:
		return "<nil>"
	case string:
		return c
	case *string:
		return *c
	case []byte:
		return string(c)
	case *[]byte:
		return string(*c)
	default:
		return fmt.Sprint(v)
	}
}
