package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newMacroCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "macro",
		Short: "Manage Zabbix user macros",
	}

	cmd.AddCommand(newMacroListCmd())

	return cmd
}

func newMacroListCmd() *cobra.Command {
	var hostName string
	var templateName string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List macros for a host or template",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": "extend",
			}

			if hostName != "" {
				params["hostids"] = []string{getHostID(client, hostName)}
			} else if templateName != "" {
				params["templateids"] = []string{getTemplateID(client, templateName)}
			} else {
				params["globalmacro"] = true
			}

			method := "usermacro.get"
			result, err := client.Call(method, params)
			handleError(err)

			var macros []map[string]interface{}
			json.Unmarshal(result, &macros)

			fmt.Printf("%-30s %-50s\n", "Macro", "Value")
			for _, m := range macros {
				fmt.Printf("%-30s %-50s\n", m["macro"], m["value"])
			}
		},
	}

	cmd.Flags().StringVarP(&hostName, "host", "H", "", "Host name to list macros for")
	cmd.Flags().StringVarP(&templateName, "template", "T", "", "Template name to list macros for")

	return cmd
}
