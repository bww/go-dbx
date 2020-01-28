package test

import (
	"database/sql"
	"fmt"
)

func createDatabase(u, n string) error {
	db, err := sql.Open("postgres", u)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	var exists int
	err = db.QueryRow(fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", n)).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if exists == 1 {
		return nil // ok
	}

	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, n))
	if err != nil {
		return err
	}

	return nil
}

func dropDatabase(u, n string) error {
	db, err := sql.Open("postgres", u)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, n))
	if err != nil {
		return err
	}

	return nil
}
