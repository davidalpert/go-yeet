package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/davidalpert/go-printers/v1"
	"github.com/davidalpert/go-yeet/internal/cmd/adapters"
	"github.com/davidalpert/go-yeet/internal/example"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ExampleOptions struct {
	*printers.PrinterOptions
	OutFolder   string
	ExampleName string
}

func NewExampleOptions(s printers.IOStreams) *ExampleOptions {
	return &ExampleOptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(s).WithDefaultOutput("yaml"),
		OutFolder:      "./example",
		ExampleName:    "default", // TODO: add more examples
	}
}

func NewCmdExample(s printers.IOStreams) *cobra.Command {
	o := NewExampleOptions(s)
	var cmd = &cobra.Command{
		Use:        "example",
		ArgAliases: []string{"ex", "sample"},
		Short:      "generate example files and templates",
		Args:       cobra.NoArgs,
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
	cmd.Flags().StringVar(&o.OutFolder, "out-dir", "./example", "file path to the folder in which to create the example")
	cmd.Flags().StringVar(&o.ExampleName, "example-name", "default", "which example to generate")

	return cmd
}

// Complete the options
func (o *ExampleOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Validate the options
func (o *ExampleOptions) Validate() error {
	if o.OutFolder == "" {
		return fmt.Errorf("OutFolder cannot be empty")
	}

	return o.PrinterOptions.Validate()
}

type writeExampleResult struct {
	OutFolder    string   `json:"out_folder"`
	ExampleName  string   `json:"example_name"`
	FilesWritten []string `json:"files_written"`
	Error        error    `json:"error,omitempty" yaml:"error,omitempty"`
}

// Run the command
func (o *ExampleOptions) Run() error {
	examplePrefix := fmt.Sprintf("%s/", o.ExampleName)
	//fmt.Fprintf(o.Out, "%s\n", examplePrefix)

	result := writeExampleResult{
		OutFolder:    o.OutFolder,
		ExampleName:  o.ExampleName,
		FilesWritten: make([]string, 0),
		Error:        nil,
	}

	if _, err := os.Stat(o.OutFolder); err == nil {
		var overwrite bool
		if err2 := adapters.AskOneWithStreams(o.IOStreams, &survey.Confirm{
			Message: fmt.Sprintf("destination folder '%s' already exists; overwrite it?", o.OutFolder),
		}, &overwrite, survey.WithValidator(survey.Required)); err2 != nil {
			return err2
		}
		if !overwrite {
			fmt.Fprintf(o.ErrOut, "chose not to overwrite existing folder %v\n", o.OutFolder)
			return nil
		}
		if existingDir, err2 := os.ReadDir(o.OutFolder); err2 != nil {
			for _, d := range existingDir {
				os.RemoveAll(path.Join([]string{o.OutFolder, d.Name()}...))
			}
		}
	}

	var exampleFS = example.DefaultFiles
	if err := fs.WalkDir(exampleFS, ".", func(path string, d fs.DirEntry, err error) error {
		if !strings.HasPrefix(path, examplePrefix) {
			return nil
		}

		destPath := filepath.Join(o.OutFolder, strings.TrimPrefix(path, examplePrefix))
		if d.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		if b, err2 := fs.ReadFile(exampleFS, path); err2 != nil {
			return err2
		} else {
			if err3 := os.WriteFile(destPath, b, os.ModePerm); err3 != nil {
				return err3
			}
			result.FilesWritten = append(result.FilesWritten, destPath)
		}
		return nil
	}); err != nil {
		return err
	}

	return o.WriteOutput(result)
}
