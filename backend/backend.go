package fattingo

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type Backend struct {
	cfg    *Config
	db     dataStore
	srv    *http.Server
	router *mux.Router
}

func NewBackend(cfg *Config) (*Backend, error) {
	db, err := newStore(cfg)
	if err != nil {
		return nil, err
	}

	bk := &Backend{
		cfg:    cfg,
		db:     db,
		router: mux.NewRouter(),
	}

	if err := bk.setupRoutes(); err != nil {
		return nil, err
	}

	return bk, nil
}

func (b *Backend) Run() error {
	b.srv = &http.Server{
		Handler:      b.router,
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
