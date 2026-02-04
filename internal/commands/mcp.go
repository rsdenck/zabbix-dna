package commands

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func newMCPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Model Context Protocol (MCP) integrations",
		Long:  `Manage and run MCP servers for Zabbix and Grafana integration with AI assistants.`,
	}

	cmd.AddCommand(newMCPZabbixCmd())
	cmd.AddCommand(newMCPGrafanaCmd())

	return cmd
}

func newMCPZabbixCmd() *cobra.Command {
	var readOnly bool
	var port int

	cmd := &cobra.Command{
		Use:   "zabbix",
		Short: "Run Zabbix MCP Server",
		Long:  `Launch the Zabbix MCP server to connect AI assistants to Zabbix monitoring.`,
		Run: func(cmd *cobra.Command, args []string) {
			uvPath, err := exec.LookPath("uv")
			if err != nil {
				handleError(fmt.Errorf("'uv' not found in PATH. Please install it from https://github.com/astral-sh/uv"))
				return
			}

			headers := []string{"Environment Check", "Result"}
			rows := [][]string{
				{"'uv' Path", uvPath},
				{"Launch Command", fmt.Sprintf("uv run python -m zabbix_mcp_server --port %d", port)},
				{"Port", fmt.Sprintf("%d", port)},
				{"Read-Only", fmt.Sprintf("%v", readOnly)},
				{"Status", "Ready to Launch"},
			}
			outputResult(cmd, nil, headers, rows)
		},
	}

	cmd.Flags().BoolVar(&readOnly, "read-only", true, "Enable read-only mode")
	cmd.Flags().IntVarP(&port, "port", "p", 8000, "Port for HTTP transport")

	return cmd
}

func newMCPGrafanaCmd() *cobra.Command {
	var host string

	cmd := &cobra.Command{
		Use:   "grafana",
		Short: "Run Grafana MCP Server",
		Long:  `Launch the Grafana MCP server to connect AI assistants to Grafana dashboards and metrics.`,
		Run: func(cmd *cobra.Command, args []string) {
			npxPath, err := exec.LookPath("npx")
			if err != nil {
				handleError(fmt.Errorf("'npx' (Node.js) not found in PATH. Please install Node.js from https://nodejs.org/"))
				return
			}

			headers := []string{"Environment Check", "Result"}
			rows := [][]string{
				{"'npx' Path", npxPath},
				{"Launch Command", fmt.Sprintf("npx @grafana/mcp-grafana --host %s", host)},
				{"Host URL", host},
				{"Status", "Ready to Launch"},
			}
			outputResult(cmd, nil, headers, rows)
		},
	}

	cmd.Flags().StringVar(&host, "host", "http://localhost:3000", "Grafana host URL")

	return cmd
}
