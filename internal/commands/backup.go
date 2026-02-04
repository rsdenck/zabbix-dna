package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"zabbix-dna/internal/api"
	"zabbix-dna/internal/config"

	"github.com/spf13/cobra"
)

func newBackupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Perform a backup of Zabbix configurations",
		Run: func(cmd *cobra.Command, args []string) {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.LoadConfig(cfgPath)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			client := api.NewClient(cfg.Zabbix.URL, cfg.Zabbix.Token, cfg.Zabbix.Timeout)

			if cfg.Zabbix.Token == "" && cfg.Zabbix.User != "" {
				fmt.Printf("Authenticating as %s...\n", cfg.Zabbix.User)
				err := client.Login(cfg.Zabbix.User, cfg.Zabbix.Password)
				if err != nil {
					fmt.Printf("Authentication Failed: %v\n", err)
					return
				}
				fmt.Println("Authenticated successfully!")
			}

			fmt.Println("Starting Zabbix backup (Version 7.0 compatibility)...")

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
				fmt.Printf("Backup failed: %v\n", err)
				fmt.Println("Attempting fallback to basic objects...")

				// Fallback for older versions or simplified export
				params["options"] = map[string]interface{}{
					"groups":    []string{},
					"hosts":     []string{},
					"templates": []string{},
				}
				result, err = client.Call("configuration.export", params)
				if err != nil {
					fmt.Printf("Fallback backup also failed: %v\n", err)
					return
				}
			}

			// Pretty print JSON
			var prettyJSON json.RawMessage
			if err := json.Unmarshal(result, &prettyJSON); err == nil {
				formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
				result = formatted
			}

			filename := fmt.Sprintf("zabbix_backup_%s.json", time.Now().Format("20060102_150405"))
			err = os.WriteFile(filename, result, 0644)
			if err != nil {
				fmt.Printf("Failed to save backup file: %v\n", err)
				return
			}

			headers := []string{"Component", "Status"}
			rows := [][]string{
				{"Zabbix API", "Connected"},
				{"Export Format", "JSON"},
				{"Configuration", "Full Export"},
				{"File Path", filename},
			}
			outputResult(cmd, map[string]string{"filename": filename}, headers, rows)
		},
	}
}
