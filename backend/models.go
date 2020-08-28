package fattingo

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type user struct {
	ID       int     `json:"id"`
	Title    *string `json:"title"`
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Address  *string `json:"address"`
	ZipCode  *string `json:"zip_code"`
	Town     *string `json:"town"`
	Province *string `json:"province"`
	Country  *string `json:"country"`
	TaxCode  *string `json:"tax_code"`
	Vat      *string `json:"vat"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
}

type customer struct {
	ID       int     `json:"id"`
	UserID   int     `json:"user_id,omitempty"`
	Title    *string `json:"title"`
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Address  *string `json:"address"`
	ZipCode  *string `json:"zip_code"`
	Town     *string `json:"town"`
	Province *string `json:"province"`
	Country  *string `json:"country"`
	TaxCode  *string `json:"tax_code"`
	Vat      *string `json:"vat"`
	Info     *string `json:"info"`
}

func (db *database) Customers() ([]*customer, error) {
	rows, err := db.Query(`
	SELECT
		id,
		title,
		name,
		surname,
		address,
		zip_code,
		town,
		province,
		country,
		tax_code,
		vat,
		info
	FROM customers
	ORDER BY id DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := make([]*customer, 0)
	for rows.Next() {
		c := &customer{}
		err = rows.Scan(
			&c.ID,
			&c.Title,
			&c.Name,
			&c.Surname,
			&c.Address,
			&c.ZipCode,
			&c.Town,
			&c.Province,
			&c.Country,
			&c.TaxCode,
			&c.Vat,
			&c.Info,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (db *database) Customer(id int) (*customer, error) {
	c := &customer{}
	err := db.QueryRow(`
	SELECT
		id,
		title,
		name,
		surname,
		address,
		zip_code,
		town,
		province,
		country,
		tax_code,
		vat,
		info
	FROM customers
	WHERE id = ?`, id).Scan(
		&c.ID,
		&c.Title,
		&c.Name,
		&c.Surname,
		&c.Address,
		&c.ZipCode,
		&c.Town,
		&c.Province,
		&c.Country,
		&c.TaxCode,
		&c.Vat,
		&c.Info,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (db *database) CreateCustomer(c *customer) (*customer, error) {
	sqlStatement := `
INSERT INTO customers (
	title,
	name,
	surname,
	address,
	zip_code,
	town,
	province,
	country,
	tax_code,
	vat,
	info,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
`
	now := time.Now()
	res, err := db.Exec(sqlStatement,
		c.Title,
		c.Name,
		c.Surname,
		c.Address,
		c.ZipCode,
		c.Town,
		c.Province,
		c.Country,
		c.TaxCode,
		c.Vat,
		c.Info,
		now,
		now,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return db.Customer(int(id))
}

func (db *database) DeleteCustomer(id int) error {
	res, err := db.Exec(`DELETE FROM customers WHERE id=?;`, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		msg := fmt.Sprintf("Can't find the customer with id %d", id)
		return &storeError{status: http.StatusNotFound, msg: msg}
	}

	if n > 1 {
		msg := fmt.Sprintf("Too many deleted customers (%d)", n)
		log.Errorf("%s with id %d", msg, id)
		return &storeError{status: http.StatusInternalServerError, msg: msg}
	}

	return nil
}
