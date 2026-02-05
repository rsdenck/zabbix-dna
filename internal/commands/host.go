package commands

import (
	"encoding/json"
	"fmt"
	"strings"
	"zabbix-dna/internal/config"

	"github.com/spf13/cobra"
)

func newHostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host",
		Short: "Manage Zabbix hosts",
	}

	// Mapeamento de comandos conforme solicitado
	cmd.AddCommand(newHostCreateCmd())  // create_host -> host create
	cmd.AddCommand(newHostDeleteCmd())  // remove_host -> host delete
	cmd.AddCommand(newHostUpdateCmd())  // update_host -> host update
	cmd.AddCommand(newHostShowCmd())    // show_host -> host show
	cmd.AddCommand(newHostListCmd())    // show_hosts -> host list
	cmd.AddCommand(newHostCloneCmd())   // clone_host -> host clone
	cmd.AddCommand(newHostEnableCmd())  // enable_host -> host enable
	cmd.AddCommand(newHostDisableCmd()) // disable_host -> host disable

	return cmd
}

func newHostListCmd() *cobra.Command {
	var limit int
	var search string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show_hosts"},
		Short:   "List Zabbix hosts",
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
		Use:     "show [host name]",
		Aliases: []string{"show_host"},
		Short:   "Show details of a Zabbix host",
		Args:    cobra.ExactArgs(1),
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
				handleError(fmt.Errorf("host not found: %s", args[0]))
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
	var groupNames []string
	var ip string
	var createInterface bool

	cmd := &cobra.Command{
		Use:     "create [host name]",
		Aliases: []string{"create_host"},
		Short:   "Create a new Zabbix host",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, _ := config.LoadConfig(cfgPath)

			// Default hostgroups from config
			if len(groupNames) == 0 && groupID == "" && cfg != nil && len(cfg.App.Commands.CreateHost.Hostgroups) > 0 {
				groupNames = cfg.App.Commands.CreateHost.Hostgroups
			}

			// Interface default from config
			if !cmd.Flags().Changed("interface") && cfg != nil {
				createInterface = cfg.App.Commands.CreateHost.CreateInterface
			}

			var groups []map[string]string
			if groupID != "" {
				groups = append(groups, map[string]string{"groupid": groupID})
			}
			if len(groupNames) > 0 {
				ids := getHostGroupsIDs(client, groupNames)
				for _, id := range ids {
					groups = append(groups, map[string]string{"groupid": id})
				}
			}

			if len(groups) == 0 {
				handleError(fmt.Errorf("at least one hostgroup must be specified (via --hostgroup, --groupid or config)"))
				return
			}

			params := map[string]interface{}{
				"host":   args[0],
				"groups": groups,
			}

			if createInterface {
				params["interfaces"] = []map[string]interface{}{
					{
						"type":  1, // Agent
						"main":  1,
						"useip": 1,
						"ip":    ip,
						"dns":   "",
						"port":  "10050",
					},
				}
			} else {
				params["interfaces"] = []interface{}{}
			}

			result, err := client.Call("host.create", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)

			hostIDs := resp["hostids"].([]interface{})
			outputResult(cmd, fmt.Sprintf("Created host %s (%v).", args[0], hostIDs[0]), nil, nil)
		},
	}

	cmd.Flags().StringVarP(&groupID, "groupid", "g", "", "Group ID for the host")
	cmd.Flags().StringSliceVar(&groupNames, "hostgroup", []string{}, "Host group names (comma-separated)")
	cmd.Flags().StringVarP(&ip, "ip", "i", "127.0.0.1", "IP address for the host interface")
	cmd.Flags().BoolVar(&createInterface, "interface", true, "Create an interface for the host")

	return cmd
}

func newHostDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete [host name]",
		Aliases: []string{"remove_host"},
		Short:   "Delete a Zabbix host",
		Args:    cobra.ExactArgs(1),
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
				handleError(fmt.Errorf("host not found: %s", args[0]))
				return
			}

			hostID := hosts[0]["hostid"].(string)

			// Delete the host
			resp, err := client.Call("host.delete", []string{hostID})
			handleError(err)

			var deleteResp map[string]interface{}
			json.Unmarshal(resp, &deleteResp)

			headers := []string{"Host", "Action", "Status"}
			rows := [][]string{{args[0], "Delete", "Success"}}
			outputResult(cmd, deleteResp, headers, rows)
		},
	}
}

func newHostUpdateCmd() *cobra.Command {
	var status string
	var name string

	cmd := &cobra.Command{
		Use:     "update [host name]",
		Aliases: []string{"update_host"},
		Short:   "Update a Zabbix host",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// Find host ID
			params := map[string]interface{}{
				"filter": map[string]interface{}{"host": args[0]},
			}
			result, err := client.Call("host.get", params)
			handleError(err)

			var hosts []map[string]interface{}
			json.Unmarshal(result, &hosts)
			if len(hosts) == 0 {
				handleError(fmt.Errorf("host not found: %s", args[0]))
				return
			}
			hostID := hosts[0]["hostid"].(string)

			updateParams := map[string]interface{}{
				"hostid": hostID,
			}
			if status != "" {
				s := "0"
				if status == "disable" || status == "off" || status == "1" {
					s = "1"
				}
				updateParams["status"] = s
			}
			if name != "" {
				updateParams["name"] = name
			}

			resp, err := client.Call("host.update", updateParams)
			handleError(err)

			var updateResp map[string]interface{}
			json.Unmarshal(resp, &updateResp)

			headers := []string{"Host", "Action", "Status"}
			rows := [][]string{{args[0], "Update", "Success"}}
			outputResult(cmd, updateResp, headers, rows)
		},
	}

	cmd.Flags().StringVarP(&status, "status", "s", "", "Set host status (enable/disable)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Set new host visible name")

	return cmd
}

func newHostEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "enable [host name]",
		Aliases: []string{"enable_host"},
		Short:   "Enable a Zabbix host",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"filter": map[string]interface{}{"host": args[0]},
			}
			result, err := client.Call("host.get", params)
			handleError(err)

			var hosts []map[string]interface{}
			json.Unmarshal(result, &hosts)
			if len(hosts) == 0 {
				handleError(fmt.Errorf("host not found: %s", args[0]))
				return
			}

			resp, err := client.Call("host.update", map[string]interface{}{
				"hostid": hosts[0]["hostid"],
				"status": "0",
			})
			handleError(err)

			var updateResp map[string]interface{}
			json.Unmarshal(resp, &updateResp)

			headers := []string{"Host", "Status", "Action"}
			rows := [][]string{{args[0], "Enabled", "Success"}}
			outputResult(cmd, updateResp, headers, rows)
		},
	}
}

func newHostDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "disable [host name]",
		Aliases: []string{"disable_host"},
		Short:   "Disable a Zabbix host",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"filter": map[string]interface{}{"host": args[0]},
			}
			result, err := client.Call("host.get", params)
			handleError(err)

			var hosts []map[string]interface{}
			json.Unmarshal(result, &hosts)
			if len(hosts) == 0 {
				handleError(fmt.Errorf("host not found: %s", args[0]))
				return
			}

			resp, err := client.Call("host.update", map[string]interface{}{
				"hostid": hosts[0]["hostid"],
				"status": "1",
			})
			handleError(err)

			var updateResp map[string]interface{}
			json.Unmarshal(resp, &updateResp)

			headers := []string{"Host", "Status", "Action"}
			rows := [][]string{{args[0], "Disabled", "Success"}}
			outputResult(cmd, updateResp, headers, rows)
		},
	}
}

func newHostCloneCmd() *cobra.Command {
	var newName string
	cmd := &cobra.Command{
		Use:     "clone [source host name]",
		Aliases: []string{"clone_host"},
		Short:   "Clone a Zabbix host",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// Get source host full details
			params := map[string]interface{}{
				"filter":           map[string]interface{}{"host": args[0]},
				"selectGroups":     "extend",
				"selectInterfaces": "extend",
				"selectTemplates":  "extend",
				"selectMacros":     "extend",
			}
			result, err := client.Call("host.get", params)
			handleError(err)

			var hosts []map[string]interface{}
			json.Unmarshal(result, &hosts)
			if len(hosts) == 0 {
				handleError(fmt.Errorf("source host not found: %s", args[0]))
				return
			}
			src := hosts[0]

			if newName == "" {
				newName = fmt.Sprintf("%s_CLONE", src["host"])
			}

			// Prepare clone params
			cloneParams := map[string]interface{}{
				"host":       newName,
				"name":       newName,
				"status":     src["status"],
				"groups":     src["groups"],
				"templates":  src["templates"],
				"interfaces": src["interfaces"],
				"macros":     src["macros"],
			}

			// Remove read-only or internal fields from select result
			if interfaces, ok := cloneParams["interfaces"].([]interface{}); ok {
				for _, iface := range interfaces {
					if m, ok := iface.(map[string]interface{}); ok {
						delete(m, "interfaceid")
						delete(m, "hostid")
					}
				}
			}

			resp, err := client.Call("host.create", cloneParams)
			handleError(err)

			var createResp map[string]interface{}
			json.Unmarshal(resp, &createResp)

			headers := []string{"Source Host", "Cloned Host", "Action", "Status"}
			rows := [][]string{{args[0], newName, "Clone", "Success"}}
			outputResult(cmd, createResp, headers, rows)
		},
	}

	cmd.Flags().StringVarP(&newName, "new-name", "n", "", "New name for the cloned host")
	return cmd
}
