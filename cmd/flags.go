package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// Registers all global flags to utility itself.
func registerFlags(command string, cmd *cobra.Command) {
	switch strings.ToLower(command) {
	case "terragen":
		cmd.PersistentFlags().StringVarP(&genin.Path, "path", "p", ".", "path where the templates has to be generated")
		cmd.PersistentFlags().BoolVarP(&genin.DryRun, "dry-run", "", false, "dry-run the process of provider scaffold creation")
	case "create":
		// cmd.PersistentFlags().StringVarP(&genin.Provider, "name", "n", "demo", "name of the provider to create scaffolds")
	case "provider":
		cmd.PersistentFlags().StringSliceVarP(&genin.DataSource, "data-source", "d", nil, "name of the data scaffold")
		cmd.PersistentFlags().StringSliceVarP(&genin.Resource, "resource", "r", nil, "name of the resource scaffold")
		cmd.PersistentFlags().StringVarP(&genin.Importer, "importer", "i", "", "name of the importer scaffold")
		cmd.PersistentFlags().StringVarP(&genin.RepoGroup, "repo-group", "g", "", "repo group to which the terraform provider to be part of")
		cmd.PersistentFlags().BoolVarP(&genin.ResourceRequired, "resource-required", "", false, "enable if resource requires scaffold")
		cmd.PersistentFlags().BoolVarP(&genin.ImporterRequired, "importer-required", "", false, "enable if importer requires scaffold")
		cmd.PersistentFlags().BoolVarP(&genin.DatasourceRequired, "datasource-required", "", false, "enable if data_source requires scaffold")
	case "datasource", "resource":
		cmd.PersistentFlags().StringVarP(&genin.Provider, "provider", "", "demo", "name of the provider for which resource/datasource to be scaffolded")
	}
}
