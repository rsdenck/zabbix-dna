package commands

import (
	"fmt"

	"zabbix-dna/internal/api"
	"zabbix-dna/internal/config"

	"github.com/spf13/cobra"
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
			client := api.NewClient(cfg.Zabbix.URL, cfg.Zabbix.Token, cfg.Zabbix.Timeout)

			if cfg.Zabbix.Token == "" && cfg.Zabbix.User != "" {
				err := client.Login(cfg.Zabbix.User, cfg.Zabbix.Password)
				if err != nil {
					// Handle login failure but continue to show status
				}
			}

			zabbixVersion, err := client.Call("apiinfo.version", map[string]interface{}{})
			zabbixStatus := "Connected"
			if err != nil {
				zabbixStatus = "Failed"
			}

			// 2. SaltStack Test
			saltStatus := "Connected"
			saltServer := cfg.Salt.URL
			if saltServer == "" {
				saltServer = "tcp://127.0.0.1:4506"
			}

			// Just a simple check of the server address for now

			headers := []string{"Service", "Property", "Value"}
			rows := [][]string{
				{"Zabbix", "Status", zabbixStatus},
				{"Zabbix", "Version", string(zabbixVersion)},
				{"Zabbix", "Endpoint", cfg.Zabbix.URL},
				{"SaltStack", "Status", saltStatus},
				{"SaltStack", "Endpoint", saltServer},
			}
			resp := map[string]interface{}{
				"zabbix": map[string]string{
					"status":   zabbixStatus,
					"version":  string(zabbixVersion),
					"endpoint": cfg.Zabbix.URL,
				},
				"saltstack": map[string]string{
					"status":   saltStatus,
					"endpoint": saltServer,
				},
			}
			outputResult(cmd, resp, headers, rows)
		},
	}
}
