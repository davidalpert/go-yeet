package resources

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/davidalpert/go-yeet/internal/diagnostics"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const DUPLICATE_TITLE = `Duplicate title found -- "%s" (%s/) matches "%s" (%s/)`

type YamlResourceLoader struct {
	Walk     func(root string, fn filepath.WalkFunc) error
	LoadYaml func(file string) ([]byte, error)
}

func DefaultYamlResourceLoader() YamlResourceLoader {
	return YamlResourceLoader{filepath.Walk, DefaultLoadYaml}
}

func LoadYamlResources(dir string) ([]*YamlResource, error) {
	return DefaultYamlResourceLoader().loadYamlResources(dir)
}

func DefaultLoadYaml(file string) ([]byte, error) {
	return os.ReadFile(file)
}

func unmarshal(data []byte) (*yaml.Node, error) {
	node := yaml.Node{}
	err := yaml.Unmarshal(data, &node)
	return &node, err
}

func (yrl YamlResourceLoader) loadYamlResources(dir string) ([]*YamlResource, error) {
	yrs := make([]*YamlResource, 0)
	parents := make(map[string]*YamlResource, 0)

	dir = filepath.Clean(dir)
	dirStringLength := len(dir)

	err := yrl.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			diagnostics.Log.WithField("path", path).Debug("load resources - walk - consider")
			if err != nil {
				return err
			}

			// skip space dir
			if path == dir {
				return nil
			}

			relPath := path[dirStringLength+1:]
			absPath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("resolving relative path '%s' to absolute", path)
			}
			diagnostics.Log.WithFields(log.Fields{
				"path":            path,
				"absPath":         absPath,
				"dir":             dir,
				"dirStringLength": dirStringLength,
				"relPath":         relPath,
			}).Debug("load resources - walk - prepare")

			if info.IsDir() {
				if ignoreDir(relPath) {
					return filepath.SkipDir
				}

				yr, loadErr := getDefaultDirYamlResource(relPath)
				if loadErr != nil {
					return fmt.Errorf("loading dir %#v: %s", relPath, loadErr)
				}

				// save a pointer to the directory YamlResource for later in case an index.yml is found
				parents[relPath] = yr
				yrs = append(yrs, yr)
			} else if IsYamlFile(path) {
				yr, loadErr := yrl.LoadYamlResource(dir, relPath)
				if loadErr != nil {
					return fmt.Errorf("loading file %#v: %s", absPath, loadErr)
				}
				if isIndexFile(path) {
					parent := parents[filepath.Dir(relPath)]
					parent.Kind = yr.Kind
					parent.Title = yr.Title
					parent.Json = yr.Json
					parent.Node = yr.Node
				} else {
					yrs = append(yrs, yr)
				}
			}

			return nil
		})

	return yrs, err
}

func LoadSingleYamlResource(file string) (*YamlResource, error) {
	fileAbs := ResolveAbsolutePathFile(file)
	yrl := YamlResourceLoader{func(root string, fn filepath.WalkFunc) error {
		fileInfo, err := os.Stat(fileAbs)
		if err != nil {
			return err
		}

		fn(fileAbs, fileInfo, nil)
		return nil
	}, DefaultLoadYaml}

	yrs, err := yrl.loadYamlResources(GetDirectoryProperties(file).SpaceDir)
	if err != nil {
		return nil, err
	}

	return yrs[0], nil
}

func (yrl YamlResourceLoader) LoadYamlResource(spaceRootDir, relFilePath string) (*YamlResource, error) {
	diagnostics.Log.Debug("LoadYamlResource - reading file")
	y, err := yrl.LoadYaml(filepath.Join(spaceRootDir, relFilePath))
	if err != nil {
		return nil, fmt.Errorf("LoadYamlResource: %s", err)
	}
	diagnostics.Log.Debug("LoadYamlResource - unmarshal from bytes")
	r, err := unmarshal(y)
	if err != nil {
		return nil, fmt.Errorf("LoadYamlResources: unmarshal: %s", err)
	}

	diagnostics.Log.Debug("LoadYamlResource - ")
	return NewYamlResource(relFilePath, r)
}

func IsYamlFile(file string) bool {
	ext := filepath.Ext(file)
	return ext == ".yml" || ext == ".yaml"
}

func isIndexFile(file string) bool {
	name := strings.Split(filepath.Base(file), ".")[0]
	return IsYamlFile(file) && (name == "index" || name == "_index")
}

func ignoreDir(path string) bool {
	return filepath.Base(path)[0:1] == "_"
}

func getDefaultDirYamlResource(relPath string) (*YamlResource, error) {
	pathTokens := strings.Split(relPath, string(os.PathSeparator))
	title := pathTokens[len(pathTokens)-1:][0]

	r, err := unmarshal([]byte(fmt.Sprintf("kind: wiki\ntitle: %s\nmarkup: \"\"", title)))
	if err != nil {
		return nil, err
	}
	return NewYamlResource(relPath, r)
}

func EnsureUniqueTitles(yrs []*YamlResource) error {
	uniqueTitle := map[string]*YamlResource{}

	for _, cur := range yrs {
		lowerTitle := strings.ToLower(cur.Title)
		if r, exists := uniqueTitle[lowerTitle]; exists {
			return errors.New(fmt.Sprintf(DUPLICATE_TITLE, r.Title, r.Path, cur.Title, cur.Path))
		} else {
			uniqueTitle[lowerTitle] = cur
		}
	}
	return nil
}

func GetAnchor(spaceDir string) string {
	data, err := os.ReadFile(filepath.Join(spaceDir, ".anchor"))
	if err != nil {
		return ""
	}

	return string(data)
}
