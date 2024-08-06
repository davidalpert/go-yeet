package atlassian

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/davidalpert/go-printers/v1"
	"github.com/davidalpert/go-yeet/internal/cmd/adapters"
)

type CloudConfig struct {
	Host     string `yaml:"host"`
	Email    string `yaml:"email"`
	APIToken string `yaml:"api_token"`
}

func (c *CloudConfig) Validate(key string) []string {
	errors := make([]string, 0)

	if c.Host == "" {
		errors = append(errors, fmt.Sprintf("%s.host is required", key))
	}

	if c.Email == "" {
		errors = append(errors, fmt.Sprintf("%s.email is required", key))
	}

	if c.APIToken == "" {
		errors = append(errors, fmt.Sprintf("%s.api_token is required", key))
	}

	return errors
}

func (c *CloudConfig) ConfigureWithSurvey(o printers.IOStreams) error {
	if err := adapters.AskOneWithStreams(o, &survey.Input{
		Message: "Cloud Host (e.g. https://<_jira_instance_>.atlassian.net)",
		Default: c.Host,
	}, &c.Host, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	if err := adapters.AskOneWithStreams(o, &survey.Input{
		Message: "User Email",
		Default: c.Email,
	}, &c.Email, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	if err := adapters.AskOneWithStreams(o, &survey.Input{
		Message: "User Token",
		Default: c.APIToken,
	}, &c.APIToken, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	return nil
}
