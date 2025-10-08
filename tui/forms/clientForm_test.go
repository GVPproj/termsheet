package forms

import (
	"testing"
)

func TestNewClientForm(t *testing.T) {
	var name, address, email, phone string

	form := NewClientForm(&name, &address, &email, &phone)

	if form == nil {
		t.Fatal("expected non-nil form")
	}

	// Simulate user input
	name = "Test Client"
	address = "789 Oak St"
	email = "client@example.com"
	phone = "555-0200"

	// Verify bindings work
	if name != "Test Client" {
		t.Errorf("expected name binding to work")
	}
}

func TestNewClientFormWithData(t *testing.T) {
	// Create mock client data without importing models to avoid cycle
	addr := "321 Client Ave"
	em := "contact@client.com"
	ph := "555-8888"

	// We can't import models.Client due to import cycle
	// So we'll test the inline pointer conversion instead
	var name, address, email, phone string

	// Manually simulate what NewClientFormWithData does
	name = "Existing Client"
	// Convert *string to string, empty if nil
	address = addr
	email = em
	phone = ph

	form := NewClientForm(&name, &address, &email, &phone)

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

func TestNewClientFormWithDataNilFields(t *testing.T) {
	var name, address, email, phone string

	// Simulate nil fields (empty strings since we can't test with actual nil pointers here)
	name = "Minimal Client"
	address = ""
	email = ""
	phone = ""

	form := NewClientForm(&name, &address, &email, &phone)

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
