package dbx

import (
	"log"
	"net/url"
	"os"

	"github.com/bww/go-upgrade/v1"
	"github.com/bww/go-upgrade/v1/driver/postgres"
	"github.com/bww/go-util/v1/debug"
	"github.com/jmoiron/sqlx"
	// "github.com/patrickmn/go-cache"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var defaultLogger = log.New(os.Stdout, "", 0)

type DB struct {
	*sqlx.DB
	log   *log.Logger
	dsn   string
	debug bool
}

func New(dsn string, opts ...Option) (*DB, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open(u.Scheme, dsn)
	if err != nil {
		return nil, err
	}

	return NewWithDB(db, opts...)
}

func NewWithDB(db *sqlx.DB, opts ...Option) (*DB, error) {
	var err error

	d := &DB{
		DB:    db,
		debug: debug.DEBUG,
		log:   defaultLogger,
	}

	for _, e := range opts {
		d, err = e(d)
		if err != nil {
			return nil, err
		}
	}

	err = d.Ping()
	if err != nil {
		return nil, err
	}

	return d, err
}

func (d *DB) Migrate(rc string) (upgrade.Results, error) {
	dr, err := postgres.New(d.dsn)
	if err != nil {
		return upgrade.Results{}, err
	}

	up, err := upgrade.New(upgrade.Config{
		Resources: rc,
		Driver:    dr,
	})
	if err != nil {
		return upgrade.Results{}, err
	}

	res, err := up.Upgrade()
	if err != nil {
		return upgrade.Results{}, err
	}

	return res, nil
}
