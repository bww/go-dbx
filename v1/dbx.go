package dbx

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/bww/go-upgrade/v1"
	"github.com/bww/go-upgrade/v1/driver/postgres"
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
	unknownDB = -1
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
	log   *log.Logger
	dsn   string
	debug bool
}

func New(dsn string, opts ...Option) (*DB, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	var drv string
	switch parseDB(u.Scheme) {
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
		DB:    x,
		dsn:   dsn,
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
	u, err := url.Parse(d.dsn)
	if err != nil {
		return upgrade.Results{}, err
	}

	var drv upgrade.Driver
	switch parseDB(u.Scheme) {
	case postgresDB:
		drv, err = postgres.New(d.dsn)
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
		return upgrade.Results{}, err
	}

	return res, nil
}
