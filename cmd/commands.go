package cmd

import "github.com/spf13/cobra"

func initCommands() {
	createRegister("provider", true, &cobra.Command{
		Use:   "provider [args] [flags]",
		Short: "Command to generate scaffolds for terraform provider",
		Long: `Command to generate scaffolds for terraform provider and its other components.
               This includes creation of provider, resource, datasource.`,
		Run:          genin.CreateProvider,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})

	createRegister("datasource", true, &cobra.Command{
		Use:          "datasource [args] [flags]",
		Short:        "Command to generate scaffolds for datasource",
		Long:         `This will help user to generate scaffolds for datasource for the chosen provider.`,
		Run:          genin.CreateDataSource,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})

	createRegister("resource", true, &cobra.Command{
		Use:          "resource [args] [flags]",
		Short:        "Command to generate scaffolds for resource",
		Long:         `This will help user to generate scaffolds for resource for the chosen provider.`,
		Run:          genin.CreateResource,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})

	editRegister("provider", false, &cobra.Command{
		Use:   "provider [args] [flags]",
		Short: "Command to edit already generated scaffolds of a provider",
		Long: `Command to edit scaffolds of terraform provider that was already generated. 
               Not all aspects of provider can be edited, it is very limited`,
		Run:          genin.CreateProvider,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})

	editRegister("datasource", false, &cobra.Command{
		Use:   "datasource [args] [flags]",
		Short: "Command to edit already generated scaffolds of a datasource",
		Long: `This will help user to edit scaffolds of datasource that was already generated.
               Not all aspects of datasource can be edited, it is very limited`,
		Run:          genin.CreateDataSource,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})

	editRegister("resource", false, &cobra.Command{
		Use:   "resource [args] [flags]",
		Short: "Command to edit already created scaffolds generated scaffolds of resource",
		Long: `This will help user to edit scaffolds of resource that was already generated.
               Not all aspects of resource can be edited, it is very limited`,
		Run:          genin.CreateResource,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})
}
