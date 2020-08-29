package fattingo

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type storeError struct {
	status int
	msg    string
}

func (err *storeError) Error() string {
	return err.msg
}

type dataStore interface {
	Customers() ([]*customer, error)
	Customer(int) (*customer, error)
	CreateCustomer(*customer) (*customer, error)
	DeleteCustomer(int) error
	Close() error
}

type database struct {
	*sql.DB
}

func newStore(cfg *Config) (dataStore, error) {
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
