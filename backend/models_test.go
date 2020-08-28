package fattingo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDB struct{}

func (mdb *mockDB) Close() error {
	return nil
}

func (mdb *mockDB) allCustomers() ([]*customer, error) {
	customers := make([]*customer, 0)
	customers = append(customers, &customer{
		ID: 1,
	})
	return customers, nil
}

func TestCustomersIndex(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/customers", nil)

	http.Handler(customersHandler(&mockDB{})).ServeHTTP(rec, req)

	want := "[{\"id\":1,\"title\":null,\"name\":null,\"surname\":null,\"address\":null,\"zip_code\":null,\"town\":null,\"province\":null,\"country\":null,\"tax_code\":null,\"vat\":null,\"info\":null}]\n"
	if want != rec.Body.String() {
		t.Errorf("want: %v, got: %v", want, rec.Body.String())
	}
}
