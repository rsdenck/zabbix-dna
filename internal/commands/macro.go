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

	cmd.AddCommand(newMacroListCmd())   // show_host_macros -> macro list
	cmd.AddCommand(newMacroCreateCmd()) // create_global_macro -> macro create

	return cmd
}

func newMacroListCmd() *cobra.Command {
	var hostName string
	var templateName string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show_host_macros"},
		Short:   "List macros for a host or template",
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

			headers := []string{"Macro", "Value"}
			var rows [][]string
			for _, m := range macros {
				rows = append(rows, []string{
					fmt.Sprintf("%v", m["macro"]),
					fmt.Sprintf("%v", m["value"]),
				})
			}

			outputResult(cmd, macros, headers, rows)
		},
	}

	cmd.Flags().StringVarP(&hostName, "host", "H", "", "Host name to list macros for")
	cmd.Flags().StringVarP(&templateName, "template", "T", "", "Template name to list macros for")

	return cmd
}

func newMacroCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "create [macro] [value]",
		Aliases: []string{"create_global_macro"},
		Short:   "Create a new global macro",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"macro": args[0],
				"value": args[1],
			}

			result, err := client.Call("usermacro.createglobal", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)

			headers := []string{"Macro", "Value", "Action", "Status"}
			rows := [][]string{{args[0], args[1], "Create Global Macro", "Success"}}
			outputResult(cmd, resp, headers, rows)
		},
	}
}
