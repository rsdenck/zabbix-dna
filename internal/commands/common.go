package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"zabbix-dna/internal/api"
	"zabbix-dna/internal/config"

	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#D20000")).
			Padding(0, 1).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D20000"))
)

type JSONResponse struct {
	Message    string      `json:"message"`
	Errors     []string    `json:"errors"`
	ReturnCode string      `json:"return_code"`
	Result     interface{} `json:"result"`
}

func outputResult(cmd *cobra.Command, result interface{}, headers []string, rows [][]string) {
	format, _ := cmd.Flags().GetString("format")
	if format == "json" {
		resp := JSONResponse{
			Message:    "",
			Errors:     []string{},
			ReturnCode: "Done",
			Result:     result,
		}
		data, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(data))
		return
	}

	// Default to table format
	table := tablewriter.NewWriter(os.Stdout)

	// Custom table style to match images
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Image-like styling
	table.SetCenterSeparator("┼")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")

	table.SetHeaderLine(true)
	table.SetBorder(true)
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(false)

	// Apply colors to header if terminal supports it
	table.SetHeaderColor(getHeaderColors(len(headers))...)

	table.AppendBulk(rows)
	table.Render()
}

func getHeaderColors(count int) []tablewriter.Colors {
	colors := make([]tablewriter.Colors, count)
	for i := range colors {
		colors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor, tablewriter.BgRedColor}
	}
	return colors
}

func getZabbixClient(cmd *cobra.Command) (*api.ZabbixClient, error) {
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := api.NewClient(cfg.Zabbix.URL, cfg.Zabbix.Token, cfg.Zabbix.Timeout)
	if cfg.Zabbix.Token == "" && cfg.Zabbix.User != "" {
		err := client.Login(cfg.Zabbix.User, cfg.Zabbix.Password)
		if err != nil {
			return nil, fmt.Errorf("login failed: %w", err)
		}
	} else if cfg.Zabbix.Token == "" && cfg.Zabbix.User == "" {
		return nil, fmt.Errorf("no authentication provided (token or user/password)")
	}

	return client, nil
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getHostID(client *api.ZabbixClient, name string) string {
	params := map[string]interface{}{
		"filter": map[string]interface{}{
			"host": name,
		},
	}
	result, err := client.Call("host.get", params)
	if err != nil {
		return ""
	}
	var hosts []map[string]interface{}
	json.Unmarshal(result, &hosts)
	if len(hosts) > 0 {
		return hosts[0]["hostid"].(string)
	}
	return ""
}

func getTemplateID(client *api.ZabbixClient, name string) string {
	params := map[string]interface{}{
		"filter": map[string]interface{}{
			"host": name,
		},
	}
	result, err := client.Call("template.get", params)
	if err != nil {
		return ""
	}
	var templates []map[string]interface{}
	json.Unmarshal(result, &templates)
	if len(templates) > 0 {
		return templates[0]["templateid"].(string)
	}
	return ""
}


