package fattingo

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

func (db *database) allCustomers() ([]*customer, error) {
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
	ORDER BY title ASC;`)
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
