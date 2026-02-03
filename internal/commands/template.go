package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage Zabbix templates",
	}

	cmd.AddCommand(newTemplateListCmd())
	cmd.AddCommand(newTemplateShowCmd())
	cmd.AddCommand(newTemplateDeleteCmd())

	return cmd
}

func newTemplateListCmd() *cobra.Command {
	var limit int
	var search string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Zabbix templates",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"output": []string{"templateid", "host", "name"},
				"limit":  limit,
			}
			if search != "" {
				params["search"] = map[string]interface{}{
					"host": search,
				}
			}

			result, err := client.Call("template.get", params)
			handleError(err)

			var templates []map[string]interface{}
			json.Unmarshal(result, &templates)

			headers := []string{"TemplateID", "Host", "Name"}
			var rows [][]string
			for _, t := range templates {
				rows = append(rows, []string{
					fmt.Sprintf("%v", t["templateid"]),
					fmt.Sprintf("%v", t["host"]),
					fmt.Sprintf("%v", t["name"]),
				})
			}

			outputResult(cmd, templates, headers, rows)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Limit the number of templates")
	cmd.Flags().StringVarP(&search, "search", "s", "", "Search for a template by name")

	return cmd
}

func newTemplateShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show [template name]",
		Short: "Show details of a Zabbix template",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			params := map[string]interface{}{
				"filter": map[string]interface{}{
					"host": args[0],
				},
				"selectGroups": "extend",
				"selectItems":  "count",
			}

			result, err := client.Call("template.get", params)
			handleError(err)

			var templates []map[string]interface{}
			json.Unmarshal(result, &templates)

			if len(templates) == 0 {
				fmt.Printf("Template not found: %s\n", args[0])
				return
			}

			t := templates[0]
			fmt.Printf("ID:         %s\n", t["templateid"])
			fmt.Printf("Host:       %s\n", t["host"])
			fmt.Printf("Name:       %s\n", t["name"])
			fmt.Printf("Items:      %s\n", t["items"])

			fmt.Println("\nGroups:")
			if groups, ok := t["groups"].([]interface{}); ok {
				for _, g := range groups {
					group := g.(map[string]interface{})
					fmt.Printf("- %s (%s)\n", group["name"], group["groupid"])
				}
			}
		},
	}
}

func newTemplateDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [template name]",
		Short: "Delete a Zabbix template",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			// First find the template ID
			params := map[string]interface{}{
				"filter": map[string]interface{}{
					"host": args[0],
				},
			}

			result, err := client.Call("template.get", params)
			handleError(err)

			var templates []map[string]interface{}
			json.Unmarshal(result, &templates)

			if len(templates) == 0 {
				fmt.Printf("Template not found: %s\n", args[0])
				return
			}

			templateID := templates[0]["templateid"].(string)

			// Delete the template
			_, err = client.Call("template.delete", []string{templateID})
			handleError(err)

			fmt.Printf("Template %s (ID: %s) deleted successfully\n", args[0], templateID)
		},
	}
}
