package dbx

import (
	"context"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bww/go-upgrade/v1"
	"github.com/bww/go-upgrade/v1/driver/postgres"
	"github.com/bww/go-upgrade/v1/driver/sqlite3"
	"github.com/bww/go-util/v1/debug"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var defaultLogger = log.New(os.Stdout, "", 0)

type database int

const (
	postgresDB database = iota
	sqliteDB
	unknownDB database = -1
)

func parseDB(v string) database {
	switch s := strings.ToLower(v); s {
	case "postgres", "postgresql":
		return postgresDB
	case "sqlite", "sqlite3":
		return sqliteDB
	default:
		return unknownDB
	}
}

type DB struct {
	*sqlx.DB
	backend database
	log     *log.Logger
	debug   bool
}

func New(dsn string, opts ...Option) (*DB, error) {
	var drv string

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	backend := parseDB(u.Scheme)
	switch backend {
	case postgresDB:
		drv, dsn = "postgres", dsn
	case sqliteDB:
		drv, dsn = "sqlite3", "file:"+u.Path+"?"+u.RawQuery
	default:
		drv, dsn = u.Scheme, dsn
	}

	x, err := sqlx.Open(drv, dsn)
	if err != nil {
		return nil, err
	}

	d := &DB{
		DB:      x,
		backend: backend,
		debug:   debug.DEBUG,
		log:     defaultLogger,
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
	var drv upgrade.Driver
	var err error

	switch d.backend {
	case postgresDB:
		drv, err = postgres.NewWithDB(d.DB.DB)
	case sqliteDB:
		drv, err = sqlite3.NewWithDB(d.DB.DB)
	default:
		err = ErrDriverNotSupported
	}
	if err != nil {
		return upgrade.Results{}, err
	}

	up, err := upgrade.New(upgrade.Config{
		Resources: rc,
		Driver:    drv,
	})
	if err != nil {
		return upgrade.Results{}, err
	}

	res, err := up.Upgrade()
	if err != nil {
		return res, err
	}

	return res, nil
}

func (d *DB) Monitor(cxt context.Context, iv time.Duration) <-chan error {
	if iv < time.Second {
		iv = time.Second
	}
	var errs chan error
	go func() {
	loop:
		for {
			select {
			case <-time.After(iv):
			case <-cxt.Done():
				break loop
			}
			if d.debug {
				d.log.Println("dbx: Polling database for connectivity")
			}
			xcx, cancel := context.WithTimeout(cxt, time.Second*10)
			err := d.PingContext(xcx)
			cancel() // clean up if we complete before the timeout
			if err != nil {
				errs <- err
			} else if d.debug {
				d.log.Println("dbx: Connection OK")
			}
		}
		close(errs)
	}()
	return errs
}
