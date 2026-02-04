package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newMonitoringCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitoring",
		Short: "Host monitoring commands",
	}

	cmd.AddCommand(newShowItemsCmd())
	cmd.AddCommand(newShowItemCmd())
	cmd.AddCommand(newShowLastValuesCmd())
	cmd.AddCommand(newShowTriggersCmd())
	cmd.AddCommand(newShowEventsCmd())
	cmd.AddCommand(newShowGraphsCmd())

	return cmd
}

func newShowItemsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "items [host name]",
		Aliases: []string{"show_items"},
		Short:   "Show items for a host",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"host":   args[0],
				"output": []string{"itemid", "name", "key_", "lastvalue", "units"},
			}

			result, err := client.Call("item.get", params)
			handleError(err)

			var items []map[string]interface{}
			json.Unmarshal(result, &items)

			headers := []string{"ItemID", "Name", "Key", "Last Value"}
			var rows [][]string
			for _, i := range items {
				val := fmt.Sprintf("%v", i["lastvalue"])
				if units, ok := i["units"].(string); ok && units != "" {
					val += " " + units
				}
				rows = append(rows, []string{
					fmt.Sprintf("%v", i["itemid"]),
					fmt.Sprintf("%v", i["name"]),
					fmt.Sprintf("%v", i["key_"]),
					val,
				})
			}

			outputResult(cmd, items, headers, rows)
		},
	}
}

func newShowItemCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "item [item name]",
		Aliases: []string{"show_item"},
		Short:   "Show details of an item",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"search": map[string]interface{}{"name": args[0]},
				"output": "extend",
			}

			result, err := client.Call("item.get", params)
			handleError(err)

			var items []map[string]interface{}
			json.Unmarshal(result, &items)

			if len(items) == 0 {
				fmt.Println("Item not found")
				return
			}

			item := items[0]
			headers := []string{"Property", "Value"}
			var rows [][]string
			rows = append(rows, []string{"ItemID", fmt.Sprintf("%v", item["itemid"])})
			rows = append(rows, []string{"Name", fmt.Sprintf("%v", item["name"])})
			rows = append(rows, []string{"Key", fmt.Sprintf("%v", item["key_"])})
			rows = append(rows, []string{"Last Value", fmt.Sprintf("%v", item["lastvalue"])})

			outputResult(cmd, item, headers, rows)
		},
	}
}

func newShowLastValuesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "last-values [host name]",
		Aliases: []string{"show_last_values"},
		Short:   "Show last values for a host",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"host":   args[0],
				"output": []string{"itemid", "name", "key_", "lastvalue", "units"},
			}

			result, err := client.Call("item.get", params)
			handleError(err)

			var items []map[string]interface{}
			json.Unmarshal(result, &items)

			headers := []string{"ItemID", "Name", "Key", "Last Value"}
			var rows [][]string
			for _, i := range items {
				val := fmt.Sprintf("%v", i["lastvalue"])
				if units, ok := i["units"].(string); ok && units != "" {
					val += " " + units
				}
				rows = append(rows, []string{
					fmt.Sprintf("%v", i["itemid"]),
					fmt.Sprintf("%v", i["name"]),
					fmt.Sprintf("%v", i["key_"]),
					val,
				})
			}

			outputResult(cmd, items, headers, rows)
		},
	}
}

func newShowTriggersCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "triggers [host name]",
		Aliases: []string{"show_triggers"},
		Short:   "Show triggers for a host",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"host":   args[0],
				"output": []string{"triggerid", "description", "priority", "value"},
			}

			result, err := client.Call("trigger.get", params)
			handleError(err)

			var triggers []map[string]interface{}
			json.Unmarshal(result, &triggers)

			headers := []string{"TriggerID", "Description", "Priority", "Status"}
			var rows [][]string
			for _, t := range triggers {
				status := "OK"
				if t["value"].(string) == "1" {
					status = "PROBLEM"
				}
				rows = append(rows, []string{
					fmt.Sprintf("%v", t["triggerid"]),
					fmt.Sprintf("%v", t["description"]),
					getPriorityName(fmt.Sprintf("%v", t["priority"])),
					status,
				})
			}

			outputResult(cmd, triggers, headers, rows)
		},
	}
}

func newShowEventsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "events [host name]",
		Aliases: []string{"show_events"},
		Short:   "Show events for a host",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"host":   args[0],
				"output": "extend",
				"limit":  10,
			}

			result, err := client.Call("event.get", params)
			handleError(err)

			var events []map[string]interface{}
			json.Unmarshal(result, &events)

			headers := []string{"EventID", "Name", "Severity", "Clock"}
			var rows [][]string
			for _, e := range events {
				rows = append(rows, []string{
					fmt.Sprintf("%v", e["eventid"]),
					fmt.Sprintf("%v", e["name"]),
					getPriorityName(fmt.Sprintf("%v", e["severity"])),
					fmt.Sprintf("%v", e["clock"]),
				})
			}

			outputResult(cmd, events, headers, rows)
		},
	}
}

func newShowGraphsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "graphs [host name]",
		Aliases: []string{"show_graphs"},
		Short:   "Show graphs for a host",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"host":   args[0],
				"output": []string{"graphid", "name", "width", "height"},
			}

			result, err := client.Call("graph.get", params)
			handleError(err)

			var graphs []map[string]interface{}
			json.Unmarshal(result, &graphs)

			headers := []string{"GraphID", "Name", "Dimensions"}
			var rows [][]string
			for _, g := range graphs {
				rows = append(rows, []string{
					fmt.Sprintf("%v", g["graphid"]),
					fmt.Sprintf("%v", g["name"]),
					fmt.Sprintf("%sx%s", g["width"], g["height"]),
				})
			}

			outputResult(cmd, graphs, headers, rows)
		},
	}
}
