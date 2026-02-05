package commands

import (
	"github.com/spf13/cobra"
)

func AddCommands(rootCmd *cobra.Command) {
	// CLI - COMANDOS GERAIS DO CLIENTE
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newTestAPICmd())
	rootCmd.AddCommand(newREPLCmd(rootCmd))

	// HOST - GERENCIAMENTO DE HOSTS
	rootCmd.AddCommand(newHostCmd())

	// HOST GROUP - GRUPOS DE HOSTS
	rootCmd.AddCommand(newHostGroupCmd())

	// TEMPLATE - TEMPLATES
	rootCmd.AddCommand(newTemplateCmd())

	// TEMPLATE GROUP - GRUPOS DE TEMPLATE
	rootCmd.AddCommand(newTemplateGroupCmd())

	// PROXY - PROXIES ZABBIX
	rootCmd.AddCommand(newProxyCmd())

	// SALTSTACK
	rootCmd.AddCommand(newSaltCmd())

	// USER
	rootCmd.AddCommand(newUserCmd())

	// USER GROUP
	rootCmd.AddCommand(newUserGroupCmd())

	// MONITORING
	rootCmd.AddCommand(newMonitoringCmd())
	rootCmd.AddCommand(newItemCmd())
	rootCmd.AddCommand(newTriggerCmd())
	rootCmd.AddCommand(newProblemCmd())
	rootCmd.AddCommand(newMaintenanceCmd())

	// MACRO
	rootCmd.AddCommand(newMacroCmd())

	// MCP - Model Context Protocol
	rootCmd.AddCommand(newMCPCmd())

	// INTERFACE
	rootCmd.AddCommand(newHostInterfaceCmd())

	// MEDIA
	rootCmd.AddCommand(newMediaCmd())

	// ACTION / SCRIPT
	rootCmd.AddCommand(newActionCmd())
	rootCmd.AddCommand(newScriptCmd())

	// CONFIG / EXPORT
	rootCmd.AddCommand(newBackupCmd())
	rootCmd.AddCommand(newExportCmd())
	rootCmd.AddCommand(newExporterCmd())

	// Alias para comandos legados ou compatibilidade se necessário
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of Zabbix-DNA",
		Run: func(cmd *cobra.Command, args []string) {
			headers := []string{"Property", "Value"}
			cgoStatus := "Disabled"
			if isCGOBuilt() {
				cgoStatus = "Enabled (SaltStack Support)"
			}
			rows := [][]string{
				{"Zabbix-DNA", "v1.0.7"},
				{"Engine", "Go 1.24.2"},
				{"CGO/SaltStack", cgoStatus},
				{"Zabbix Compatibility", "6.4, 7.0, 7.2, 8.0"},
				{"Features", "SaltStack, OTLP, TUI"},
			}
			outputResult(cmd, nil, headers, rows)
		},
	}
}
