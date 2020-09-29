package fattingo

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tommyblue/fattingo/backend/model"
)

type mockDB struct {
	customers []*model.Customer
}

func (mdb *mockDB) Close() error {
	return nil
}
func (mdb *mockDB) Customer(id int) (*model.Customer, error) {
	return mdb.customers[0], nil
}

func (mdb *mockDB) Customers() ([]*model.Customer, error) {
	return mdb.customers, nil
}

func (mdb *mockDB) CreateCustomer(c *model.Customer) (*model.Customer, error) {
	c.ID = 9999
	mdb.customers = append(mdb.customers, c)
	return c, nil
}

func (mdb *mockDB) UpdateCustomer(id int, c *model.Customer) (*model.Customer, error) {
	for i, cust := range mdb.customers {
		if cust.ID == id {
			mdb.customers[i].Title = c.Title
			mdb.customers[i].Name = c.Name
			mdb.customers[i].Surname = c.Surname
			mdb.customers[i].Address = c.Address
			mdb.customers[i].ZipCode = c.ZipCode
			mdb.customers[i].Town = c.Town
			mdb.customers[i].Province = c.Province
			mdb.customers[i].Country = c.Country
			mdb.customers[i].TaxCode = c.TaxCode
			mdb.customers[i].Vat = c.Vat
			mdb.customers[i].Info = c.Info

			return mdb.customers[i], nil
		}
	}
	return nil, errors.New("Cannot find the customer")
}

func (mdb *mockDB) DeleteCustomer(id int) error {
	for i, c := range mdb.customers {
		if c.ID == id {
			copy(mdb.customers[i:], mdb.customers[i+1:])
			mdb.customers[len(mdb.customers)-1] = nil
			mdb.customers = mdb.customers[:len(mdb.customers)-1]
			return nil
		}
	}
	return errors.New("Cannot find the customer")
}

func TestCustomers(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/customers", nil)

	customers := make([]*model.Customer, 0)
	n := "NameTest"
	s := "SurnameTest"
	customers = append(customers, &model.Customer{
		ID:      1,
		Name:    &n,
		Surname: &s,
	})
	b := &Backend{
		db: &mockDB{
			customers: customers,
		},
	}
	http.Handler(b.customersHandler()).ServeHTTP(rec, req)

	var c []*model.Customer
	json.Unmarshal(rec.Body.Bytes(), &c)

	if !reflect.DeepEqual(customers, c) {
		t.Errorf("want: %v, got: %v", customers, c)
	}
}

func TestCustomer(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/customer/1", nil)

	customers := make([]*model.Customer, 0)
	n := "NameTest"
	s := "SurnameTest"
	customers = append(customers, &model.Customer{
		ID:      1,
		Name:    &n,
		Surname: &s,
	})

	b := &Backend{
		db: &mockDB{
			customers: customers,
		},
		router: mux.NewRouter(),
	}

	if err := b.setupRoutes(); err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/customer/{id:[0-9]+}", b.customerHandler())
	router.ServeHTTP(rec, req)

	var c *model.Customer
	json.Unmarshal(rec.Body.Bytes(), &c)

	if !reflect.DeepEqual(customers[0], c) {
		t.Errorf("want: %v, got: %v", customers[0], c)
	}
}

func TestCreateCustomer(t *testing.T) {
	rec := httptest.NewRecorder()
	tl := "CreateTitleTest"
	n := "CreateNameTest"
	s := "CreateSurnameTest"
	want := &model.Customer{
		ID:      9999,
		Title:   &tl,
		Name:    &n,
		Surname: &s,
	}
	jsonStr := []byte(fmt.Sprintf(`{"title":"%s", "name":"%s", "surname":"%s"}`, tl, n, s))
	req, _ := http.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(jsonStr))

	b := &Backend{
		db: &mockDB{
			customers: make([]*model.Customer, 0),
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/customers", b.createCustomerHandler())
	router.ServeHTTP(rec, req)

	var c *model.Customer
	json.Unmarshal(rec.Body.Bytes(), &c)

	if !reflect.DeepEqual(want, c) {
		t.Errorf("want: %v, got: %v", want, c)
	}
}

func TestDeleteCustomer(t *testing.T) {
	customers := make([]*model.Customer, 0)
	for i := 0; i < 3; i++ {
		n := fmt.Sprintf("NameTest%d", i)
		s := fmt.Sprintf("SurnameTest%d", i)
		customers = append(customers, &model.Customer{
			ID:      i,
			Name:    &n,
			Surname: &s,
		})
	}
	b := &Backend{
		db: &mockDB{
			customers: customers,
		},
	}

	// Load customers
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/customers", nil)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/customers", b.customersHandler())
	router.ServeHTTP(rec, req)

	var oldCustomers []*model.Customer
	json.Unmarshal(rec.Body.Bytes(), &oldCustomers)

	// Delete 1 customer
	delRec := httptest.NewRecorder()
	delReq, _ := http.NewRequest("DELETE", "/api/v1/customer/2", nil)
	router.HandleFunc("/api/v1/customer/{id:[0-9]+}", b.deleteCustomerHandler())
	router.ServeHTTP(delRec, delReq)

	// Reload customers
	newRec := httptest.NewRecorder()
	newReq, _ := http.NewRequest("GET", "/api/v1/customers", nil)
	router.HandleFunc("/api/v1/customers", b.customersHandler())
	router.ServeHTTP(newRec, newReq)

	var newCustomers []*model.Customer
	json.Unmarshal(newRec.Body.Bytes(), &newCustomers)

	// Check
	if len(newCustomers) != len(oldCustomers)-1 {
		t.Errorf("want: %d, got: %d", len(oldCustomers)-1, len(newCustomers))
	}
}

func TestUpdateCustomer(t *testing.T) {
	customers := make([]*model.Customer, 0)
	n := "NameTest"
	s := "SurnameTest"
	customers = append(customers, &model.Customer{
		ID:      1,
		Name:    &n,
		Surname: &s,
	})
	b := &Backend{
		db: &mockDB{
			customers: customers,
		},
	}

	tl1 := "UpdatedTitleTest"
	n1 := "UpdatedNameTest"
	s1 := "UpdatedSurnameTest"
	jsonStr := []byte(fmt.Sprintf(`{"title":"%s", "name":"%s", "surname":"%s"}`, tl1, n1, s1))
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/customer/%d", customers[0].ID), bytes.NewBuffer(jsonStr))

	if err != nil {
		t.Fatal(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/customer/{id:[0-9]+}", b.updateCustomerHandler())
	router.ServeHTTP(rec, req)

	var c *model.Customer
	json.Unmarshal(rec.Body.Bytes(), &c)

	if *c.Title != tl1 {
		t.Fatalf("want: %v, got: %v", tl1, *c.Title)
	}
	if *c.Name != n1 {
		t.Fatalf("want: %v, got: %v", n1, *c.Name)
	}
	if *c.Surname != s1 {
		t.Fatalf("want: %v, got: %v", s1, *c.Surname)
	}
}

func TestCustomerLifeCycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip in short mode")
	}

	dbPath := "./test-db.sqlite"
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", dbPath))
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dbPath)

	_, err = db.Exec(`CREATE TABLE customers (
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		title TEXT,
		name TEXT,
		surname TEXT,
		address TEXT,
		zip_code TEXT,
		town TEXT,
		province TEXT,
		country TEXT,
		tax_code TEXT,
		vat TEXT,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		info TEXT)
	`)
	if err != nil {
		t.Fatal(err)
	}
	bk := &Backend{
		db: &model.Database{db},
	}

	customers, err := bk.db.Customers()
	if err != nil {
		t.Fatal(err)
	}

	if len(customers) != 0 {
		t.Fatalf("customers, want: 0, got: %d", len(customers))
	}

	tl := "TitleTest"
	n := "NameTest"
	s := "SurnameTest"
	c := &model.Customer{
		Title:   &tl,
		Name:    &n,
		Surname: &s,
	}

	_, err = bk.db.CreateCustomer(c)
	if err != nil {
		t.Fatal(err)
	}

	customers, err = bk.db.Customers()
	if err != nil {
		t.Fatal(err)
	}

	if len(customers) != 1 {
		t.Fatalf("customers, want: 1, got: %d", len(customers))
	}

	newCustomer, err := bk.db.Customer(customers[0].ID)

	if *newCustomer.Name != n {
		t.Fatalf("Name, want: %s, got: %s", n, *newCustomer.Name)
	}
	if *newCustomer.Surname != s {
		t.Fatalf("Surname, want: %s, got: %s", s, *newCustomer.Surname)
	}
	if *newCustomer.Title != tl {
		t.Fatalf("Title, want: %s, got: %s", tl, *newCustomer.Title)
	}

	tl1 := "UpdatedTitleTest"
	n1 := "UpdatedNameTest"
	s1 := "UpdatedSurnameTest"
	c1 := &model.Customer{
		Title:   &tl1,
		Name:    &n1,
		Surname: &s1,
	}
	updatedCustomer, err := bk.db.UpdateCustomer(customers[0].ID, c1)
	if err != nil {
		t.Fatal(err)
	}

	if *updatedCustomer.Title != tl1 {
		t.Fatalf("updated title, want %s, got: %s", tl1, *updatedCustomer.Title)
	}
	if *updatedCustomer.Name != n1 {
		t.Fatalf("updated name, want %s, got: %s", n1, *updatedCustomer.Name)
	}
	if *updatedCustomer.Surname != s1 {
		t.Fatalf("updated surname, want %s, got: %s", s1, *updatedCustomer.Surname)
	}

	if err := bk.db.DeleteCustomer(updatedCustomer.ID); err != nil {
		t.Fatal(err)
	}

	customers, err = bk.db.Customers()
	if err != nil {
		t.Fatal(err)
	}

	if len(customers) != 0 {
		t.Fatalf("customers, want: 0, got: %d", len(customers))
	}

	if err := bk.db.DeleteCustomer(updatedCustomer.ID); err == nil {
		t.Fatal("Deleting not existing customer, want error, got nil")
	}
}
