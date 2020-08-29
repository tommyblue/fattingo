package fattingo

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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
	http.HandleFunc("/api/v1/customers", withLogs(withMetrics(customersHandler(b.db))))
	http.HandleFunc("/api/v1/customer", withLogs(withMetrics(customerHandler(b.db))))
	http.HandleFunc("/", withLogs(withMetrics(rootHandler())))

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

func withLogs(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("[%s] %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	}
}

func withMetrics(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()
		next.ServeHTTP(w, r)
		log.Debugf("[%s] %s took %s", r.Method, r.URL, time.Since(began))
	}
}

func getURLQueryParam(key string, w http.ResponseWriter, r *http.Request) (int, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("(400) Missing '%s' param\n", key)))
		return 0, fmt.Errorf("[%s] %s - (400) Missing '%s' param", r.Method, r.URL, key)
	}

	value, err := strconv.Atoi(keys[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("(400) Wrong '%s' param %s\n", key, keys[0])))
		return 0, fmt.Errorf("[%s] %s - (400) Missing '%s' param", r.Method, r.URL, key)
	}

	return value, nil
}
