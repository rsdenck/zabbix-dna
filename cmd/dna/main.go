package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"zabbix-dna/internal/commands"
	"zabbix-dna/internal/tui"
	"zabbix-dna/internal/wizard"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	bannerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#D20000")).
		Padding(1, 2).
		Bold(true).
		MarginBottom(1)
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "zabbix-dna",
		Short: "Zabbix CLI | Enterprise Observability",
		Long: `ZABBIX-DNA is a high-performance CLI for Zabbix, 
written 100% in Go with a focus on observability and automation.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			batchFile, _ := cmd.Flags().GetString("batch")
			if batchFile != "" {
				runBatch(cmd, batchFile)
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println(bannerStyle.Render("ZABBIX-DNA CLI | v1.0.0"))
				choice, _ := tui.Start()
				if choice != "" {
					fmt.Printf("\n> Executando: %s\n\n", choice)
					args := strings.Fields(choice)
					cmd.SetArgs(args)
					cmd.Execute()
				}
			}
		},
	}

	// Wizard command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "wizard",
		Short: "Start the configuration wizard",
		Run: func(cmd *cobra.Command, args []string) {
			wizard.Start()
		},
	})

	// TUI command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gui",
		Short: "Start the TUI interface",
		Run: func(cmd *cobra.Command, args []string) {
			choice, _ := tui.Start()
			if choice != "" {
				fmt.Printf("\n> Executando: %s\n\n", choice)
				args := strings.Fields(choice)
				cmd.Root().SetArgs(args)
				cmd.Root().Execute()
			}
		},
	})

	// Persistent flags
	rootCmd.PersistentFlags().StringP("config", "c", "zabbix-dna.toml", "config file")
	rootCmd.PersistentFlags().StringP("format", "f", "table", "output format (table, json)")
	rootCmd.PersistentFlags().String("batch", "", "run commands from a file")

	// Add commands
	commands.AddCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runBatch(rootCmd *cobra.Command, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening batch file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		args := strings.Fields(line)
		rootCmd.SetArgs(args)
		if err := rootCmd.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing command '%s': %v\n", line, err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading batch file: %v\n", err)
	}
}
