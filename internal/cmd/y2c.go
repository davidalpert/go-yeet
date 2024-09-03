package cmd

import (
	"github.com/davidalpert/go-printers/v1"
	"github.com/davidalpert/go-yeet/internal/atlassian"
	"github.com/davidalpert/go-yeet/internal/cfg"
	"github.com/davidalpert/go-yeet/internal/y2c"
	"github.com/spf13/cobra"
)

type Y2COptions struct {
	*printers.PrinterOptions
	cfg            *cfg.Config
	Client         *atlassian.Confluence
	DocumentRootID int
	YamlSourceDir  string
	DryRun         bool
}

func NewY2COptions(s printers.IOStreams) *Y2COptions {
	return &Y2COptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(s).WithDefaultOutput("yaml"),
		cfg:            &cfg.Config{},
	}
}

func NewCmdY2C(s printers.IOStreams) *cobra.Command {
	o := NewY2COptions(s)
	var cmd = &cobra.Command{
		Use:     "y2c <source-dir>",
		Aliases: []string{"yaml2confluence"},
		Short:   "synchronize a remote confluence workspace to match a local yaml source dir",
		Args:    cobra.ExactArgs(1),
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

	cmd.Flags().IntVar(&o.DocumentRootID, "document-root-id", 0, "Document ID of the root documentation page")
	cmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "Describe what would happen without making any changes to confluence")

	o.AddPrinterFlags(cmd.Flags())

	return cmd
}

// Complete the options
func (o *Y2COptions) Complete(cmd *cobra.Command, args []string) error {
	o.YamlSourceDir = args[0]

	if err := cfg.ReadMergedInto(o.cfg); err != nil {
		return err
	}
	if err := o.cfg.Validate(); err != nil {
		return err
	}

	if c, err := atlassian.New(&o.cfg.AtlassianCloud); err != nil {
		return err
	} else {
		o.Client = c
	}

	return nil
}

// Validate the options
func (o *Y2COptions) Validate() error {
	return o.PrinterOptions.Validate()
}

// Run the command
func (o *Y2COptions) Run() error {
	if o.FormatCategory() == "csv" {
		o.WithDefaultOutput("json")
	}

	y2c.LoadPages(o.YamlSourceDir)

	return nil
}
