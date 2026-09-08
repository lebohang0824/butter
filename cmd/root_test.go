package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"butter/pkg/lexer"
	"butter/pkg/parser"
	"butter/pkg/semantic"
)

func TestRepoSpecsParse(t *testing.T) {
	specDirs := []string{"specs", ".butter", "docs", "butter-extension"}
	seen := map[string]bool{}
	for _, dir := range specDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".butter") {
				continue
			}
			path := filepath.Join(dir, e.Name())
			if seen[path] {
				continue
			}
			seen[path] = true
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			l := lexer.NewLexer(string(content))
			p := parser.NewParser(l)
			spec, err := p.Parse()
			if err != nil {
				t.Errorf("parse %s: %v", path, err)
				continue
			}
			for _, d := range semantic.Analyze(spec) {
				if d.Severity == semantic.SemError {
					t.Errorf("%s: %v", path, d)
				}
			}
		}
	}
}
