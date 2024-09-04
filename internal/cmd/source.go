package cmd

import (
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
)

func NewCmdSource(s printers.IOStreams) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "source",
		Aliases: []string{"src", "s"},
		Short:   "source file subcommands",
		//Args:    cobra.NoArgs,
	}

	cmd.AddCommand(NewCmdSourceValidate(s))

	return cmd
}
