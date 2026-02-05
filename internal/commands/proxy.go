package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newProxyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Manage Zabbix proxies",
	}

	cmd.AddCommand(newProxyListCmd())

	return cmd
}

func newProxyListCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Zabbix proxies",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": []string{"proxyid", "name", "operating_mode", "address", "version", "compatibility"},
				"limit":  limit,
			}
			result, err := client.Call("proxy.get", params)
			handleError(err)

			var proxies []map[string]interface{}
			json.Unmarshal(result, &proxies)

			headers := []string{"Name", "Address", "Mode", "Version", "Compatibility"}
			var rows [][]string

			for _, p := range proxies {
				name := fmt.Sprintf("%v", p["name"])
				address := fmt.Sprintf("%v", p["address"])
				if address == "" || address == "<nil>" {
					address = "127.0.0.1"
				}

				mode := "Active"
				if m, ok := p["operating_mode"].(string); ok && m == "1" {
					mode = "Passive"
				}

				version := fmt.Sprintf("%v", p["version"])
				if version == "" || version == "<nil>" {
					version = "0"
				}

				comp := "Undefined"
				if c, ok := p["compatibility"].(string); ok {
					switch c {
					case "1":
						comp = "Compatible"
					case "2":
						comp = "Incompatible"
					}
				}

				rows = append(rows, []string{name, address, mode, version, comp})
			}

			outputResult(cmd, proxies, headers, rows)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of proxies")

	return cmd
}
