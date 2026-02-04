package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newHostGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hostgroup",
		Short: "Manage Zabbix host groups",
	}

	cmd.AddCommand(newHostGroupCreateCmd())      // create_hostgroup -> hostgroup create
	cmd.AddCommand(newHostGroupDeleteCmd())      // remove_hostgroup -> hostgroup delete
	cmd.AddCommand(newHostGroupShowCmd())        // show_hostgroup -> hostgroup show
	cmd.AddCommand(newHostGroupListCmd())        // show_hostgroups -> hostgroup list
	cmd.AddCommand(newAddHostToGroupCmd())       // add_host_to_hostgroup -> hostgroup add-host
	cmd.AddCommand(newRemoveHostFromGroupCmd())  // remove_host_from_hostgroup -> hostgroup remove-host
	cmd.AddCommand(newHostGroupPermissionsCmd()) // show_hostgroup_permissions -> hostgroup permissions

	return cmd
}

func newHostGroupListCmd() *cobra.Command {
	var limit int
	var search string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show_hostgroups"},
		Short:   "List Zabbix host groups",
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

			result, err := client.Call("hostgroup.get", params)
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

func newHostGroupCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [group name]",
		Short: "Create a new Zabbix host group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"name": args[0],
			}

			result, err := client.Call("hostgroup.create", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)

			groupIDs := resp["groupids"].([]interface{})
			headers := []string{"Host Group", "Action", "Status", "ID"}
			rows := [][]string{{args[0], "Create", "Success", fmt.Sprintf("%v", groupIDs[0])}}
			outputResult(cmd, resp, headers, rows)
		},
	}
}

func newHostGroupDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete [group name]",
		Aliases: []string{"remove_hostgroup"},
		Short:   "Delete a Zabbix host group",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// First find the group ID
			params := map[string]interface{}{
				"filter": map[string]interface{}{
					"name": args[0],
				},
			}

			result, err := client.Call("hostgroup.get", params)
			handleError(err)

			var groups []map[string]interface{}
			json.Unmarshal(result, &groups)

			if len(groups) == 0 {
				fmt.Printf("Host group not found: %s\n", args[0])
				return
			}

			groupID := groups[0]["groupid"].(string)

			// Delete the group
			resp, err := client.Call("hostgroup.delete", []string{groupID})
			handleError(err)

			var deleteResp map[string]interface{}
			json.Unmarshal(resp, &deleteResp)

			headers := []string{"Host Group", "Action", "Status"}
			rows := [][]string{{args[0], "Delete", "Success"}}
			outputResult(cmd, deleteResp, headers, rows)
		},
	}
}

func newHostGroupShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show [group name]",
		Aliases: []string{"show_hostgroup"},
		Short:   "Show details of a Zabbix host group",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"filter": map[string]interface{}{
					"name": args[0],
				},
				"selectHosts": "extend",
			}

			result, err := client.Call("hostgroup.get", params)
			handleError(err)

			var groups []map[string]interface{}
			json.Unmarshal(result, &groups)

			if len(groups) == 0 {
				fmt.Println("Host group not found")
				return
			}

			g := groups[0]
			headers := []string{"Property", "Value"}
			var rows [][]string

			rows = append(rows, []string{"GroupID", fmt.Sprintf("%v", g["groupid"])})
			rows = append(rows, []string{"Name", fmt.Sprintf("%v", g["name"])})

			if hosts, ok := g["hosts"].([]interface{}); ok {
				for i, h := range hosts {
					if host, ok := h.(map[string]interface{}); ok {
						label := "Host"
						if i > 0 {
							label = ""
						}
						rows = append(rows, []string{label, fmt.Sprintf("%v", host["name"])})
					}
				}
			}

			outputResult(cmd, g, headers, rows)
		},
	}
}

func newAddHostToGroupCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "add-host [host names] [group names]",
		Aliases: []string{"add_host_to_hostgroup"},
		Short:   "Add hosts to host groups",
		Long:    "Add one or more hosts (comma-separated) to one or more host groups (comma-separated)",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			hostNames := strings.Split(args[0], ",")
			groupNames := strings.Split(args[1], ",")

			var hostIDs []string
			var groupIDs []string

			// Resolve Host IDs
			for _, hn := range hostNames {
				hn = strings.TrimSpace(hn)
				res, err := client.Call("host.get", map[string]interface{}{
					"filter": map[string]interface{}{"host": hn},
				})
				handleError(err)
				var hosts []map[string]interface{}
				json.Unmarshal(res, &hosts)
				if len(hosts) == 0 {
					fmt.Printf("Host not found: %s\n", hn)
					continue
				}
				hostIDs = append(hostIDs, hosts[0]["hostid"].(string))
			}

			// Resolve Group IDs
			for _, gn := range groupNames {
				gn = strings.TrimSpace(gn)
				res, err := client.Call("hostgroup.get", map[string]interface{}{
					"filter": map[string]interface{}{"name": gn},
				})
				handleError(err)
				var groups []map[string]interface{}
				json.Unmarshal(res, &groups)
				if len(groups) == 0 {
					fmt.Printf("Host group not found: %s\n", gn)
					continue
				}
				groupIDs = append(groupIDs, groups[0]["groupid"].(string))
			}

			if len(hostIDs) == 0 || len(groupIDs) == 0 {
				fmt.Println("No valid hosts or groups found to process.")
				return
			}

			// Add hosts to groups (massadd)
			var groupsParam []map[string]string
			for _, id := range groupIDs {
				groupsParam = append(groupsParam, map[string]string{"groupid": id})
			}
			var hostsParam []map[string]string
			for _, id := range hostIDs {
				hostsParam = append(hostsParam, map[string]string{"hostid": id})
			}

			resp, err := client.Call("hostgroup.massadd", map[string]interface{}{
				"groups": groupsParam,
				"hosts":  hostsParam,
			})
			handleError(err)

			var massResp map[string]interface{}
			json.Unmarshal(resp, &massResp)

			headers := []string{"Hosts", "Groups", "Action", "Status"}
			rows := [][]string{{args[0], args[1], "Add to Group", "Success"}}
			outputResult(cmd, massResp, headers, rows)
		},
	}
}

func newRemoveHostFromGroupCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove-host [host names] [group names]",
		Aliases: []string{"remove_host_from_hostgroup"},
		Short:   "Remove hosts from host groups",
		Long:    "Remove one or more hosts (comma-separated) from one or more host groups (comma-separated)",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			hostNames := strings.Split(args[0], ",")
			groupNames := strings.Split(args[1], ",")

			var hostIDs []string
			var groupIDs []string

			// Resolve Host IDs
			for _, hn := range hostNames {
				hn = strings.TrimSpace(hn)
				res, err := client.Call("host.get", map[string]interface{}{
					"filter": map[string]interface{}{"host": hn},
				})
				handleError(err)
				var hosts []map[string]interface{}
				json.Unmarshal(res, &hosts)
				if len(hosts) == 0 {
					fmt.Printf("Host not found: %s\n", hn)
					continue
				}
				hostIDs = append(hostIDs, hosts[0]["hostid"].(string))
			}

			// Resolve Group IDs
			for _, gn := range groupNames {
				gn = strings.TrimSpace(gn)
				res, err := client.Call("hostgroup.get", map[string]interface{}{
					"filter": map[string]interface{}{"name": gn},
				})
				handleError(err)
				var groups []map[string]interface{}
				json.Unmarshal(res, &groups)
				if len(groups) == 0 {
					fmt.Printf("Host group not found: %s\n", gn)
					continue
				}
				groupIDs = append(groupIDs, groups[0]["groupid"].(string))
			}

			if len(hostIDs) == 0 || len(groupIDs) == 0 {
				fmt.Println("No valid hosts or groups found to process.")
				return
			}

			// Remove hosts from groups (massremove)
			resp, err := client.Call("hostgroup.massremove", map[string]interface{}{
				"groupids": groupIDs,
				"hostids":  hostIDs,
			})
			handleError(err)

			var massResp map[string]interface{}
			json.Unmarshal(resp, &massResp)

			headers := []string{"Hosts", "Groups", "Action", "Status"}
			rows := [][]string{{args[0], args[1], "Remove from Group", "Success"}}
			outputResult(cmd, massResp, headers, rows)
		},
	}
}

func newHostGroupPermissionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "permissions [group names]",
		Aliases: []string{"show_hostgroup_permissions"},
		Short:   "Show permissions for host groups",
		Long:    "Show which user groups have access to the specified host groups",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			groupNames := strings.Split(args[0], ",")
			var groupIDs []string

			// Resolve Group IDs
			for _, gn := range groupNames {
				gn = strings.TrimSpace(gn)
				res, err := client.Call("hostgroup.get", map[string]interface{}{
					"filter": map[string]interface{}{"name": gn},
				})
				handleError(err)
				var groups []map[string]interface{}
				json.Unmarshal(res, &groups)
				if len(groups) == 0 {
					fmt.Printf("Host group not found: %s\n", gn)
					continue
				}
				groupIDs = append(groupIDs, groups[0]["groupid"].(string))
			}

			if len(groupIDs) == 0 {
				return
			}

			// Get user groups and their permissions
			res, err := client.Call("usergroup.get", map[string]interface{}{
				"selectRights": "extend",
				"output":       []string{"usrgrpid", "name"},
			})
			handleError(err)
			var userGroups []map[string]interface{}
			json.Unmarshal(res, &userGroups)

			headers := []string{"Host Group", "User Group", "Permission"}
			var rows [][]string

			permissionMap := map[string]string{
				"0": "None",
				"2": "Read-only",
				"3": "Read-write",
			}

			for _, gn := range groupNames {
				gn = strings.TrimSpace(gn)
				// Find group ID again or use a map
				res, _ := client.Call("hostgroup.get", map[string]interface{}{"filter": map[string]interface{}{"name": gn}})
				var groups []map[string]interface{}
				json.Unmarshal(res, &groups)
				if len(groups) == 0 {
					continue
				}
				gid := groups[0]["groupid"].(string)

				for _, ug := range userGroups {
					if rights, ok := ug["rights"].([]interface{}); ok {
						for _, r := range rights {
							right := r.(map[string]interface{})
							if right["id"].(string) == gid {
								perm := permissionMap[right["permission"].(string)]
								rows = append(rows, []string{gn, ug["name"].(string), perm})
							}
						}
					}
				}
			}

			if len(rows) == 0 {
				fmt.Println("No specific permissions found for these groups.")
				return
			}

			outputResult(cmd, userGroups, headers, rows)
		},
	}
}


