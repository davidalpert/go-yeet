package cfg

import (
	"fmt"
	"github.com/davidalpert/go-printers/v1"
	"github.com/davidalpert/go-yeet/internal/atlassian"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"path/filepath"
	"strings"
)

var File string
var Dir string

type Config struct {
	AtlassianCloud atlassian.CloudConfig `yaml:"atlassian_cloud"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) String() string {
	s, _, _ := printers.NewPrinterOptions().WithDefaultOutput("yaml").FormatOutput(c)
	return s
}

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}
	errors := make([]string, 0)

	errors = append(errors, c.AtlassianCloud.Validate("atlassian_cloud")...)

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	return nil
}

func ReadMergedInto(c *Config) error {
	if err := ensureFileExists(File, c.String()); err != nil {
		return err
	}
	return cleanenv.ReadConfig(File, c)
}

func (c *Config) Write() error {
	return c.WriteToFile(File)
}

func (c *Config) WriteToFile(file string) error {
	configDir, _ := filepath.Split(file)

	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return fmt.Errorf("make %s: %v", configDir, err)
	}

	y, _, err := printers.NewPrinterOptions().WithDefaultOutput("yaml").FormatOutput(*c)
	if err != nil {
		return fmt.Errorf("formatting %#v: %v", c, err)
	}

	if err = os.WriteFile(file, []byte(y), os.ModePerm); err != nil {
		return fmt.Errorf("write %s: %v", file, err)
	}

	return nil
}

func ensureFileExists(path string, defaultContent string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("making %s; %v", filepath.Dir(path), err)
	}

	return os.WriteFile(path, []byte(defaultContent), os.ModePerm)
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	Dir = filepath.Join(home, ".yeet")
	File = filepath.Join(Dir, "config.yaml")
}
