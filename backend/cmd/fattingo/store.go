package main

import "database/sql"

// store mimics the go-sql-driver features. The primary target
// of using an interface here instead of the final type (that
// doesn't change) is to be able to mock the db when testing
type store interface {
	Close() error
	Ping() error
	Query(string, ...interface{}) (*sql.Rows, error)
}
