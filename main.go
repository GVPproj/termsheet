package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GVPproj/termsheet/models"
	"github.com/GVPproj/termsheet/types"
	"github.com/GVPproj/termsheet/views"
	"github.com/GVPproj/termsheet/views/forms"
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

	// Provider form fields
	providerName    string
	providerAddress string
	providerEmail   string
	providerPhone   string
	selectedProviderID string
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
				m.currentView = types.ProvidersListView
				// Create provider list form
				providerForm, err := views.CreateProviderListForm(&m.selection)
				if err != nil {
					log.Printf("Error creating provider form: %v", err)
					return m, nil
				}
				m.form = providerForm
				return m, m.form.Init()
			case "Clients":
				m.currentView = types.ClientsListView
			case "Invoices":
				m.currentView = types.InvoicesListView
			}
		}

		return m, cmd
	}

	// When in provider list view, let the form handle messages
	if m.currentView == types.ProvidersListView {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		// Check if form is completed
		if m.form.State == huh.StateCompleted {
			// Handle selection
			if m.selection == "CREATE_NEW" {
				// Navigate to create provider view
				m.currentView = types.ProviderCreateView
				m.providerName = ""
				m.providerAddress = ""
				m.providerEmail = ""
				m.providerPhone = ""
				m.form = forms.NewProviderForm(&m.providerName, &m.providerAddress, &m.providerEmail, &m.providerPhone)
				return m, m.form.Init()
			} else {
				// Navigate to edit provider view
				m.selectedProviderID = m.selection
				// Get provider data
				providers, err := models.ListProviders()
				if err != nil {
					log.Printf("Error loading provider: %v", err)
					return m, nil
				}
				var selectedProvider *models.Provider
				for _, p := range providers {
					if p.ID == m.selectedProviderID {
						selectedProvider = &p
						break
					}
				}
				if selectedProvider == nil {
					log.Printf("Provider not found: %s", m.selectedProviderID)
					return m, nil
				}
				m.currentView = types.ProviderEditView
				m.form = forms.NewProviderFormWithData(*selectedProvider, &m.providerName, &m.providerAddress, &m.providerEmail, &m.providerPhone)
				return m, m.form.Init()
			}
		}

		return m, cmd
	}

	// When in provider create/edit view, let the form handle messages
	if m.currentView == types.ProviderCreateView || m.currentView == types.ProviderEditView {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		// Check if form is completed
		if m.form.State == huh.StateCompleted {
			// Handle provider creation/update
			addressPtr := forms.StringToStringPtr(m.providerAddress)
			emailPtr := forms.StringToStringPtr(m.providerEmail)
			phonePtr := forms.StringToStringPtr(m.providerPhone)

			if m.currentView == types.ProviderCreateView {
				_, err := models.CreateProvider(m.providerName, addressPtr, emailPtr, phonePtr)
				if err != nil {
					log.Printf("Error creating provider: %v", err)
					return m, nil
				}
			} else {
				err := models.UpdateProvider(m.selectedProviderID, m.providerName, addressPtr, emailPtr, phonePtr)
				if err != nil {
					log.Printf("Error updating provider: %v", err)
					return m, nil
				}
			}

			// Return to provider list
			m.currentView = types.ProvidersListView
			m.selection = ""
			providerForm, err := views.CreateProviderListForm(&m.selection)
			if err != nil {
				log.Printf("Error creating provider form: %v", err)
				return m, nil
			}
			m.form = providerForm
			return m, m.form.Init()
		}

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
	case types.ProviderCreateView:
		return views.RenderProviders(m.form) // Reuse the same renderer with the form
	case types.ProviderEditView:
		return views.RenderProviders(m.form) // Reuse the same renderer with the form
	case types.ClientsListView:
		return views.RenderClients()
	case types.InvoicesListView:
		return views.RenderInvoices()
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
