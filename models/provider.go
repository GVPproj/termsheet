package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

func CreateProvider(name string, address, email, phone *string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", errors.New("provider name is required")
	}

	providerID := uuid.New().String()

	_, err := db.Exec(
		"INSERT INTO provider (id, name, address, email, phone) VALUES (?, ?, ?, ?, ?)",
		providerID,
		strings.TrimSpace(name),
		address,
		email,
		phone,
	)
	if err != nil {
		return "", err
	}

	return providerID, nil
}

func ListProviders() ([]Provider, error) {
	rows, err := db.Query("SELECT id, name, address, email, phone FROM provider")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []Provider
	for rows.Next() {
		var p Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.Address, &p.Email, &p.Phone); err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}

	return providers, rows.Err()
}

func UpdateProvider(providerID, name string, address, email, phone *string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("provider name is required")
	}

	result, err := db.Exec(
		"UPDATE provider SET name = ?, address = ?, email = ?, phone = ? WHERE id = ?",
		strings.TrimSpace(name),
		address,
		email,
		phone,
		providerID,
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

func DeleteProvider(providerID string) error {
	result, err := db.Exec("DELETE FROM provider WHERE id = ?", providerID)
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
