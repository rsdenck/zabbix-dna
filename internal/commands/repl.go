package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newREPLCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Start an interactive shell",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)

			bannerStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#D20000")).
				Padding(0, 1).
				Bold(true)

			promptStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D20000")).
				Bold(true)

			fmt.Println(bannerStyle.Render("Zabbix-DNA Interactive Shell"))
			fmt.Println("Type 'exit' or 'quit' to leave")

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

				if input == "exit" || input == "quit" {
					break
				}

				args := strings.Split(input, " ")
				rootCmd.SetArgs(args)
				rootCmd.Execute()
			}
		},
	}
}
