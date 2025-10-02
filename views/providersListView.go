package views

import (
	"fmt"
	"strings"

	"github.com/GVPproj/termsheet/models"
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
func CreateProviderListForm(selection *string) (*huh.Form, error) {
	providers, err := models.ListProviders()
	if err != nil {
		return nil, err
	}

	// Build options for the select form
	options := make([]huh.Option[string], 0)

	// Add existing providers
	for _, p := range providers {
		label := fmt.Sprintf("%s", p.Name)
		if p.Email != nil {
			label += fmt.Sprintf(" (%s)", *p.Email)
		}
		options = append(options, huh.NewOption(label, p.ID))
	}

	// Add "Create New Provider" option
	options = append(options, huh.NewOption("+ Create New Provider", "CREATE_NEW"))

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
	b.WriteString(helpStyle.Render("\n\nPress ESC to return to menu"))

	// Wrap in container
	return providerContainerStyle.Render(b.String())
}
