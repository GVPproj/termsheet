package views

import (
	"fmt"
	"strings"

	"github.com/GVPproj/termsheet/storage"
	"github.com/charmbracelet/huh"
)

// CreateClientListForm creates a form for selecting or creating providers
func CreateClientListForm(selection *string) (*huh.Form, error) {
	return CreateClientListFormWithError(selection, "")
}

// CreateClientListFormWithError creates a form with an optional error message
func CreateClientListFormWithError(selection *string, errorMsg string) (*huh.Form, error) {
	clients, err := storage.ListClients()
	if err != nil {
		return nil, err
	}

	// Build options for the select form
	options := make([]huh.Option[string], 0)

	// Add existing clients
	for _, c := range clients {
		label := c.Name
		if c.Email != nil {
			label += fmt.Sprintf(" (%s)", *c.Email)
		}
		options = append(options, huh.NewOption(label, c.ID))
	}

	// Add "Create New Provider" option
	options = append(options, huh.NewOption("+ Create New Client", "CREATE_NEW"))

	title := "Select a client or create a new one"
	if errorMsg != "" {
		title = "⚠️  " + errorMsg + "\n\nSelect a client or create a new one"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(options...).
				Value(selection),
		),
	).WithTheme(GetMenuTheme())

	return form, nil
}

func RenderClients(form *huh.Form) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Clients"))
	b.WriteString("\n\n")

	// Render the form
	b.WriteString(form.View())

	// Render help text
	b.WriteString(helpStyle.Render("\n\nPress 'd' to delete | ESC to return to menu"))
	return containerStyle.Render(b.String())
}

func DeleteSelectedClient(clientID string) (int, error) {
	if clientID == "" || clientID == "CREATE_NEW" {
		return -1, fmt.Errorf("invalid client selection")
	}

	// Get current providers to find the index before deletion
	clients, err := storage.ListClients()
	if err != nil {
		return -1, err
	}

	// Find index of provider being deleted
	deletedIndex := -1
	for i, p := range clients {
		if p.ID == clientID {
			deletedIndex = i
			break
		}
	}

	// Delete the provider
	if err := storage.DeleteClient(clientID); err != nil {
		return -1, err
	}

	// Return the index where cursor should be after deletion
	// (same index will select next item, or previous if last was deleted)
	return deletedIndex, nil
}
