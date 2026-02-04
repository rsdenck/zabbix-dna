package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newScriptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "script",
		Short: "Manage Zabbix scripts",
	}

	cmd.AddCommand(newScriptListCmd())
	cmd.AddCommand(newScriptExecuteCmd())

	return cmd
}

func newScriptListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Zabbix scripts",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": []string{"scriptid", "name", "command"},
			}

			result, err := client.Call("script.get", params)
			handleError(err)

			var scripts []map[string]interface{}
			json.Unmarshal(result, &scripts)

			headers := []string{"ID", "Name", "Command"}
			var rows [][]string
			for _, s := range scripts {
				rows = append(rows, []string{
					fmt.Sprintf("%v", s["scriptid"]),
					fmt.Sprintf("%v", s["name"]),
					fmt.Sprintf("%v", s["command"]),
				})
			}

			outputResult(cmd, scripts, headers, rows)
		},
	}

	return cmd
}

func newScriptExecuteCmd() *cobra.Command {
	var hostID string

	cmd := &cobra.Command{
		Use:   "execute [script id]",
		Short: "Execute a Zabbix script",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"scriptid": args[0],
				"hostid":   hostID,
			}

			result, err := client.Call("script.execute", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)

			if value, ok := resp["value"].(string); ok {
				headers := []string{"Script execution result"}
				rows := [][]string{{value}}
				outputResult(cmd, resp, headers, rows)
			}
		},
	}

	cmd.Flags().StringVarP(&hostID, "hostid", "H", "", "Host ID to execute script on")
	cmd.MarkFlagRequired("hostid")

	return cmd
}
