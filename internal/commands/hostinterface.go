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

	cmd.AddCommand(newHostInterfaceListCmd())

	return cmd
}

func newHostInterfaceListCmd() *cobra.Command {
	var hostName string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List interfaces for a host",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": "extend",
			}

			if hostName != "" {
				hostID := getHostID(client, hostName)
				if hostID == "" {
					fmt.Printf("Host not found: %s\n", hostName)
					return
				}
				params["hostids"] = []string{hostID}
			}

			result, err := client.Call("hostinterface.get", params)
			handleError(err)

			var interfaces []map[string]interface{}
			json.Unmarshal(result, &interfaces)

			fmt.Printf("%-10s %-20s %-10s %-10s %-10s\n", "ID", "IP", "DNS", "Port", "Type")
			for _, i := range interfaces {
				iType := getInterfaceTypeName(i["type"].(string))
				fmt.Printf("%-10s %-20s %-10s %-10s %-10s\n", i["interfaceid"], i["ip"], i["dns"], i["port"], iType)
			}
		},
	}

	cmd.Flags().StringVarP(&hostName, "host", "H", "", "Host name to list interfaces for")

	return cmd
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
