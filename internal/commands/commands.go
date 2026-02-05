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
	hostCmd := newHostCmd()
	rootCmd.AddCommand(hostCmd)
	// Aliases para compatibilidade com zabbix-cli
	rootCmd.AddCommand(aliasCmd(hostCmd, "create_host", "create"))
	rootCmd.AddCommand(aliasCmd(hostCmd, "remove_host", "delete"))
	rootCmd.AddCommand(aliasCmd(hostCmd, "show_host", "show"))
	rootCmd.AddCommand(aliasCmd(hostCmd, "show_hosts", "list"))

	// HOST GROUP - GRUPOS DE HOSTS
	hostGroupCmd := newHostGroupCmd()
	rootCmd.AddCommand(hostGroupCmd)
	rootCmd.AddCommand(aliasCmd(hostGroupCmd, "create_hostgroup", "create"))
	rootCmd.AddCommand(aliasCmd(hostGroupCmd, "remove_hostgroup", "delete"))
	rootCmd.AddCommand(aliasCmd(hostGroupCmd, "show_hostgroup", "show"))
	rootCmd.AddCommand(aliasCmd(hostGroupCmd, "show_hostgroups", "list"))

	// TEMPLATE - TEMPLATES
	templateCmd := newTemplateCmd()
	rootCmd.AddCommand(templateCmd)
	rootCmd.AddCommand(aliasCmd(templateCmd, "create_template", "create"))
	rootCmd.AddCommand(aliasCmd(templateCmd, "remove_template", "delete"))
	rootCmd.AddCommand(aliasCmd(templateCmd, "show_template", "show"))
	rootCmd.AddCommand(aliasCmd(templateCmd, "show_templates", "list"))

	// TEMPLATE GROUP - GRUPOS DE TEMPLATE
	templateGroupCmd := newTemplateGroupCmd()
	rootCmd.AddCommand(templateGroupCmd)
	rootCmd.AddCommand(aliasCmd(templateGroupCmd, "create_templategroup", "create"))
	rootCmd.AddCommand(aliasCmd(templateGroupCmd, "remove_templategroup", "delete"))
	rootCmd.AddCommand(aliasCmd(templateGroupCmd, "show_templategroup", "show"))
	rootCmd.AddCommand(aliasCmd(templateGroupCmd, "show_templategroups", "list"))

	// PROXY - PROXIES ZABBIX
	proxyCmd := newProxyCmd()
	rootCmd.AddCommand(proxyCmd)
	rootCmd.AddCommand(aliasCmd(proxyCmd, "show_proxy", "show"))
	rootCmd.AddCommand(aliasCmd(proxyCmd, "show_proxies", "list"))

	// SALTSTACK
	rootCmd.AddCommand(newSaltCmd())

	// USER
	userCmd := newUserCmd()
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(aliasCmd(userCmd, "create_user", "create"))
	rootCmd.AddCommand(aliasCmd(userCmd, "remove_user", "delete"))
	rootCmd.AddCommand(aliasCmd(userCmd, "show_user", "show"))
	rootCmd.AddCommand(aliasCmd(userCmd, "show_users", "list"))

	// USER GROUP
	userGroupCmd := newUserGroupCmd()
	rootCmd.AddCommand(userGroupCmd)
	rootCmd.AddCommand(aliasCmd(userGroupCmd, "create_usergroup", "create"))
	rootCmd.AddCommand(aliasCmd(userGroupCmd, "remove_usergroup", "delete"))
	rootCmd.AddCommand(aliasCmd(userGroupCmd, "show_usergroup", "show"))
	rootCmd.AddCommand(aliasCmd(userGroupCmd, "show_usergroups", "list"))

	// MONITORING
	rootCmd.AddCommand(newMonitoringCmd())

	itemCmd := newItemCmd()
	rootCmd.AddCommand(itemCmd)
	rootCmd.AddCommand(aliasCmd(itemCmd, "show_items", "list"))

	triggerCmd := newTriggerCmd()
	rootCmd.AddCommand(triggerCmd)
	rootCmd.AddCommand(aliasCmd(triggerCmd, "show_triggers", "list"))

	problemCmd := newProblemCmd()
	rootCmd.AddCommand(problemCmd)
	rootCmd.AddCommand(aliasCmd(problemCmd, "show_problems", "list"))
	rootCmd.AddCommand(aliasCmd(problemCmd, "acknowledge_event", "acknowledge"))
	rootCmd.AddCommand(aliasCmd(problemCmd, "acknowledge_events", "acknowledge"))
	rootCmd.AddCommand(aliasCmd(problemCmd, "acknowledge_trigger_last_event", "acknowledge_trigger"))
	rootCmd.AddCommand(aliasCmd(problemCmd, "show_trigger_events", "events"))
	rootCmd.AddCommand(aliasCmd(problemCmd, "show_alarms", "alarms"))

	maintenanceCmd := newMaintenanceCmd()
	rootCmd.AddCommand(maintenanceCmd)
	rootCmd.AddCommand(aliasCmd(maintenanceCmd, "show_maintenance", "show"))
	rootCmd.AddCommand(aliasCmd(maintenanceCmd, "show_maintenances", "list"))
	rootCmd.AddCommand(aliasCmd(maintenanceCmd, "create_maintenance_definition", "create"))

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

func aliasCmd(parent *cobra.Command, use string, subCommandName string) *cobra.Command {
	child, _, _ := parent.Find([]string{subCommandName})
	if child == nil {
		return &cobra.Command{Use: use, Hidden: true}
	}

	alias := &cobra.Command{
		Use:                use,
		Short:              child.Short,
		Long:               child.Long,
		Example:            child.Example,
		Args:               child.Args,
		Hidden:             true,
		DisableFlagParsing: false,
		Run:                child.Run,
	}

	alias.Flags().AddFlagSet(child.Flags())
	return alias
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
