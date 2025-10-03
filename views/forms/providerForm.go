package forms

import (
	"errors"

	"github.com/GVPproj/termsheet/models"
	"github.com/charmbracelet/huh"
)

// NewProviderForm creates a new form for provider input
// Takes pointers to string variables that will be bound to form fields
func NewProviderForm(name, address, email, phone *string) *huh.Form {
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
			huh.NewInput().
				Title("Phone").
				Value(phone),
		),
	)
}

// NewProviderFormWithData creates a form pre-populated with existing provider data
// Useful for editing an existing provider
func NewProviderFormWithData(provider models.Provider, name, address, email, phone *string) *huh.Form {
	// Pre-populate the bound variables with existing data
	*name = provider.Name

	// Convert *string to string, empty if nil
	if provider.Address != nil {
		*address = *provider.Address
	} else {
		*address = ""
	}

	if provider.Email != nil {
		*email = *provider.Email
	} else {
		*email = ""
	}

	if provider.Phone != nil {
		*phone = *provider.Phone
	} else {
		*phone = ""
	}

	return NewProviderForm(name, address, email, phone)
}
