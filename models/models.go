// Package models defines the core data structures for invoicing entities.
package models

import "time"

// Provider represents a service provider (sender of invoices)
type Provider struct {
	ID      string
	Name    string
	Address *string
	Email   *string
	Phone   *string
}

// Client represents a client receiving invoices
type Client struct {
	ID      string
	Name    string
	Address *string
	Email   *string
}

// Invoice represents an invoice with provider, client, and payment status
type Invoice struct {
	ID          int
	ProviderID  string
	ClientID    string
	Paid        bool
	DateCreated time.Time
}

// InvoiceItem represents a line item on an invoice
type InvoiceItem struct {
	ID          int
	InvoiceID   int
	ItemName    string
	Amount      float64
	CostPerUnit float64
}
