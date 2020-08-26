package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type backend struct {
	db store
}

func NewBackend(cfg *config) (*backend, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	bk := &backend{
		db: db,
	}
	return bk, nil
}

func (b *backend) Run() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/customers", b.customersHandler).Methods("GET")
	r.HandleFunc("/", b.rootHandler).Methods("GET")
	r.HandleFunc("/{path}", b.catchAllHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func (b *backend) Stop() {
	b.db.Close()
}

func (b *backend) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Fatt-in-Go!\n"))
}

type Customer struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (b *backend) customersHandler(w http.ResponseWriter, r *http.Request) {
	query, err := b.db.Query("SELECT id, title FROM customers ORDER BY title ASC")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// defer cancel()
	var customers []Customer
	for query.Next() {
		var c Customer
		err = query.Scan(&c.ID, &c.Title)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		customers = append(customers, c)
	}

	json.NewEncoder(w).Encode(customers)
}

func (b *backend) catchAllHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	vars := mux.Vars(r)
	w.Write([]byte(fmt.Sprintf("404 - Path '/%s' not found\n", vars["path"])))
}

func main() {
	time.Sleep(5 * time.Second)
	cfg, err := readConf()
	if err != nil {
		log.Fatal(err)
	}
	bk, err := NewBackend(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Add goroutines for recurring jobs
	log.Printf("Serving...")
	bk.Run()

}
