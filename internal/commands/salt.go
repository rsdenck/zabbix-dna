package commands

import (
	"fmt"
	"strings"

	"zabbix-dna/internal/config"

	"github.com/spf13/cobra"
	salt "github.com/tsaridas/salt-golang/lib/client"
)

func newSaltCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "salt",
		Short: "SaltStack integration for Zabbix Proxies",
		Long:  `Manage Zabbix Proxies and other infrastructure using SaltStack.`,
	}

	cmd.AddCommand(newSaltPingCmd())
	cmd.AddCommand(newSaltRunCmd())

	return cmd
}

func newSaltPingCmd() *cobra.Command {
	var target string
	var targetType string

	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping minions (proxies)",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getSaltClient(cmd)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			jid := client.GetJid()

			headers := []string{"Property", "Value"}
			rows := [][]string{
				{"Target", target},
				{"Target Type", targetType},
				{"JID", jid},
				{"Command", "test.ping"},
			}
			outputResult(cmd, nil, headers, rows)

			err = client.SendCommand(jid, target, targetType, "test.ping")
			if err != nil {
				if strings.Contains(err.Error(), "root_key") {
					fmt.Println("\nError: SaltStack root_key not found. This command must be run on the Salt Master.")
				} else {
					fmt.Printf("\nError sending command: %v\n", err)
				}
				return
			}

			fmt.Println("\nCommand published successfully. Check Salt events for returns.")
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "*", "Target minions")
	cmd.Flags().StringVarP(&targetType, "type", "T", "glob", "Target type (glob, list, pcre)")

	return cmd
}

func newSaltRunCmd() *cobra.Command {
	var target string
	var targetType string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a module on minions",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			module := args[0]
			client, err := getSaltClient(cmd)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			jid := client.GetJid()

			headers := []string{"Property", "Value"}
			rows := [][]string{
				{"Target", target},
				{"Target Type", targetType},
				{"JID", jid},
				{"Module", module},
			}
			outputResult(cmd, nil, headers, rows)

			err = client.SendCommand(jid, target, targetType, module)
			if err != nil {
				if strings.Contains(err.Error(), "root_key") {
					fmt.Println("\nError: SaltStack root_key not found. This command must be run on the Salt Master.")
				} else {
					fmt.Printf("\nError running command: %v\n", err)
				}
				return
			}

			fmt.Println("\nCommand published successfully. Check Salt events for returns.")
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "*", "Target minions")
	cmd.Flags().StringVarP(&targetType, "type", "T", "glob", "Target type (glob, list, pcre)")

	return cmd
}

func getSaltClient(cmd *cobra.Command) (*salt.Client, error) {
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		// If config fails, try default values
		return &salt.Client{
			Server: "tcp://127.0.0.1:4506",
		}, nil
	}

	server := cfg.Salt.URL
	if server == "" {
		server = "tcp://127.0.0.1:4506"
	}

	return &salt.Client{
		Server:  server,
		Verbose: false,
	}, nil
}
