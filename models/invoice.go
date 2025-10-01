package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

// InvoiceData contains complete invoice information including provider and client details
type InvoiceData struct {
	InvoiceID       int
	DateCreated     time.Time
	Paid            bool
	ProviderName    string
	ProviderAddress *string
	ProviderEmail   *string
	ProviderPhone   *string
	ClientName      string
	ClientAddress   *string
	ClientEmail     *string
	Items           []InvoiceItem
}

func CreateInvoice(providerID, clientID string, paid bool) (int, error) {
	result, err := db.Exec(
		"INSERT INTO invoice (provider_id, client_id, paid) VALUES (?, ?, ?)",
		providerID,
		clientID,
		paid,
	)
	if err != nil {
		return 0, err
	}

	invoiceID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(invoiceID), nil
}

func UpdateInvoice(invoiceID int, providerID, clientID string, paid bool) error {
	result, err := db.Exec(
		"UPDATE invoice SET provider_id = ?, client_id = ?, paid = ? WHERE id = ?",
		providerID,
		clientID,
		paid,
		invoiceID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func ListInvoices() ([]struct {
	ID           int
	ProviderName string
	ClientName   string
	DateCreated  time.Time
	Paid         bool
}, error,
) {
	rows, err := db.Query(`
		SELECT
			i.id,
			p.name as provider_name,
			c.name as client_name,
			i.date_created,
			i.paid
		FROM invoice i
		LEFT JOIN provider p ON i.provider_id = p.id
		LEFT JOIN client c ON i.client_id = c.id
		ORDER BY i.date_created DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []struct {
		ID           int
		ProviderName string
		ClientName   string
		DateCreated  time.Time
		Paid         bool
	}

	for rows.Next() {
		var inv struct {
			ID           int
			ProviderName string
			ClientName   string
			DateCreated  time.Time
			Paid         bool
		}
		if err := rows.Scan(&inv.ID, &inv.ProviderName, &inv.ClientName, &inv.DateCreated, &inv.Paid); err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}

	return invoices, rows.Err()
}

func GetInvoiceData(invoiceID int) (*InvoiceData, error) {
	var data InvoiceData

	err := db.QueryRow(`
		SELECT
			i.id,
			i.date_created,
			i.paid,
			p.name, p.address, p.email, p.phone,
			c.name, c.address, c.email
		FROM invoice i
		LEFT JOIN provider p ON i.provider_id = p.id
		LEFT JOIN client c ON i.client_id = c.id
		WHERE i.id = ?
	`, invoiceID).Scan(
		&data.InvoiceID,
		&data.DateCreated,
		&data.Paid,
		&data.ProviderName,
		&data.ProviderAddress,
		&data.ProviderEmail,
		&data.ProviderPhone,
		&data.ClientName,
		&data.ClientAddress,
		&data.ClientEmail,
	)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT id, invoice_id, item_name, amount, cost_per_unit
		FROM invoice_item
		WHERE invoice_id = ?
	`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item InvoiceItem
		if err := rows.Scan(&item.ID, &item.InvoiceID, &item.ItemName, &item.Amount, &item.CostPerUnit); err != nil {
			return nil, err
		}
		data.Items = append(data.Items, item)
	}

	return &data, rows.Err()
}

func AddInvoiceItem(invoiceID int, itemName string, amount, costPerUnit float64) (int, error) {
	if strings.TrimSpace(itemName) == "" {
		return 0, errors.New("item name is required")
	}
	if amount <= 0 {
		return 0, errors.New("amount must be positive")
	}
	if costPerUnit <= 0 {
		return 0, errors.New("cost per unit must be positive")
	}

	result, err := db.Exec(
		"INSERT INTO invoice_item (invoice_id, item_name, amount, cost_per_unit) VALUES (?, ?, ?, ?)",
		invoiceID,
		strings.TrimSpace(itemName),
		amount,
		costPerUnit,
	)
	if err != nil {
		return 0, err
	}

	itemID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(itemID), nil
}
