package storage

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/GVPproj/termsheet/models"
	"github.com/google/uuid"
)

func CreateClient(name string, address, email *string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", errors.New("client name is required")
	}

	clientID := uuid.New().String()

	_, err := db.Exec(
		"INSERT INTO client (id, name, address, email) VALUES (?, ?, ?, ?)",
		clientID,
		strings.TrimSpace(name),
		address,
		email,
	)
	if err != nil {
		return "", err
	}

	return clientID, nil
}

func ListClients() ([]models.Client, error) {
	rows, err := db.Query("SELECT id, name, address, email FROM client")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		if err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.Email); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}

	return clients, rows.Err()
}

func UpdateClient(clientID, name string, address, email *string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("client name is required")
	}

	result, err := db.Exec(
		"UPDATE client SET name = ?, address = ?, email = ? WHERE id = ?",
		strings.TrimSpace(name),
		address,
		email,
		clientID,
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

func DeleteClient(clientID string) error {
	result, err := db.Exec("DELETE FROM client WHERE id = ?", clientID)
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
