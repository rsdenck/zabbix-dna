package output

import (
	"os"

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

	// Strictly aligned columns and fixed headers
	table.SetHeader(r.Headers)
	table.SetAutoWrapText(true) // Enable wrap for better responsiveness
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetReflowDuringAutoWrap(true) // Ensure it wraps correctly

	// Visual consistency according to reference
	// We use the "rounded" style or a very clean ASCII style
	table.SetCenterSeparator("┼")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")

	table.SetHeaderLine(true)
	table.SetBorder(true)
	table.SetTablePadding(" ") // Clean padding
	table.SetNoWhiteSpace(false)

	// Responsiveness: auto-resize based on terminal width
	// tablewriter handles this by default with SetAutoWrapText(true)

	// Refined styling for NOC/SOC/SRE environments
	// White on Black/Transparent with clean borders
	table.SetHeaderColor(getHeaderColors(len(r.Headers))...)

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
