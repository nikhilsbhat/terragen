package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nikhilsbhat/terragen/pkg/gen"
	"github.com/nikhilsbhat/terragen/version"
	"github.com/spf13/cobra"
)

var (
	createCommands map[string]*cobra.Command
	editCommands   map[string]*cobra.Command
	generateInput  gen.Input
)

type terragenCommands struct {
	commands []*cobra.Command
}

// SetTerragenCommands helps in gathering all the subcommands so that it can be used while registering it with main command.
func SetTerragenCommands() *cobra.Command {
	return getTerragenCommands()
}

// Add an entry in below function to register new command.
func getTerragenCommands() *cobra.Command {
	command := new(terragenCommands)
	command.commands = append(command.commands, getCreateCommand())
	command.commands = append(command.commands, getVersionCommand())
	command.commands = append(command.commands, getEditCommand())

	return command.prepareCommands()
}

func (c *terragenCommands) prepareCommands() *cobra.Command {
	rootCmd := getRootCommand()
	for _, cmnd := range c.commands {
		rootCmd.AddCommand(cmnd)
	}
	registerFlags("terragen", rootCmd)

	return rootCmd
}

func getRootCommand() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:     "terragen [command]",
		Short:   "Utility that helps in generating scaffolds for terraform provider",
		Long:    `Terragen helps user to create custom terraform provider and its components by generating scaffolds.`,
		PreRunE: setCLI,
		RunE:    echoTerragen,
	}
	rootCommand.SetUsageTemplate(getUsageTemplate())

	return rootCommand
}

func getCreateCommand() *cobra.Command {
	createCommand := &cobra.Command{
		Use:     "create [command] [flags]",
		Short:   "Command to scaffold provider and other components of terraform provider",
		Long:    `This will help user to generate the initial components of terraform provider.`,
		PreRunE: setCLI,
		RunE:    echoTerragen,
	}
	registerFlags("create", createCommand)
	for _, command := range createCommands {
		createCommand.AddCommand(command)
	}

	return createCommand
}

func getEditCommand() *cobra.Command {
	editCommand := &cobra.Command{
		Use:     "edit [command] [flags]",
		Short:   "Command to edit the scaffold created for a provider",
		Long:    `This will help user to edit the scaffolds generated for terraform provider and other components of them.`,
		PreRunE: setCLI,
		RunE:    echoTerragen,
	}
	registerFlags("edit", editCommand)
	for _, command := range editCommands {
		command.SilenceUsage = true
		editCommand.AddCommand(command)
	}

	return editCommand
}

func getVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version [flags]",
		Short: "Command to fetch the version of terragen installed",
		Long:  `This will help user to find what version of terragen he/she installed in her machine.`,
		RunE:  versionConfig,
	}
}

func echoTerragen(cmd *cobra.Command, _ []string) error {
	return cmd.Usage()
}

func setCLI(_ *cobra.Command, _ []string) error {
	InitLogger(cliLogLevel)

	return nil
}

func versionConfig(_ *cobra.Command, _ []string) error {
	buildInfo, err := json.Marshal(version.GetBuildInfo())
	if err != nil {
		cliLogger.Errorf("fetching version information of terragen errored with: %s", err.Error())
		os.Exit(1)
	}
	fmt.Printf("terragen version: %s\n", string(buildInfo))

	return nil
}

// The function that helps in registering the subcommands with the respective main command.
// Make sure you call this, and this is the only way to register the subcommands.
func createRegister(name string, flagsRequired bool, fn *cobra.Command) {
	if createCommands == nil {
		createCommands = make(map[string]*cobra.Command)
	}

	if createCommands[name] != nil {
		panic(fmt.Sprintf("Command %s is already registered", name))
	}

	if flagsRequired {
		registerFlags(name, fn)
	}
	createCommands[name] = fn
}

func editRegister(name string, flagsRequired bool, fn *cobra.Command) {
	if editCommands == nil {
		editCommands = make(map[string]*cobra.Command)
	}

	if editCommands[name] != nil {
		panic(fmt.Sprintf("Command %s is already registered", name))
	}

	if flagsRequired {
		registerFlags(name, fn)
	}
	editCommands[name] = fn
}

// This function will return the custom template for usage function,
// only functions/methods inside this package can call this.

func getUsageTemplate() string {
	return `Usage:{{if .Runnable}}
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
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}
{{if .HasAvailableSubCommands}}{{printf "\n"}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
{{printf "\n"}}`
}
