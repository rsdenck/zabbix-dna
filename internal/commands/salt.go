//go:build cgo

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
	cmd.AddCommand(newSaltDeployAgentCmd())

	return cmd
}

func newSaltDeployAgentCmd() *cobra.Command {
	var target string
	var targetType string
	var osType string

	cmd := &cobra.Command{
		Use:   "deploy_agent",
		Short: "Deploy Dimo Zabbix Agent via SaltStack",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getSaltClient(cmd)
			handleError(err)

			fmt.Printf("Starting deployment for target: %s (%s)\n", target, osType)

			var commands []string
			if osType == "linux" {
				commands = []string{
					"file.mkdir /opt/dimo/",
					"cmd.run 'curl -o /opt/dimo/zabbix_agent_dimo https://repo.dimo.com/zabbix/agent_linux'",
					"cmd.run 'chmod +x /opt/dimo/zabbix_agent_dimo'",
					"service.restart zabbix-agent2",
				}
			} else {
				commands = []string{
					"file.mkdir C:\\Dimo\\",
					"cmd.run 'powershell Invoke-WebRequest -Uri https://repo.dimo.com/zabbix/agent_win -OutFile C:\\Dimo\\zabbix_agent_dimo.exe'",
					"service.restart zabbix-agent2",
				}
			}

			for _, saltCmd := range commands {
				jid := client.GetJid()
				fmt.Printf("Executing: %s (JID: %s)... ", saltCmd, jid)
				err = client.SendCommand(jid, target, targetType, saltCmd)
				if err != nil {
					fmt.Printf("FAILED: %v\n", err)
					handleError(fmt.Errorf("deployment failed at step: %s", saltCmd))
					return
				}
				fmt.Println("SUCCESS")
			}

			outputResult(cmd, "Dimo Agent deployed successfully via SaltStack.", nil, nil)
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "*", "Target minions")
	cmd.Flags().StringVarP(&targetType, "type", "T", "glob", "Target type")
	cmd.Flags().StringVarP(&osType, "os", "o", "linux", "Operating system (linux/windows)")

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
			handleError(err)

			jid := client.GetJid()

			err = client.SendCommand(jid, target, targetType, "test.ping")
			if err != nil {
				if strings.Contains(err.Error(), "root_key") {
					handleError(fmt.Errorf("SaltStack root_key not found. This command must be run on the Salt Master"))
				} else {
					handleError(err)
				}
				return
			}

			headers := []string{"Target", "JID", "Command", "Status"}
			rows := [][]string{
				{target, jid, "test.ping", "Published Successfully"},
			}
			outputResult(cmd, nil, headers, rows)
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
			handleError(err)

			jid := client.GetJid()

			err = client.SendCommand(jid, target, targetType, module)
			if err != nil {
				if strings.Contains(err.Error(), "root_key") {
					handleError(fmt.Errorf("SaltStack root_key not found. This command must be run on the Salt Master"))
				} else {
					handleError(err)
				}
				return
			}

			headers := []string{"Target", "JID", "Module", "Status"}
			rows := [][]string{
				{target, jid, module, "Published Successfully"},
			}
			outputResult(cmd, nil, headers, rows)
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
