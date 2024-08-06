package cmd

import (
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
)

func NewCmdConfluence(s printers.IOStreams) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "confluence",
		Aliases: []string{"cn", "conf"},
		Short:   "confluence subcommands",
		//Args:    cobra.NoArgs,
	}

	cmd.AddCommand(NewCmdConfluencePage(s))

	return cmd
}
