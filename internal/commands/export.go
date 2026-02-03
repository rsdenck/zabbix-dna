package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export Zabbix configurations",
	}

	cmd.AddCommand(newExportHostCmd())
	cmd.AddCommand(newExportTemplateCmd())

	return cmd
}

func newExportHostCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "host [host name]",
		Short: "Export a Zabbix host configuration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			hostID := getHostID(client, args[0])
			if hostID == "" {
				fmt.Printf("Host not found: %s\n", args[0])
				return
			}

			params := map[string]interface{}{
				"options": map[string]interface{}{
					"hosts": []string{hostID},
				},
				"format": "json",
			}

			result, err := client.Call("configuration.export", params)
			handleError(err)

			var prettyJSON json.RawMessage
			if err := json.Unmarshal(result, &prettyJSON); err == nil {
				formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
				fmt.Println(string(formatted))
			} else {
				fmt.Println(string(result))
			}
		},
	}
}

func newExportTemplateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "template [template name]",
		Short: "Export a Zabbix template configuration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			templateID := getTemplateID(client, args[0])
			if templateID == "" {
				fmt.Printf("Template not found: %s\n", args[0])
				return
			}

			params := map[string]interface{}{
				"options": map[string]interface{}{
					"templates": []string{templateID},
				},
				"format": "json",
			}

			result, err := client.Call("configuration.export", params)
			handleError(err)

			var prettyJSON json.RawMessage
			if err := json.Unmarshal(result, &prettyJSON); err == nil {
				formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
				fmt.Println(string(formatted))
			} else {
				fmt.Println(string(result))
			}
		},
	}
}
