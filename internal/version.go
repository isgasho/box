package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version release version
	Version = "0.0.1"

	// Build will be overwritten automatically by the build system
	Build = "dev"

	// GitCommit will be overwritten automatically by the build system
	GitCommit = "HEAD"
)

// FullVersion returns the full version, build and commit hash
func FullVersion() string {
	return fmt.Sprintf("%s-%s@%s", Version, Build, GitCommit)
}

// DisplayVersion displays the version of box and exits
func DisplayVersion(_ *cobra.Command, _ []string) error {
	fmt.Printf("box v%s\n", FullVersion())
	return nil
}
