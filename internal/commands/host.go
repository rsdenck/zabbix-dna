package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newHostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host",
		Short: "Manage Zabbix hosts",
	}

	cmd.AddCommand(newHostListCmd())
	cmd.AddCommand(newHostShowCmd())
	cmd.AddCommand(newHostCreateCmd())
	cmd.AddCommand(newHostDeleteCmd())

	return cmd
}

func newHostListCmd() *cobra.Command {
	var limit int
	var search string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Zabbix hosts",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output":           []string{"hostid", "host", "name", "status", "maintenance_status", "proxy_hostid"},
				"selectGroups":     []string{"name"},
				"selectTemplates":  []string{"name"},
				"selectInterfaces": []string{"type", "available"},
				"limit":            limit,
			}
			if search != "" {
				params["search"] = map[string]interface{}{
					"host": search,
				}
			}

			result, err := client.Call("host.get", params)
			handleError(err)

			var hosts []map[string]interface{}
			json.Unmarshal(result, &hosts)

			// Fetch proxies to resolve names
			proxyResult, err := client.Call("proxy.get", map[string]interface{}{
				"output": []string{"proxyid", "name"},
			})
			var proxies []map[string]interface{}
			if err == nil {
				json.Unmarshal(proxyResult, &proxies)
			}
			proxyMap := make(map[string]string)
			for _, p := range proxies {
				id := fmt.Sprintf("%v", p["proxyid"])
				name := fmt.Sprintf("%v", p["name"])
				proxyMap[id] = name
			}

			headers := []string{"HostID", "Name", "Host groups", "Templates", "Zabbix agent", "Maintenance", "Status", "Proxy"}
			var rows [][]string

			for _, h := range hosts {
				hostID := fmt.Sprintf("%v", h["hostid"])
				name := fmt.Sprintf("%v", h["name"])

				// Groups
				var groups []string
				if g, ok := h["groups"].([]interface{}); ok {
					for _, item := range g {
						if group, ok := item.(map[string]interface{}); ok {
							groups = append(groups, fmt.Sprintf("%v", group["name"]))
						}
					}
				}
				groupsStr := strings.Join(groups, ", ")

				// Templates
				var templates []string
				if t, ok := h["templates"].([]interface{}); ok {
					for _, item := range t {
						if tmpl, ok := item.(map[string]interface{}); ok {
							templates = append(templates, fmt.Sprintf("%v", tmpl["name"]))
						}
					}
				}
				templatesStr := strings.Join(templates, ", ")

				// Agent availability
				agentStatus := "Unknown"
				if interfaces, ok := h["interfaces"].([]interface{}); ok {
					for _, item := range interfaces {
						if iface, ok := item.(map[string]interface{}); ok {
							if iface["type"].(string) == "1" { // Agent
								switch iface["available"].(string) {
								case "1":
									agentStatus = "Available"
								case "2":
									agentStatus = "Unavailable"
								}
							}
						}
					}
				}

				// Maintenance
				maintenance := "Off"
				if h["maintenance_status"].(string) == "1" {
					maintenance = "On"
				}

				// Status
				status := "On"
				if h["status"].(string) == "1" {
					status = "Off"
				}

				// Proxy
				proxyName := "None"
				if pID, ok := h["proxy_hostid"].(string); ok && pID != "0" {
					if name, ok := proxyMap[pID]; ok {
						proxyName = name
					} else {
						proxyName = pID
					}
				}

				rows = append(rows, []string{hostID, name, groupsStr, templatesStr, agentStatus, maintenance, status, proxyName})
			}

			outputResult(cmd, hosts, headers, rows)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of hosts")
	cmd.Flags().StringVarP(&search, "search", "s", "", "Search for a host by name")

	return cmd
}

func newHostShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show [host name]",
		Short: "Show details of a Zabbix host",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"filter": map[string]interface{}{
					"host": args[0],
				},
				"selectGroups":     "extend",
				"selectInterfaces": "extend",
				"selectTemplates":  "extend",
			}

			result, err := client.Call("host.get", params)
			handleError(err)

			var hosts []map[string]interface{}
			json.Unmarshal(result, &hosts)

			if len(hosts) == 0 {
				fmt.Println("Host not found")
				return
			}

			h := hosts[0]
			headers := []string{"Property", "Value"}
			var rows [][]string

			rows = append(rows, []string{"HostID", fmt.Sprintf("%v", h["hostid"])})
			rows = append(rows, []string{"Name", fmt.Sprintf("%v", h["name"])})
			rows = append(rows, []string{"Host", fmt.Sprintf("%v", h["host"])})

			// Status
			status := "On"
			if h["status"].(string) == "1" {
				status = "Off"
			}
			rows = append(rows, []string{"Status", status})

			// Groups
			var groups []string
			if g, ok := h["groups"].([]interface{}); ok {
				for _, item := range g {
					if group, ok := item.(map[string]interface{}); ok {
						groups = append(groups, fmt.Sprintf("%v", group["name"]))
					}
				}
			}
			for i, g := range groups {
				label := "Group"
				if i > 0 {
					label = ""
				}
				rows = append(rows, []string{label, g})
			}

			// Interfaces
			if interfaces, ok := h["interfaces"].([]interface{}); ok {
				for i, item := range interfaces {
					if iface, ok := item.(map[string]interface{}); ok {
						label := "Interface"
						if i > 0 {
							label = ""
						}
						val := fmt.Sprintf("%s:%s", iface["ip"], iface["port"])
						rows = append(rows, []string{label, val})
					}
				}
			}

			outputResult(cmd, h, headers, rows)
		},
	}
}

func newHostCreateCmd() *cobra.Command {
	var groupID string
	var ip string

	cmd := &cobra.Command{
		Use:   "create [host name]",
		Short: "Create a new Zabbix host",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"host": args[0],
				"groups": []map[string]string{
					{"groupid": groupID},
				},
				"interfaces": []map[string]interface{}{
					{
						"type":  1, // Agent
						"main":  1,
						"useip": 1,
						"ip":    ip,
						"dns":   "",
						"port":  "10050",
					},
				},
			}

			result, err := client.Call("host.create", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)

			format, _ := cmd.Flags().GetString("format")
			if format == "json" {
				outputResult(cmd, resp, nil, nil)
				return
			}

			hostIDs := resp["hostids"].([]interface{})
			fmt.Printf("Host created successfully with ID: %s\n", hostIDs[0])
		},
	}

	cmd.Flags().StringVarP(&groupID, "groupid", "g", "", "Group ID for the host")
	cmd.Flags().StringVarP(&ip, "ip", "i", "127.0.0.1", "IP address for the host interface")
	cmd.MarkFlagRequired("groupid")

	return cmd
}

func newHostDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [host name]",
		Short: "Delete a Zabbix host",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// First find the host ID
			params := map[string]interface{}{
				"filter": map[string]interface{}{
					"host": args[0],
				},
			}

			result, err := client.Call("host.get", params)
			handleError(err)

			var hosts []map[string]interface{}
			json.Unmarshal(result, &hosts)

			if len(hosts) == 0 {
				fmt.Printf("Host not found: %s\n", args[0])
				return
			}

			hostID := hosts[0]["hostid"].(string)

			// Delete the host
			_, err = client.Call("host.delete", []string{hostID})
			handleError(err)

			fmt.Printf("Host %s (ID: %s) deleted successfully\n", args[0], hostID)
		},
	}
}
