package views

import (
	"testing"
)

func TestCreateInvoiceActionForm(t *testing.T) {
	var selection string
	form := CreateInvoiceActionForm(&selection)

	if form == nil {
		t.Fatal("CreateInvoiceActionForm returned nil")
	}

	// Verify form is properly created by checking it's not nil
	// We can't access internal fields, but we can verify it exists
	if form.View() == "" {
		t.Error("Form should render something")
	}
}

func TestRenderInvoiceActionMenu(t *testing.T) {
	var selection string
	form := CreateInvoiceActionForm(&selection)

	rendered := RenderInvoiceActionMenu(form)

	if rendered == "" {
		t.Error("RenderInvoiceActionMenu should return non-empty string")
	}

	// Check for expected content
	if !contains(rendered, "Invoice Actions") {
		t.Error("Rendered output should contain title 'Invoice Actions'")
	}

	if !contains(rendered, "ESC to return to invoice list") {
		t.Error("Rendered output should contain help text")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
