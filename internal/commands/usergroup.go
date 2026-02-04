package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newUserGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "usergroup",
		Short: "Manage Zabbix user groups",
	}

	cmd.AddCommand(newUserGroupListCmd())
	cmd.AddCommand(newUserGroupDeleteCmd())

	return cmd
}

func newUserGroupListCmd() *cobra.Command {
	var limit int
	var search string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Zabbix user groups",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": []string{"usrgrpid", "name"},
				"limit":  limit,
			}
			if search != "" {
				params["search"] = map[string]interface{}{
					"name": search,
				}
			}

			result, err := client.Call("usergroup.get", params)
			handleError(err)

			var groups []map[string]interface{}
			json.Unmarshal(result, &groups)

			headers := []string{"UsrGrpID", "Name"}
			var rows [][]string
			for _, g := range groups {
				rows = append(rows, []string{
					fmt.Sprintf("%v", g["usrgrpid"]),
					fmt.Sprintf("%v", g["name"]),
				})
			}

			outputResult(cmd, groups, headers, rows)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of user groups")
	cmd.Flags().StringVarP(&search, "search", "s", "", "Search for a user group by name")

	return cmd
}

func newUserGroupDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [group name]",
		Short: "Delete a Zabbix user group",
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

			result, err := client.Call("usergroup.get", params)
			handleError(err)

			var groups []map[string]interface{}
			json.Unmarshal(result, &groups)

			if len(groups) == 0 {
				fmt.Printf("User group not found: %s\n", args[0])
				return
			}

			groupID := groups[0]["usrgrpid"].(string)

			// Delete the group
			resp, err := client.Call("usergroup.delete", []string{groupID})
			handleError(err)

			var deleteResp map[string]interface{}
			json.Unmarshal(resp, &deleteResp)

			headers := []string{"User Group", "Action", "Status"}
			rows := [][]string{{args[0], "Delete", "Success"}}
			outputResult(cmd, deleteResp, headers, rows)
		},
	}
}
