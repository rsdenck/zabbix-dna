package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newMaintenanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maintenance",
		Short: "Manage Zabbix maintenance periods",
	}

	cmd.AddCommand(newMaintenanceListCmd())
	cmd.AddCommand(newMaintenanceCreateCmd())
	cmd.AddCommand(newMaintenanceDeleteCmd())

	return cmd
}

func newMaintenanceListCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Zabbix maintenance periods",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": []string{"maintenanceid", "name", "maintenance_type", "active_since", "active_till"},
				"limit":  limit,
			}

			result, err := client.Call("maintenance.get", params)
			handleError(err)

			var periods []map[string]interface{}
			json.Unmarshal(result, &periods)

			headers := []string{"MaintenanceID", "Name", "Type", "Since", "Till"}
			var rows [][]string
			for _, p := range periods {
				mType := "With data"
				if p["maintenance_type"].(string) == "1" {
					mType = "No data"
				}
				rows = append(rows, []string{
					fmt.Sprintf("%v", p["maintenanceid"]),
					fmt.Sprintf("%v", p["name"]),
					mType,
					fmt.Sprintf("%v", p["active_since"]),
					fmt.Sprintf("%v", p["active_till"]),
				})
			}

			outputResult(cmd, periods, headers, rows)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of periods")

	return cmd
}

func newMaintenanceCreateCmd() *cobra.Command {
	var hostID string
	var groupID string
	var activeSince string
	var activeTill string

	cmd := &cobra.Command{
		Use:   "create [maintenance name]",
		Short: "Create a new Zabbix maintenance period",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"name":         args[0],
				"active_since": activeSince,
				"active_till":  activeTill,
				"timeperiods": []map[string]interface{}{
					{
						"timeperiod_type": 0, // One time only
						"start_date":      activeSince,
						"period":          3600, // 1 hour default
					},
				},
			}

			if hostID != "" {
				params["hostids"] = []string{hostID}
			}
			if groupID != "" {
				params["groupids"] = []string{groupID}
			}

			result, err := client.Call("maintenance.create", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)
			maintenanceIDs := resp["maintenanceids"].([]interface{})

			fmt.Printf("Maintenance period created successfully with ID: %s\n", maintenanceIDs[0])
		},
	}

	cmd.Flags().StringVarP(&hostID, "hostid", "H", "", "Host ID for maintenance")
	cmd.Flags().StringVarP(&groupID, "groupid", "g", "", "Host group ID for maintenance")
	cmd.Flags().StringVarP(&activeSince, "since", "s", "", "Active since (Unix timestamp)")
	cmd.Flags().StringVarP(&activeTill, "till", "t", "", "Active till (Unix timestamp)")
	cmd.MarkFlagRequired("since")
	cmd.MarkFlagRequired("till")

	return cmd
}

func newMaintenanceDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [maintenance name]",
		Short: "Delete a Zabbix maintenance period",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// First find the maintenance ID
			params := map[string]interface{}{
				"filter": map[string]interface{}{
					"name": args[0],
				},
			}

			result, err := client.Call("maintenance.get", params)
			handleError(err)

			var periods []map[string]interface{}
			json.Unmarshal(result, &periods)

			if len(periods) == 0 {
				fmt.Printf("Maintenance period not found: %s\n", args[0])
				return
			}

			maintenanceID := periods[0]["maintenanceid"].(string)

			// Delete the period
			_, err = client.Call("maintenance.delete", []string{maintenanceID})
			handleError(err)

			fmt.Printf("Maintenance period %s (ID: %s) deleted successfully\n", args[0], maintenanceID)
		},
	}
}
