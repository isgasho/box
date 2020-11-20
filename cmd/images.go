package cmd

import (
	"github.com/spf13/cobra"

	"github.com/prologic/box/internal"
)

// NewImagesCommand implements and returns the images command.
func NewImagesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "images",
		Short:                 "List local images",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.NoArgs,
		PreRunE:               isRoot,
		RunE:                  internal.Images,
	}

	return cmd
}
