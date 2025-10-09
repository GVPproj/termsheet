package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/GVPproj/termsheet/models"
	"github.com/google/uuid"
)

// CreateEntity creates a new entity in the specified table
// In Go, we use string pointers (*string) for optional fields to distinguish between three states:
// 1. Field not provided - pointer is nil
// 2. Field provided but empty - pointer points to an empty string ""
// 3. Field has a value - pointer points to the actual string value
func CreateEntity(tableName, name string, address, email, phone *string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("%s name is required", tableName)
	}

	entityID := uuid.New().String()

	_, err := db.Exec(
		fmt.Sprintf("INSERT INTO %s (id, name, address, email, phone) VALUES (?, ?, ?, ?, ?)", tableName),
		entityID,
		strings.TrimSpace(name),
		address,
		email,
		phone,
	)
	if err != nil {
		return "", err
	}

	return entityID, nil
}

// ListEntities retrieves all entities from the specified table
func ListEntities(tableName string) ([]models.Entity, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT id, name, address, email, phone FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []models.Entity
	for rows.Next() {
		var e models.Entity
		if err := rows.Scan(&e.ID, &e.Name, &e.Address, &e.Email, &e.Phone); err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	return entities, rows.Err()
}

// UpdateEntity updates an entity in the specified table
func UpdateEntity(tableName, entityID, name string, address, email, phone *string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%s name is required", tableName)
	}

	result, err := db.Exec(
		fmt.Sprintf("UPDATE %s SET name = ?, address = ?, email = ?, phone = ? WHERE id = ?", tableName),
		strings.TrimSpace(name),
		address,
		email,
		phone,
		entityID,
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

// HasInvoiceReferences checks if an entity (client or provider) has any associated invoices
func HasInvoiceReferences(tableName, entityID string) (bool, error) {
	var columnName string
	if tableName == "client" {
		columnName = "client_id"
	} else if tableName == "provider" {
		columnName = "provider_id"
	} else {
		// For other entity types, no invoice check needed
		return false, nil
	}

	var count int
	err := db.QueryRow(
		fmt.Sprintf("SELECT COUNT(*) FROM invoice WHERE %s = ?", columnName),
		entityID,
	).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// DeleteEntity deletes an entity from the specified table
func DeleteEntity(tableName, entityID string) error {
	// Check for invoice references before deleting
	hasRefs, err := HasInvoiceReferences(tableName, entityID)
	if err != nil {
		return err
	}
	if hasRefs {
		return fmt.Errorf("cannot delete %s: it has associated invoices", tableName)
	}

	result, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName), entityID)
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

// ValidateEntityName validates that an entity name is not empty
func ValidateEntityName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}
	return nil
}
