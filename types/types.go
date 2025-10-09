// Package types defines the view constants used for navigation in the termsheet application.
package types

import "github.com/charmbracelet/huh"

type View int

const (
	MenuView View = iota
	ProvidersListView
	ProviderCreateView
	ProviderEditView
	ProviderDeleteConfirmView
	ClientsListView
	ClientCreateView
	ClientEditView
	ClientDeleteConfirmView
	InvoicesListView
	InvoiceCreateView
	InvoiceEditView
)

// ViewTransition represents a request to change views
type ViewTransition struct {
	NewView View
	Form    *huh.Form
}
