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

			// No prints here according to rule 6 (Shell Interativo must use same renderer)
			// But the prompt and banner are allowed as they are part of the UI, 
			// though the output of commands MUST be tabular.
			// The Execute() call will trigger the command's Run which uses outputResult.
			
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
				
				// Ensure that even in shell mode, the output is directed to our table renderer
				rootCmd.Execute()
			}
		},
	}
}


