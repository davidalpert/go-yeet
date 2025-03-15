package cmd

import (
	"fmt"
	"github.com/davidalpert/go-printers/v1"
	"github.com/davidalpert/go-yeet/internal/cfg"
	"github.com/davidalpert/go-yeet/internal/resources"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type SourceRenderOptions struct {
	*printers.PrinterOptions
	cfg          *cfg.Config
	SourcePath   string
	TemplatesDir string
	OutputPath   string
}

func NewSourceRenderOptions(s printers.IOStreams) *SourceRenderOptions {
	return &SourceRenderOptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(s).WithDefaultTableWriter(),
		cfg:            &cfg.Config{},
	}
}

func NewCmdSourceRender(s printers.IOStreams) *cobra.Command {
	o := NewSourceRenderOptions(s)
	var cmd = &cobra.Command{
		Use:          "render [sourcePath] [outputPath]",
		Short:        "render source documents",
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
	cmd.Flags().StringVar(&o.OutputPath, "output-path", "STDOUT", "write to path")
	cmd.Flags().StringVar(&o.TemplatesDir, "template-dir", "./.templates", "a folder where markdown rendering templates can be located")

	return cmd
}

// Complete the options
func (o *SourceRenderOptions) Complete(cmd *cobra.Command, args []string) error {
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
func (o *SourceRenderOptions) Validate() error {
	return o.PrinterOptions.Validate()
}

// Run the command
func (o *SourceRenderOptions) Run() error {
	if o.FormatCategory() == "table" || o.FormatCategory() == "csv" {
		o.WithDefaultOutput("json")
	}

	if isDir, err := afero.IsDir(Fs, o.SourcePath); err != nil {
		return fmt.Errorf("test path: %#v: %s", o.SourcePath, err)
	} else if isDir {
		return o.renderDirectory()
	} else {
		return fmt.Errorf("TODO: validate file: %#v", o.SourcePath)
	}
}

func (o *SourceRenderOptions) renderDirectory() error {
	yrs, err := resources.LoadYamlResources(o.SourcePath)
	if err != nil {
		return fmt.Errorf("load yaml resources: %s", err)
	}

	for _, yr := range yrs {
		if resultBytes, renderErr := yr.Render(o.TemplatesDir); renderErr != nil {
			fmt.Fprintf(o.Out, "---\n%s\n-v-\nERROR rendering: %s\n", yr.Path, renderErr)
		} else if o.OutputPath == "STDOUT" {
			fmt.Fprintf(o.Out, "---\n%s\n-v-\n%s\n", yr.Path, string(resultBytes))
		} else {
			outFilePath := filepath.Join(o.OutputPath, fmt.Sprintf("%s.md", filepath.Base(yr.Path)))
			outFileDir := filepath.Dir(outFilePath)
			if err = os.MkdirAll(outFileDir, os.ModePerm); err != nil {
				return fmt.Errorf("creating dir '%s': %s", outFileDir, err)
			}
			if err = os.WriteFile(outFilePath, resultBytes, os.ModePerm); err != nil {
				return fmt.Errorf("writing to '%s': %s", outFilePath, err)
			}
		}
	}

	return nil
}
