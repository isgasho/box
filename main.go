package main

import (
	"github.com/prologic/box/cmd"
)

func main() {
	rootCmd := cmd.NewBoxCommand()
	rootCmd.AddCommand(cmd.NewRunCommand())
	rootCmd.AddCommand(cmd.NewForkCommand())
	rootCmd.AddCommand(cmd.NewExecCommand())
	rootCmd.AddCommand(cmd.NewPsCommand())
	rootCmd.AddCommand(cmd.NewImagesCommand())
	rootCmd.AddCommand(cmd.NewVersionCommand())
	rootCmd.Execute()
}
