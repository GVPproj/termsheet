package provider

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

// Controller manages provider-related state and behavior
type Controller struct {
	// Form state
	form      *huh.Form
	selection string

	// Provider form fields
	name    string
	address string
	email   string
	phone   string

	// Edit state
	selectedID string
}

// NewController creates a new provider controller
func NewController() *Controller {
	return &Controller{}
}

// InitListView initializes the provider list view
func (c *Controller) InitListView() (*huh.Form, error) {
	c.selection = ""
	providerForm, err := views.CreateProviderListForm(&c.selection, -1)
	if err != nil {
		return nil, err
	}
	c.form = providerForm
	return c.form, nil
}

// Update handles provider-related messages and returns view transition if needed
func (c *Controller) Update(msg tea.Msg, currentView types.View) (*types.ViewTransition, tea.Cmd) {
	switch currentView {
	case types.ProvidersListView:
		return c.handleListView(msg)
	case types.ProviderCreateView, types.ProviderEditView:
		return c.handleFormView(msg, currentView)
	}
	return nil, nil
}

// handleListView manages the provider list view logic
func (c *Controller) handleListView(msg tea.Msg) (*types.ViewTransition, tea.Cmd) {
	// Handle delete key before passing to form
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "d" {
		if c.selection != "" && c.selection != "CREATE_NEW" {
			preserveIndex, err := views.DeleteSelectedProvider(c.selection)
			if err != nil {
				log.Printf("Error deleting provider: %v", err)
			} else {
				// Refresh the provider list, preserving cursor position
				c.selection = ""
				providerForm, err := views.CreateProviderListForm(&c.selection, preserveIndex)
				if err != nil {
					log.Printf("Error refreshing provider list: %v", err)
					return nil, nil
				}
				c.form = providerForm
				return nil, c.form.Init()
			}
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
			// Navigate to create provider view
			c.resetFormFields()
			c.form = forms.NewProviderForm(&c.name, &c.address, &c.email, &c.phone)
			return &types.ViewTransition{
				NewView: types.ProviderCreateView,
				Form:    c.form,
			}, c.form.Init()
		} else {
			// Navigate to edit provider view
			c.selectedID = c.selection
			// Get provider data
			providers, err := storage.ListProviders()
			if err != nil {
				log.Printf("Error loading provider: %v", err)
				return nil, nil
			}
			var selectedProvider *models.Provider
			for _, p := range providers {
				if p.ID == c.selectedID {
					selectedProvider = &p
					break
				}
			}
			if selectedProvider == nil {
				log.Printf("Provider not found: %s", c.selectedID)
				return nil, nil
			}
			c.form = forms.NewProviderFormWithData(*selectedProvider, &c.name, &c.address, &c.email, &c.phone)
			return &types.ViewTransition{
				NewView: types.ProviderEditView,
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
		// Handle provider creation/update
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

		if currentView == types.ProviderCreateView {
			_, err := storage.CreateProvider(c.name, addressPtr, emailPtr, phonePtr)
			if err != nil {
				log.Printf("Error creating provider: %v", err)
				return nil, nil
			}
		} else {
			err := storage.UpdateProvider(c.selectedID, c.name, addressPtr, emailPtr, phonePtr)
			if err != nil {
				log.Printf("Error updating provider: %v", err)
				return nil, nil
			}
		}

		// Return to provider list
		c.selection = ""
		providerForm, err := views.CreateProviderListForm(&c.selection, -1)
		if err != nil {
			log.Printf("Error creating provider form: %v", err)
			return nil, nil
		}
		c.form = providerForm
		return &types.ViewTransition{
			NewView: types.ProvidersListView,
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
