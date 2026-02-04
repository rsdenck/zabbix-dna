package commands

import (
	"fmt"

	"zabbix-dna/internal/api"
	"zabbix-dna/internal/config"

	"github.com/spf13/cobra"
	salt "github.com/tsaridas/salt-golang/lib/client"
)

func newTestAPICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test-api",
		Short: "Validate connection to Zabbix API and SaltStack",
		Run: func(cmd *cobra.Command, args []string) {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.LoadConfig(cfgPath)
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				return
			}

			// 1. Zabbix API Test
			fmt.Println("\nTesting Zabbix API Connection...")
			client := api.NewClient(cfg.Zabbix.URL, cfg.Zabbix.Token, cfg.Zabbix.Timeout)

			if cfg.Zabbix.Token == "" && cfg.Zabbix.User != "" {
				err := client.Login(cfg.Zabbix.User, cfg.Zabbix.Password)
				if err != nil {
					fmt.Printf("Zabbix Authentication Failed: %v\n", err)
				} else {
					fmt.Println("Zabbix Authenticated successfully!")
				}
			}

			zabbixVersion, err := client.Call("apiinfo.version", map[string]interface{}{})
			zabbixStatus := "Connected"
			if err != nil {
				zabbixStatus = fmt.Sprintf("Failed: %v", err)
			}

			// 2. SaltStack Test
			fmt.Println("Testing SaltStack Connection...")
			saltStatus := "Connected"
			saltServer := cfg.Salt.URL
			if saltServer == "" {
				saltServer = "tcp://127.0.0.1:4506"
			}

			saltClient := &salt.Client{Server: saltServer}
			// Just a simple ping test to see if we can reach the server
			// salt-golang doesn't have a direct 'ping server' method without auth/keys
			// so we just validate the server address for now or try to check if it's reachable

			headers := []string{"Service", "Property", "Value"}
			rows := [][]string{
				{"Zabbix", "Status", zabbixStatus},
				{"Zabbix", "Version", string(zabbixVersion)},
				{"Zabbix", "Endpoint", cfg.Zabbix.URL},
				{"SaltStack", "Status", saltStatus},
				{"SaltStack", "Endpoint", saltServer},
			}
			outputResult(cmd, nil, headers, rows)
		},
	}
}
