package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newTemplateGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templategroup",
		Short: "Manage Zabbix template groups",
	}

	cmd.AddCommand(newTemplateGroupListCmd())
	cmd.AddCommand(newTemplateGroupCreateCmd())
	cmd.AddCommand(newTemplateGroupDeleteCmd())

	return cmd
}

func newTemplateGroupListCmd() *cobra.Command {
	var limit int
	var search string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Zabbix template groups",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": []string{"groupid", "name"},
				"limit":  limit,
			}
			if search != "" {
				params["search"] = map[string]interface{}{
					"name": search,
				}
			}

			result, err := client.Call("templategroup.get", params)
			handleError(err)

			var groups []map[string]interface{}
			json.Unmarshal(result, &groups)

			headers := []string{"ID", "Name"}
			var rows [][]string
			for _, g := range groups {
				rows = append(rows, []string{
					fmt.Sprintf("%v", g["groupid"]),
					fmt.Sprintf("%v", g["name"]),
				})
			}

			outputResult(cmd, groups, headers, rows)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of groups")
	cmd.Flags().StringVarP(&search, "search", "s", "", "Search for a group by name")

	return cmd
}

func newTemplateGroupCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [group name]",
		Short: "Create a new Zabbix template group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"name": args[0],
			}

			result, err := client.Call("templategroup.create", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)
			groupIDs := resp["groupids"].([]interface{})

			fmt.Printf("Template group created successfully with ID: %s\n", groupIDs[0])
		},
	}
}

func newTemplateGroupDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [group name]",
		Short: "Delete a Zabbix template group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// First find the group ID
			params := map[string]interface{}{
				"filter": map[string]interface{}{
					"name": args[0],
				},
			}

			result, err := client.Call("templategroup.get", params)
			handleError(err)

			var groups []map[string]interface{}
			json.Unmarshal(result, &groups)

			if len(groups) == 0 {
				fmt.Printf("Template group not found: %s\n", args[0])
				return
			}

			groupID := groups[0]["groupid"].(string)

			// Delete the group
			_, err = client.Call("templategroup.delete", []string{groupID})
			handleError(err)

			fmt.Printf("Template group %s (ID: %s) deleted successfully\n", args[0], groupID)
		},
	}
}
