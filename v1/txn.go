package dbx

type TransactionHandler func(cxt Context) error

// Execute in a transaction. A transaction is created and the handler is invoked.
// If the handler returns a non-nil error the transaction is rolled back, otherwise
// the transaction is committed.
func (d *DB) Transaction(h TransactionHandler) error {
	tx, err := d.Beginx()
	if err != nil {
		return err
	}

	err = h(d.wrapTx(tx))
	if err == nil {
		err = tx.Commit()
	} else if terr := tx.Rollback(); terr != nil {
		d.log.Printf("Could not roll back transaction: %v", terr)
	}

	return err
}
