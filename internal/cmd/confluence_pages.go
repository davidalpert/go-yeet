package cmd

import (
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
)

func NewCmdConfluencePage(s printers.IOStreams) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "page",
		Aliases: []string{"p", "pg", "pages"},
		Short:   "page subcommands",
		//Args:    cobra.NoArgs,
	}

	cmd.AddCommand(NewCmdConfluencePageGet(s))

	return cmd
}
