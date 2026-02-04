package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newActionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "action",
		Short: "Manage Zabbix actions",
	}

	cmd.AddCommand(newActionListCmd())

	return cmd
}

func newActionListCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Zabbix actions",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": []string{"actionid", "name", "eventsource", "status"},
				"limit":  limit,
			}

			result, err := client.Call("action.get", params)
			handleError(err)

			var actions []map[string]interface{}
			json.Unmarshal(result, &actions)

			headers := []string{"ID", "Name", "Source", "Status"}
			var rows [][]string
			for _, a := range actions {
				source := getEventSourceName(a["eventsource"].(string))
				status := "Enabled"
				if a["status"].(string) == "1" {
					status = "Disabled"
				}
				rows = append(rows, []string{
					fmt.Sprintf("%v", a["actionid"]),
					fmt.Sprintf("%v", a["name"]),
					source,
					status,
				})
			}

			outputResult(cmd, actions, headers, rows)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of actions")

	return cmd
}

func getEventSourceName(s string) string {
	switch s {
	case "0":
		return "Triggers"
	case "1":
		return "Discovery"
	case "2":
		return "Autoregistration"
	case "3":
		return "Internal"
	case "4":
		return "Service"
	default:
		return "Unknown"
	}
}
