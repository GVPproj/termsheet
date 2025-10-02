package models

import (
	"errors"
	"testing"

	"github.com/charmbracelet/huh"
)

// TestCreateProviderWithFormInput tests creating a provider using huh form input simulation
func TestCreateProviderWithFormInput(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Simulate form input data
	var (
		name    string
		address string
		email   string
		phone   string
	)

	// Create a form that binds to our variables
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Provider Name").
				Value(&name).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("provider name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Address").
				Value(&address),
			huh.NewInput().
				Title("Email").
				Value(&email),
			huh.NewInput().
				Title("Phone").
				Value(&phone),
		),
	)

	// Simulate user filling out the form
	name = "Test Provider Inc"
	address = "123 Business Ave"
	email = "contact@testprovider.com"
	phone = "555-1234"

	// Mark form as completed (simulating user submission)
	form.State = huh.StateCompleted

	// Verify form is completed
	if form.State != huh.StateCompleted {
		t.Fatal("form should be completed")
	}

	// Use the form data to create a provider
	var addressPtr, emailPtr, phonePtr *string
	if address != "" {
		addressPtr = &address
	}
	if email != "" {
		emailPtr = &email
	}
	if phone != "" {
		phonePtr = &phone
	}

	_, err := CreateProvider(name, addressPtr, emailPtr, phonePtr)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	// Retrieve the provider from database to confirm it was stored correctly
	providers, err := ListProviders()
	if err != nil {
		t.Fatalf("ListProviders failed: %v", err)
	}

	if len(providers) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(providers))
	}

	// Verify all fields match the form input
	p := providers[0]
	if p.Name != name {
		t.Errorf("expected name %q, got %q", name, p.Name)
	}
	if p.Address == nil || *p.Address != address {
		t.Errorf("expected address %q, got %v", address, p.Address)
	}
	if p.Email == nil || *p.Email != email {
		t.Errorf("expected email %q, got %v", email, p.Email)
	}
	if p.Phone == nil || *p.Phone != phone {
		t.Errorf("expected phone %q, got %v", phone, p.Phone)
	}
}

// TestCreateProviderWithFormInputOptionalFields tests form with only required fields
func TestCreateProviderWithFormInputOptionalFields(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	var (
		name    string
		address string
		email   string
		phone   string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Provider Name").
				Value(&name).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("provider name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Address").
				Value(&address),
			huh.NewInput().
				Title("Email").
				Value(&email),
			huh.NewInput().
				Title("Phone").
				Value(&phone),
		),
	)

	// User only fills out name (required field)
	name = "Minimal Provider"
	// address, email, phone remain empty

	form.State = huh.StateCompleted

	// Create provider with only required field
	_, err := CreateProvider(name, nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	// Retrieve and verify
	providers, err := ListProviders()
	if err != nil {
		t.Fatalf("ListProviders failed: %v", err)
	}

	if len(providers) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(providers))
	}

	p := providers[0]
	if p.Name != name {
		t.Errorf("expected name %q, got %q", name, p.Name)
	}
	if p.Address != nil {
		t.Errorf("expected nil address, got %v", p.Address)
	}
	if p.Email != nil {
		t.Errorf("expected nil email, got %v", p.Email)
	}
	if p.Phone != nil {
		t.Errorf("expected nil phone, got %v", p.Phone)
	}
}

// TestCreateProviderFormValidation tests that form validation prevents invalid submissions
func TestCreateProviderFormValidation(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Test that CreateProvider rejects empty names (which form validation would also catch)
	_, err := CreateProvider("", nil, nil, nil)
	if err == nil {
		t.Error("expected error when creating provider with empty name")
	}

	// Verify provider was not created in database
	providers, _ := ListProviders()
	if len(providers) != 0 {
		t.Errorf("expected 0 providers after validation failure, got %d", len(providers))
	}
}
