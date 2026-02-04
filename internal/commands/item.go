package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newItemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item",
		Short: "Manage Zabbix items",
	}

	cmd.AddCommand(newItemListCmd())
	cmd.AddCommand(newItemCreateCmd())

	return cmd
}

func newItemListCmd() *cobra.Command {
	var hostName string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List items for a host",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output":    []string{"itemid", "name", "key_", "lastvalue", "units"},
				"limit":     limit,
				"sortfield": "name",
			}

			if hostName != "" {
				params["host"] = hostName
			}

			result, err := client.Call("item.get", params)
			handleError(err)

			var items []map[string]interface{}
			json.Unmarshal(result, &items)

			headers := []string{"ItemID", "Name", "Key", "Last Value"}
			var rows [][]string
			for _, i := range items {
				lastValue := i["lastvalue"].(string)
				if units, ok := i["units"].(string); ok && units != "" {
					lastValue += " " + units
				}
				rows = append(rows, []string{
					fmt.Sprintf("%v", i["itemid"]),
					fmt.Sprintf("%v", i["name"]),
					fmt.Sprintf("%v", i["key_"]),
					lastValue,
				})
			}

			outputResult(cmd, items, headers, rows)
		},
	}

	cmd.Flags().StringVarP(&hostName, "host", "H", "", "Host name to list items for")
	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of items")

	return cmd
}

func newItemCreateCmd() *cobra.Command {
	var hostID string
	var key string
	var itemType int
	var valueType int
	var interfaceID string
	var delay string

	cmd := &cobra.Command{
		Use:   "create [item name]",
		Short: "Create a new Zabbix item",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"name":        args[0],
				"key_":        key,
				"hostid":      hostID,
				"type":        itemType,
				"value_type":  valueType,
				"interfaceid": interfaceID,
				"delay":       delay,
			}

			result, err := client.Call("item.create", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)
			itemIDs := resp["itemids"].([]interface{})

			headers := []string{"Item Name", "Action", "Status", "ID"}
			rows := [][]string{{args[0], "Create", "Success", fmt.Sprintf("%v", itemIDs[0])}}
			outputResult(cmd, resp, headers, rows)
		},
	}

	cmd.Flags().StringVarP(&hostID, "hostid", "H", "", "Host ID for the item")
	cmd.Flags().StringVarP(&key, "key", "k", "", "Key for the item")
	cmd.Flags().IntVarP(&itemType, "type", "t", 0, "Item type (default: 0 - Zabbix agent)")
	cmd.Flags().IntVarP(&valueType, "value-type", "v", 3, "Value type (default: 3 - Numeric unsigned)")
	cmd.Flags().StringVarP(&interfaceID, "interfaceid", "i", "", "Interface ID for the item")
	cmd.Flags().StringVarP(&delay, "delay", "d", "1m", "Update interval (default: 1m)")
	cmd.MarkFlagRequired("hostid")
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("interfaceid")

	return cmd
}

