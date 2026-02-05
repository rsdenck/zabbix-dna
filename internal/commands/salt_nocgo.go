//go:build !cgo

package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newSaltCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "salt",
		Short: "SaltStack integration (disabled in this build)",
		Long:  `SaltStack integration requires CGO and is disabled in this binary.`,
		Run: func(cmd *cobra.Command, args []string) {
			handleError(fmt.Errorf("SaltStack support was not included in this build (CGO disabled)"))
		},
	}
}
