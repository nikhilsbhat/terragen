// Package cli will initialize cli for neuron.
package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	cmd *cobra.Command
)

func init() {
	cmd = SetTerragenCmds()
}

// Main will take the workload of executing/starting the cli, when the command is passed to it.
func Main() {
	err := Execute(os.Args[1:])
	if err != nil {
		cm.NeuronSaysItsError(err.Error())
		os.Exit(1)
	}
}

// Execute will actually execute the cli by taking the arguments passed to cli.
func Execute(args []string) error {

	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	if err != nil {
		return err
	}
	return nil
}
