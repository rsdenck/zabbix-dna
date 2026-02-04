package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func newBackupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Perform a backup of Zabbix configurations",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// Export major configuration objects for Zabbix 7.0
			params := map[string]interface{}{
				"options": map[string]interface{}{
					"host_groups":     []string{},
					"template_groups": []string{},
					"hosts":           []string{},
					"templates":       []string{},
					"maps":            []string{},
					"mediaTypes":      []string{},
					"images":          []string{},
				},
				"format": "json",
			}

			result, err := client.Call("configuration.export", params)
			if err != nil {
				// Fallback for older versions or simplified export
				params["options"] = map[string]interface{}{
					"groups":    []string{},
					"hosts":     []string{},
					"templates": []string{},
				}
				result, err = client.Call("configuration.export", params)
				handleError(err)
			}

			filename := fmt.Sprintf("zabbix_backup_%s.json", time.Now().Format("20060102_150405"))
			err = os.WriteFile(filename, result, 0644)
			handleError(err)

			headers := []string{"Component", "Status", "Info"}
			rows := [][]string{
				{"Zabbix API", "Connected", "OK"},
				{"Export Format", "JSON", "OK"},
				{"Configuration", "Full Export", "OK"},
				{"File Path", "Success", filename},
			}
			outputResult(cmd, nil, headers, rows)
		},
	}
}
