package prompt

import (
	"strings"
	"testing"

	"butter/pkg/ast"
)

func TestSerializeIncludesRules(t *testing.T) {
	spec := &ast.AppSpec{
		App:   "TestApp",
		Rules: []ast.RuleSpec{{Statement: "Users may only access their own tasks"}, {Statement: "A task cannot be completed until it is assigned to a user"}},
	}

	out, err := (promptExt{}).Serialize(spec)
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	got := string(out)
	for _, want := range []string{"### Rules", "* Users may only access their own tasks", "* A task cannot be completed until it is assigned to a user"} {
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