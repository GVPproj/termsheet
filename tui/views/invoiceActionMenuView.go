package views

import (
	"strings"

	"github.com/charmbracelet/huh"
)

// InvoiceActionOption represents the action to take on an invoice
type InvoiceActionOption string

const (
	ActionView   InvoiceActionOption = "view"
	ActionEdit   InvoiceActionOption = "edit"
	ActionPDF    InvoiceActionOption = "pdf"
	ActionCancel InvoiceActionOption = "cancel"
)

// CreateInvoiceActionForm creates a form for selecting an action on an invoice
func CreateInvoiceActionForm(selection *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
				Options(
					huh.NewOption("View Invoice", string(ActionView)),
					huh.NewOption("Edit Invoice", string(ActionEdit)),
					huh.NewOption("Output PDF", string(ActionPDF)),
				).
				Value(selection),
		),
	).WithTheme(GetMenuTheme())
}

// RenderInvoiceActionMenu renders the invoice action menu view
func RenderInvoiceActionMenu(form *huh.Form) string {
	var b strings.Builder

	// Render title
	b.WriteString(titleStyle.Render("Invoice Actions"))
	b.WriteString("\n\n")

	// Render the form
	b.WriteString(form.View())

	// Render help text
	b.WriteString(helpStyle.Render("\n\nESC to return to invoice list"))

	// Wrap in container
	return containerStyle.Render(b.String())
}
