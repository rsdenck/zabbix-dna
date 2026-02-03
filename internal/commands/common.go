package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"zabbix-dna/internal/api"
	"zabbix-dna/internal/config"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
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
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	table.SetHeaderLine(true)
	table.SetBorder(true)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(false)

	table.AppendBulk(rows)
	table.Render()
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
