package forms

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/GVPproj/termsheet/models"
	"github.com/GVPproj/termsheet/storage"
	"github.com/charmbracelet/huh"
)

// NewProviderSelectForm creates a form for selecting a provider
func NewProviderSelectForm(providerID *string) (*huh.Form, error) {
	providers, err := storage.ListProviders()
	if err != nil {
		return nil, err
	}

	providerOptions := make([]huh.Option[string], 0, len(providers))
	for _, p := range providers {
		label := p.Name
		if p.Email != nil {
			label += fmt.Sprintf(" (%s)", *p.Email)
		}
		providerOptions = append(providerOptions, huh.NewOption(label, p.ID))
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Provider").
				Options(providerOptions...).
				Value(providerID).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("provider is required")
					}
					return nil
				}),
		),
	), nil
}

// NewClientSelectForm creates a form for selecting a client
func NewClientSelectForm(clientID *string) (*huh.Form, error) {
	clients, err := storage.ListClients()
	if err != nil {
		return nil, err
	}

	clientOptions := make([]huh.Option[string], 0, len(clients))
	for _, c := range clients {
		label := c.Name
		if c.Email != nil {
			label += fmt.Sprintf(" (%s)", *c.Email)
		}
		clientOptions = append(clientOptions, huh.NewOption(label, c.ID))
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Client").
				Options(clientOptions...).
				Value(clientID).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("client is required")
					}
					return nil
				}),
		),
	), nil
}

// NewInvoiceItemForm creates a form for entering a single invoice item
func NewInvoiceItemForm(itemName, itemAmount, itemCostPerUnit *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Item Name").
				Value(itemName).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("item name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Amount").
				Value(itemAmount).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("amount is required")
					}
					amount, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return errors.New("amount must be a number")
					}
					if amount <= 0 {
						return errors.New("amount must be positive")
					}
					return nil
				}),
			huh.NewInput().
				Title("Cost Per Unit").
				Value(itemCostPerUnit).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("cost per unit is required")
					}
					cost, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return errors.New("cost per unit must be a number")
					}
					if cost <= 0 {
						return errors.New("cost per unit must be positive")
					}
					return nil
				}),
		),
	)
}

// NewAddAnotherItemForm creates a confirmation form for adding another item
func NewAddAnotherItemForm(addAnother *bool) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add another item?").
				Affirmative("Yes").
				Negative("No").
				Value(addAnother),
		),
	)
}

// NewMarkPaidForm creates a form for marking invoice as paid
func NewMarkPaidForm(paid *bool) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Mark as Paid?").
				Value(paid),
		),
	)
}

// NewProviderSelectFormWithData creates a provider select form with pre-populated data
func NewProviderSelectFormWithData(providerID *string, existingProviderID string) (*huh.Form, error) {
	*providerID = existingProviderID
	return NewProviderSelectForm(providerID)
}

// NewClientSelectFormWithData creates a client select form with pre-populated data
func NewClientSelectFormWithData(clientID *string, existingClientID string) (*huh.Form, error) {
	*clientID = existingClientID
	return NewClientSelectForm(clientID)
}

// NewInvoiceItemFormWithData creates an item form with pre-populated data
func NewInvoiceItemFormWithData(itemName, itemAmount, itemCostPerUnit *string, item models.InvoiceItem) *huh.Form {
	*itemName = item.ItemName
	*itemAmount = fmt.Sprintf("%.2f", item.Amount)
	*itemCostPerUnit = fmt.Sprintf("%.2f", item.CostPerUnit)
	return NewInvoiceItemForm(itemName, itemAmount, itemCostPerUnit)
}

// NewMarkPaidFormWithData creates a mark paid form with pre-populated data
func NewMarkPaidFormWithData(paid *bool, existingPaid bool) *huh.Form {
	*paid = existingPaid
	return NewMarkPaidForm(paid)
}

// NewDeleteConfirmForm creates a confirmation form for deleting an item
func NewDeleteConfirmForm(confirmed *bool) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are You Sure?").
				Value(confirmed),
		),
	)
}
