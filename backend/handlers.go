package fattingo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (b *Backend) rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("404 - Path '%s' not found\n", r.URL.Path)))
			return
		}
		w.Write([]byte("Fatt-in-Go!\n"))
	}
}

func (b *Backend) customersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			b.createCustomerHandler(w, r)
			return
		} else if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			log.Warnf("[%s] %s - Method not allowed", r.Method, r.URL)
			return
		}

		customers, err := b.db.Customers()

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Error(err)
			return
		}

		json.NewEncoder(w).Encode(customers)
	}
}

func (b *Backend) customerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerID, err := getURLQueryParam("id", w, r)
		if err != nil {
			log.Warn(err.Error())
			return
		}

		switch r.Method {
		case "DELETE":
			if err := b.db.DeleteCustomer(customerID); err != nil {
				var sErr *storeError
				if errors.As(err, &sErr) {
					http.Error(w, sErr.msg, sErr.status)
				} else {
					log.Error(err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		case "PUT":
			c, ok := b.populateCustomer(w, r)
			if !ok {
				return
			}

			customer, err := b.db.UpdateCustomer(customerID, c)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Error(err)
				return
			}
			json.NewEncoder(w).Encode(customer)
		case "GET":
			customer, err := b.db.Customer(customerID)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Error(err)
				return
			}

			json.NewEncoder(w).Encode(customer)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			log.Warnf("[%s] %s - Method not allowed", r.Method, r.URL)
			return
		}
	}
}

func (b *Backend) createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	c, ok := b.populateCustomer(w, r)
	if !ok {
		return
	}

	customer, err := b.db.CreateCustomer(c)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func (b *Backend) populateCustomer(w http.ResponseWriter, r *http.Request) (*customer, bool) {
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
		return nil, false
	}
	return &c, true
}
