package resources

import (
	"bytes"
	"encoding/json"
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
	Kind  string
	Title string
	Path  string
	Node  *yaml.Node
	Json  string
}

var jsonEncoder yqlib.Encoder = yqlib.NewJSONEncoder(yqlib.JsonPreferences{Indent: 0, ColorsEnabled: false, UnwrapScalar: false})

func NewYamlResource(path string, node *yaml.Node) *YamlResource {
	setHeadComment(path, node)
	node.FootComment = "V2"

	yr := &YamlResource{
		Path: path,
		Node: node,
	}

	yr.UpdateJson()

	return yr
}

func (yr *YamlResource) GetParentPath() string {
	return filepath.Dir(yr.Path)
}

func (yr *YamlResource) UpdateJson() {
	var buf bytes.Buffer
	var yqNode yqlib.CandidateNode
	yqNode.UnmarshalYAML(yr.Node, make(map[string]*yqlib.CandidateNode, 0))
	jsonEncoder.Encode(&buf, &yqNode)
	yr.Json = buf.String()
	yr.UpdateKindAndTitle()
}

// this is ugly, but it works for now
func (yr *YamlResource) UpdateKindAndTitle() {
	kindAndTitle := &KindAndTitle{}
	if err := json.Unmarshal([]byte(yr.Json), &kindAndTitle); err != nil {
		panic(err)
	}

	yr.Kind = kindAndTitle.Kind
	yr.Title = kindAndTitle.Title
}

func (yr *YamlResource) GetLabels() []string {
	labels := &Labels{}
	if err := json.Unmarshal([]byte(yr.Json), &labels); err != nil {
		panic(err)
	}

	return labels.Labels
}

func (yr *YamlResource) ToObject() map[string]interface{} {
	var obj map[string]interface{}

	if err := json.Unmarshal([]byte(yr.Json), &obj); err != nil {
		panic(err)
	}

	return obj
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
