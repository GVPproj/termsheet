// Package types defines the view constants used for navigation in the termsheet application.
package types

type View int

const (
	MenuView View = iota
	ProvidersListView
	ProviderCreateView
	ProviderEditView
	ClientsListView
	InvoicesListView
)
