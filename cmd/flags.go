package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// Registers all global flags to utility itself.
func registerFlags(command string, cmd *cobra.Command) {
	switch strings.ToLower(command) {
	case "terragen":
		cmd.PersistentFlags().StringVarP(&generateInput.Path, "path", "p", ".", "path where the templates has to be generated")
		cmd.PersistentFlags().BoolVarP(&generateInput.DryRun, "dry-run", "", false, "dry-run the process of provider scaffold creation")
		cmd.PersistentFlags().BoolVarP(&generateInput.Force, "force", "f",
			false, "enable this to forcefully create resource/datasource/importers (this might tamper the scaffold)")
		cmd.PersistentFlags().BoolVarP(&generateInput.TerraformPluginFramework, "use-plugin-framework", "",
			false, "enable this to generate scaffolds with terraform-plugin-framework(https://github.com/hashicorp/terraform-plugin-framework)")
		cmd.PersistentFlags().BoolVarP(&generateInput.SkipProviderUpdate, "skip-provider-update", "", false,
			"when enabled, updating provider.go with newly created datasource/resource would be skipped")
		cmd.PersistentFlags().StringVarP(&cliLogLevel, "log-level", "l", "info",
			"log level for terragen, log levels supported by [https://github.com/sirupsen/logrus] will work")
	case "create":
		// cmd.PersistentFlags().StringVarP(&generateInput.Provider, "name", "n", "demo", "name of the provider to create scaffolds")
	case "provider":
		cmd.PersistentFlags().StringSliceVarP(&generateInput.DataSource, "data-source", "d", nil, "name of the data scaffold")
		cmd.PersistentFlags().StringSliceVarP(&generateInput.Resource, "resource", "r", nil, "name of the resource scaffold")
		cmd.PersistentFlags().StringVarP(&generateInput.Importer, "importer", "i", "", "name of the importer scaffold")
		cmd.PersistentFlags().StringVarP(&generateInput.RepoGroup, "repo-group", "g", "",
			"repo group to which the terraform provider to be part of")
		cmd.PersistentFlags().BoolVarP(&generateInput.ResourceRequired, "resource-required", "", false,
			"enable if resource requires scaffold")
		cmd.PersistentFlags().BoolVarP(&generateInput.ImporterRequired, "importer-required", "", false,
			"enable if importer requires scaffold")
		cmd.PersistentFlags().BoolVarP(&generateInput.DatasourceRequired, "datasource-required", "", false,
			"enable if data_source requires scaffold")
		cmd.PersistentFlags().BoolVarP(&generateInput.SkipValidation, "skip-validation", "",
			false, "enable if prerequisite validation needs to be skipped")
	case "datasource", "resource":
		cmd.PersistentFlags().StringVarP(&generateInput.Provider, "provider", "", "demo",
			"name of the provider for which resource/datasource to be scaffolded")
	}
}
