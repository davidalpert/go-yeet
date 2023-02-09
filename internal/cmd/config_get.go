package cmd

import (
	"github.com/davidalpert/go-yeet/internal/cfg"
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
)

type ConfigGetOptions struct {
	*printers.PrinterOptions
	Values *cfg.Config
}

func NewConfigGetOptions(s printers.IOStreams) *ConfigGetOptions {
	return &ConfigGetOptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(s).WithDefaultOutput("yaml"),
		Values:         &cfg.Config{},
	}
}

func NewCmdConfigGet(s printers.IOStreams) *cobra.Command {
	o := NewConfigGetOptions(s)
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "show current configuration values",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			return o.Run()
		},
	}

	o.AddPrinterFlags(cmd.Flags())

	return cmd
}

// Complete the options
func (o *ConfigGetOptions) Complete(cmd *cobra.Command, args []string) error {
	if err := cfg.ReadMergedInto(o.Values); err != nil {
		return err
	}

	return nil
}

// Validate the options
func (o *ConfigGetOptions) Validate() error {
	return o.PrinterOptions.Validate()
}

// Run the command
func (o *ConfigGetOptions) Run() error {
	if o.FormatCategory() == "table" || o.FormatCategory() == "csv" {
		o.WithDefaultOutput("json")
	}

	return o.WriteOutput(*o.Values)
}
