package cmd

import (
	"github.com/prologic/box/internal"
	"github.com/spf13/cobra"
)

// NewPsCommand implements and returns the ps command.
func NewPsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ps",
		Short:                 "List Containers",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.NoArgs,
		PreRunE:               isRoot,
		RunE:                  internal.Ps,
	}

	return cmd
}
