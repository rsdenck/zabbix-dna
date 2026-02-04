package commands

import (
	"fmt"
	"os"

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
				handleError(fmt.Errorf("host not found: %s", args[0]))
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

			// According to absolute requirements, we must NOT output JSON directly.
			// We will save to a file and show a summary table.
			filename := fmt.Sprintf("export_host_%s_%s.json", args[0], fmt.Sprintf("%v", hostID))
			err = os.WriteFile(filename, result, 0644)
			handleError(err)

			headers := []string{"Host", "HostID", "Export File", "Status"}
			rows := [][]string{{args[0], hostID, filename, "Success"}}
			outputResult(cmd, nil, headers, rows)
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
				handleError(fmt.Errorf("template not found: %s", args[0]))
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

			// According to absolute requirements, we must NOT output JSON directly.
			// We will save to a file and show a summary table.
			filename := fmt.Sprintf("export_template_%s_%s.json", args[0], fmt.Sprintf("%v", templateID))
			err = os.WriteFile(filename, result, 0644)
			handleError(err)

			headers := []string{"Template", "TemplateID", "Export File", "Status"}
			rows := [][]string{{args[0], templateID, filename, "Success"}}
			outputResult(cmd, nil, headers, rows)
		},
	}
}
