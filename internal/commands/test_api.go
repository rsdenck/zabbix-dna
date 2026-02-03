package commands

import (
	"fmt"
	"os"

	"zabbix-dna/internal/api"
	"zabbix-dna/internal/config"
	"github.com/spf13/cobra"
)

func newTestAPICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test-api",
		Short: "Validate connection to Zabbix API",
		Run: func(cmd *cobra.Command, args []string) {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.LoadConfig(cfgPath)
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				fmt.Println("Attempting with default environment values...")
				// Fallback or manual entry could go here
				return
			}

			client := api.NewClient(cfg.Zabbix.URL, cfg.Zabbix.Token, cfg.Zabbix.Timeout)

			if cfg.Zabbix.Token == "" && cfg.Zabbix.User != "" {
				fmt.Printf("Authenticating as %s...\n", cfg.Zabbix.User)
				err := client.Login(cfg.Zabbix.User, cfg.Zabbix.Password)
				if err != nil {
					fmt.Printf("Authentication Failed: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("Authenticated successfully!")
			}

			// Test with apiinfo.version
			result, err := client.Call("apiinfo.version", map[string]interface{}{})
			if err != nil {
				fmt.Printf("API Connection Failed: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Successfully connected to Zabbix API!\nVersion: %s\n", string(result))
		},
	}
}
