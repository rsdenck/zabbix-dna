package commands

import (
	"zabbix-dna/internal/config"

	"github.com/spf13/cobra"
)

func newTestAPICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test-api",
		Short: "Validate connection to Zabbix API and SaltStack",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getZabbixClient(cmd)
			handleError(err)

			zabbixVersion, err := client.Call("apiinfo.version", map[string]interface{}{})
			zabbixStatus := "Connected"
			if err != nil {
				zabbixStatus = "Failed"
			}

			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, _ := config.LoadConfig(cfgPath)

			// 2. SaltStack Test
			saltStatus := "Connected"
			saltServer := cfg.Salt.URL
			if saltServer == "" {
				saltServer = "tcp://127.0.0.1:4506"
			}

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
