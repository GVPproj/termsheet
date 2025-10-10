package views

import (
	"testing"

	"github.com/GVPproj/termsheet/storage"
)

// TestInvoiceActionMenuIntegration tests the full flow from invoice list to action menu
func TestInvoiceActionMenuIntegration(t *testing.T) {
	// Initialize test database
	if err := storage.InitDB(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer storage.CloseDB()

	// Create test data
	providerID, err := storage.CreateProvider("Test Provider", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	clientID, err := storage.CreateClient("Test Client", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	invoiceID, err := storage.CreateInvoice(providerID, clientID, false)
	if err != nil {
		t.Fatalf("Failed to create invoice: %v", err)
	}

	_, err = storage.AddInvoiceItem(invoiceID, "Test Item", 10, 100.00)
	if err != nil {
		t.Fatalf("Failed to add invoice item: %v", err)
	}

	// Test 1: Invoice list form can be created
	var selection string
	invoiceListForm, err := CreateInvoiceListForm(&selection)
	if err != nil {
		t.Fatalf("Failed to create invoice list form: %v", err)
	}
	if invoiceListForm == nil {
		t.Fatal("Invoice list form should not be nil")
	}

	// Test 2: Action menu form can be created
	var actionSelection string
	actionForm := CreateInvoiceActionForm(&actionSelection)
	if actionForm == nil {
		t.Fatal("Action form should not be nil")
	}

	// Test 3: Action menu renders correctly
	rendered := RenderInvoiceActionMenu(actionForm)
	if rendered == "" {
		t.Error("Action menu should render non-empty string")
	}

	// Test 4: Invoice data can be retrieved and displayed
	invoiceData, err := storage.GetInvoiceData(invoiceID)
	if err != nil {
		t.Fatalf("Failed to get invoice data: %v", err)
	}

	viewRendered := RenderInvoiceView(invoiceData)
	if viewRendered == "" {
		t.Error("Invoice view should render non-empty string")
	}

	// Verify the rendered view contains expected data
	if !contains(viewRendered, "Test Provider") {
		t.Error("Invoice view should contain provider name")
	}
	if !contains(viewRendered, "Test Client") {
		t.Error("Invoice view should contain client name")
	}
	if !contains(viewRendered, "Test Item") {
		t.Error("Invoice view should contain item name")
	}
}

// TestActionMenuOptions tests that all action options are available
func TestActionMenuOptions(t *testing.T) {
	var selection string
	form := CreateInvoiceActionForm(&selection)

	if form == nil {
		t.Fatal("Form should not be nil")
	}

	// Verify the form renders (can't access options directly)
	rendered := form.View()
	if rendered == "" {
		t.Error("Form should render something")
	}
}

// TestInvoiceViewWithEmptyData tests rendering with minimal data
func TestInvoiceViewWithEmptyData(t *testing.T) {
	// Initialize test database
	if err := storage.InitDB(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer storage.CloseDB()

	// Create minimal invoice
	providerID, err := storage.CreateProvider("Provider", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	clientID, err := storage.CreateClient("Client", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	invoiceID, err := storage.CreateInvoice(providerID, clientID, false)
	if err != nil {
		t.Fatalf("Failed to create invoice: %v", err)
	}

	// Get and render invoice with no items
	invoiceData, err := storage.GetInvoiceData(invoiceID)
	if err != nil {
		t.Fatalf("Failed to get invoice data: %v", err)
	}

	rendered := RenderInvoiceView(invoiceData)
	if rendered == "" {
		t.Error("Should render even with no items")
	}

	// Should show "No items" message
	if !contains(rendered, "No items") {
		t.Error("Should display 'No items' message when invoice has no items")
	}
}
