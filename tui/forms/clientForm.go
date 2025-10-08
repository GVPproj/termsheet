package forms

import (
	"errors"

	"github.com/GVPproj/termsheet/models"
	"github.com/charmbracelet/huh"
)

// NewClientForm creates a new form for client input
// Takes pointers to string variables that will be bound to form fields
func NewClientForm(name, address, email, phone *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Client Name").
				Value(name).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("client name is required")
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

// NewClientFormWithData creates a form pre-populated with existing client data
// Useful for editing an existing client
func NewClientFormWithData(client models.Client, name, address, email, phone *string) *huh.Form {
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

	if client.Phone != nil {
		*phone = *client.Phone
	} else {
		*phone = ""
	}

	return NewClientForm(name, address, email, phone)
}
