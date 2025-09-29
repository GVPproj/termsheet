// Package types defines the view constants used for navigation in the invogo terminal application.
package types

type View int

const (
	MenuView View = iota
	ProvidersView
	ClientsView
	InvoicesView
)
