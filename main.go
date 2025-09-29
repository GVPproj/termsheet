package main

import (
	"fmt"
	"os"

	"github.com/GVPproj/termsheet/types"
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
		return m.menuView()
	case types.ProvidersView:
		return m.providersView()
	default:
		return "View not implemented yet\n\nPress ESC to return to menu"
	}
}

func (m model) menuView() string {
	s := "Please select a view\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit."
	return s
}

func (m model) providersView() string {
	s := "PROVIDERS VIEW\n\n"
	s += "• Provider 1\n"
	s += "• Provider 2\n"
	s += "• Provider 3\n\n"
	s += "Press ESC to return to menu"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
