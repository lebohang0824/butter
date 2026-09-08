package prompt

import (
	"strings"
	"testing"

	"butter/pkg/ast"
)

func TestSerializeIncludesRules(t *testing.T) {
	spec := &ast.AppSpec{
		App:   "TestApp",
		Rules: []ast.RuleSpec{{Statement: "Use TypeScript"}, {Statement: "Use React"}},
	}

	out, err := (promptExt{}).Serialize(spec)
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	got := string(out)
	for _, want := range []string{"### Rules", "* Use TypeScript", "* Use React"} {
		if !strings.Contains(got, want) {
			t.Errorf("Serialize() output missing %q; got:\n%s", want, got)
		}
	}
}

func TestSerializeOmitsRulesWhenEmpty(t *testing.T) {
	spec := &ast.AppSpec{App: "TestApp"}

	out, err := (promptExt{}).Serialize(spec)
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	if strings.Contains(string(out), "### Rules") {
		t.Errorf("Serialize() should omit Rules section when no rules present; got:\n%s", out)
	}
}