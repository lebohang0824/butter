package yaml

import (
	"butter/pkg/ast"
	"butter/pkg/output"

	"gopkg.in/yaml.v3"
)

func init() {
	output.Register(yamlExt{})
}

type yamlExt struct{}

func (yamlExt) Name() string          { return "yaml" }
func (yamlExt) FileExtension() string { return ".yaml" }

func (yamlExt) Serialize(spec *ast.AppSpec) ([]byte, error) {
	return yaml.Marshal(spec)
}
