// Package cli will initialize cli of terragen.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cmd *cobra.Command

func init() {
	initCommands()
	cmd = SetTerragenCommands()
}

// Main will take the workload of executing/starting the cli, when the command is passed to it.
func Main() {
	if err := Execute(os.Args[1:]); err != nil {
		cm.NeuronSaysItsError(err.Error())
		os.Exit(1)
	}
}

// Execute will actually execute the cli by taking the arguments passed to cli.
func Execute(args []string) error {
	cmd.SetArgs(args)
	if _, err := cmd.ExecuteC(); err != nil {
		return err
	}

	return nil
}
