package cmd

import "github.com/spf13/cobra"

func initCommands() {
	createRegister("provider", true, &cobra.Command{
		Use:   "provider [flags]",
		Short: "Command to generate scaffolds for terraform provider",
		Long: `command to generate scaffolds for terraform provider and its other components,
                       this includes creation of provider, resource, datasource.`,
		Run:          genin.CreateProvider,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})

	createRegister("datasource", true, &cobra.Command{
		Use:          "datasource [flags]",
		Short:        "Command to generate scaffolds for datasource",
		Long:         `This will help user to generate scaffolds for datasource for the chosen provider.`,
		Run:          genin.CreateDataSource,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})

	createRegister("resource", true, &cobra.Command{
		Use:          "resource [flags]",
		Short:        "Command to generate scaffolds for resource",
		Long:         `This will help user to generate scaffolds for resource for the chosen provider.`,
		Run:          genin.CreateResource,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})
}
