package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newMediaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "media",
		Short: "Manage Zabbix media and media types",
	}

	cmd.AddCommand(newMediaTypeListCmd())

	return cmd
}

func newMediaTypeListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "type-list",
		Short: "List Zabbix media types",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": []string{"mediatypeid", "name", "type", "status"},
			}

			result, err := client.Call("mediatype.get", params)
			handleError(err)

			var types []map[string]interface{}
			json.Unmarshal(result, &types)

			headers := []string{"ID", "Name", "Type", "Status"}
			var rows [][]string
			for _, t := range types {
				mType := getMediaTypeName(t["type"].(string))
				status := "Enabled"
				if t["status"].(string) == "1" {
					status = "Disabled"
				}
				rows = append(rows, []string{
					fmt.Sprintf("%v", t["mediatypeid"]),
					fmt.Sprintf("%v", t["name"]),
					mType,
					status,
				})
			}

			outputResult(cmd, types, headers, rows)
		},
	}

	return cmd
}

func getMediaTypeName(t string) string {
	switch t {
	case "0":
		return "Email"
	case "1":
		return "Script"
	case "2":
		return "SMS"
	case "4":
		return "WebHook"
	default:
		return "Other"
	}
}


