// Package models defines the core data structures for invoicing entities.
package models

import "time"

type Provider struct {
	ID   string
	Name string
	// we use pointers for optional fields so that they can be NULL
	// If the database column is NULL, the *string will be nil.
	// If the database column has a string value, the *string will point to that string.
	Address *string
	Email   *string
	Phone   *string
}

type Client struct {
	ID      string
	Name    string
	Address *string
	Email   *string
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
