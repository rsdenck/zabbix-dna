package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newTriggerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trigger",
		Short: "Manage Zabbix triggers",
	}

	cmd.AddCommand(newTriggerListCmd())
	cmd.AddCommand(newTriggerCreateCmd())

	return cmd
}

func newTriggerListCmd() *cobra.Command {
	var hostName string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List triggers for a host",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output":    []string{"triggerid", "description", "priority", "value", "status"},
				"limit":     limit,
				"sortfield": "priority",
				"sortorder": "DESC",
			}

			if hostName != "" {
				params["host"] = hostName
			}

			result, err := client.Call("trigger.get", params)
			handleError(err)

			var triggers []map[string]interface{}
			json.Unmarshal(result, &triggers)

			headers := []string{"TriggerID", "Description", "Priority", "Status"}
			var rows [][]string
			for _, t := range triggers {
				priority := getPriorityName(t["priority"].(string))
				status := "OK"
				if t["value"].(string) == "1" {
					status = "PROBLEM"
				}
				if t["status"].(string) == "1" {
					status = "DISABLED"
				}
				rows = append(rows, []string{
					fmt.Sprintf("%v", t["triggerid"]),
					fmt.Sprintf("%v", t["description"]),
					priority,
					status,
				})
			}

			outputResult(cmd, triggers, headers, rows)
		},
	}

	cmd.Flags().StringVarP(&hostName, "host", "H", "", "Host name to list triggers for")
	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of triggers")

	return cmd
}

func newTriggerCreateCmd() *cobra.Command {
	var expression string
	var priority int

	cmd := &cobra.Command{
		Use:   "create [trigger description]",
		Short: "Create a new Zabbix trigger",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"description": args[0],
				"expression":  expression,
				"priority":    priority,
			}

			result, err := client.Call("trigger.create", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)
			triggerIDs := resp["triggerids"].([]interface{})

			fmt.Printf("Trigger created successfully with ID: %s\n", triggerIDs[0])
		},
	}

	cmd.Flags().StringVarP(&expression, "expression", "e", "", "Trigger expression")
	cmd.Flags().IntVarP(&priority, "priority", "p", 0, "Trigger priority (0-5)")
	cmd.MarkFlagRequired("expression")

	return cmd
}

func getPriorityName(p string) string {
	switch p {
	case "0":
		return "Not classified"
	case "1":
		return "Information"
	case "2":
		return "Warning"
	case "3":
		return "Average"
	case "4":
		return "High"
	case "5":
		return "Disaster"
	default:
		return "Unknown"
	}
}
