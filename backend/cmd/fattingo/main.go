package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Fatt-in-Go!\n"))
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	vars := mux.Vars(r)
	w.Write([]byte(fmt.Sprintf("404 - Path '/%s' not found\n", vars["path"])))
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/{path}", catchAllHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
