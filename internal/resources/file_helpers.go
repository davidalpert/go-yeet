package resources

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thanhpk/randstr"
)

const DEFAULT_SECRET_FILENAME = ".secret"

type DirectoryProperties struct {
	ConfigPath   string
	SpaceDir     string
	TemplatesDir string
	HooksDir     string
	SpaceKey     string
}

func GetDefaultY2cHomeDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(dirname, ".y2c")
}

func GetDefaultSecretPath() string {
	return filepath.Join(GetDefaultY2cHomeDir(), DEFAULT_SECRET_FILENAME)
}

func GetSecret() (string, error) {
	secret, err := os.ReadFile(GetDefaultSecretPath())
	if err != nil {
		return "", err
	}

	return string(secret), nil
}

func GetSecretAndGenerateIfMissing() string {
	secret, err := GetSecret()
	if err == nil {
		return secret
	}

	return GenerateSecret()
}

func GenerateSecret() string {
	secret := []byte(randstr.String(24))

	dir := GetDefaultY2cHomeDir()
	secretPath := GetDefaultSecretPath()
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(secretPath, secret, 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated secret key: " + secretPath)

	return string(secret)
}

func GetDirectoryProperties(path string) DirectoryProperties {
	dirTokens := strings.Split(ResolveAbsolutePathFile(path), "spaces/")
	baseDir := dirTokens[0]
	spaceKey := ""
	if len(dirTokens) == 2 {
		spaceKey = strings.Split(dirTokens[1], "/")[0]
	} else if len(dirTokens) == 1 {
		// it's assumed the space key was passed, treat the current directory as baseDir
		// TODO this is ugly, clean up
		current, err := fs.Getwd()
		if err != nil {
			panic(err)
		}
		baseDir = current
		spaceKey = path
	}

	props := DirectoryProperties{}
	props.ConfigPath = filepath.Join(baseDir, "config.yml")
	props.SpaceDir = filepath.Join(baseDir, "spaces", spaceKey)
	props.SpaceKey = spaceKey
	props.TemplatesDir = filepath.Join(baseDir, "templates")
	props.HooksDir = filepath.Join(baseDir, "hooks")

	os.Setenv("SPACE_DIR", props.SpaceDir)

	if _, err := os.Stat(props.ConfigPath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Could not find config.yml")
		os.Exit(1)
	}

	if stat, err := os.Stat(props.SpaceDir); errors.Is(err, os.ErrNotExist) || !stat.IsDir() {
		fmt.Printf("Could not find '%s' space directory", props.SpaceKey)
		os.Exit(1)
	}

	if stat, err := os.Stat(props.TemplatesDir); errors.Is(err, os.ErrNotExist) || !stat.IsDir() {
		fmt.Println("Could not find templates directory")
		os.Exit(1)
	}

	return props
}

func CreateInstanceDirectory(instanceDir string, config string) {
	configFile := filepath.Join(instanceDir, "config.yml")

	// create instance directory
	os.Mkdir(instanceDir, 0755)

	// create spaces, templates, and hooks directories
	createDirectories(instanceDir, []string{"spaces", "templates", "hooks"})

	// write config.yml
	err := os.WriteFile(configFile, []byte(config), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created file " + configFile)

}

func createDirectories(instanceDir string, dirs []string) {
	for _, dir := range dirs {
		fullDir := filepath.Join(instanceDir, dir)

		err := os.Mkdir(fullDir, 0755)
		if err != nil {
			fmt.Printf("Failed to create directory %s\n%s\n", fullDir, err.Error())
			os.Exit(1)
		}
		fmt.Println("Created directory " + fullDir)
	}
}

func ResolveAbsolutePathDir(dir string) string {
	if filepath.IsAbs(dir) {
		return dir
	}

	current, err := fs.Getwd()
	if err != nil {
		panic(err)
	}

	if dir == "" {
		return current
	}

	absPath, err := filepath.Abs(filepath.Join(current, dir))
	if err != nil {
		panic(err)
	}

	return absPath
}

func ResolveAbsolutePathFile(file string) string {
	if file == "" {
		panic("no file provided")
	}

	return ResolveAbsolutePathDir(file)
}
