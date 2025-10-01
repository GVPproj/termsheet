package main

import (
	"fmt"
	"os"

	"github.com/GVPproj/termsheet/types"
	"github.com/GVPproj/termsheet/views"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	currentView types.View
	cursor      int
	choices     []string
}

func initialModel() model {
	return model{
		currentView: types.MenuView,
		choices:     []string{"Providers", "Clients", "Invoices"},
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("termsheet")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.currentView = types.MenuView
		case "up", "k":
			if m.currentView == types.MenuView && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.currentView == types.MenuView && m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.currentView == types.MenuView {
				m.currentView = types.View(m.cursor + 1)
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.currentView {
	case types.MenuView:
		return views.RenderMenu(m.cursor, m.choices)
	case types.ProvidersView:
		return views.RenderProviders()
	default:
		return "View not implemented yet\n\nPress ESC to return to menu"
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
