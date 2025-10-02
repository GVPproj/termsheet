package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GVPproj/termsheet/models"
	"github.com/GVPproj/termsheet/types"
	"github.com/GVPproj/termsheet/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type model struct {
	currentView types.View
	cursor      int
	choices     []string
	form        *huh.Form
	selection   string
}

func (m *model) createMenuForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Please select a view").
				Options(
					huh.NewOption("Providers", "Providers"),
					huh.NewOption("Clients", "Clients"),
					huh.NewOption("Invoices", "Invoices"),
				).
				Value(&m.selection),
		),
	)
}

func initialModel() model {
	m := model{
		currentView: types.MenuView,
		choices:     []string{"Providers", "Clients", "Invoices"},
	}

	m.form = m.createMenuForm()
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("termsheet"), m.form.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.currentView != types.MenuView {
				// Reset form when returning to menu
				m.currentView = types.MenuView
				m.selection = ""
				m.form = m.createMenuForm()
				return m, m.form.Init()
			}
		}
	}

	// When in menu view, let the form handle messages
	if m.currentView == types.MenuView {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		// Check if form is completed
		if m.form.State == huh.StateCompleted {
			// Map selection to view
			switch m.selection {
			case "Providers":
				m.currentView = types.ProvidersView
			case "Clients":
				m.currentView = types.ClientsView
			case "Invoices":
				m.currentView = types.InvoicesView
			}
		}

		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.currentView {
	case types.MenuView:
		return views.RenderMenu(m.form)
	case types.ProvidersView:
		return views.RenderProviders()
	default:
		return "View not implemented yet\n\nPress ESC to return to menu"
	}
}

func main() {
	// Initialize database
	if err := models.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer models.CloseDB()

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
