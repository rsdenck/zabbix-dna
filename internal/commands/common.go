package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"zabbix-dna/internal/api"
	"zabbix-dna/internal/config"
	"zabbix-dna/internal/output"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)
)

type Result struct {
	ReturnCode string      `json:"return_code"`
	Errors     []string    `json:"errors"`
	Message    string      `json:"message"`
	Result     interface{} `json:"result"`
}

func outputResult(cmd *cobra.Command, data interface{}, headers []string, rows [][]string) {
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, _ := config.LoadConfig(cfgPath)

	format := "table"
	if cfg != nil && cfg.App.Output.Format != "" {
		format = cfg.App.Output.Format
	}

	// Always wrap in Result for JSON
	res := Result{
		ReturnCode: "Done",
		Errors:     []string{},
		Message:    "",
		Result:     data,
	}

	// If it's a message-only result (like "Created host...")
	if msg, ok := data.(string); ok && headers == nil && rows == nil {
		res.Message = msg
		res.Result = nil
	}

	if format == "json" {
		if res.Result == nil && rows != nil {
			res.Result = rows
		}
		jsonData, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(jsonData))
		return
	}

	// Default to Table
	if res.Message != "" {
		fmt.Println(res.Message)
		return
	}

	if len(rows) == 0 {
		fmt.Println("No results found.")
		return
	}

	renderer := output.NewTableRenderer(headers, rows)
	renderer.Render()
}

func getZabbixClient(cmd *cobra.Command) (*api.ZabbixClient, error) {
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := api.NewClient(cfg.API.URL, cfg.API.AuthToken, cfg.API.Timeout)
	if cfg.API.AuthToken == "" && cfg.API.Username != "" {
		err := client.Login(cfg.API.Username, cfg.API.Password)
		if err != nil {
			return nil, fmt.Errorf("login failed: %w", err)
		}
	} else if cfg.API.AuthToken == "" && cfg.API.Username == "" {
		return nil, fmt.Errorf("no authentication provided (token or username/password)")
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

func getHostGroupsIDs(client *api.ZabbixClient, names []string) []string {
	params := map[string]interface{}{
		"filter": map[string]interface{}{
			"name": names,
		},
	}
	result, err := client.Call("hostgroup.get", params)
	if err != nil {
		return nil
	}
	var groups []map[string]interface{}
	json.Unmarshal(result, &groups)
	var ids []string
	for _, g := range groups {
		ids = append(ids, g["groupid"].(string))
	}
	return ids
}

func getHostsIDs(client *api.ZabbixClient, names []string) []string {
	params := map[string]interface{}{
		"filter": map[string]interface{}{
			"host": names,
		},
	}
	result, err := client.Call("host.get", params)
	if err != nil {
		return nil
	}
	var hosts []map[string]interface{}
	json.Unmarshal(result, &hosts)
	var ids []string
	for _, h := range hosts {
		ids = append(ids, h["hostid"].(string))
	}
	return ids
}

func getPriorityName(p string) string {
	switch p {
	case "0":
		return "Not classified"
	case "1":
		return "Information"
	case "2":
		return "Warning"
	case "3":
		return "Average"
	case "4":
		return "High"
	case "5":
		return "Disaster"
	default:
		return "Unknown"
	}
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

func getEventForTrigger(client *api.ZabbixClient, triggerID string) string {
	params := map[string]interface{}{
		"output":    []string{"eventid"},
		"objectids": []string{triggerID},
		"sortfield": "clock",
		"sortorder": "DESC",
		"limit":     1,
	}
	result, err := client.Call("event.get", params)
	if err != nil {
		return ""
	}
	var events []map[string]interface{}
	json.Unmarshal(result, &events)
	if len(events) > 0 {
		return events[0]["eventid"].(string)
	}
	return ""
}
