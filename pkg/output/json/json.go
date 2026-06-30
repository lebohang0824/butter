package json

import (
	"encoding/json"

	"butter/pkg/ast"
	"butter/pkg/output"
)

func init() {
	output.Register(jsonExt{})
}

type jsonExt struct{}

func (jsonExt) Name() string          { return "json" }
func (jsonExt) FileExtension() string { return ".json" }

func (jsonExt) Serialize(spec *ast.AppSpec) ([]byte, error) {
	return json.MarshalIndent(spec, "", "  ")
}
