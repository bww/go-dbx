package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	assert.Equal(t, "ASC", Ascending.String())
	assert.Equal(t, "ASC", Order(9).String())
	assert.Equal(t, "DESC", Descending.String())
}
