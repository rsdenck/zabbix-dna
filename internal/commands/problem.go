package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newProblemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "problem",
		Short: "Manage Zabbix problems",
	}

	cmd.AddCommand(newProblemListCmd())
	cmd.AddCommand(newProblemAcknowledgeCmd())
	cmd.AddCommand(newProblemAcknowledgeTriggerCmd())
	cmd.AddCommand(newProblemShowEventsCmd())
	cmd.AddCommand(newProblemShowAlarmsCmd())

	return cmd
}

func newProblemShowEventsCmd() *cobra.Command {
	var triggerIDs []string
	var hostGroups []string
	var hosts []string
	var limit int

	cmd := &cobra.Command{
		Use:     "events",
		Aliases: []string{"show_trigger_events"},
		Short:   "Show recent events for triggers, hosts, or hostgroups",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// Support legacy positional arguments: show_trigger_events [triggerid] [limit]
			if len(args) > 0 {
				triggerIDs = append(triggerIDs, args[0])
				if len(args) > 1 {
					fmt.Sscanf(args[1], "%d", &limit)
				}
			}

			params := map[string]interface{}{
				"output":    "extend",
				"sortfield": "clock",
				"sortorder": "DESC",
				"limit":     limit,
			}

			if len(triggerIDs) > 0 {
				params["objectids"] = triggerIDs
			}
			if len(hostGroups) > 0 {
				params["groupids"] = getHostGroupsIDs(client, hostGroups)
			}
			if len(hosts) > 0 {
				params["hostids"] = getHostsIDs(client, hosts)
			}

			result, err := client.Call("event.get", params)
			handleError(err)

			var events []map[string]interface{}
			json.Unmarshal(result, &events)

			headers := []string{"EventID", "ObjectID", "Name", "Severity", "Time", "Ack"}
			var rows [][]string
			for _, e := range events {
				ack := "No"
				if e["acknowledged"].(string) == "1" {
					ack = "Yes"
				}
				sev := getPriorityName(e["severity"].(string))
				rows = append(rows, []string{
					fmt.Sprintf("%v", e["eventid"]),
					fmt.Sprintf("%v", e["objectid"]),
					fmt.Sprintf("%v", e["name"]),
					sev,
					fmt.Sprintf("%v", e["clock"]),
					ack,
				})
			}

			outputResult(cmd, events, headers, rows)
		},
	}

	cmd.Flags().StringSliceVar(&triggerIDs, "trigger-id", []string{}, "Trigger ID(s)")
	cmd.Flags().StringSliceVar(&hostGroups, "hostgroup", []string{}, "Host group name(s)")
	cmd.Flags().StringSliceVar(&hosts, "host", []string{}, "Host name(s)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 10, "Limit the number of events")

	return cmd
}

func newProblemShowAlarmsCmd() *cobra.Command {
	var description string
	var priority int
	var hostGroups []string
	var unacknowledged bool

	cmd := &cobra.Command{
		Use:     "alarms",
		Aliases: []string{"show_alarms"},
		Short:   "Show active alarms/triggers",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// Support legacy positional arguments: show_alarms [description] [priority] [hostgroup] [unack]
			if len(args) > 0 {
				description = args[0]
				if len(args) > 1 {
					fmt.Sscanf(args[1], "%d", &priority)
				}
				if len(args) > 2 {
					hostGroups = append(hostGroups, args[2])
				}
				if len(args) > 3 {
					if args[3] == "*" {
						unacknowledged = false
					} else {
						fmt.Sscanf(args[3], "%t", &unacknowledged)
					}
				}
			}

			params := map[string]interface{}{
				"output":            "extend",
				"selectHosts":       "extend",
				"monitored":         true,
				"active":            true,
				"expandDescription": true,
				"filter": map[string]interface{}{
					"value": 1, // Problem state
				},
				"skipDependent": true,
			}

			if description != "" {
				params["search"] = map[string]interface{}{
					"description": description,
				}
			}
			if priority >= 0 {
				params["filter"].(map[string]interface{})["priority"] = priority
			}
			if len(hostGroups) > 0 {
				params["groupids"] = getHostGroupsIDs(client, hostGroups)
			}
			if unacknowledged {
				params["withLastEventUnacknowledged"] = true
			}

			result, err := client.Call("trigger.get", params)
			handleError(err)

			var triggers []map[string]interface{}
			json.Unmarshal(result, &triggers)

			headers := []string{"TriggerID", "Host", "Description", "Priority", "Last Change"}
			var rows [][]string
			for _, t := range triggers {
				hostName := ""
				if hosts, ok := t["hosts"].([]interface{}); ok && len(hosts) > 0 {
					hostName = hosts[0].(map[string]interface{})["name"].(string)
				}
				rows = append(rows, []string{
					fmt.Sprintf("%v", t["triggerid"]),
					hostName,
					fmt.Sprintf("%v", t["description"]),
					getPriorityName(t["priority"].(string)),
					fmt.Sprintf("%v", t["lastchange"]),
				})
			}

			outputResult(cmd, triggers, headers, rows)
		},
	}

	cmd.Flags().StringVar(&description, "description", "", "Filter by description")
	cmd.Flags().IntVar(&priority, "priority", -1, "Filter by priority")
	cmd.Flags().StringSliceVar(&hostGroups, "hostgroup", []string{}, "Filter by host group name(s)")
	cmd.Flags().BoolVar(&unacknowledged, "unack", true, "Show only unacknowledged alarms")

	return cmd
}

func newProblemAcknowledgeTriggerCmd() *cobra.Command {
	var message string
	var close bool

	cmd := &cobra.Command{
		Use:     "acknowledge_trigger [triggerid]",
		Aliases: []string{"acknowledge_trigger_last_event"},
		Short:   "Acknowledge the last event for one or more triggers",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			var eventIDs []string
			for _, arg := range args {
				// Handle comma-separated IDs
				tids := strings.Split(arg, ",")
				for _, tid := range tids {
					tid = strings.TrimSpace(tid)
					if tid == "" {
						continue
					}
					eid := getEventForTrigger(client, tid)
					if eid != "" {
						eventIDs = append(eventIDs, eid)
					}
				}
			}

			if len(eventIDs) == 0 {
				handleError(fmt.Errorf("no events found for triggers: %v", args))
				return
			}

			if message == "" {
				message = "[Zabbix-DNA] Acknowledged via CLI"
			}

			action := 2 // Acknowledge
			if close {
				action |= 1 // Close
			}

			params := map[string]interface{}{
				"eventids": eventIDs,
				"message":  message,
				"action":   action,
			}

			result, err := client.Call("event.acknowledge", params)
			handleError(err)

			var ackResult map[string]interface{}
			json.Unmarshal(result, &ackResult)

			outputResult(cmd, "Event(s) acknowledged successfully.", nil, nil)
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Acknowledgement message")
	cmd.Flags().BoolVarP(&close, "close", "c", false, "Close the event(s)")

	return cmd
}

func newProblemAcknowledgeCmd() *cobra.Command {
	var message string
	var close bool

	cmd := &cobra.Command{
		Use:   "acknowledge [eventid]",
		Short: "Acknowledge one or more events",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			var eventIDs []string
			for _, arg := range args {
				ids := strings.Split(arg, ",")
				for _, id := range ids {
					id = strings.TrimSpace(id)
					if id != "" {
						eventIDs = append(eventIDs, id)
					}
				}
			}

			if message == "" {
				message = "[Zabbix-DNA] Acknowledged via CLI"
			}

			action := 2 // Acknowledge
			if close {
				action |= 1 // Close
			}

			params := map[string]interface{}{
				"eventids": eventIDs,
				"message":  message,
				"action":   action,
			}

			result, err := client.Call("event.acknowledge", params)
			handleError(err)

			var ackResult map[string]interface{}
			json.Unmarshal(result, &ackResult)

			outputResult(cmd, "Event(s) acknowledged successfully.", nil, nil)
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Acknowledgement message")
	cmd.Flags().BoolVarP(&close, "close", "c", false, "Close the event(s)")

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
