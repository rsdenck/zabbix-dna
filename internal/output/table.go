package output

import (
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/olekukonko/tablewriter"
)

// TableRenderer is the centralized renderer for all CLI outputs.
// It strictly follows the enterprise-grade visual standard.
type TableRenderer struct {
	Headers []string
	Rows    [][]string
}

// NewTableRenderer creates a new instance of the centralized renderer.
func NewTableRenderer(headers []string, rows [][]string) *TableRenderer {
	return &TableRenderer{
		Headers: headers,
		Rows:    rows,
	}
}

// Render outputs the data to stdout in a strictly tabular format.
// Any other format (JSON, YAML, etc.) is prohibited.
func (r *TableRenderer) Render() {
	table := tablewriter.NewWriter(os.Stdout)

	// Get terminal width for responsiveness
	width, _, err := term.GetSize(uintptr(os.Stdout.Fd()))
	if err == nil && width > 0 {
		table.SetColWidth(width / (len(r.Headers) + 1)) // Distribute width
	}

	// Strictly aligned columns and fixed headers
	table.SetHeader(r.Headers)
	table.SetAutoWrapText(true) // Enable wrap for better responsiveness
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetReflowDuringAutoWrap(true) // Ensure it wraps correctly

	// Visual consistency: Rounded-like clean ASCII style
	// We use Unicode box-drawing characters for a professional look
	table.SetCenterSeparator("┼")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")

	// Set header styling
	table.SetHeaderLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Borders
	table.SetBorder(true)
	table.SetTablePadding(" ") // Clean padding
	table.SetNoWhiteSpace(false)

	// Auto-wrap for responsiveness
	table.SetAutoWrapText(true)
	table.SetReflowDuringAutoWrap(true)

	table.AppendBulk(r.Rows)
	table.Render()
}

func getHeaderColors(count int) []tablewriter.Colors {
	colors := make([]tablewriter.Colors, count)
	for i := range colors {
		colors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor}
	}
	return colors
}
