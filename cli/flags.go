package cli

import (
	"github.com/spf13/cobra"
)

// Registering all the flags to the command neuron itself.
func registerFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&genin.Package, "name", "n", "", "name of the provider that needs templates")
	cmd.PersistentFlags().StringVarP(&genin.Path, "path", "p", "", "path where the templates has to be generated")
}
