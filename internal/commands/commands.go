package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func AddCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newBackupCmd())
	rootCmd.AddCommand(newExportCmd())
	rootCmd.AddCommand(newExporterCmd())
	rootCmd.AddCommand(newTestAPICmd())
	rootCmd.AddCommand(newHostCmd())
	rootCmd.AddCommand(newHostGroupCmd())
	rootCmd.AddCommand(newTemplateCmd())
	rootCmd.AddCommand(newTemplateGroupCmd())
	rootCmd.AddCommand(newProxyCmd())
	rootCmd.AddCommand(newSaltCmd()) // Add SaltStack command
	rootCmd.AddCommand(newUserCmd())
	rootCmd.AddCommand(newUserGroupCmd())
	rootCmd.AddCommand(newItemCmd())
	rootCmd.AddCommand(newTriggerCmd())
	rootCmd.AddCommand(newProblemCmd())
	rootCmd.AddCommand(newMaintenanceCmd())
	rootCmd.AddCommand(newMacroCmd())
	rootCmd.AddCommand(newHostInterfaceCmd())
	rootCmd.AddCommand(newMediaCmd())
	rootCmd.AddCommand(newActionCmd())
	rootCmd.AddCommand(newScriptCmd())
	rootCmd.AddCommand(newREPLCmd(rootCmd))
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of Zabbix-DNA",
		Run: func(cmd *cobra.Command, args []string) {
			// Version will be injected at build time in main.go
			// but we can also show some static info here
			fmt.Println("ZABBIX-DNA CLI | v1.0.6")
			fmt.Println("Engine: Go 1.24.2")
			fmt.Println("Zabbix Compatibility: 6.4, 7.0, 7.2, 8.0")
			fmt.Println("Features: SaltStack Integration, OTLP, TUI")
		},
	}
}
