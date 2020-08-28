package fattingo

import (
	"database/sql"
	"fmt"
	"time"
)

// store mimics the go-sql-driver features. The primary target
// of using an interface here instead of the final type (that
// doesn't change) is to be able to mock the db when testing
type dataStore interface {
	Customers() ([]*customer, error)
	Customer(int) (*customer, error)
	Close() error
}

type database struct {
	*sql.DB
}

func newStore(cfg *Config) (*database, error) {
	connString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		cfg.dbUser,
		cfg.dbPassword,
		cfg.dbHost,
		cfg.dbPort,
		cfg.dbName)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &database{db}, nil
}
