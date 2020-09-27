package fattingo

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type frontendHandler struct {
	staticPath string
	indexPath  string
}

func (h frontendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func (b *Backend) setupRoutes() error {

	r := b.router.PathPrefix("/api/v1").Subrouter()

	r.HandleFunc("/customers", b.customersHandler()).Methods(http.MethodGet)
	r.HandleFunc("/customers", b.createCustomerHandler()).Methods(http.MethodPost).HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/customer/{id:[0-9]+}", b.customerHandler()).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id:[0-9]+}", b.updateCustomerHandler()).Methods(http.MethodPut).HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/customer/{id:[0-9]+}", b.deleteCustomerHandler()).Methods(http.MethodDelete).HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	r.Use(jsonResponseMiddleware)

	spa := frontendHandler{staticPath: "../frontend/build", indexPath: "index.html"}
	b.router.PathPrefix("/").Handler(spa)

	b.router.HandleFunc("/", b.rootHandler())

	b.router.Use(loggingMiddleware)
	b.router.Use(metricsMiddleware)
	b.router.Use(mux.CORSMethodMiddleware(b.router))

	http.Handle("/", b.router)
	return nil
}

func jsonResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
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
