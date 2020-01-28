package test

import (
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-util/debug"
)

const initDB = "postgres"

func dburl(n string) string {
	return fmt.Sprintf("postgres://postgres@localhost/%s?sslmode=disable", url.PathEscape(n))
}

var (
	initOnce sync.Once
	sharedDB *dbx.DB
)

func Init(name string, opts ...Option) {
	initOnce.Do(func() {
		var err error
		conf := &Config{
			Name: name,
		}

		for _, o := range opts {
			conf, err = o(conf)
			if err != nil {
				panic(err)
			}
		}

		err = teardown(conf)
		if err != nil {
			panic(err)
		}

		err = setup(conf)
		if err != nil {
			panic(err)
		}
	})
}

func DB() *dbx.DB {
	if sharedDB == nil {
		panic("Test database not initialized; did you call test.Init() in this package?")
	}
	return sharedDB
}

func setup(conf *Config) error {
	debug.DEBUG = debug.DEBUG || istrue(os.Getenv("DBX_DEBUG"))
	debug.VERBOSE = debug.VERBOSE || istrue(os.Getenv("DBX_VERBOSE"))
	debug.TRACE = debug.TRACE || istrue(os.Getenv("DBX_TRACE"))

	err := createDatabase(dburl(initDB), conf.Name)
	if err != nil {
		return fmt.Errorf("Creating %s (from %s): %v", conf.Name, initDB, err)
	}

	sharedDB, err = dbx.New(dburl(conf.Name))
	if err != nil {
		return err
	}

	if conf.Migrations != "" {
		err = sharedDB.Migrate(conf.Migrations)
		if err != nil {
			return err
		}
	}

	return nil
}

func teardown(conf *Config) error {
	err := dropDatabase(dburl(initDB), conf.Name)
	if err != nil {
		return err
	}
	return nil
}

func istrue(s string) bool {
	return s == "true" || s == "t" || s == "yes" || s == "y"
}
