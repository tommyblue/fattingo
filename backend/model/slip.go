package model

func (db *Database) Slips(customerId int) ([]*Slip, error) {
	rows, err := db.Query(`
	SELECT
		id,
		customer_id,
		estimate_id,
		invoice_id,
		invoice_project_id,
		name,
		rate,
		created_at,
		updated_at
	FROM slips
	WHERE customer_id = ?
	ORDER BY id DESC;`, customerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slips := make([]*Slip, 0)
	for rows.Next() {
		s := &Slip{}
		err = rows.Scan(
			&s.ID,
			&s.CustomerID,
			&s.EstimateID,
			&s.InvoiceID,
			&s.InvoiceProjectID,
			&s.Name,
			&s.Rate,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		slips = append(slips, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return slips, nil
}
