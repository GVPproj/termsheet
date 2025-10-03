package storage

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

// setupTestDB initializes an in-memory SQLite database for testing
func setupTestDB(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	if err := createTables(); err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}
}

// teardownTestDB closes the test database connection
func teardownTestDB(t *testing.T) {
	if db != nil {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close test database: %v", err)
		}
	}
}

// TestCreateClient tests creating a client
func TestCreateClient(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	address := "123 Main St"
	email := "test@example.com"

	id, err := CreateClient("Acme Corp", &address, &email)
	if err != nil {
		t.Fatalf("CreateClient failed: %v", err)
	}

	if id == "" {
		t.Error("expected non-empty client ID")
	}
}

// TestCreateClientEmptyName tests that empty names are rejected
func TestCreateClientEmptyName(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	testCases := []struct {
		name     string
		input    string
		wantErr  bool
	}{
		{"empty string", "", true},
		{"whitespace only", "   ", true},
		{"valid name", "Valid Client", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CreateClient(tc.input, nil, nil)
			if (err != nil) != tc.wantErr {
				t.Errorf("CreateClient(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			}
		})
	}
}

// TestListClients tests listing clients
func TestListClients(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Test empty list
	clients, err := ListClients()
	if err != nil {
		t.Fatalf("ListClients failed: %v", err)
	}
	if len(clients) != 0 {
		t.Errorf("expected empty list, got %d clients", len(clients))
	}

	// Create test clients
	address1 := "123 Main St"
	email1 := "client1@example.com"
	_, err = CreateClient("Client One", &address1, &email1)
	if err != nil {
		t.Fatalf("CreateClient failed: %v", err)
	}

	_, err = CreateClient("Client Two", nil, nil)
	if err != nil {
		t.Fatalf("CreateClient failed: %v", err)
	}

	// List clients
	clients, err = ListClients()
	if err != nil {
		t.Fatalf("ListClients failed: %v", err)
	}

	if len(clients) != 2 {
		t.Errorf("expected 2 clients, got %d", len(clients))
	}

	// Verify first client
	if clients[0].Name != "Client One" {
		t.Errorf("expected name 'Client One', got %q", clients[0].Name)
	}
	if clients[0].Address == nil || *clients[0].Address != address1 {
		t.Error("expected address to match")
	}
	if clients[0].Email == nil || *clients[0].Email != email1 {
		t.Error("expected email to match")
	}

	// Verify second client has nil optional fields
	if clients[1].Address != nil {
		t.Error("expected nil address")
	}
	if clients[1].Email != nil {
		t.Error("expected nil email")
	}
}

// TestUpdateClient tests updating a client
func TestUpdateClient(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create a client
	address := "123 Main St"
	id, err := CreateClient("Original Name", &address, nil)
	if err != nil {
		t.Fatalf("CreateClient failed: %v", err)
	}

	// Update the client
	newAddress := "456 Elm St"
	newEmail := "new@example.com"
	err = UpdateClient(id, "Updated Name", &newAddress, &newEmail)
	if err != nil {
		t.Fatalf("UpdateClient failed: %v", err)
	}

	// Verify update
	clients, _ := ListClients()
	if len(clients) != 1 {
		t.Fatalf("expected 1 client, got %d", len(clients))
	}

	if clients[0].Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got %q", clients[0].Name)
	}
	if clients[0].Address == nil || *clients[0].Address != newAddress {
		t.Error("expected updated address")
	}
	if clients[0].Email == nil || *clients[0].Email != newEmail {
		t.Error("expected updated email")
	}
}

// TestUpdateClientNonExistent tests updating a non-existent client
func TestUpdateClientNonExistent(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	err := UpdateClient("non-existent-id", "Name", nil, nil)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

// TestUpdateClientEmptyName tests that empty names are rejected in updates
func TestUpdateClientEmptyName(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	id, _ := CreateClient("Original", nil, nil)

	err := UpdateClient(id, "", nil, nil)
	if err == nil {
		t.Error("expected error for empty name")
	}

	err = UpdateClient(id, "   ", nil, nil)
	if err == nil {
		t.Error("expected error for whitespace-only name")
	}
}

// TestCreateProvider tests creating a provider
func TestCreateProvider(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	address := "789 Oak Ave"
	email := "provider@example.com"
	phone := "555-1234"

	id, err := CreateProvider("Test Provider", &address, &email, &phone)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	if id == "" {
		t.Error("expected non-empty provider ID")
	}
}

// TestCreateProviderEmptyName tests that empty names are rejected
func TestCreateProviderEmptyName(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	_, err := CreateProvider("", nil, nil, nil)
	if err == nil {
		t.Error("expected error for empty name")
	}

	_, err = CreateProvider("   ", nil, nil, nil)
	if err == nil {
		t.Error("expected error for whitespace-only name")
	}
}

// TestListProviders tests listing providers
func TestListProviders(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Test empty list
	providers, err := ListProviders()
	if err != nil {
		t.Fatalf("ListProviders failed: %v", err)
	}
	if len(providers) != 0 {
		t.Errorf("expected empty list, got %d providers", len(providers))
	}

	// Create test providers
	phone := "555-1234"
	_, err = CreateProvider("Provider One", nil, nil, &phone)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	_, err = CreateProvider("Provider Two", nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	// List providers
	providers, err = ListProviders()
	if err != nil {
		t.Fatalf("ListProviders failed: %v", err)
	}

	if len(providers) != 2 {
		t.Errorf("expected 2 providers, got %d", len(providers))
	}
}

// TestUpdateProvider tests updating a provider
func TestUpdateProvider(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	id, err := CreateProvider("Original Provider", nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	newPhone := "555-9999"
	err = UpdateProvider(id, "Updated Provider", nil, nil, &newPhone)
	if err != nil {
		t.Fatalf("UpdateProvider failed: %v", err)
	}

	providers, _ := ListProviders()
	if len(providers) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(providers))
	}

	if providers[0].Name != "Updated Provider" {
		t.Errorf("expected name 'Updated Provider', got %q", providers[0].Name)
	}
	if providers[0].Phone == nil || *providers[0].Phone != newPhone {
		t.Error("expected updated phone")
	}
}

// TestUpdateProviderNonExistent tests updating a non-existent provider
func TestUpdateProviderNonExistent(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	err := UpdateProvider("non-existent-id", "Name", nil, nil, nil)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

// TestCreateInvoice tests creating an invoice
func TestCreateInvoice(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create provider and client first
	providerID, _ := CreateProvider("Test Provider", nil, nil, nil)
	clientID, _ := CreateClient("Test Client", nil, nil)

	invoiceID, err := CreateInvoice(providerID, clientID, false)
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	if invoiceID == 0 {
		t.Error("expected non-zero invoice ID")
	}
}

// TestUpdateInvoice tests updating an invoice
func TestUpdateInvoice(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Setup
	providerID, _ := CreateProvider("Provider", nil, nil, nil)
	clientID, _ := CreateClient("Client", nil, nil)
	invoiceID, _ := CreateInvoice(providerID, clientID, false)

	// Update
	err := UpdateInvoice(invoiceID, providerID, clientID, true)
	if err != nil {
		t.Fatalf("UpdateInvoice failed: %v", err)
	}

	// Verify
	invoices, _ := ListInvoices()
	if len(invoices) != 1 {
		t.Fatalf("expected 1 invoice, got %d", len(invoices))
	}

	if !invoices[0].Paid {
		t.Error("expected invoice to be marked as paid")
	}
}

// TestUpdateInvoiceNonExistent tests updating a non-existent invoice
func TestUpdateInvoiceNonExistent(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	err := UpdateInvoice(999, "provider-id", "client-id", false)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

// TestListInvoices tests listing invoices
func TestListInvoices(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Test empty list
	invoices, err := ListInvoices()
	if err != nil {
		t.Fatalf("ListInvoices failed: %v", err)
	}
	if len(invoices) != 0 {
		t.Errorf("expected empty list, got %d invoices", len(invoices))
	}

	// Create test data
	providerID, _ := CreateProvider("Test Provider", nil, nil, nil)
	clientID, _ := CreateClient("Test Client", nil, nil)
	_, err = CreateInvoice(providerID, clientID, false)
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	// List invoices
	invoices, err = ListInvoices()
	if err != nil {
		t.Fatalf("ListInvoices failed: %v", err)
	}

	if len(invoices) != 1 {
		t.Errorf("expected 1 invoice, got %d", len(invoices))
	}

	if invoices[0].ProviderName != "Test Provider" {
		t.Errorf("expected provider name 'Test Provider', got %q", invoices[0].ProviderName)
	}
	if invoices[0].ClientName != "Test Client" {
		t.Errorf("expected client name 'Test Client', got %q", invoices[0].ClientName)
	}
}

// TestAddInvoiceItem tests adding items to an invoice
func TestAddInvoiceItem(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Setup
	providerID, _ := CreateProvider("Provider", nil, nil, nil)
	clientID, _ := CreateClient("Client", nil, nil)
	invoiceID, _ := CreateInvoice(providerID, clientID, false)

	// Add item
	itemID, err := AddInvoiceItem(invoiceID, "Widget", 10.0, 25.50)
	if err != nil {
		t.Fatalf("AddInvoiceItem failed: %v", err)
	}

	if itemID == 0 {
		t.Error("expected non-zero item ID")
	}
}

// TestAddInvoiceItemValidation tests invoice item validation
func TestAddInvoiceItemValidation(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	providerID, _ := CreateProvider("Provider", nil, nil, nil)
	clientID, _ := CreateClient("Client", nil, nil)
	invoiceID, _ := CreateInvoice(providerID, clientID, false)

	testCases := []struct {
		name        string
		itemName    string
		amount      float64
		costPerUnit float64
		wantErr     bool
	}{
		{"valid item", "Widget", 5.0, 10.0, false},
		{"empty name", "", 5.0, 10.0, true},
		{"whitespace name", "   ", 5.0, 10.0, true},
		{"zero amount", "Widget", 0, 10.0, true},
		{"negative amount", "Widget", -5.0, 10.0, true},
		{"zero cost", "Widget", 5.0, 0, true},
		{"negative cost", "Widget", 5.0, -10.0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := AddInvoiceItem(invoiceID, tc.itemName, tc.amount, tc.costPerUnit)
			if (err != nil) != tc.wantErr {
				t.Errorf("AddInvoiceItem() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// TestGetInvoiceData tests retrieving complete invoice data
func TestGetInvoiceData(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Setup
	providerAddress := "123 Provider St"
	providerEmail := "provider@test.com"
	providerPhone := "555-0000"
	providerID, _ := CreateProvider("My Provider", &providerAddress, &providerEmail, &providerPhone)

	clientAddress := "456 Client Ave"
	clientEmail := "client@test.com"
	clientID, _ := CreateClient("My Client", &clientAddress, &clientEmail)

	invoiceID, _ := CreateInvoice(providerID, clientID, true)

	// Add items
	_, _ = AddInvoiceItem(invoiceID, "Item 1", 2.0, 50.0)
	_, _ = AddInvoiceItem(invoiceID, "Item 2", 1.0, 75.0)

	// Get invoice data
	data, err := GetInvoiceData(invoiceID)
	if err != nil {
		t.Fatalf("GetInvoiceData failed: %v", err)
	}

	// Verify invoice data
	if data.InvoiceID != invoiceID {
		t.Errorf("expected invoice ID %d, got %d", invoiceID, data.InvoiceID)
	}
	if !data.Paid {
		t.Error("expected invoice to be paid")
	}
	if data.ProviderName != "My Provider" {
		t.Errorf("expected provider name 'My Provider', got %q", data.ProviderName)
	}
	if data.ClientName != "My Client" {
		t.Errorf("expected client name 'My Client', got %q", data.ClientName)
	}

	// Verify items
	if len(data.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(data.Items))
	}

	if data.Items[0].ItemName != "Item 1" {
		t.Errorf("expected item name 'Item 1', got %q", data.Items[0].ItemName)
	}
	if data.Items[0].Amount != 2.0 {
		t.Errorf("expected amount 2.0, got %f", data.Items[0].Amount)
	}
	if data.Items[0].CostPerUnit != 50.0 {
		t.Errorf("expected cost per unit 50.0, got %f", data.Items[0].CostPerUnit)
	}
}

// TestGetInvoiceDataNonExistent tests retrieving non-existent invoice
func TestGetInvoiceDataNonExistent(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	_, err := GetInvoiceData(999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}
