package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GVPproj/termsheet/storage"
	"github.com/GVPproj/termsheet/tui/components/client"
	"github.com/GVPproj/termsheet/tui/components/provider"
	"github.com/GVPproj/termsheet/tui/views"
	"github.com/GVPproj/termsheet/types"
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

	// Provider component
	providerComponent *provider.Controller
	// Client component
	clientComponent *client.Controller
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
	).WithTheme(views.GetMenuTheme())
}

func initialModel() *model {
	m := &model{
		currentView:       types.MenuView,
		choices:           []string{"Providers", "Clients", "Invoices"},
		providerComponent: provider.NewController(),
		clientComponent:   client.NewController(),
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
			if m.currentView == types.MenuView {
				return m, tea.Quit
			}
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
				m.currentView = types.ProvidersListView
				// Initialize provider list form
				providerForm, err := m.providerComponent.InitListView()
				if err != nil {
					log.Printf("Error creating provider form: %v", err)
					return m, nil
				}
				m.form = providerForm
				return m, m.form.Init()
			case "Clients":
				m.currentView = types.ClientsListView
				clientForm, err := m.clientComponent.InitListView()
				if err != nil {
					log.Printf("Error creating client form: %v", err)
					return m, nil
				}
				m.form = clientForm
				return m, m.form.Init()
			case "Invoices":
				m.currentView = types.InvoicesListView
			}
		}

		return m, cmd
	}

	// Delegate to provider component for provider-related views
	if m.currentView == types.ProvidersListView ||
		m.currentView == types.ProviderCreateView ||
		m.currentView == types.ProviderEditView {
		transition, cmd := m.providerComponent.Update(msg, m.currentView)
		if transition != nil {
			m.currentView = transition.NewView
			m.form = transition.Form
			return m, cmd
		}
		// Update form reference from component
		m.form = m.providerComponent.GetForm()
		return m, cmd
	}

	// Delegate to client component for client-related views
	if m.currentView == types.ClientsListView ||
		m.currentView == types.ClientCreateView ||
		m.currentView == types.ClientEditView {
		transition, cmd := m.clientComponent.Update(msg, m.currentView)
		if transition != nil {
			m.currentView = transition.NewView
			m.form = transition.Form
			return m, cmd
		}
		// Update form reference from component
		m.form = m.clientComponent.GetForm()
		return m, cmd
	}

	return m, nil
}

func (m *model) View() string {
	switch m.currentView {
	case types.MenuView:
		return views.RenderMenu(m.form)
	case types.ProvidersListView:
		return views.RenderProviders(m.form)
	case types.ProviderCreateView, types.ProviderEditView:
		return views.RenderProviders(m.form)
	case types.ClientsListView:
		return views.RenderClients(m.form)
	case types.ClientCreateView, types.ClientEditView:
		return views.RenderClients(m.form)
	case types.InvoicesListView:
		return views.RenderInvoices()
	default:
		return "View not implemented yet\n\nPress ESC to return to menu"
	}
}

func main() {
	// Initialize database
	if err := storage.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer storage.CloseDB()

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
