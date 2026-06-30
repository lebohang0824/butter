package output

import (
	"fmt"
	"sort"

	"butter/pkg/ast"
)

type Extension interface {
	Name() string
	FileExtension() string
	Serialize(spec *ast.AppSpec) ([]byte, error)
}

var registry = map[string]Extension{}

func Register(ext Extension) {
	name := ext.Name()
	if _, ok := registry[name]; ok {
		panic(fmt.Sprintf("output extension %q already registered", name))
	}
	registry[name] = ext
}

func Get(name string) (Extension, bool) {
	ext, ok := registry[name]
	return ext, ok
}

func Names() []string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
