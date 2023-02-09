package cmd

import (
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
)

func NewCmdConfig(s printers.IOStreams) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg", "c"},
		Short:   "configuration subcommands",
		//Args:    cobra.NoArgs,
	}

	cmd.AddCommand(NewCmdConfigGet(s))
	cmd.AddCommand(NewCmdConfigSetup(s))

	return cmd
}
