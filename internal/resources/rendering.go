package resources

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var templatesByKind = make(map[string]*template.Template)

func (yr *YamlResource) Render(templatesDir string) ([]byte, error) {
	if tmpl, err := loadTemplate(templatesDir, yr.Kind); err != nil {
		return nil, err
	} else {
		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, yr); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

func loadTemplate(tmplDir, kind string) (*template.Template, error) {
	if _, found := templatesByKind[kind]; !found {
		templateSrcPath := filepath.Join(tmplDir, fmt.Sprintf("kind.%s.tmpl", strings.ToLower(kind)))
		if templateSrcBytes, err := os.ReadFile(templateSrcPath); err != nil {
			return nil, fmt.Errorf("reading template: '%s': %s'", templateSrcPath, err)
		} else {
			t := template.New(kind)
			if tmpl, terr := t.Parse(string(templateSrcBytes)); terr != nil {
				return nil, fmt.Errorf("parsing template: '%s': %s", templateSrcPath, terr)
			} else {
				templatesByKind[kind] = tmpl

			}
		}
	}

	return templatesByKind[kind], nil
}
