package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"zabbix-dna/internal/config"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newREPLCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:     "shell",
		Aliases: []string{"repl"},
		Short:   "Start an interactive shell (Zabbix-CLI style)",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)

			promptStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF00")).
				Bold(true)

			// Intro Panel Style (zabbix-cli style)
			infoStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF00")).
				Bold(true)

			panelStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#00FF00")).
				Padding(0, 1).
				MarginBottom(1)

			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, _ := config.LoadConfig(cfgPath)
			serverURL := "unknown"
			if cfg != nil {
				serverURL = cfg.API.URL
			}

			welcomeMsg := fmt.Sprintf("Welcome to the Zabbix command-line interface (v1.0.7)\nConnected to server %s", serverURL)
			fmt.Println(panelStyle.Render(infoStyle.Render(welcomeMsg)))
			fmt.Println("Type --help to list commands, :h for REPL help, :q to exit.\n")

			for {
				fmt.Print(promptStyle.Render("zabbix-dna> "))
				input, err := reader.ReadString('\n')
				if err != nil {
					break
				}

				input = strings.TrimSpace(input)
				if input == "" {
					continue
				}

				// Internal Commands (zabbix-cli style)
				if strings.HasPrefix(input, ":") {
					handleInternalReplCommand(input, rootCmd)
					continue
				}

				// System Commands (zabbix-cli style)
				if strings.HasPrefix(input, "!") {
					handleSystemCommand(input[1:])
					continue
				}

				if input == "exit" || input == "quit" {
					break
				}

				args := strings.Fields(input)
				rootCmd.SetArgs(args)
				if err := rootCmd.Execute(); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
			}
		},
	}
}

func handleInternalReplCommand(input string, rootCmd *cobra.Command) {
	cmd := strings.TrimPrefix(input, ":")
	switch cmd {
	case "help", "h", "?":
		fmt.Println("\nInternal Commands:")
		fmt.Println("  :help, :h, :?    Show this help")
		fmt.Println("  :exit, :q        Exit the shell")
		fmt.Println("  :clear           Clear the screen")
		fmt.Println("\nSystem Commands:")
		fmt.Println("  ! <command>      Execute a system command (e.g., !ls -la)\n")
	case "exit", "q", "quit":
		os.Exit(0)
	case "clear":
		fmt.Print("\033[H\033[2J")
	default:
		fmt.Printf("Unknown internal command: %s\n", cmd)
	}
}

func handleSystemCommand(input string) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "System command failed: %v\n", err)
	}
}
