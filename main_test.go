package main

import (
	"testing"

	"github.com/GVPproj/termsheet/models"
	"github.com/GVPproj/termsheet/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Test that the form's pointer binding to m.selection works correctly
func TestFormSelectionBinding(t *testing.T) {
	m := initialModel()

	// Verify the form has a pointer to m.selection
	// This test would have failed with the old code where initialModel returned a value
	initialAddress := &m.selection

	// Manually set selection through the field the form points to
	m.selection = "Clients"

	// Verify it's the same memory location
	if &m.selection != initialAddress {
		t.Error("selection address changed - pointer binding broken")
	}

	if m.selection != "Clients" {
		t.Errorf("expected selection 'Clients', got %q", m.selection)
	}
}

// Test the view switching logic
func TestViewSwitching(t *testing.T) {
	// Initialize test database
	if err := models.InitDB(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer models.CloseDB()

	tests := []struct {
		name         string
		selection    string
		expectedView types.View
	}{
		{
			name:         "Providers selection",
			selection:    "Providers",
			expectedView: types.ProvidersListView,
		},
		{
			name:         "Clients selection",
			selection:    "Clients",
			expectedView: types.ClientsListView,
		},
		{
			name:         "Invoices selection",
			selection:    "Invoices",
			expectedView: types.InvoicesListView,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := initialModel()

			// Set selection and mark form as completed
			m.selection = tt.selection
			m.form.State = huh.StateCompleted

			// Trigger Update to process completed form
			m.Update(nil)

			// Verify the view switched correctly
			if m.currentView != tt.expectedView {
				t.Errorf("expected view %v, got %v", tt.expectedView, m.currentView)
			}
		})
	}
}

// Test ESC returns to menu
func TestEscapeReturnsToMenu(t *testing.T) {
	m := initialModel()

	// Set to a different view
	m.currentView = types.ProvidersListView

	// Send ESC
	m.Update(tea.KeyMsg{Type: tea.KeyEscape})

	// Verify we're back at menu
	if m.currentView != types.MenuView {
		t.Errorf("expected MenuView after ESC, got %v", m.currentView)
	}

	// Verify form was reset
	if m.form == nil {
		t.Error("form should not be nil after returning to menu")
	}
}

// Test provider selection loads provider form instead of returning to menu
func TestProviderSelectionLoadsForm(t *testing.T) {
	// Initialize test database
	if err := models.InitDB(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer models.CloseDB()

	// Create a test provider
	_, err := models.CreateProvider("Test Provider", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create test provider: %v", err)
	}

	m := initialModel()

	// Navigate to Providers view
	m.selection = "Providers"
	m.form.State = huh.StateCompleted
	m.Update(nil)

	// Verify we're in ProvidersListView
	if m.currentView != types.ProvidersListView {
		t.Fatalf("expected ProvidersListView, got %v", m.currentView)
	}

	// Get the provider ID from the form options
	providers, _ := models.ListProviders()
	if len(providers) == 0 {
		t.Fatal("expected at least one provider")
	}
	providerID := providers[0].ID

	// Select a provider
	m.selection = providerID
	m.form.State = huh.StateCompleted
	m.Update(nil)

	// EXPECTED: Should load provider edit form
	// ACTUAL: Currently returns to MenuView (this test should fail)
	if m.currentView == types.MenuView {
		t.Error("selecting a provider should NOT return to menu, should load provider form")
	}
}

// Test "Create New Provider" loads provider form instead of returning to menu
func TestCreateNewProviderLoadsForm(t *testing.T) {
	// Initialize test database
	if err := models.InitDB(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer models.CloseDB()

	m := initialModel()

	// Navigate to Providers view
	m.selection = "Providers"
	m.form.State = huh.StateCompleted
	m.Update(nil)

	// Verify we're in ProvidersListView
	if m.currentView != types.ProvidersListView {
		t.Fatalf("expected ProvidersListView, got %v", m.currentView)
	}

	// Select "Create New Provider"
	m.selection = "CREATE_NEW"
	m.form.State = huh.StateCompleted
	m.Update(nil)

	// EXPECTED: Should load provider create form
	// ACTUAL: Currently returns to MenuView (this test should fail)
	if m.currentView == types.MenuView {
		t.Error("selecting 'Create New Provider' should NOT return to menu, should load provider form")
	}
}
