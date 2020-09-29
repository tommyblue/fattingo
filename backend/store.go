package fattingo

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tommyblue/fattingo/backend/model"
)

type dataStore interface {
	Customers() ([]*model.Customer, error)
	Customer(int) (*model.Customer, error)
	CreateCustomer(*model.Customer) (*model.Customer, error)
	UpdateCustomer(int, *model.Customer) (*model.Customer, error)
	DeleteCustomer(int) error
	CustomerInfo(int) (*model.CustomerInfo, error)
	Close() error
}

func newStore(cfg *Config) (dataStore, error) {
	var db *sql.DB
	var err error

	switch cfg.dbType {
	case "mysql":
		connString := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.dbUser,
			cfg.dbPassword,
			cfg.dbHost,
			cfg.dbPort,
			cfg.dbName)

		db, err = sql.Open("mysql", connString)
		if err != nil {
			return nil, err
		}

		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
	case "sqlite3":
		db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s", cfg.dbPath))
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unknown db type %s", cfg.dbType)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &model.Database{db}, nil
}
