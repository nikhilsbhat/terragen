package cmd

import "github.com/spf13/cobra"

func initCommands() {
	createRegister("provider", true, &cobra.Command{
		Use:   "provider [args] [flags]",
		Short: "Command to generate scaffolds for terraform provider",
		Long: `Command to generate scaffolds for terraform provider and its other components.
               This includes creation of provider, resource, datasource.`,
		PreRunE: setInputLogger,
		Args:    cobra.RangeArgs(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateInput.Generate(args[0])
		},
	})

	createRegister("datasource", true, &cobra.Command{
		Use:     "datasource [args] [flags]",
		Short:   "Command to generate scaffolds for datasource",
		Long:    `This will help user to generate scaffolds for datasource of chosen provider.`,
		PreRunE: setInputLogger,
		Args:    cobra.RangeArgs(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateInput.GenerateDataSource(args)
		},
	})

	createRegister("resource", true, &cobra.Command{
		Use:     "resource [args] [flags]",
		Short:   "Command to generate scaffolds for resource",
		Long:    `This will help user to generate scaffolds for resource of chosen provider.`,
		PreRunE: setInputLogger,
		Args:    cobra.RangeArgs(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateInput.GenerateResource(args)
		},
	})

	editRegister("provider", false, &cobra.Command{
		Use:   "provider [args] [flags]",
		Short: "Command to edit already generated scaffolds of a provider",
		Long: `Command to edit scaffolds of terraform provider that was already generated. 
               Not all aspects of provider can be edited, it is very limited`,
		PreRunE: setInputLogger,
		Args:    cobra.RangeArgs(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateInput.Generate(args[0])
		},
	})

	editRegister("datasource", false, &cobra.Command{
		Use:   "datasource [args] [flags]",
		Short: "Command to edit already generated scaffolds of a datasource",
		Long: `This will help user to edit scaffolds of datasource that was already generated.
               Not all aspects of datasource can be edited, it is very limited`,
		PreRunE: setInputLogger,
		Args:    cobra.RangeArgs(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateInput.GenerateDataSource(args)
		},
	})

	editRegister("resource", false, &cobra.Command{
		Use:   "resource [args] [flags]",
		Short: "Command to edit already created scaffolds generated scaffolds of resource",
		Long: `This will help user to edit scaffolds of resource that was already generated.
               Not all aspects of resource can be edited, it is very limited`,
		Args: cobra.RangeArgs(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateInput.GenerateResource(args)
		},
	})
}

func setInputLogger(_ *cobra.Command, _ []string) error {
	InitLogger(cliLogLevel)
	generateInput.SetLogger(cliLogger)

	return nil
}
