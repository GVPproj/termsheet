package views

import (
	"strings"
	"testing"
	"time"

	"github.com/GVPproj/termsheet/models"
)

func TestRenderInvoiceView(t *testing.T) {
	address := "123 Main St"
	email := "test@example.com"
	phone := "555-1234"

	data := &models.InvoiceData{
		InvoiceID:   1,
		DateCreated: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Paid:        false,
		Provider: models.Entity{
			ID:      "p1",
			Name:    "Test Provider",
			Address: &address,
			Email:   &email,
			Phone:   &phone,
		},
		Client: models.Entity{
			ID:      "c1",
			Name:    "Test Client",
			Address: &address,
			Email:   &email,
			Phone:   &phone,
		},
		Items: []models.InvoiceItem{
			{
				ID:          1,
				InvoiceID:   1,
				ItemName:    "Consulting",
				Amount:      10,
				CostPerUnit: 100.00,
			},
			{
				ID:          2,
				InvoiceID:   1,
				ItemName:    "Development",
				Amount:      5,
				CostPerUnit: 150.00,
			},
		},
	}

	rendered := RenderInvoiceView(data)

	// Basic checks
	if rendered == "" {
		t.Fatal("RenderInvoiceView should return non-empty string")
	}

	// Check for invoice number
	if !strings.Contains(rendered, "Invoice #1") {
		t.Error("Rendered output should contain invoice number")
	}

	// Check for status
	if !strings.Contains(rendered, "Unpaid") {
		t.Error("Rendered output should contain payment status")
	}

	// Check for provider name
	if !strings.Contains(rendered, "Test Provider") {
		t.Error("Rendered output should contain provider name")
	}

	// Check for client name
	if !strings.Contains(rendered, "Test Client") {
		t.Error("Rendered output should contain client name")
	}

	// Check for items
	if !strings.Contains(rendered, "Consulting") {
		t.Error("Rendered output should contain item name")
	}

	// Check for total
	if !strings.Contains(rendered, "1750.00") {
		t.Error("Rendered output should contain correct total (10*100 + 5*150 = 1750)")
	}
}

func TestRenderInvoiceViewPaid(t *testing.T) {
	data := &models.InvoiceData{
		InvoiceID:   2,
		DateCreated: time.Now(),
		Paid:        true,
		Provider: models.Entity{
			ID:   "p1",
			Name: "Provider",
		},
		Client: models.Entity{
			ID:   "c1",
			Name: "Client",
		},
		Items: []models.InvoiceItem{},
	}

	rendered := RenderInvoiceView(data)

	if !strings.Contains(rendered, "Paid") {
		t.Error("Rendered output should contain 'Paid' status")
	}
}

func TestCalculateTotal(t *testing.T) {
	tests := []struct {
		name     string
		items    []models.InvoiceItem
		expected float64
	}{
		{
			name:     "empty items",
			items:    []models.InvoiceItem{},
			expected: 0.0,
		},
		{
			name: "single item",
			items: []models.InvoiceItem{
				{ItemName: "Item1", Amount: 5, CostPerUnit: 10.0},
			},
			expected: 50.0,
		},
		{
			name: "multiple items",
			items: []models.InvoiceItem{
				{ItemName: "Item1", Amount: 5, CostPerUnit: 10.0},
				{ItemName: "Item2", Amount: 3, CostPerUnit: 20.0},
			},
			expected: 110.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateTotal(tt.items)
			if result != tt.expected {
				t.Errorf("calculateTotal() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRenderItemsTable(t *testing.T) {
	items := []models.InvoiceItem{
		{ItemName: "Test Item", Amount: 2, CostPerUnit: 50.0},
	}

	rendered := renderItemsTable(items)

	if rendered == "" {
		t.Error("renderItemsTable should return non-empty string")
	}

	if !strings.Contains(rendered, "Test Item") {
		t.Error("Rendered table should contain item name")
	}
}

func TestRenderItemsTableEmpty(t *testing.T) {
	items := []models.InvoiceItem{}

	rendered := renderItemsTable(items)

	if !strings.Contains(rendered, "No items") {
		t.Error("Empty items should render 'No items' message")
	}
}

func TestRenderEntity(t *testing.T) {
	address := "123 Main St"
	email := "test@example.com"
	phone := "555-1234"

	entity := &models.Entity{
		ID:      "e1",
		Name:    "Test Entity",
		Address: &address,
		Email:   &email,
		Phone:   &phone,
	}

	rendered := renderEntity(entity)

	if !strings.Contains(rendered, "Test Entity") {
		t.Error("Rendered entity should contain name")
	}

	if !strings.Contains(rendered, "123 Main St") {
		t.Error("Rendered entity should contain address")
	}

	if !strings.Contains(rendered, "test@example.com") {
		t.Error("Rendered entity should contain email")
	}

	if !strings.Contains(rendered, "555-1234") {
		t.Error("Rendered entity should contain phone")
	}
}

func TestRenderEntityMinimal(t *testing.T) {
	entity := &models.Entity{
		ID:   "e1",
		Name: "Minimal Entity",
	}

	rendered := renderEntity(entity)

	if !strings.Contains(rendered, "Minimal Entity") {
		t.Error("Rendered entity should contain name")
	}

	// Should not contain field labels for missing optional fields
	if strings.Count(rendered, "Address:") > 0 {
		t.Error("Rendered entity should not contain address label when address is nil")
	}
}

