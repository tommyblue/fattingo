package fattingo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func rootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("404 - Path '%s' not found\n", r.URL.Path)))
			return
		}
		w.Write([]byte("Fatt-in-Go!\n"))
	})
}

func customersHandler(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			log.Warnf("[%s] %s - Method not allowed", r.Method, r.URL)
			return
		}

		customers, err := db.Customers()

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Error(err)
			return
		}

		json.NewEncoder(w).Encode(customers)
	})
}

func customerHandler(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			log.Warnf("[%s] %s - (405) Method not allowed", r.Method, r.URL)
			return
		}

		keys, ok := r.URL.Query()["id"]

		if !ok || len(keys[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("(400) Missing 'id' param\n"))
			log.Warnf("[%s] %s - (400) Missing 'id' param", r.Method, r.URL)
			return
		}

		customerID, err := strconv.Atoi(keys[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("(400) Wrong 'id' param %s\n", keys[0])))
			log.Warnf("[%s] %s - (400) Missing 'id' param", r.Method, r.URL)
			return
		}
		customer, err := db.Customer(customerID)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Error(err)
			return
		}

		json.NewEncoder(w).Encode(customer)
	})
}

func withLogs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("[%s] %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func withMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()
		next.ServeHTTP(w, r)
		log.Debugf("[%s] %s took %s", r.Method, r.URL, time.Since(began))
	})
}
