package dbx

import (
	"log"
	"net/url"
	"os"

	"github.com/bww/go-upgrade/v1"
	"github.com/bww/go-upgrade/v1/driver/postgres"
	"github.com/bww/go-util/debug"
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

	x, err := sqlx.Open(u.Scheme, dsn)
	if err != nil {
		return nil, err
	}

	d := &DB{
		DB:    x,
		dsn:   dsn,
		debug: debug.DEBUG,
	}

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

func (d *DB) Migrate(rc string) (int, error) {
	up := upgrade.New(upgrade.Config{
		Resources: rc,
		Driver:    postgres.New(d.dsn),
	})

	rev, err := up.Upgrade()
	if err != nil {
		return -1, err
	}

	return rev, nil
}
