package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/davidalpert/go-printers/v1"
)

// askOneWithStreams is a convenience method adapting survey to use go-printers for STDIO
func askOneWithStreams(s printers.IOStreams, p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	opts = append(opts, survey.WithStdio(s.In.(terminal.FileReader), s.Out.(terminal.FileWriter), s.ErrOut))
	return survey.AskOne(p, response, opts...)
}
