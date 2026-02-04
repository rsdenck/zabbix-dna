//go:build nosalt
// +build nosalt

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSaltCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "salt",
		Short: "SaltStack integration (disabled in this build)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("SaltStack integration is disabled in this build because of missing CGO/ZeroMQ dependencies.")
		},
	}
}
