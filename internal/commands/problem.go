package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newProblemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "problem",
		Short: "Manage Zabbix problems",
	}

	cmd.AddCommand(newProblemListCmd())

	return cmd
}

func newProblemListCmd() *cobra.Command {
	var limit int
	var severity int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List active Zabbix problems",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output":    []string{"eventid", "name", "severity", "clock", "objectid"},
				"limit":     limit,
				"sortfield": "eventid",
				"sortorder": "DESC",
			}

			if severity >= 0 {
				params["severities"] = []int{severity}
			}

			result, err := client.Call("problem.get", params)
			handleError(err)

			var problems []map[string]interface{}
			json.Unmarshal(result, &problems)

			headers := []string{"EventID", "Problem", "Severity", "Time"}
			var rows [][]string
			for _, p := range problems {
				sev := getPriorityName(p["severity"].(string))
				clockStr := p["clock"].(string)

				rows = append(rows, []string{
					fmt.Sprintf("%v", p["eventid"]),
					fmt.Sprintf("%v", p["name"]),
					sev,
					clockStr,
				})
			}

			outputResult(cmd, problems, headers, rows)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of problems")
	cmd.Flags().IntVarP(&severity, "severity", "s", -1, "Filter by severity (0-5)")

	return cmd
}


