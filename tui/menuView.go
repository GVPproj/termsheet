// Package tui collects the various app state views in the Termsheet app
package tui

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// GetMenuTheme returns a customized theme for the menu form
func GetMenuTheme() *huh.Theme {
	t := huh.ThemeCharm()
	t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color("#61AFEF"))
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(lipgloss.Color("#D33682"))
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(lipgloss.Color("#F255A1"))
	return t
}

var (
	// Adaptive colors for light/dark terminals
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	// Title styling
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#61AFEF")).
			MarginBottom(1)

	// Help text styling
	helpStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1)

	// Container styling
	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#61AFEF")).
			Padding(1, 2).
			Width(60)
)

func RenderMenu(form *huh.Form) string {
	var b strings.Builder

	// Render title
	b.WriteString(titleStyle.Render("Termsheet"))
	b.WriteString("\n\n")

	// Render the form
	// TODO: pass styles into here
	b.WriteString(form.View())

	// Render help text
	b.WriteString(helpStyle.Render("\nPress q to quit"))

	// Wrap in container
	return containerStyle.Render(b.String())
}
