package cmd

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	imagesPath     = "/var/lib/box/images"
	containersPath = "/var/lib/box/containers"
	netnsPath      = "/var/lib/box/netns"
)

var ErrNotPermitted = errors.New("operation not permitted")

// Make box directories first.
func init() {
	os.MkdirAll(netnsPath, 0700)
	os.MkdirAll(imagesPath, 0700)
	os.MkdirAll(containersPath, 0700)
}

// NewBoxCommand returns the root cobra.Command for box.
func NewBoxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "box [OPTIONS] COMMAND",
		Short:                 "A tiny tool for managing containers and sandbox processes",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// set logging level
			isDebug, _ := cmd.Flags().GetBool("debug")
			if isDebug {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}
		},
	}

	cmd.PersistentFlags().BoolP(
		"debug", "D", false,
		"Enable debug logging",
	)

	return cmd
}

// isRoot implements a cobra acceptable function and
// returns ErrNotPermitted if user is not root.
func isRoot(_ *cobra.Command, _ []string) error {
	if os.Getuid() != 0 {
		return ErrNotPermitted
	}
	return nil
}
