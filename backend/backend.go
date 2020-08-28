package fattingo

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Backend struct {
	cfg *Config
	db  dataStore
	srv *http.Server
}

func NewBackend(cfg *Config) (*Backend, error) {
	db, err := newStore(cfg)
	if err != nil {
		return nil, err
	}

	bk := &Backend{
		cfg: cfg,
		db:  db,
	}

	return bk, nil
}

func (b *Backend) Run() error {
	http.Handle("/customers", withLogs(withMetrics(customersHandler(b.db))))
	http.Handle("/customer", withLogs(withMetrics(customerHandler(b.db))))
	http.Handle("/", withLogs(withMetrics(rootHandler())))

	b.srv = &http.Server{
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Info("starting HTTP server...")
	return b.srv.ListenAndServe()
}

func (b *Backend) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	log.Info("stopping http server...")
	if err := b.srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %+v", err)
	}

	log.Info("closing db connection...")
	if err := b.db.Close(); err != nil {
		log.Fatalf("DB connection closing failed: %+v", err)
	}

	log.Info("exited properly")
}
