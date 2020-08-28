package fattingo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	db := &mockDB{
		customers: customers,
	}
	http.Handler(customersHandler(db)).ServeHTTP(rec, req)

	want := "[{\"id\":1,\"title\":null,\"name\":\"NameTest\",\"surname\":\"SurnameTest\",\"address\":null,\"zip_code\":null,\"town\":null,\"province\":null,\"country\":null,\"tax_code\":null,\"vat\":null,\"info\":null}]\n"
	if want != rec.Body.String() {
		t.Errorf("want: %v, got: %v", want, rec.Body.String())
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
	db := &mockDB{
		customers: customers,
	}

	http.Handler(customerHandler(db)).ServeHTTP(rec, req)

	want := "{\"id\":1,\"title\":null,\"name\":\"NameTest\",\"surname\":\"SurnameTest\",\"address\":null,\"zip_code\":null,\"town\":null,\"province\":null,\"country\":null,\"tax_code\":null,\"vat\":null,\"info\":null}\n"
	if want != rec.Body.String() {
		t.Errorf("want: %v, got: %v", want, rec.Body.String())
	}
}

func TestCreateCustomer(t *testing.T) {
	rec := httptest.NewRecorder()
	jsonStr := []byte(`{"title":"CreateTitleTest", "name":"CreateNameTest", "surname":"CreateSurnameTest"}`)
	req, _ := http.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(jsonStr))

	db := &mockDB{
		customers: make([]*customer, 0),
	}

	http.Handler(customersHandler(db)).ServeHTTP(rec, req)

	want := "{\"id\":9999,\"title\":\"CreateTitleTest\",\"name\":\"CreateNameTest\",\"surname\":\"CreateSurnameTest\",\"address\":null,\"zip_code\":null,\"town\":null,\"province\":null,\"country\":null,\"tax_code\":null,\"vat\":null,\"info\":null}\n"
	if want != rec.Body.String() {
		t.Errorf("want: %v, got: %v", want, rec.Body.String())
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
	db := &mockDB{
		customers: customers,
	}

	// Load customers
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/customers", nil)
	http.Handler(customersHandler(db)).ServeHTTP(rec, req)

	var oldCustomers []*customer
	json.Unmarshal(rec.Body.Bytes(), &oldCustomers)

	// Delete 1 customer
	delRec := httptest.NewRecorder()
	delReq, _ := http.NewRequest("DELETE", "/api/v1/customers?id=2", nil)
	http.Handler(customerHandler(db)).ServeHTTP(delRec, delReq)

	// Reload customers
	newRec := httptest.NewRecorder()
	newReq, _ := http.NewRequest("GET", "/api/v1/customers", nil)
	http.Handler(customersHandler(db)).ServeHTTP(newRec, newReq)

	var newCustomers []*customer
	json.Unmarshal(newRec.Body.Bytes(), &newCustomers)

	// Check
	if len(newCustomers) != len(oldCustomers)-1 {
		t.Errorf("want: %d, got: %d", len(oldCustomers)-1, len(newCustomers))
	}
}
