package y2c

import (
	"fmt"
	"github.com/spf13/afero"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var Fs = afero.NewOsFs()

type PageNode struct {
	parent *PageNode // used only during construction to match index.yaml with

	IsIndex  bool
	Name     string
	Path     string
	FilePath string
	Children []*PageNode
}

func (n *PageNode) addPage(filepath string) {

}

func newPageNode(filePath string, info fs.FileInfo) *PageNode {
	//filePath := strings.TrimPrefix(path, strings.TrimSuffix(sourceRoot, string(os.PathSeparator))+string(os.PathSeparator))
	isIndex := strings.EqualFold(info.Name(), "index.yaml")
	path := filepath.Dir(filePath)

	return &PageNode{
		parent:   nil,
		IsIndex:  isIndex,
		Name:     info.Name(),
		Path:     path,
		FilePath: filePath,
		Children: make([]*PageNode,0),
	}
}

func LoadPages(sourceRoot string) error {
	var rootNode *PageNode = newPageNode("index.yaml")
		parent:   nil,
		IsIndex:  true,
		Name:     "index.yaml",
		Path:     ".",
		FilePath: "index.yaml",
		Children: nil,
	}

	afero.Walk(Fs, sourceRoot, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		filePath := strings.TrimPrefix(path, strings.TrimSuffix(sourceRoot, string(os.PathSeparator))+string(os.PathSeparator))
		node := newPageNode(filePath, info)

		//fmt.Printf("%s [%s] (%s)\n", filePath, path, strconv.FormatBool(isIndex))

		return nil
	})
	return nil
}

//b, err := afero.ReadFile(Fs, path)
//if err != nil {
//	return nil, fmt.Errorf("LoadPages: reading %s: %s", path, err)
//}
//var p PageYaml
//if err = yaml.Unmarshal(b, &p); err != nil {
//	return nil, fmt.Errorf("LoadPages: parsing %s: %s", path, err)
//}
//fmt.Printf("sourcePath: %s\n", sourcePath)
//fmt.Printf("            %s\n", strings.TrimSuffix(sourcePath, string(os.PathSeparator)))
//fmt.Printf("            %s\n", strings.TrimSuffix(sourcePath, string(os.PathSeparator))+string(os.PathSeparator))
//fmt.Printf("            %s\n", path)
//fmt.Printf("            %s\n---\n", strings.TrimPrefix(path, strings.TrimSuffix(sourcePath, string(os.PathSeparator))+string(os.PathSeparator)))
