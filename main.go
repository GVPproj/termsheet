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
	// this Pointer allows form state to be modified (selection, focus, etc.)
	// this means form as value doesn't get recreated each time
	form      *huh.Form
	selection string
}

// createMenuForm is a method on the model struct
// The (m *model) part is called a receiver - it makes createMenuForm() a method on the model struct
func (m *model) createMenuForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				// .Title() - Sets the prompt text shown to the user
				Title("Please select a view").
				// .Options() - Defines the selectable choices (label + value pairs)
				Options(
					huh.NewOption("Providers", "Providers"),
					huh.NewOption("Clients", "Clients"),
					huh.NewOption("Invoices", "Invoices"),
				).
				// .Value(&m.selection) - Binds the selected value to the m.selection field on the model struct
				// using a pointer to that variable.
				// Why a Pointer is Needed for Modification: If Value() were given just m.selection (the value itself, not its address),
				// it would be working with a copy of the string. Any changes it made to that copy
				// would not affect the original m.selection in your model.
				// By giving it a pointer (&m.selection), huh.Form knows exactly where in memory to store
				// the user's selection, directly updating m.selection in your model struct.
				Value(&m.selection),
		),
	)
}

func initialModel() *model {
	m := &model{
		currentView: types.MenuView,
		choices:     []string{"Providers", "Clients", "Invoices"},
	}

	m.form = m.createMenuForm()
	return m
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("termsheet"), m.form.Init())
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *model) View() string {
	switch m.currentView {
	case types.MenuView:
		return views.RenderMenu(m.form)
	case types.ProvidersView:
		return views.RenderProviders()
	case types.ClientsView:
		return views.RenderClients()
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
