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
			fmt.Println("Checking environment for Zabbix MCP Server...")

			uvPath, err := exec.LookPath("uv")
			if err != nil {
				fmt.Println("Error: 'uv' not found in PATH. Please install it from https://github.com/astral-sh/uv")
				return
			}

			fmt.Printf("Found 'uv' at: %s\n", uvPath)
			fmt.Printf("Preparing to launch Zabbix MCP Server on port %d (Read-Only: %v)...\n", port, readOnly)
			fmt.Println("\nExecution command:")
			fmt.Printf("uv run python -m zabbix_mcp_server --port %d\n", port)

			fmt.Println("\nNote: Ensure you have the zabbix-mcp-server package installed or available in your uv project.")
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
			fmt.Println("Checking environment for Grafana MCP Server...")

			npxPath, err := exec.LookPath("npx")
			if err != nil {
				fmt.Println("Error: 'npx' (Node.js) not found in PATH. Please install Node.js from https://nodejs.org/")
				return
			}

			fmt.Printf("Found 'npx' at: %s\n", npxPath)
			fmt.Printf("Preparing to launch Grafana MCP Server for host %s...\n", host)
			fmt.Println("\nExecution command:")
			fmt.Printf("npx @grafana/mcp-grafana --host %s\n", host)
		},
	}

	cmd.Flags().StringVar(&host, "host", "http://localhost:3000", "Grafana host URL")

	return cmd
}
