package main

import (
	"database/sql"
	"fmt"
	"time"
)

// store mimics the go-sql-driver features. The primary target
// of using an interface here instead of the final type (that
// doesn't change) is to be able to mock the db when testing
type store interface {
	Close() error
	Ping() error
	Query(string, ...interface{}) (*sql.Rows, error)
}

func NewStore(cfg *config) (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
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

	return db, nil
}
