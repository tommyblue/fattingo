package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type backend struct {
	cfg *config
	db  store
}

func NewBackend(cfg *config) (*backend, error) {
	db, err := NewStore(cfg)
	if err != nil {
		return nil, err
	}

	bk := &backend{
		cfg: cfg,
		db:  db,
	}

	return bk, nil
}

func (b *backend) Run() {
	http.Handle("/", withLogs(withMetrics(rootHandler())))
	http.Handle("/customers", withLogs(withMetrics(customersHandler(b.db))))

	srv := &http.Server{
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func (b *backend) Stop() {
	b.db.Close()
}
