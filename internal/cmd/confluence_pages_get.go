package cmd

import (
	"fmt"
	"github.com/davidalpert/go-printers/v1"
	"github.com/davidalpert/go-yeet/internal/atlassian"
	"github.com/davidalpert/go-yeet/internal/cfg"
	"github.com/spf13/cobra"
)

type ConfluencePageGetOptions struct {
	*printers.PrinterOptions
	cfg            *cfg.Config
	Client         *atlassian.Confluence
	DocumentRootID int
	WithChildren   bool
}

func NewConfluencePageGetOptions(s printers.IOStreams) *ConfluencePageGetOptions {
	return &ConfluencePageGetOptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(s).WithDefaultOutput("yaml"),
		cfg:            &cfg.Config{},
	}
}

func NewCmdConfluencePageGet(s printers.IOStreams) *cobra.Command {
	o := NewConfluencePageGetOptions(s)
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "get pages",
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

	cmd.Flags().IntVar(&o.DocumentRootID, "document-root-id", 0, "Document ID of the root documentation page")
	cmd.Flags().BoolVar(&o.WithChildren, "with-children", false, "Include children of Document ID")

	o.AddPrinterFlags(cmd.Flags())

	return cmd
}

// Complete the options
func (o *ConfluencePageGetOptions) Complete(cmd *cobra.Command, args []string) error {
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
func (o *ConfluencePageGetOptions) Validate() error {
	if o.WithChildren && o.DocumentRootID == 0 {
		return fmt.Errorf("with-children is only valid with document-root-id")
	}

	return o.PrinterOptions.Validate()
}

// Run the command
func (o *ConfluencePageGetOptions) Run() error {
	if o.FormatCategory() == "table" || o.FormatCategory() == "csv" {
		o.WithDefaultOutput("json")
	}

	if o.DocumentRootID == 0 {
		if result, err := o.Client.GetPages(); err != nil {
			return err
		} else {
			return o.WriteOutput(result)
		}
	}

	if result, err := o.Client.GetPageByID(o.DocumentRootID, o.WithChildren); err != nil {
		return err
	} else {
		return o.WriteOutput(result)
	}
}
