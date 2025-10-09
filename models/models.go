// Package models defines the core data structures for invoicing entities.
package models

import "time"

// Entity represents a generic contact entity (client or provider)
type Entity struct {
	ID   string
	Name string
	// we use pointers for optional fields so that they can be NULL
	// If the database column is NULL, the *string will be nil.
	// If the database column has a string value, the *string will point to that string.
	Address *string
	Email   *string
	Phone   *string
}

type Invoice struct {
	ID          int
	ProviderID  string
	ClientID    string
	Paid        bool
	DateCreated time.Time
}

type InvoiceItem struct {
	ID          int
	InvoiceID   int
	ItemName    string
	Amount      float64
	CostPerUnit float64
}

// InvoiceData contains complete invoice information including provider and client details
type InvoiceData struct {
	InvoiceID   int
	DateCreated time.Time
	Paid        bool
	Provider    Entity
	Client      Entity
	Items       []InvoiceItem
}
