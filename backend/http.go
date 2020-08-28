package fattingo

import (
	"encoding/json"
	"errors"
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
		if r.Method == "POST" {
			createCustomerHandler(db, w, r)
			return
		} else if r.Method != "GET" {
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
			log.Warnf("[%s] %s - Method not allowed", r.Method, r.URL)
			return
		}

		customerID, err := getURLQueryParam("id", w, r)
		if err != nil {
			log.Warn(err.Error())
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

func createCustomerHandler(db dataStore, w http.ResponseWriter, r *http.Request) {
	var c customer
	err := decodeJSONBody(w, r, &c)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	customer, err := db.CreateCustomer(&c)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
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
