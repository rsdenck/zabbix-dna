package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newHostInterfaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hostinterface",
		Short: "Manage Zabbix host interfaces",
	}

	cmd.AddCommand(newHostInterfaceListCmd())   // show_host_interfaces -> hostinterface list
	cmd.AddCommand(newHostInterfaceCreateCmd()) // create_host_interface -> hostinterface create
	cmd.AddCommand(newHostInterfaceUpdateCmd()) // update_host_interface -> hostinterface update
	cmd.AddCommand(newHostInterfaceDeleteCmd()) // remove_host_interface -> hostinterface delete

	return cmd
}

func newHostInterfaceListCmd() *cobra.Command {
	var hostName string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show_host_interfaces"},
		Short:   "List interfaces for a host",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": "extend",
			}

			if hostName != "" {
				// Search by name
				res, err := client.Call("host.get", map[string]interface{}{
					"filter": map[string]interface{}{"host": hostName},
				})
				handleError(err)
				var hosts []map[string]interface{}
				json.Unmarshal(res, &hosts)
				if len(hosts) == 0 {
					fmt.Printf("Host not found: %s\n", hostName)
					return
				}
				params["hostids"] = []string{hosts[0]["hostid"].(string)}
			} else if len(args) > 0 {
				hostName = args[0]
				res, err := client.Call("host.get", map[string]interface{}{
					"filter": map[string]interface{}{"host": hostName},
				})
				handleError(err)
				var hosts []map[string]interface{}
				json.Unmarshal(res, &hosts)
				if len(hosts) == 0 {
					fmt.Printf("Host not found: %s\n", hostName)
					return
				}
				params["hostids"] = []string{hosts[0]["hostid"].(string)}
			}

			result, err := client.Call("hostinterface.get", params)
			handleError(err)

			var interfaces []map[string]interface{}
			json.Unmarshal(result, &interfaces)

			headers := []string{"ID", "IP", "DNS", "Port", "Type"}
			var rows [][]string
			for _, i := range interfaces {
				iType := getInterfaceTypeName(fmt.Sprintf("%v", i["type"]))
				rows = append(rows, []string{
					fmt.Sprintf("%v", i["interfaceid"]),
					fmt.Sprintf("%v", i["ip"]),
					fmt.Sprintf("%v", i["dns"]),
					fmt.Sprintf("%v", i["port"]),
					iType,
				})
			}

			outputResult(cmd, interfaces, headers, rows)
		},
	}

	cmd.Flags().StringVarP(&hostName, "host", "H", "", "Host name to list interfaces for")

	return cmd
}

func newHostInterfaceCreateCmd() *cobra.Command {
	var hostName string
	var ip string
	var dns string
	var port string
	var main bool
	var iType string

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"create_host_interface"},
		Short:   "Create a new host interface",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			hostID := getHostID(client, hostName)
			if hostID == "" {
				fmt.Printf("Host not found: %s\n", hostName)
				return
			}

			params := map[string]interface{}{
				"hostid": hostID,
				"ip":     ip,
				"dns":    dns,
				"port":   port,
				"main":   "0",
				"type":   iType,
				"useip":  "1",
			}
			if main {
				params["main"] = "1"
			}
			if dns != "" && ip == "" {
				params["useip"] = "0"
			}

			result, err := client.Call("hostinterface.create", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)

			headers := []string{"Host", "IP", "Port", "Action", "Status"}
			rows := [][]string{{hostName, ip, port, "Create Interface", "Success"}}
			outputResult(cmd, resp, headers, rows)
		},
	}

	cmd.Flags().StringVarP(&hostName, "host", "H", "", "Host name")
	cmd.Flags().StringVarP(&ip, "ip", "i", "", "IP address")
	cmd.Flags().StringVarP(&dns, "dns", "d", "", "DNS name")
	cmd.Flags().StringVarP(&port, "port", "p", "10050", "Port number")
	cmd.Flags().BoolVarP(&main, "main", "m", false, "Set as main interface")
	cmd.Flags().StringVarP(&iType, "type", "t", "1", "Interface type (1: Agent, 2: SNMP, 3: IPMI, 4: JMX)")

	cmd.MarkFlagRequired("host")

	return cmd
}

func newHostInterfaceUpdateCmd() *cobra.Command {
	var interfaceID string
	var ip string
	var dns string
	var port string

	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"update_host_interface"},
		Short:   "Update a host interface",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"interfaceid": interfaceID,
			}
			if ip != "" {
				params["ip"] = ip
			}
			if dns != "" {
				params["dns"] = dns
			}
			if port != "" {
				params["port"] = port
			}

			result, err := client.Call("hostinterface.update", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)

			headers := []string{"InterfaceID", "Action", "Status"}
			rows := [][]string{{interfaceID, "Update Interface", "Success"}}
			outputResult(cmd, resp, headers, rows)
		},
	}

	cmd.Flags().StringVarP(&interfaceID, "id", "i", "", "Interface ID")
	cmd.Flags().StringVar(&ip, "ip", "", "New IP address")
	cmd.Flags().StringVar(&dns, "dns", "", "New DNS name")
	cmd.Flags().StringVar(&port, "port", "", "New port number")

	cmd.MarkFlagRequired("id")

	return cmd
}

func newHostInterfaceDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete [interface id]",
		Aliases: []string{"remove_host_interface"},
		Short:   "Remove a host interface",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			interfaceID := args[0]
			result, err := client.Call("hostinterface.delete", []string{interfaceID})
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)

			headers := []string{"InterfaceID", "Action", "Status"}
			rows := [][]string{{interfaceID, "Delete Interface", "Success"}}
			outputResult(cmd, resp, headers, rows)
		},
	}
}

func getInterfaceTypeName(t string) string {
	switch t {
	case "1":
		return "Agent"
	case "2":
		return "SNMP"
	case "3":
		return "IPMI"
	case "4":
		return "JMX"
	default:
		return "Unknown"
	}
}

