package prompt

import (
	"fmt"
	"strings"

	"butter/pkg/ast"
	"butter/pkg/output"
)

func init() {
	output.Register(promptExt{})
}

type promptExt struct{}

func (promptExt) Name() string          { return "prompt" }
func (promptExt) FileExtension() string { return ".prompt.md" }

func (promptExt) Serialize(spec *ast.AppSpec) ([]byte, error) {
	var b strings.Builder

	fmt.Fprintf(&b, "# [SYSTEM SPEC] %s\n", spec.App)
	if spec.Version != "" {
		fmt.Fprintf(&b, "> **Version:** %s\n", spec.Version)
	}
	if spec.Description != "" {
		fmt.Fprintf(&b, "> **Description:** %s\n", spec.Description)
	}
	if spec.Version != "" || spec.Description != "" {
		b.WriteString("\n")
	}

	for i, feat := range spec.Features {
		if i > 0 {
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, "## Feature: %s\n", feat.Name)
		if feat.Version != "" {
			fmt.Fprintf(&b, "**Version:** %s\n", feat.Version)
		}
		if feat.Description != "" {
			fmt.Fprintf(&b, "%s\n", feat.Description)
		}

		if len(feat.Params) > 0 {
			b.WriteString("\n### Params\n")
			for _, p := range feat.Params {
				b.WriteString(formatParam(p))
			}
		}

		if len(feat.Actions) > 0 {
			b.WriteString("\n### Execution Sequence\n")
			b.WriteString("**CRITICAL:** Execute the following steps strictly in order. Do not proceed to the next step until the current one is complete.\n\n")
			for j, a := range feat.Actions {
				b.WriteString(formatAction(a, j+1))
			}
		}
	}

	return []byte(b.String()), nil
}

func formatParam(p ast.ParamSpec) string {
	var sb strings.Builder
	req := "optional"
	if p.Required {
		req = "required"
	}

	typeStr := p.Type
	if len(p.Validate) > 0 {
		typeStr += " (" + strings.Join(p.Validate, ", ") + ")"
	}

	sb.WriteString(fmt.Sprintf("* `%s` (%s, %s)", p.Name, typeStr, req))

	if p.Default != nil {
		defaultStr := fmt.Sprintf("%v", p.Default)
		if p.Type == "string" || strings.HasPrefix(p.Type, "enum") {
			defaultStr = fmt.Sprintf("%q", p.Default)
		}
		sb.WriteString(fmt.Sprintf(" default: %s", defaultStr))
	}

	sb.WriteString("\n")
	return sb.String()
}

func formatAction(a ast.ActionSpec, num int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d. **%s**\n", num, a.Statement))

	for _, e := range a.Enforce {
		sb.WriteString(fmt.Sprintf("    * **ENFORCE:** %s\n", e))
	}

	if a.Condition != nil {
		keyword := strings.ToUpper(a.Condition.Type)
		sb.WriteString(fmt.Sprintf("    * **%s:** `%s`\n", keyword, a.Condition.Expression))
	}

	return sb.String()
}
