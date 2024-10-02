package cmd

import (
	"fmt"
	"github.com/davidalpert/go-printers/v1"
	"github.com/davidalpert/go-yeet/internal/cfg"
	"github.com/davidalpert/go-yeet/internal/resources"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"strconv"
)

type SourceValidateOptions struct {
	*printers.PrinterOptions
	cfg        *cfg.Config
	SourcePath string
}

func NewSourceValidateOptions(s printers.IOStreams) *SourceValidateOptions {
	return &SourceValidateOptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(s).WithDefaultTableWriter(),
		cfg:            &cfg.Config{},
	}
}

func NewCmdSourceValidate(s printers.IOStreams) *cobra.Command {
	o := NewSourceValidateOptions(s)
	var cmd = &cobra.Command{
		Use:          "validate [sourcePath]",
		Short:        "validate source documents",
		Args:         cobra.RangeArgs(0, 1),
		SilenceUsage: true,
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
func (o *SourceValidateOptions) Complete(cmd *cobra.Command, args []string) error {
	if err := cfg.ReadMergedInto(o.cfg); err != nil {
		return err
	}
	if err := o.cfg.Validate(); err != nil {
		return err
	}

	if len(args) > 0 {
		o.SourcePath = args[0]
	} else {
		o.SourcePath = "."
	}

	return nil
}

// Validate the options
func (o *SourceValidateOptions) Validate() error {
	return o.PrinterOptions.Validate()
}

// Run the command
func (o *SourceValidateOptions) Run() error {
	if o.FormatCategory() == "table" || o.FormatCategory() == "csv" {
		o.WithDefaultOutput("json")
	}

	if isDir, err := afero.IsDir(Fs, o.SourcePath); err != nil {
		return fmt.Errorf("test path: %#v: %s", o.SourcePath, err)
	} else if isDir {
		return o.validateDirectory()
	} else {
		return fmt.Errorf("TODO: validate file: %#v", o.SourcePath)
	}
}

func (o *SourceValidateOptions) validateDirectory() error {
	yrs, err := resources.LoadYamlResources(o.SourcePath)
	if err != nil {
		return fmt.Errorf("load yaml resources: %s", err)
	}

	return o.WithTableWriter(o.SourcePath, func(t *tablewriter.Table) {
		t.SetHeader([]string{
			"#",
			"Kind",
			"Path",
			"Title",
		})
		for i, r := range yrs {
			t.Append([]string{
				strconv.Itoa(i),
				r.Kind,
				r.Path,
				r.Title,
			})
		}
	}).WriteOutput(yrs)
}
