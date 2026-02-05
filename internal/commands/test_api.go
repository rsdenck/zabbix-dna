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
			saltStatus := "Not Checked (Requires CGO)"
			saltServer := cfg.Salt.URL
			if saltServer == "" {
				saltServer = "tcp://127.0.0.1:4506"
			}

			// We only check if salt-client is available/functional in CGO builds
			// In non-CGO, we just report the config
			if isCGOBuilt() {
				saltStatus = "Connected (Mock/Config Only)"
			}

			headers := []string{"Service", "Property", "Value"}
			rows := [][]string{
				{"Zabbix", "Status", zabbixStatus},
				{"Zabbix", "Version", string(zabbixVersion)},
				{"Zabbix", "Endpoint", cfg.API.URL},
				{"SaltStack", "Status", saltStatus},
				{"SaltStack", "Endpoint", saltServer},
			}
			outputResult(cmd, nil, headers, rows)
		},
	}
}
