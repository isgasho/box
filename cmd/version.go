package cmd

import (
	"github.com/spf13/cobra"

	"github.com/prologic/box/internal"
)

// NewVersionCommand implements and returns the version command.
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		Short:                 "Display the version of box and exit",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.NoArgs,
		RunE:                  internal.DisplayVersion,
	}

	return cmd
}
