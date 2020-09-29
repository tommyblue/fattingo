package model

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func (db *Database) Customers() ([]*Customer, error) {
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

	customers := make([]*Customer, 0)
	for rows.Next() {
		c := &Customer{}
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

func (db *Database) Customer(id int) (*Customer, error) {
	c := &Customer{}
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

func (db *Database) CreateCustomer(c *Customer) (*Customer, error) {
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

func (db *Database) UpdateCustomer(id int, c *Customer) (*Customer, error) {
	sqlStatement := `
UPDATE customers SET
	title = ?,
	name = ?,
	surname = ?,
	address = ?,
	zip_code = ?,
	town = ?,
	province = ?,
	country = ?,
	tax_code = ?,
	vat = ?,
	info = ?,
	updated_at = ?
WHERE id = ?;
`
	now := time.Now()
	_, err := db.Exec(sqlStatement,
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
		id,
	)

	if err != nil {
		return nil, err
	}

	return db.Customer(id)
}

func (db *Database) DeleteCustomer(id int) error {
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
		return &DbError{Status: http.StatusNotFound, Msg: msg}
	}

	if n > 1 {
		msg := fmt.Sprintf("Too many deleted customers (%d)", n)
		log.Errorf("%s with id %d", msg, id)
		return &DbError{Status: http.StatusInternalServerError, Msg: msg}
	}

	return nil
}

func (db *Database) CustomerInfo(id int) (*CustomerInfo, error) {
	c, err := db.Customer(id)
	if err != nil {
		return nil, err
	}

	s, err := db.Slips(id)
	if err != nil {
		return nil, err
	}

	return &CustomerInfo{
		Customer: c,
		Slips:    s,
	}, nil
}
