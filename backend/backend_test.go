package fattingo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockDB struct {
	customers []*customer
}

func (mdb *mockDB) Close() error {
	return nil
}
func (mdb *mockDB) Customer(id int) (*customer, error) {
	return mdb.customers[0], nil
}

func (mdb *mockDB) Customers() ([]*customer, error) {
	return mdb.customers, nil
}

func (mdb *mockDB) CreateCustomer(c *customer) (*customer, error) {
	c.ID = 9999
	mdb.customers = append(mdb.customers, c)
	return c, nil
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

	customers := make([]*customer, 0)
	n := "NameTest"
	s := "SurnameTest"
	customers = append(customers, &customer{
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

	var c []*customer
	json.Unmarshal(rec.Body.Bytes(), &c)

	if !reflect.DeepEqual(customers, c) {
		t.Errorf("want: %v, got: %v", customers, c)
	}
}

func TestCustomer(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/customer?id=1", nil)

	customers := make([]*customer, 0)
	n := "NameTest"
	s := "SurnameTest"
	customers = append(customers, &customer{
		ID:      1,
		Name:    &n,
		Surname: &s,
	})
	b := &Backend{
		db: &mockDB{
			customers: customers,
		},
	}

	http.Handler(b.customerHandler()).ServeHTTP(rec, req)

	var c *customer
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
	want := &customer{
		ID:      9999,
		Title:   &tl,
		Name:    &n,
		Surname: &s,
	}
	jsonStr := []byte(fmt.Sprintf(`{"title":"%s", "name":"%s", "surname":"%s"}`, tl, n, s))
	req, _ := http.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(jsonStr))

	b := &Backend{
		db: &mockDB{
			customers: make([]*customer, 0),
		},
	}
	http.Handler(b.customersHandler()).ServeHTTP(rec, req)

	var c *customer
	json.Unmarshal(rec.Body.Bytes(), &c)

	if !reflect.DeepEqual(want, c) {
		t.Errorf("want: %v, got: %v", want, c)
	}
}

func TestDeleteCustomer(t *testing.T) {
	customers := make([]*customer, 0)
	for i := 0; i < 3; i++ {
		n := fmt.Sprintf("NameTest%d", i)
		s := fmt.Sprintf("SurnameTest%d", i)
		customers = append(customers, &customer{
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
	http.Handler(b.customersHandler()).ServeHTTP(rec, req)

	var oldCustomers []*customer
	json.Unmarshal(rec.Body.Bytes(), &oldCustomers)

	// Delete 1 customer
	delRec := httptest.NewRecorder()
	delReq, _ := http.NewRequest("DELETE", "/api/v1/customers?id=2", nil)
	http.Handler(b.customerHandler()).ServeHTTP(delRec, delReq)

	// Reload customers
	newRec := httptest.NewRecorder()
	newReq, _ := http.NewRequest("GET", "/api/v1/customers", nil)
	http.Handler(b.customersHandler()).ServeHTTP(newRec, newReq)

	var newCustomers []*customer
	json.Unmarshal(newRec.Body.Bytes(), &newCustomers)

	// Check
	if len(newCustomers) != len(oldCustomers)-1 {
		t.Errorf("want: %d, got: %d", len(oldCustomers)-1, len(newCustomers))
	}
}
