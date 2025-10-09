// Package invoice
package invoice

import (
	"fmt"
	"log"
	"strconv"

	"github.com/GVPproj/termsheet/models"
	"github.com/GVPproj/termsheet/storage"
	"github.com/GVPproj/termsheet/tui/forms"
	"github.com/GVPproj/termsheet/tui/views"
	"github.com/GVPproj/termsheet/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// InvoiceFormStep represents the current step in the invoice form flow
type InvoiceFormStep int

const (
	StepSelectProvider InvoiceFormStep = iota
	StepSelectClient
	StepAddItem
	StepAskForMore
	StepMarkPaid
)

// InvoiceItem represents a single item being added to an invoice
type InvoiceItem struct {
	Name        string
	Amount      float64
	CostPerUnit float64
}

// Controller manages invoice-related state and behavior
type Controller struct {
	// Form state
	form      *huh.Form
	selection string

	// Invoice form fields
	providerID      string
	clientID        string
	itemName        string
	itemAmount      string
	itemCostPerUnit string
	paid            bool
	addAnother      bool

	// Multi-step flow
	currentStep InvoiceFormStep
	items       []InvoiceItem

	// Edit state
	selectedID    string
	invoiceID     int
	existingItems []models.InvoiceItem
	isEditMode    bool
	currentItemIndex int
}

// NewController creates a new invoice controller
func NewController() *Controller {
	return &Controller{}
}

// InitListView initializes the invoice list view
func (c *Controller) InitListView() (*huh.Form, error) {
	c.selection = ""
	invoiceForm, err := views.CreateInvoiceListForm(&c.selection)
	if err != nil {
		return nil, err
	}
	c.form = invoiceForm
	return c.form, nil
}

// Update handles invoice-related messages and returns view transition if needed
func (c *Controller) Update(msg tea.Msg, currentView types.View) (*types.ViewTransition, tea.Cmd) {
	switch currentView {
	case types.InvoicesListView:
		return c.handleListView(msg)
	case types.InvoiceCreateView, types.InvoiceEditView:
		return c.handleFormView(msg, currentView)
	}
	return nil, nil
}

// handleListView manages the invoice list view logic
func (c *Controller) handleListView(msg tea.Msg) (*types.ViewTransition, tea.Cmd) {
	// Handle delete key before passing to form
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "d" {
		if c.selection != "" && c.selection != "CREATE_NEW" {
			// Parse invoice ID
			var invoiceID int
			_, err := fmt.Sscanf(c.selection, "%d", &invoiceID)
			if err != nil {
				log.Printf("Error parsing invoice ID: %v", err)
				return nil, nil
			}

			// Delete the invoice
			err = storage.DeleteInvoice(invoiceID)
			if err != nil {
				log.Printf("Error deleting invoice: %v", err)
				return nil, nil
			}

			// Refresh the invoice list
			c.selection = ""
			invoiceForm, err := views.CreateInvoiceListForm(&c.selection)
			if err != nil {
				log.Printf("Error refreshing invoice list: %v", err)
				return nil, nil
			}
			c.form = invoiceForm
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
			// Navigate to create invoice view - start with provider selection
			c.resetFormFields()
			c.currentStep = StepSelectProvider
			c.isEditMode = false
			invoiceForm, err := forms.NewProviderSelectForm(&c.providerID)
			if err != nil {
				log.Printf("Error creating invoice form: %v", err)
				return nil, nil
			}
			c.form = invoiceForm
			return &types.ViewTransition{
				NewView: types.InvoiceCreateView,
				Form:    c.form,
			}, c.form.Init()
		} else {
			// Navigate to edit invoice view
			c.selectedID = c.selection

			// Parse invoice ID
			var invoiceID int
			_, err := fmt.Sscanf(c.selectedID, "%d", &invoiceID)
			if err != nil {
				log.Printf("Error parsing invoice ID: %v", err)
				return nil, nil
			}
			c.invoiceID = invoiceID

			// Get invoice data
			invoiceData, err := storage.GetInvoiceData(invoiceID)
			if err != nil {
				log.Printf("Error loading invoice: %v", err)
				return nil, nil
			}

			// Create Invoice model from InvoiceData
			invoice := models.Invoice{
				ID:         invoiceData.InvoiceID,
				ProviderID: "", // We'll need to fetch this separately
				ClientID:   "", // We'll need to fetch this separately
				Paid:       invoiceData.Paid,
			}

			// Get provider and client IDs by name
			providers, err := storage.ListProviders()
			if err != nil {
				log.Printf("Error loading providers: %v", err)
				return nil, nil
			}
			for _, p := range providers {
				if p.Name == invoiceData.ProviderName {
					invoice.ProviderID = p.ID
					break
				}
			}

			clients, err := storage.ListClients()
			if err != nil {
				log.Printf("Error loading clients: %v", err)
				return nil, nil
			}
			for _, cl := range clients {
				if cl.Name == invoiceData.ClientName {
					invoice.ClientID = cl.ID
					break
				}
			}

			c.existingItems = invoiceData.Items
			c.isEditMode = true
			c.currentStep = StepSelectProvider
			c.currentItemIndex = 0

			// Pre-load existing items into items slice
			c.items = make([]InvoiceItem, 0, len(invoiceData.Items))
			for _, item := range invoiceData.Items {
				c.items = append(c.items, InvoiceItem{
					Name:        item.ItemName,
					Amount:      item.Amount,
					CostPerUnit: item.CostPerUnit,
				})
			}

			invoiceForm, err := forms.NewProviderSelectFormWithData(&c.providerID, invoice.ProviderID)
			if err != nil {
				log.Printf("Error creating invoice form: %v", err)
				return nil, nil
			}
			c.form = invoiceForm
			return &types.ViewTransition{
				NewView: types.InvoiceEditView,
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
		return c.handleStepComplete(currentView)
	}

	return nil, cmd
}

// handleStepComplete handles the completion of each form step
func (c *Controller) handleStepComplete(currentView types.View) (*types.ViewTransition, tea.Cmd) {
	switch c.currentStep {
	case StepSelectProvider:
		// Move to client selection
		c.currentStep = StepSelectClient
		var nextForm *huh.Form
		var err error

		if c.isEditMode && c.clientID != "" {
			nextForm, err = forms.NewClientSelectFormWithData(&c.clientID, c.clientID)
		} else {
			nextForm, err = forms.NewClientSelectForm(&c.clientID)
		}

		if err != nil {
			log.Printf("Error creating client form: %v", err)
			return nil, nil
		}
		c.form = nextForm
		return nil, c.form.Init()

	case StepSelectClient:
		// Move to item entry
		c.currentStep = StepAddItem
		var nextForm *huh.Form

		// In edit mode, pre-populate with first existing item
		if c.isEditMode && c.currentItemIndex < len(c.items) {
			item := c.items[c.currentItemIndex]
			c.itemName = item.Name
			c.itemAmount = fmt.Sprintf("%.2f", item.Amount)
			c.itemCostPerUnit = fmt.Sprintf("%.2f", item.CostPerUnit)
			nextForm = forms.NewInvoiceItemForm(&c.itemName, &c.itemAmount, &c.itemCostPerUnit)
		} else {
			// Clear item fields for new item
			c.itemName = ""
			c.itemAmount = ""
			c.itemCostPerUnit = ""
			nextForm = forms.NewInvoiceItemForm(&c.itemName, &c.itemAmount, &c.itemCostPerUnit)
		}

		c.form = nextForm
		return nil, c.form.Init()

	case StepAddItem:
		// Save the item
		amount, err := strconv.ParseFloat(c.itemAmount, 64)
		if err != nil {
			log.Printf("Error parsing amount: %v", err)
			return nil, nil
		}

		costPerUnit, err := strconv.ParseFloat(c.itemCostPerUnit, 64)
		if err != nil {
			log.Printf("Error parsing cost per unit: %v", err)
			return nil, nil
		}

		// Store or update item in items slice
		newItem := InvoiceItem{
			Name:        c.itemName,
			Amount:      amount,
			CostPerUnit: costPerUnit,
		}

		if c.isEditMode && c.currentItemIndex < len(c.items) {
			// Update existing item
			c.items[c.currentItemIndex] = newItem
			c.currentItemIndex++
		} else {
			// Add new item
			c.items = append(c.items, newItem)
		}

		// Check if there are more existing items to edit, or ask if user wants to add more
		if c.isEditMode && c.currentItemIndex < len(c.existingItems) {
			// Move to next existing item
			c.currentStep = StepAddItem
			item := c.items[c.currentItemIndex]
			c.itemName = item.Name
			c.itemAmount = fmt.Sprintf("%.2f", item.Amount)
			c.itemCostPerUnit = fmt.Sprintf("%.2f", item.CostPerUnit)
			nextForm := forms.NewInvoiceItemForm(&c.itemName, &c.itemAmount, &c.itemCostPerUnit)
			c.form = nextForm
			return nil, c.form.Init()
		}

		// Ask if user wants to add another item
		c.currentStep = StepAskForMore
		c.addAnother = false
		nextForm := forms.NewAddAnotherItemForm(&c.addAnother)
		c.form = nextForm
		return nil, c.form.Init()

	case StepAskForMore:
		if c.addAnother {
			// Add another item
			c.currentStep = StepAddItem
			c.itemName = ""
			c.itemAmount = ""
			c.itemCostPerUnit = ""
			nextForm := forms.NewInvoiceItemForm(&c.itemName, &c.itemAmount, &c.itemCostPerUnit)
			c.form = nextForm
			return nil, c.form.Init()
		}

		// Move to mark paid
		c.currentStep = StepMarkPaid
		var nextForm *huh.Form

		if c.isEditMode {
			nextForm = forms.NewMarkPaidFormWithData(&c.paid, c.paid)
		} else {
			c.paid = false
			nextForm = forms.NewMarkPaidForm(&c.paid)
		}

		c.form = nextForm
		return nil, c.form.Init()

	case StepMarkPaid:
		// Save the invoice
		return c.saveInvoice(currentView)
	}

	return nil, nil
}

// saveInvoice saves the invoice and all items to the database
func (c *Controller) saveInvoice(currentView types.View) (*types.ViewTransition, tea.Cmd) {
	if currentView == types.InvoiceCreateView {
		// Create invoice
		invoiceID, err := storage.CreateInvoice(c.providerID, c.clientID, c.paid)
		if err != nil {
			log.Printf("Error creating invoice: %v", err)
			return nil, nil
		}

		// Add all invoice items
		for _, item := range c.items {
			_, err = storage.AddInvoiceItem(invoiceID, item.Name, item.Amount, item.CostPerUnit)
			if err != nil {
				log.Printf("Error adding invoice item: %v", err)
				return nil, nil
			}
		}
	} else {
		// Update invoice
		err := storage.UpdateInvoice(c.invoiceID, c.providerID, c.clientID, c.paid)
		if err != nil {
			log.Printf("Error updating invoice: %v", err)
			return nil, nil
		}

		// Delete all existing items and re-add
		// This is simpler than trying to diff and update individual items
		for _, existingItem := range c.existingItems {
			err = storage.DeleteInvoiceItem(existingItem.ID)
			if err != nil {
				log.Printf("Error deleting invoice item: %v", err)
				// Continue anyway
			}
		}

		// Add all items
		for _, item := range c.items {
			_, err = storage.AddInvoiceItem(c.invoiceID, item.Name, item.Amount, item.CostPerUnit)
			if err != nil {
				log.Printf("Error adding invoice item: %v", err)
				return nil, nil
			}
		}
	}

	// Return to invoice list
	c.selection = ""
	invoiceForm, err := views.CreateInvoiceListForm(&c.selection)
	if err != nil {
		log.Printf("Error creating invoice form: %v", err)
		return nil, nil
	}
	c.form = invoiceForm
	return &types.ViewTransition{
		NewView: types.InvoicesListView,
		Form:    c.form,
	}, c.form.Init()
}

// resetFormFields clears all form field values
func (c *Controller) resetFormFields() {
	c.providerID = ""
	c.clientID = ""
	c.itemName = ""
	c.itemAmount = ""
	c.itemCostPerUnit = ""
	c.paid = false
	c.addAnother = false
	c.existingItems = nil
	c.items = nil
	c.currentStep = StepSelectProvider
	c.isEditMode = false
	c.currentItemIndex = 0
}

// GetForm returns the current form
func (c *Controller) GetForm() *huh.Form {
	return c.form
}
