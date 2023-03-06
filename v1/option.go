package dbx

import (
	"log"
	"time"
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

func WithConnMaxLifetime(v time.Duration) Option {
	return func(d *DB) (*DB, error) {
		d.DB.SetConnMaxLifetime(v)
		return d, nil
	}
}

func WithConnMaxIdleTime(v time.Duration) Option {
	return func(d *DB) (*DB, error) {
		d.DB.SetConnMaxIdleTime(v)
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
