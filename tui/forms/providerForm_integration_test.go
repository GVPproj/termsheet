package forms_test

import (
	"testing"

	"github.com/GVPproj/termsheet/models"
	"github.com/GVPproj/termsheet/tui/forms"
)

// TestNewProviderFormWithDataIntegration tests the integration between models.Provider and forms
func TestNewProviderFormWithDataIntegration(t *testing.T) {
	// This is a unit test, so we'll use mock data instead of requiring DB setup
	// Create a provider with all fields populated
	addr := "456 Business Ave"
	em := "contact@business.com"
	ph := "555-9999"

	provider := models.Provider{
		Name:    "Existing Provider",
		Address: &addr,
		Email:   &em,
		Phone:   &ph,
	}

	var name, address, email, phone string
	form := forms.NewProviderFormWithData(provider, &name, &address, &email, &phone)

	if form == nil {
		t.Fatal("expected non-nil form")
	}

	// Verify data was pre-populated
	if name != "Existing Provider" {
		t.Errorf("expected name %q, got %q", "Existing Provider", name)
	}
	if address != "456 Business Ave" {
		t.Errorf("expected address %q, got %q", "456 Business Ave", address)
	}
	if email != "contact@business.com" {
		t.Errorf("expected email %q, got %q", "contact@business.com", email)
	}
	if phone != "555-9999" {
		t.Errorf("expected phone %q, got %q", "555-9999", phone)
	}
}

// TestNewProviderFormWithDataNilFieldsIntegration tests the form with a minimal provider
func TestNewProviderFormWithDataNilFieldsIntegration(t *testing.T) {
	provider := models.Provider{
		Name:    "Minimal Provider",
		Address: nil,
		Email:   nil,
		Phone:   nil,
	}

	var name, address, email, phone string
	form := forms.NewProviderFormWithData(provider, &name, &address, &email, &phone)

	if form == nil {
		t.Fatal("expected non-nil form")
	}

	// Verify nil fields become empty strings
	if name != "Minimal Provider" {
		t.Errorf("expected name %q, got %q", "Minimal Provider", name)
	}
	if address != "" {
		t.Errorf("expected empty address, got %q", address)
	}
	if email != "" {
		t.Errorf("expected empty email, got %q", email)
	}
	if phone != "" {
		t.Errorf("expected empty phone, got %q", phone)
	}
}
