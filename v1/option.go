package dbx

import (
	"log"
)

type Option func(d *DB) (*DB, error)

func WithMaxOpenConns(v int) Option {
	return func(d *DB) (*DB, error) {
		d.DB.SetMaxOpenConns(v)
		return d, nil
	}
}

func WithMaxIdleConns(v int) Option {
	return func(d *DB) (*DB, error) {
		d.DB.SetMaxIdleConns(v)
		return d, nil
	}
}

func WithDebug(on bool) Option {
	return func(d *DB) (*DB, error) {
		d.debug = on
		return d, nil
	}
}

func WithLogger(l *log.Logger) Option {
	return func(d *DB) (*DB, error) {
		d.log = l
		return d, nil
	}
}
