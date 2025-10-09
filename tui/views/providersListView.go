package views

import (
	"fmt"
	"strings"

	"github.com/GVPproj/termsheet/storage"
	"github.com/charmbracelet/huh"
)

// CreateProviderListForm creates a form for selecting or creating providers
// If preserveIndex >= 0, it attempts to select the item at that index
func CreateProviderListForm(selection *string) (*huh.Form, error) {
	return CreateProviderListFormWithError(selection, "")
}

// CreateProviderListFormWithError creates a form with an optional error message
func CreateProviderListFormWithError(selection *string, errorMsg string) (*huh.Form, error) {
	providers, err := storage.ListProviders()
	if err != nil {
		return nil, err
	}

	// Build options for the select form
	options := make([]huh.Option[string], 0)

	// Add existing providers
	for _, p := range providers {
		label := p.Name
		if p.Email != nil {
			label += fmt.Sprintf(" (%s)", *p.Email)
		}
		options = append(options, huh.NewOption(label, p.ID))
	}

	// Add "Create New Provider" option
	options = append(options, huh.NewOption("+ Create New Provider", "CREATE_NEW"))

	title := "Select a provider or create a new one"
	if errorMsg != "" {
		title = "⚠️  " + errorMsg + "\n\nSelect a provider or create a new one"
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

// RenderProviders renders the provider list view with the given form
func RenderProviders(form *huh.Form) string {
	var b strings.Builder

	// Render title
	b.WriteString(titleStyle.Render("Providers"))
	b.WriteString("\n\n")

	// Render the form
	b.WriteString(form.View())

	// Render help text
	b.WriteString(helpStyle.Render("\n\nPress 'd' to delete | ESC to return to menu"))

	// Wrap in container
	return containerStyle.Render(b.String())
}

// DeleteSelectedProvider deletes the provider with the given ID and returns the index to preserve
func DeleteSelectedProvider(providerID string) (int, error) {
	if providerID == "" || providerID == "CREATE_NEW" {
		return -1, fmt.Errorf("invalid provider selection")
	}

	// Get current providers to find the index before deletion
	providers, err := storage.ListProviders()
	if err != nil {
		return -1, err
	}

	// Find index of provider being deleted
	deletedIndex := -1
	for i, p := range providers {
		if p.ID == providerID {
			deletedIndex = i
			break
		}
	}

	// Delete the provider
	if err := storage.DeleteProvider(providerID); err != nil {
		return -1, err
	}

	// Return the index where cursor should be after deletion
	// (same index will select next item, or previous if last was deleted)
	return deletedIndex, nil
}
