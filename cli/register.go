package cli

import (
	"fmt"

	gen "github.com/nikhilsbhat/terragen/gen"
	"github.com/nikhilsbhat/terragen/version"
	"github.com/spf13/cobra"
)

var (
	cmds  map[string]*cobra.Command
	genin gen.GenInput
)

type confcmds struct {
	commands []*cobra.Command
}

// SetTerragenCmds helps in gathering all the subcommands so that it can be used while registering it with main command.
func SetTerragenCmds() *cobra.Command {
	cmd := getTerragenCmds()
	return cmd
}

func getTerragenCmds() *cobra.Command {

	var terragenCmd = &cobra.Command{
		Use:   "terragen [command]",
		Short: "command to create files/folder for terraform provider",
		Long:  `Terragen helps user to create custom terraform provider by generating templates for it.`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  cm.echoTerragen,
	}
	terragenCmd.SetUsageTemplate(getUsageTemplate())

	var setCmd = &cobra.Command{
		Use:          "generate [flags]",
		Short:        "command to generate the initial components for terraform provider",
		Long:         `This will help user to generate the initial components of terraform provider.`,
		Run:          genin.Generate,
		SilenceUsage: true,
	}

	// fetching "version" will be done here.
	var versionCmd = &cobra.Command{
		Use:   "version [flags]",
		Short: "command to fetch the version of terragen installed",
		Long:  `This will help user to find what version of terragen he/she installed in her machine.`,
		RunE:  versionConfig,
	}

	terragenCmd.AddCommand(setCmd)
	terragenCmd.AddCommand(versionCmd)
	registerFlags(terragenCmd)
	return terragenCmd
}

func (cm *cliMeta) echoTerragen(cmd *cobra.Command, args []string) error {
	cmd.Usage()
	return nil
}

func versionConfig(cmd *cobra.Command, args []string) error {
	fmt.Println("terragen", version.GetVersion())
	return nil
}

// This function will return the custom template for usage function,
// only functions/methods inside this package can call this.

func getUsageTemplate() string {
	return `{{printf "\n"}}Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if gt (len .Aliases) 0}}{{printf "\n" }}
Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}{{printf "\n" }}
Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{printf "\n"}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}{{printf "\n"}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}{{printf "\n"}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}{{printf "\n"}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}{{printf "\n"}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}"
{{printf "\n"}}`
}
