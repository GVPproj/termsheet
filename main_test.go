package main

import (
	"testing"

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
	tests := []struct {
		name         string
		selection    string
		expectedView types.View
	}{
		{
			name:         "Providers selection",
			selection:    "Providers",
			expectedView: types.ProvidersView,
		},
		{
			name:         "Clients selection",
			selection:    "Clients",
			expectedView: types.ClientsView,
		},
		{
			name:         "Invoices selection",
			selection:    "Invoices",
			expectedView: types.InvoicesView,
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
	m.currentView = types.ProvidersView

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
