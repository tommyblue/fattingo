package fattingo

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (b *Backend) setupRoutes() error {

	r := b.router.PathPrefix("/api/v1").Subrouter()

	r.HandleFunc("/customers", b.customersHandler()).Methods("GET")
	r.HandleFunc("/customers", b.createCustomerHandler()).Methods("POST").HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/customer/{id:[0-9]+}", b.customerHandler()).Methods("GET")
	r.HandleFunc("/customer/{id:[0-9]+}", b.updateCustomerHandler()).Methods("PUT").HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/customer/{id:[0-9]+}", b.deleteCustomerHandler()).Methods("DELETE").HeadersRegexp("Content-Type", "application/json")

	b.router.HandleFunc("/", b.rootHandler())

	b.router.Use(loggingMiddleware)
	b.router.Use(metricsMiddleware)

	http.Handle("/", b.router)
	return nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("[%s] %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()
		next.ServeHTTP(w, r)
		log.Debugf("[%s] %s took %s", r.Method, r.URL, time.Since(began))
	})
}

func getURLVar(key string, r *http.Request) (string, bool) {
	vars := mux.Vars(r)
	value, ok := vars[key]
	return value, ok
}

func getIDVar(key string, r *http.Request) (int, error) {
	idStr, ok := getURLVar(key, r)
	if !ok {
		return 0, fmt.Errorf("[%s] %s - (400) Cannot get '%s' from URL path", r.Method, r.URL, key)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("[%s] %s - (400) '%s' isn't a number", r.Method, r.URL, key)
	}

	return id, nil
}
