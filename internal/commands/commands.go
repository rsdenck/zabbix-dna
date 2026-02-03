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
			fmt.Println("Zabbix-DNA v1.0.0")
			fmt.Println("Engine: Go 1.25")
			fmt.Println("Zabbix Compatibility: 7.0, 7.2, 8.0")
		},
	}
}
