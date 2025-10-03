package views

import (
	"fmt"
	"strings"

	"github.com/GVPproj/termsheet/storage"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	providerTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#61AFEF")).
				MarginBottom(1)

	providerItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#98C379"))

	providerContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#61AFEF")).
				Padding(1, 2).
				Width(60)
)

// CreateProviderListForm creates a form for selecting or creating providers
// If preserveIndex >= 0, it attempts to select the item at that index
func CreateProviderListForm(selection *string, preserveIndex int) (*huh.Form, error) {
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

	// If preserveIndex is specified, set selection to that index
	// Adjust if out of bounds (deleted last item)
	if preserveIndex >= 0 {
		if preserveIndex >= len(options) {
			preserveIndex = len(options) - 1
		}
		if preserveIndex >= 0 && preserveIndex < len(options) {
			*selection = options[preserveIndex].Value
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a provider or create a new one").
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
	b.WriteString(providerTitleStyle.Render("Providers"))
	b.WriteString("\n\n")

	// Render the form
	b.WriteString(form.View())

	// Render help text
	b.WriteString(helpStyle.Render("\n\nPress 'd' to delete | ESC to return to menu"))

	// Wrap in container
	return providerContainerStyle.Render(b.String())
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
