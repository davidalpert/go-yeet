package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidalpert/go-yeet/internal/diagnostics"
	"path/filepath"

	"github.com/aybabtme/orderedjson"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/yaml.v3"
)

type KindAndTitle struct {
	Kind  string `json:"kind"`
	Title string `json:"title"`
}

type Labels struct {
	Labels []string `json:"labels"`
}

type YamlResource struct {
	Kind  string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Title string `json:"title" yaml:"title"`
	Path  string `json:"path" yaml:"path"`
	Node  *yaml.Node
	Json  string `json:"json,omitempty" yaml:"json,omitempty"`
}

var jsonEncoder yqlib.Encoder = yqlib.NewJSONEncoder(yqlib.JsonPreferences{Indent: 0, ColorsEnabled: false, UnwrapScalar: false})

func NewYamlResource(path string, node *yaml.Node) (*YamlResource, error) {
	setHeadComment(path, node)
	node.FootComment = "V2"

	yr := &YamlResource{
		Path: path,
		Node: node,
	}

	err := yr.UpdateJson()

	return yr, err
}

func (yr *YamlResource) GetParentPath() string {
	return filepath.Dir(yr.Path)
}

func (yr *YamlResource) UpdateJson() error {
	var buf bytes.Buffer
	var yqNode yqlib.CandidateNode
	diagnostics.Log.Debug("unmarshalYAML")
	if err := yqNode.UnmarshalYAML(yr.Node, make(map[string]*yqlib.CandidateNode, 0)); err != nil {
		return fmt.Errorf("UpdateJson: UnmarshalYAML: %s: %#v", err, yr.Node)
	}
	diagnostics.Log.Debug("encoding as JSON")
	if err := jsonEncoder.Encode(&buf, &yqNode); err != nil {
		return fmt.Errorf("UpdateJson: EncodeJSON: %s", err)
	}
	yr.Json = buf.String()
	diagnostics.Log.Debug("updating kind and title")
	return yr.UpdateKindAndTitle()
}

// this is ugly, but it works for now
func (yr *YamlResource) UpdateKindAndTitle() error {
	if yr == nil {
		return fmt.Errorf("UpdateKindAndTitle: cannot update a nil resource")
	}

	kindAndTitle := &KindAndTitle{}
	diagnostics.Log.Debug("unmarshalJSON")
	if err := json.Unmarshal([]byte(yr.Json), &kindAndTitle); err != nil {
		return fmt.Errorf("UpdateKindAndTitle: %s", err)
	}

	yr.Kind = kindAndTitle.Kind
	yr.Title = kindAndTitle.Title

	return nil
}

func (yr *YamlResource) GetLabels() ([]string, error) {
	if yr == nil {
		return nil, fmt.Errorf("update kind and title: cannot update a nil resource")
	}

	labels := &Labels{}
	err := json.Unmarshal([]byte(yr.Json), &labels)

	return labels.Labels, err
}

func (yr *YamlResource) ToObject() (map[string]interface{}, error) {
	var obj map[string]interface{}

	if yr == nil {
		return nil, fmt.Errorf("update kind and title: cannot update a nil resource")
	}

	if err := json.Unmarshal([]byte(yr.Json), &obj); err != nil {
		panic(err)
	}

	return obj, nil
}

func (yr *YamlResource) ToOrderedMap() orderedjson.Map {
	var object orderedjson.Map
	err := json.Unmarshal([]byte(yr.Json), &object)
	if err != nil {
		panic(err)
	}

	return object
}

func setHeadComment(path string, node *yaml.Node) {
	if node.Kind == 0 {
		node.Kind = yaml.MappingNode
	}

	if node.HeadComment == "" {
		node.HeadComment = path
	} else {
		node.HeadComment = path + "\n" + node.HeadComment
	}

}
