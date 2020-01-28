package dbx

import (
	"log"
	"net/url"
	"os"

	"github.com/bww/go-util/debug"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	// "github.com/patrickmn/go-cache"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var defaultLogger = log.New(os.Stdout, "", 0)

type DB struct {
	*sqlx.DB
	log   *log.Logger
	debug bool
}

func New(dsn string, opts ...Option) (*DB, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	x, err := sqlx.Open(u.Scheme, dsn)
	if err != nil {
		return nil, err
	}

	d := &DB{DB: x, debug: debug.DEBUG}

	for _, e := range opts {
		d, err = e(d)
		if err != nil {
			return nil, err
		}
	}

	if d.log == nil {
		d.log = defaultLogger
	}

	err = d.Ping()
	if err != nil {
		return nil, err
	}

	return d, err
}

func (d *DB) Migrate(rc string) error {
	v, err := postgres.WithInstance(d.DB.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(rc, "postgres", v)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
