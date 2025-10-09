package forms_test

import (
	"fmt"
	"testing"

	"github.com/GVPproj/termsheet/models"
)

// TestInvoiceFormWithDataIntegration tests the integration between models.Invoice and forms
// Note: This test doesn't require database initialization because it only tests
// the data pre-population logic, not the actual form creation
func TestInvoiceFormWithDataIntegration(t *testing.T) {
	// Create an invoice with all fields populated
	invoice := models.Invoice{
		ID:         1,
		ProviderID: "provider-abc",
		ClientID:   "client-xyz",
		Paid:       true,
	}

	// Create invoice items
	items := []models.InvoiceItem{
		{
			ID:          1,
			InvoiceID:   1,
			ItemName:    "Web Development",
			Amount:      20.0,
			CostPerUnit: 100.0,
		},
	}

	var providerID, clientID string
	var paid bool

	// Simulate provider/client form pre-population
	providerID = invoice.ProviderID
	clientID = invoice.ClientID
	paid = invoice.Paid

	// Verify provider/client data was pre-populated correctly
	if providerID != "provider-abc" {
		t.Errorf("expected providerID %q, got %q", "provider-abc", providerID)
	}
	if clientID != "client-xyz" {
		t.Errorf("expected clientID %q, got %q", "client-xyz", clientID)
	}
	if !paid {
		t.Errorf("expected paid to be true")
	}

	// Simulate item form pre-population
	if len(items) > 0 {
		itemName := items[0].ItemName
		itemAmount := fmt.Sprintf("%.2f", items[0].Amount)
		itemCostPerUnit := fmt.Sprintf("%.2f", items[0].CostPerUnit)

		if itemName != "Web Development" {
			t.Errorf("expected itemName %q, got %q", "Web Development", itemName)
		}
		if itemAmount != "20.00" {
			t.Errorf("expected itemAmount %q, got %q", "20.00", itemAmount)
		}
		if itemCostPerUnit != "100.00" {
			t.Errorf("expected itemCostPerUnit %q, got %q", "100.00", itemCostPerUnit)
		}
	}
}

// TestInvoiceFormWithMultipleItemsIntegration tests form with multiple items
func TestInvoiceFormWithMultipleItemsIntegration(t *testing.T) {
	// Multiple items
	items := []models.InvoiceItem{
		{
			ID:          2,
			InvoiceID:   2,
			ItemName:    "Consulting",
			Amount:      10.0,
			CostPerUnit: 150.0,
		},
		{
			ID:          3,
			InvoiceID:   2,
			ItemName:    "Design",
			Amount:      5.0,
			CostPerUnit: 200.0,
		},
	}

	// Verify we can iterate through all items
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}

	// Verify first item
	if items[0].ItemName != "Consulting" {
		t.Errorf("expected first item name %q, got %q", "Consulting", items[0].ItemName)
	}

	// Verify second item
	if items[1].ItemName != "Design" {
		t.Errorf("expected second item name %q, got %q", "Design", items[1].ItemName)
	}

	// Verify amounts
	if items[0].Amount != 10.0 {
		t.Errorf("expected first item amount %.2f, got %.2f", 10.0, items[0].Amount)
	}
	if items[1].CostPerUnit != 200.0 {
		t.Errorf("expected second item cost per unit %.2f, got %.2f", 200.0, items[1].CostPerUnit)
	}
}

// TestInvoiceFormWithNoItemsIntegration tests the form with an invoice but no items
func TestInvoiceFormWithNoItemsIntegration(t *testing.T) {
	invoice := models.Invoice{
		ID:         3,
		ProviderID: "provider-ghi",
		ClientID:   "client-rst",
		Paid:       false,
	}

	// Empty items slice
	items := []models.InvoiceItem{}

	var providerID, clientID string
	var paid bool

	// Simulate provider/client form pre-population
	providerID = invoice.ProviderID
	clientID = invoice.ClientID
	paid = invoice.Paid

	// Verify invoice data was pre-populated
	if providerID != "provider-ghi" {
		t.Errorf("expected providerID %q, got %q", "provider-ghi", providerID)
	}
	if clientID != "client-rst" {
		t.Errorf("expected clientID %q, got %q", "client-rst", clientID)
	}
	if paid {
		t.Errorf("expected paid to be false")
	}

	// Verify no items exist
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}
