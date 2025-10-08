package forms_test

import (
	"testing"

	"github.com/GVPproj/termsheet/models"
	"github.com/GVPproj/termsheet/tui/forms"
)

// TestNewClientFormWithDataIntegration tests the integration between models.Client and forms
func TestNewClientFormWithDataIntegration(t *testing.T) {
	// This is a unit test, so we'll use mock data instead of requiring DB setup
	// Create a client with all fields populated
	addr := "321 Client Ave"
	em := "contact@client.com"
	ph := "555-8888"

	client := models.Client{
		Name:    "Existing Client",
		Address: &addr,
		Email:   &em,
		Phone:   &ph,
	}

	var name, address, email, phone string
	form := forms.NewClientFormWithData(client, &name, &address, &email, &phone)

	if form == nil {
		t.Fatal("expected non-nil form")
	}

	// Verify data was pre-populated
	if name != "Existing Client" {
		t.Errorf("expected name %q, got %q", "Existing Client", name)
	}
	if address != "321 Client Ave" {
		t.Errorf("expected address %q, got %q", "321 Client Ave", address)
	}
	if email != "contact@client.com" {
		t.Errorf("expected email %q, got %q", "contact@client.com", email)
	}
	if phone != "555-8888" {
		t.Errorf("expected phone %q, got %q", "555-8888", phone)
	}
}

// TestNewClientFormWithDataNilFieldsIntegration tests the form with a minimal client
func TestNewClientFormWithDataNilFieldsIntegration(t *testing.T) {
	client := models.Client{
		Name:    "Minimal Client",
		Address: nil,
		Email:   nil,
		Phone:   nil,
	}

	var name, address, email, phone string
	form := forms.NewClientFormWithData(client, &name, &address, &email, &phone)

	if form == nil {
		t.Fatal("expected non-nil form")
	}

	// Verify nil fields become empty strings
	if name != "Minimal Client" {
		t.Errorf("expected name %q, got %q", "Minimal Client", name)
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

// TestNewClientFormWithDataMixedFields tests the form with some nil and some populated fields
func TestNewClientFormWithDataMixedFields(t *testing.T) {
	addr := "999 Mixed St"

	client := models.Client{
		Name:    "Partial Client",
		Address: &addr,
		Email:   nil,
		Phone:   nil,
	}

	var name, address, email, phone string
	form := forms.NewClientFormWithData(client, &name, &address, &email, &phone)

	if form == nil {
		t.Fatal("expected non-nil form")
	}

	// Verify mixed fields
	if name != "Partial Client" {
		t.Errorf("expected name %q, got %q", "Partial Client", name)
	}
	if address != "999 Mixed St" {
		t.Errorf("expected address %q, got %q", "999 Mixed St", address)
	}
	if email != "" {
		t.Errorf("expected empty email, got %q", email)
	}
	if phone != "" {
		t.Errorf("expected empty phone, got %q", phone)
	}
}
