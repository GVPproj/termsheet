package views

import (
	"fmt"
	"strings"

	"github.com/GVPproj/termsheet/storage"
	"github.com/charmbracelet/huh"
)

// CreateInvoiceListForm creates a form for selecting or creating invoices
func CreateInvoiceListForm(selection *string) (*huh.Form, error) {
	invoices, err := storage.ListInvoices()
	if err != nil {
		return nil, err
	}

	// Build options for the select form
	options := make([]huh.Option[string], 0)

	// Add existing invoices
	for _, inv := range invoices {
		paidStatus := "Unpaid"
		if inv.Paid {
			paidStatus = "Paid"
		}
		label := fmt.Sprintf("#%d - %s → %s (%s)", inv.ID, inv.ProviderName, inv.ClientName, paidStatus)
		// Store invoice ID as string for selection
		options = append(options, huh.NewOption(label, fmt.Sprintf("%d", inv.ID)))
	}

	// Add "Create New Invoice" option
	options = append(options, huh.NewOption("+ Create New Invoice", "CREATE_NEW"))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select an invoice or create a new one").
				Options(options...).
				Value(selection),
		),
	).WithTheme(GetMenuTheme())

	return form, nil
}

// RenderInvoices renders the invoice list view with the given form
func RenderInvoices(form *huh.Form) string {
	var b strings.Builder

	// Render title
	b.WriteString(titleStyle.Render("Invoices"))
	b.WriteString("\n\n")

	// Render the form
	b.WriteString(form.View())

	// Render help text
	b.WriteString(helpStyle.Render("\n\nPress 'd' to delete | ESC to return to menu"))

	// Wrap in container
	return containerStyle.Render(b.String())
}

// DeleteSelectedInvoice deletes the invoice with the given ID and returns the index to preserve
func DeleteSelectedInvoice(invoiceIDStr string) (int, error) {
	if invoiceIDStr == "" || invoiceIDStr == "CREATE_NEW" {
		return -1, fmt.Errorf("invalid invoice selection")
	}

	// Parse invoice ID
	var invoiceID int
	_, err := fmt.Sscanf(invoiceIDStr, "%d", &invoiceID)
	if err != nil {
		return -1, fmt.Errorf("invalid invoice ID: %v", err)
	}

	// Get current invoices to find the index before deletion
	invoices, err := storage.ListInvoices()
	if err != nil {
		return -1, err
	}

	// Find index of invoice being deleted
	deletedIndex := -1
	for i, inv := range invoices {
		if inv.ID == invoiceID {
			deletedIndex = i
			break
		}
	}

	// Delete the invoice
	if err := storage.DeleteInvoice(invoiceID); err != nil {
		return -1, err
	}

	// Return the index where cursor should be after deletion
	// (same index will select next item, or previous if last was deleted)
	return deletedIndex, nil
}
