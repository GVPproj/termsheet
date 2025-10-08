// Package client
package client

import (
	"log"

	"github.com/GVPproj/termsheet/models"
	"github.com/GVPproj/termsheet/storage"
	"github.com/GVPproj/termsheet/tui/forms"
	"github.com/GVPproj/termsheet/tui/views"
	"github.com/GVPproj/termsheet/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Controller manages client-related state and behavior
type Controller struct {
	// Form state
	form      *huh.Form
	selection string

	// Client form fields
	name    string
	address string
	email   string
	phone   string

	// Edit state
	selectedID string
}

// NewController creates a new client controller
func NewController() *Controller {
	return &Controller{}
}

// InitListView initializes the client list view
func (c *Controller) InitListView() (*huh.Form, error) {
	c.selection = ""
	clientListForm, err := views.CreateClientListForm(&c.selection)
	if err != nil {
		return nil, err
	}
	c.form = clientListForm
	return c.form, nil
}

// Update handles client-related messages and returns view transition if needed
func (c *Controller) Update(msg tea.Msg, currentView types.View) (*types.ViewTransition, tea.Cmd) {
	switch currentView {
	case types.ClientsListView:
		return c.handleListView(msg)
	case types.ClientCreateView, types.ClientEditView:
		return c.handleFormView(msg, currentView)
	}
	return nil, nil
}

// handleListView manages the client list view logic
func (c *Controller) handleListView(msg tea.Msg) (*types.ViewTransition, tea.Cmd) {
	// Handle delete key before passing to form
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "d" {
		if c.selection != "" && c.selection != "CREATE_NEW" {
			// Delete the client
			err := storage.DeleteClient(c.selection)
			if err != nil {
				log.Printf("Error deleting client: %v", err)
				return nil, nil
			}
			// Refresh the client list
			c.selection = ""
			clientListForm, err := views.CreateClientListForm(&c.selection)
			if err != nil {
				log.Printf("Error refreshing client list: %v", err)
				return nil, nil
			}
			c.form = clientListForm
			return nil, c.form.Init()
		}
	}

	// Update form
	form, cmd := c.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		c.form = f
	}

	// Check if form is completed
	if c.form.State == huh.StateCompleted {
		if c.selection == "CREATE_NEW" {
			// Navigate to create client view
			c.resetFormFields()
			c.form = forms.NewClientForm(&c.name, &c.address, &c.email, &c.phone)
			return &types.ViewTransition{
				NewView: types.ClientCreateView,
				Form:    c.form,
			}, c.form.Init()
		} else {
			// Navigate to edit client view
			c.selectedID = c.selection
			// Get client data
			clients, err := storage.ListClients()
			if err != nil {
				log.Printf("Error loading client: %v", err)
				return nil, nil
			}
			var selectedClient *models.Client
			for _, p := range clients {
				if p.ID == c.selectedID {
					selectedClient = &p
					break
				}
			}
			if selectedClient == nil {
				log.Printf("Client not found: %s", c.selectedID)
				return nil, nil
			}
			c.form = forms.NewClientFormWithData(*selectedClient, &c.name, &c.address, &c.email, &c.phone)
			return &types.ViewTransition{
				NewView: types.ClientEditView,
				Form:    c.form,
			}, c.form.Init()
		}
	}

	return nil, cmd
}

// handleFormView manages create and edit form views
func (c *Controller) handleFormView(msg tea.Msg, currentView types.View) (*types.ViewTransition, tea.Cmd) {
	// Update form
	form, cmd := c.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		c.form = f
	}

	// Check if form is completed
	if c.form.State == huh.StateCompleted {
		// Handle client creation/update
		// Convert empty strings to nil pointers, otherwise return pointer to value
		var addressPtr, emailPtr, phonePtr *string
		if c.address != "" {
			addressPtr = &c.address
		}
		if c.email != "" {
			emailPtr = &c.email
		}
		if c.phone != "" {
			phonePtr = &c.phone
		}

		if currentView == types.ClientCreateView {
			_, err := storage.CreateClient(c.name, addressPtr, emailPtr, phonePtr)
			if err != nil {
				log.Printf("Error creating client: %v", err)
				return nil, nil
			}
		} else {
			err := storage.UpdateClient(c.selectedID, c.name, addressPtr, emailPtr, phonePtr)
			if err != nil {
				log.Printf("Error updating client: %v", err)
				return nil, nil
			}
		}

		// Return to client list
		c.selection = ""
		clientForm, err := views.CreateClientListForm(&c.selection)
		if err != nil {
			log.Printf("Error creating client form: %v", err)
			return nil, nil
		}
		c.form = clientForm
		return &types.ViewTransition{
			NewView: types.ClientsListView,
			Form:    c.form,
		}, c.form.Init()
	}

	return nil, cmd
}

// resetFormFields clears all form field values
func (c *Controller) resetFormFields() {
	c.name = ""
	c.address = ""
	c.email = ""
	c.phone = ""
}

// GetForm returns the current form
func (c *Controller) GetForm() *huh.Form {
	return c.form
}
