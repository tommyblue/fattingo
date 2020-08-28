package fattingo

import (
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

func TestCustomers(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/customers", nil)

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
	req, _ := http.NewRequest("GET", "/customer?id=1", nil)

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
