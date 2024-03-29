package cmd

import (
	"github.com/prologic/box/internal"
	"github.com/spf13/cobra"
)

// NewExecCommand implements and returns the exec command.
func NewExecCommand() *cobra.Command {
	var detach bool
	cmd := &cobra.Command{
		Use:                   "exec [OPTIONS] CONTAINER COMMAND [ARG...]",
		Short:                 "Run a command inside a existing Container.",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.MinimumNArgs(2),
		PreRunE:               isRoot,
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.Exec(args[0], args[1:], detach)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&detach, "detach", "d", false, "run command in the background")

	return cmd
}
