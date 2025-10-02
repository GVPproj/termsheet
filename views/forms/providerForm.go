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
	*address = StringPtrToString(provider.Address)
	*email = StringPtrToString(provider.Email)
	*phone = StringPtrToString(provider.Phone)

	return NewProviderForm(name, address, email, phone)
}

// StringPtrToString converts a *string to string, returning empty string if nil
func StringPtrToString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// StringToStringPtr converts a non-empty string to *string, returning nil for empty strings
func StringToStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
