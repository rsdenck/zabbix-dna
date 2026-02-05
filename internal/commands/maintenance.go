package commands

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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
	cmd.AddCommand(newMaintenanceRemoveCmd())

	return cmd
}

func newMaintenanceRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove [maintenance id(s)]",
		Aliases: []string{"remove_maintenance_definition"},
		Short:   "Remove a maintenance definition by ID(s)",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// Delete the maintenance(s)
			result, err := client.Call("maintenance.delete", args)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)

			outputResult(cmd, "Removed maintenance definition(s).", nil, nil)
		},
	}
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

				sinceSec, _ := strconv.ParseInt(p["active_since"].(string), 10, 64)
				tillSec, _ := strconv.ParseInt(p["active_till"].(string), 10, 64)

				rows = append(rows, []string{
					fmt.Sprintf("%v", p["maintenanceid"]),
					fmt.Sprintf("%v", p["name"]),
					mType,
					time.Unix(sinceSec, 0).Format("2006-01-02 15:04:05"),
					time.Unix(tillSec, 0).Format("2006-01-02 15:04:05"),
				})
			}

			outputResult(cmd, periods, headers, rows)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of periods")

	return cmd
}

func newMaintenanceCreateCmd() *cobra.Command {
	var hosts []string
	var hostgroups []string
	var activeSince int64
	var activeTill int64
	var period int

	cmd := &cobra.Command{
		Use:   "create [maintenance name]",
		Short: "Create a new Zabbix maintenance period",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			if activeSince == 0 {
				activeSince = time.Now().Unix()
			}
			if activeTill == 0 {
				if period == 0 {
					period = 3600 // 1 hour default
				}
				activeTill = activeSince + int64(period)
			}

			params := map[string]interface{}{
				"name":         args[0],
				"active_since": activeSince,
				"active_till":  activeTill,
				"timeperiods": []map[string]interface{}{
					{
						"timeperiod_type": 0, // One time only
						"start_date":      activeSince,
						"period":          activeTill - activeSince,
					},
				},
			}

			if len(hosts) > 0 {
				params["hostids"] = getHostsIDs(client, hosts)
			}
			if len(hostgroups) > 0 {
				params["groupids"] = getHostGroupsIDs(client, hostgroups)
			}

			result, err := client.Call("maintenance.create", params)
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)
			maintenanceIDs := resp["maintenanceids"].([]interface{})

			outputResult(cmd, fmt.Sprintf("Created maintenance definition (%v).", maintenanceIDs[0]), nil, nil)
		},
	}

	cmd.Flags().StringSliceVar(&hosts, "host", []string{}, "Host names (comma-separated)")
	cmd.Flags().StringSliceVar(&hostgroups, "hostgroup", []string{}, "Host group names (comma-separated)")
	cmd.Flags().Int64Var(&activeSince, "since", 0, "Active since (Unix timestamp, defaults to now)")
	cmd.Flags().Int64Var(&activeTill, "till", 0, "Active till (Unix timestamp)")
	cmd.Flags().IntVar(&period, "period", 3600, "Period in seconds (default 1 hour if till is not set)")

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
				handleError(fmt.Errorf("maintenance period not found: %s", args[0]))
				return
			}

			maintenanceID := periods[0]["maintenanceid"].(string)

			// Delete the period
			result, err = client.Call("maintenance.delete", []string{maintenanceID})
			handleError(err)

			var resp map[string]interface{}
			json.Unmarshal(result, &resp)

			headers := []string{"Maintenance", "Action", "Status", "ID"}
			rows := [][]string{{args[0], "Delete", "Success", maintenanceID}}
			outputResult(cmd, resp, headers, rows)
		},
	}
}
