package dbx

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

const testDB = "dbx_v1_test"

func TestPersist(t *testing.T) {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=postgres sslmode=disable")
	if assert.NoError(t, err) {
		defer db.Close()
	}

	tx, err := db.Beginx()
	assert.NoError(t, err)
	assert.Equal(t, true, IsTx(tx))

}
