package forms

import (
	"errors"

	"github.com/GVPproj/termsheet/models"
	"github.com/charmbracelet/huh"
)

// NewProviderForm creates a new form for provider input
// Takes pointers to string variables that will be bound to form fields
func NewClientForm(name, address, email *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Provider Name").
				Value(name).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("provider name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Address").
				Value(address),
			huh.NewInput().
				Title("Email").
				Value(email),
		),
	)
}

// NewProviderFormWithData creates a form pre-populated with existing provider data
// Useful for editing an existing provider
func NewClientFormWithData(client models.Client, name, address, email *string) *huh.Form {
	// Pre-populate the bound variables with existing data
	*name = client.Name

	// Convert *string to string, empty if nil
	if client.Address != nil {
		*address = *client.Address
	} else {
		*address = ""
	}

	if client.Email != nil {
		*email = *client.Email
	} else {
		*email = ""
	}

	return NewClientForm(name, address, email)
}
