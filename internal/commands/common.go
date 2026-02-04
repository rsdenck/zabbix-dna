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

func outputResult(cmd *cobra.Command, result interface{}, headers []string, rows [][]string) {
	// Prohibited: JSON, XML, YAML, etc.
	// Only centralized TableRenderer is allowed.
	renderer := output.NewTableRenderer(headers, rows)
	renderer.Render()
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
