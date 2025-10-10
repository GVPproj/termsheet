package views

import (
	"fmt"
	"strings"

	"github.com/GVPproj/termsheet/models"
	"github.com/GVPproj/termsheet/utils"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Table styles
	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#61AFEF")).
				Padding(0, 1)

	tableCellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98C379")).
			Padding(0, 1)

	tableRowStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle)

	sectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#D33682")).
				MarginTop(1).
				MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#61AFEF"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98C379"))
)

// RenderInvoiceView renders a read-only view of an invoice
func RenderInvoiceView(data *models.InvoiceData) string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render(fmt.Sprintf("Invoice #%d", data.InvoiceID)))
	b.WriteString("\n\n")

	// Date and status
	dateStr := data.DateCreated.Format("2006-01-02")
	status := "Unpaid"
	if data.Paid {
		status = "Paid"
	}
	b.WriteString(fmt.Sprintf("%s %s  |  %s %s\n\n",
		labelStyle.Render("Date:"),
		valueStyle.Render(dateStr),
		labelStyle.Render("Status:"),
		valueStyle.Render(status),
	))

	// Provider section
	b.WriteString(sectionTitleStyle.Render("Provider"))
	b.WriteString("\n")
	b.WriteString(renderEntity(&data.Provider))
	b.WriteString("\n\n")

	// Client section
	b.WriteString(sectionTitleStyle.Render("Client"))
	b.WriteString("\n")
	b.WriteString(renderEntity(&data.Client))
	b.WriteString("\n\n")

	// Items section
	b.WriteString(sectionTitleStyle.Render("Items"))
	b.WriteString("\n")
	b.WriteString(renderItemsTable(data.Items))
	b.WriteString("\n\n")

	// Total
	total := calculateTotal(data.Items)
	b.WriteString(labelStyle.Render("Total: "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("$%.2f", total)))

	// Help text
	b.WriteString(helpStyle.Render("\n\n\nESC to return"))

	// Wrap in container
	return containerStyle.Render(b.String())
}

// renderEntity renders provider or client information
func renderEntity(entity *models.Entity) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s %s\n",
		labelStyle.Render("Name:"),
		valueStyle.Render(entity.Name),
	))

	if entity.Address != nil && *entity.Address != "" {
		b.WriteString(fmt.Sprintf("%s %s\n",
			labelStyle.Render("Address:"),
			valueStyle.Render(*entity.Address),
		))
	}

	if entity.Email != nil && *entity.Email != "" {
		b.WriteString(fmt.Sprintf("%s %s\n",
			labelStyle.Render("Email:"),
			valueStyle.Render(*entity.Email),
		))
	}

	if entity.Phone != nil && *entity.Phone != "" {
		b.WriteString(fmt.Sprintf("%s %s\n",
			labelStyle.Render("Phone:"),
			valueStyle.Render(*entity.Phone),
		))
	}

	return b.String()
}

// renderItemsTable renders the invoice items in a table format
func renderItemsTable(items []models.InvoiceItem) string {
	if len(items) == 0 {
		return valueStyle.Render("No items")
	}

	var b strings.Builder

	// Header row
	headerRow := lipgloss.JoinHorizontal(
		lipgloss.Left,
		tableHeaderStyle.Width(20).Render("Item"),
		tableHeaderStyle.Width(12).Render("Quantity"),
		tableHeaderStyle.Width(15).Render("Cost/Unit"),
		tableHeaderStyle.Width(15).Render("Total"),
	)
	b.WriteString(headerRow)
	b.WriteString("\n")

	// Data rows
	for _, item := range items {
		itemTotal := item.Amount * item.CostPerUnit
		row := lipgloss.JoinHorizontal(
			lipgloss.Left,
			tableCellStyle.Width(20).Render(utils.TruncateText(item.ItemName, 18)),
			tableCellStyle.Width(12).Render(fmt.Sprintf("%.2f", item.Amount)),
			tableCellStyle.Width(15).Render(fmt.Sprintf("$%.2f", item.CostPerUnit)),
			tableCellStyle.Width(15).Render(fmt.Sprintf("$%.2f", itemTotal)),
		)
		b.WriteString(row)
		b.WriteString("\n")
	}

	return b.String()
}

// calculateTotal calculates the total cost of all items
func calculateTotal(items []models.InvoiceItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Amount * item.CostPerUnit
	}
	return total
}
