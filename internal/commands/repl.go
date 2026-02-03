package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newREPLCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Start an interactive shell",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Zabbix-DNA Interactive Shell")
			fmt.Println("Type 'exit' or 'quit' to leave")

			for {
				fmt.Print("zabbix-dna> ")
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
