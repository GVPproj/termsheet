// Package views collects the various app state views in the Termsheet app
package views

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Adaptive colors for light/dark terminals
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	// Title styling
	titleStyle = lipgloss.NewStyle().
			Bold(true).
		// gross green test here
		Foreground(lipgloss.Color("#00FF00")).
		MarginBottom(1)

	// Help text styling
	helpStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1)

	// Container styling
	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(60)
)

func RenderMenu(form *huh.Form) string {
	var b strings.Builder

	// Render title
	b.WriteString(titleStyle.Render("Termsheet"))
	b.WriteString("\n\n")

	// Render the form
	b.WriteString(form.View())

	// Render help text
	b.WriteString(helpStyle.Render("\nPress q to quit"))

	// Wrap in container
	return containerStyle.Render(b.String())
}
