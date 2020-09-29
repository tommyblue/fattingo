package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Database struct {
	*sql.DB
}

type DbError struct {
	Status int
	Msg    string
}

func (err *DbError) Error() string {
	return err.Msg
}

// +------------------------+--------------+------+-----+---------+----------------+
// | Field                  | Type         | Null | Key | Default | Extra          |
// +------------------------+--------------+------+-----+---------+----------------+
// | id                     | int          | NO   | PRI | NULL    | auto_increment |
// | title                  | varchar(255) | YES  |     | NULL    |                |
// | name                   | varchar(255) | NO   |     | NULL    |                |
// | surname                | varchar(255) | NO   |     | NULL    |                |
// | address                | varchar(255) | NO   |     | NULL    |                |
// | zip_code               | varchar(255) | NO   |     | NULL    |                |
// | town                   | varchar(255) | NO   |     | NULL    |                |
// | province               | varchar(255) | NO   |     | NULL    |                |
// | country                | varchar(255) | YES  |     | NULL    |                |
// | tax_code               | varchar(255) | NO   |     | NULL    |                |
// | vat                    | varchar(255) | NO   |     | NULL    |                |
// | phone                  | varchar(255) | NO   |     | NULL    |                |
// | created_at             | datetime     | NO   |     | NULL    |                |
// | updated_at             | datetime     | NO   |     | NULL    |                |
// | email                  | varchar(255) | NO   | UNI |         |                |
// | encrypted_password     | varchar(128) | NO   |     |         |                |
// | reset_password_token   | varchar(255) | YES  | UNI | NULL    |                |
// | reset_password_sent_at | datetime     | YES  |     | NULL    |                |
// | remember_created_at    | datetime     | YES  |     | NULL    |                |
// | sign_in_count          | int          | YES  |     | 0       |                |
// | current_sign_in_at     | datetime     | YES  |     | NULL    |                |
// | last_sign_in_at        | datetime     | YES  |     | NULL    |                |
// | current_sign_in_ip     | varchar(255) | YES  |     | NULL    |                |
// | last_sign_in_ip        | varchar(255) | YES  |     | NULL    |                |
// | bank_coordinates       | text         | YES  |     | NULL    |                |
// | language               | varchar(255) | YES  |     | it      |                |
// | logo_file_name         | varchar(255) | YES  |     | NULL    |                |
// | logo_content_type      | varchar(255) | YES  |     | NULL    |                |
// | logo_file_size         | int          | YES  |     | NULL    |                |
// | logo_updated_at        | datetime     | YES  |     | NULL    |                |
// +------------------------+--------------+------+-----+---------+----------------+
type User struct {
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

// +------------+--------------+------+-----+---------+----------------+
// | Field      | Type         | Null | Key | Default | Extra          |
// +------------+--------------+------+-----+---------+----------------+
// | id         | int          | NO   | PRI | NULL    | auto_increment |
// | user_id    | int          | YES  | MUL | NULL    |                |
// | title      | varchar(255) | YES  | MUL | NULL    |                |
// | name       | varchar(255) | YES  | MUL | NULL    |                |
// | surname    | varchar(255) | YES  | MUL | NULL    |                |
// | address    | varchar(255) | YES  |     | NULL    |                |
// | zip_code   | varchar(255) | YES  |     | NULL    |                |
// | town       | varchar(255) | YES  |     | NULL    |                |
// | province   | varchar(255) | YES  |     | NULL    |                |
// | country    | varchar(255) | YES  |     | NULL    |                |
// | tax_code   | varchar(255) | YES  |     | NULL    |                |
// | vat        | varchar(255) | YES  |     | NULL    |                |
// | created_at | datetime     | NO   |     | NULL    |                |
// | updated_at | datetime     | NO   |     | NULL    |                |
// | info       | text         | YES  |     | NULL    |                |
// +------------+--------------+------+-----+---------+----------------+
type Customer struct {
	ID       int     `json:"id"`
	UserID   int     `json:"-"`
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

type CustomerInfo struct {
	Customer *Customer
	Slips    []*Slip
}

// +--------------------+--------------+------+-----+---------+----------------+
// | Field              | Type         | Null | Key | Default | Extra          |
// +--------------------+--------------+------+-----+---------+----------------+
// | id                 | int          | NO   | PRI | NULL    | auto_increment |
// | customer_id        | int          | NO   | MUL | NULL    |                |
// | estimate_id        | int          | YES  | MUL | NULL    |                |
// | invoice_id         | int          | YES  | MUL | NULL    |                |
// | invoice_project_id | int          | YES  | MUL | NULL    |                |
// | name               | varchar(255) | NO   | MUL | NULL    |                |
// | rate               | decimal(8,2) | NO   |     | NULL    |                |
// | timed              | tinyint(1)   | YES  |     | 0       |                |
// | duration           | int          | YES  |     | NULL    |                |
// | created_at         | datetime     | NO   |     | NULL    |                |
// | updated_at         | datetime     | NO   |     | NULL    |                |
// +--------------------+--------------+------+-----+---------+----------------+
type Slip struct {
	ID               JsonNullInt64 `json:"id"`
	CustomerID       JsonNullInt64 `json:"customer_id"`
	EstimateID       JsonNullInt64 `json:"estimate_id"`
	InvoiceID        JsonNullInt64 `json:"invoice_id"`
	InvoiceProjectID JsonNullInt64 `json:"invoice_project_id"`
	Name             *string       `json:"name"`
	Rate             *float64      `json:"rate"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

// +---------------------+--------------+------+-----+---------+----------------+
// | Field               | Type         | Null | Key | Default | Extra          |
// +---------------------+--------------+------+-----+---------+----------------+
// | id                  | int          | NO   | PRI | NULL    | auto_increment |
// | consolidated_tax_id | int          | YES  | MUL | NULL    |                |
// | order               | int          | YES  |     | NULL    |                |
// | name                | varchar(255) | YES  | MUL | NULL    |                |
// | rate                | int          | YES  |     | NULL    |                |
// | fixed_rate          | tinyint(1)   | YES  |     | 0       |                |
// | compound            | tinyint(1)   | YES  |     | NULL    |                |
// | withholding         | tinyint(1)   | YES  |     | 0       |                |
// +---------------------+--------------+------+-----+---------+----------------+
type tax struct{}

// +---------+--------------+------+-----+---------+----------------+
// | Field   | Type         | Null | Key | Default | Extra          |
// +---------+--------------+------+-----+---------+----------------+
// | id      | int          | NO   | PRI | NULL    | auto_increment |
// | user_id | int          | YES  | MUL | NULL    |                |
// | name    | varchar(255) | YES  | MUL | NULL    |                |
// | notes   | text         | YES  |     | NULL    |                |
// +---------+--------------+------+-----+---------+----------------+
type consolidatedTax struct{}

// +---------------------+------------+------+-----+---------+----------------+
// | Field               | Type       | Null | Key | Default | Extra          |
// +---------------------+------------+------+-----+---------+----------------+
// | id                  | int        | NO   | PRI | NULL    | auto_increment |
// | customer_id         | int        | YES  | MUL | NULL    |                |
// | consolidated_tax_id | int        | YES  | MUL | NULL    |                |
// | date                | date       | NO   | MUL | NULL    |                |
// | number              | int        | NO   | MUL | NULL    |                |
// | payment_date        | date       | YES  |     | NULL    |                |
// | created_at          | datetime   | NO   |     | NULL    |                |
// | updated_at          | datetime   | NO   |     | NULL    |                |
// | invoice_project_id  | int        | YES  | UNI | NULL    |                |
// | downloaded          | tinyint(1) | YES  |     | 0       |                |
// +---------------------+------------+------+-----+---------+----------------+
type invoice struct{}

// +---------------------+------------+------+-----+---------+----------------+
// | Field               | Type       | Null | Key | Default | Extra          |
// +---------------------+------------+------+-----+---------+----------------+
// | id                  | int        | NO   | PRI | NULL    | auto_increment |
// | customer_id         | int        | YES  | MUL | NULL    |                |
// | consolidated_tax_id | int        | YES  | MUL | NULL    |                |
// | date                | date       | NO   | MUL | NULL    |                |
// | number              | int        | NO   | MUL | NULL    |                |
// | invoiced            | tinyint(1) | YES  | MUL | 0       |                |
// | created_at          | datetime   | NO   |     | NULL    |                |
// | updated_at          | datetime   | NO   |     | NULL    |                |
// +---------------------+------------+------+-----+---------+----------------+
type estimate struct{}

// +-----------------+--------------+------+-----+---------+----------------+
// | Field           | Type         | Null | Key | Default | Extra          |
// +-----------------+--------------+------+-----+---------+----------------+
// | id              | int          | NO   | PRI | NULL    | auto_increment |
// | customer_id     | int          | NO   | MUL | NULL    |                |
// | schedule        | varchar(255) | NO   |     | NULL    |                |
// | last_occurrence | datetime     | YES  |     | NULL    |                |
// | next_occurrence | datetime     | NO   |     | NULL    |                |
// | name            | varchar(255) | NO   | MUL | NULL    |                |
// | rate            | decimal(8,2) | NO   |     | NULL    |                |
// +-----------------+--------------+------+-----+---------+----------------+
type recurringSlip struct{}

// +---------+--------------+------+-----+---------+----------------+
// | Field   | Type         | Null | Key | Default | Extra          |
// +---------+--------------+------+-----+---------+----------------+
// | id      | int          | NO   | PRI | NULL    | auto_increment |
// | user_id | int          | YES  | MUL | NULL    |                |
// | name    | varchar(255) | YES  | MUL | NULL    |                |
// | value   | varchar(255) | YES  |     | NULL    |                |
// | integer | tinyint(1)   | YES  |     | 0       |                |
// +---------+--------------+------+-----+---------+----------------+
type option struct{}

// +---------------------+------------+------+-----+---------+----------------+
// | Field               | Type       | Null | Key | Default | Extra          |
// +---------------------+------------+------+-----+---------+----------------+
// | id                  | int        | NO   | PRI | NULL    | auto_increment |
// | customer_id         | int        | YES  | MUL | NULL    |                |
// | consolidated_tax_id | int        | YES  | MUL | NULL    |                |
// | date                | date       | NO   | MUL | NULL    |                |
// | number              | int        | NO   | MUL | NULL    |                |
// | invoiced            | tinyint(1) | YES  | MUL | 0       |                |
// | created_at          | datetime   | NO   |     | NULL    |                |
// | updated_at          | datetime   | NO   |     | NULL    |                |
// | downloaded          | tinyint(1) | YES  |     | 0       |                |
// +---------------------+------------+------+-----+---------+----------------+
type invoiceProject struct{}

// +-------------------------+--------------+------+-----+---------+----------------+
// | Field                   | Type         | Null | Key | Default | Extra          |
// +-------------------------+--------------+------+-----+---------+----------------+
// | id                      | int          | NO   | PRI | NULL    | auto_increment |
// | user_id                 | int          | YES  |     | NULL    |                |
// | customer_id             | int          | YES  |     | NULL    |                |
// | year                    | int          | YES  |     | NULL    |                |
// | received_at             | date         | YES  |     | NULL    |                |
// | attachment_file_name    | varchar(255) | YES  |     | NULL    |                |
// | attachment_content_type | varchar(255) | YES  |     | NULL    |                |
// | attachment_file_size    | int          | YES  |     | NULL    |                |
// | attachment_updated_at   | datetime     | YES  |     | NULL    |                |
// | rate                    | decimal(8,2) | YES  |     | NULL    |                |
// | created_at              | datetime     | NO   |     | NULL    |                |
// | updated_at              | datetime     | NO   |     | NULL    |                |
// +-------------------------+--------------+------+-----+---------+----------------+
type certification struct{}

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}
