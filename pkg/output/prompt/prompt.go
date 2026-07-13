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

	for i, ep := range spec.Endpoints {
		if i > 0 || len(spec.Features) > 0 {
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, "## Endpoint: %s\n", ep.Name)
		if ep.Version != "" {
			fmt.Fprintf(&b, "**Version:** %s\n", ep.Version)
		}
		if ep.Description != "" {
			fmt.Fprintf(&b, "%s\n", ep.Description)
		}
		fmt.Fprintf(&b, "**%s** `%s`\n", ep.Method, ep.Route)

		if len(ep.Params) > 0 {
			b.WriteString("\n### Params\n")
			for _, p := range ep.Params {
				b.WriteString(formatParam(p))
			}
		}

		if len(ep.Responses) > 0 {
			b.WriteString("\n### Response Schemas\n")
			for _, r := range ep.Responses {
				b.WriteString(formatResponse(r))
			}
		}

		if len(ep.Actions) > 0 {
			b.WriteString("\n### Execution Sequence\n")
			b.WriteString("**CRITICAL:** Execute the following steps strictly in order. Do not proceed to the next step until the current one is complete.\n\n")
			for j, a := range ep.Actions {
				b.WriteString(formatAction(a, j+1))
			}
		}

		if len(ep.Returns) > 0 {
			b.WriteString("\n### Return Mapping\n")
			for _, r := range ep.Returns {
				b.WriteString(formatReturn(r))
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

func formatResponse(r ast.ResponseSpec) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("* **%s**\n", r.Name))
	for _, f := range r.Fields {
		sb.WriteString(formatField(f, 1))
	}
	return sb.String()
}

func formatField(f ast.FieldSpec, depth int) string {
	var sb strings.Builder
	indent := strings.Repeat("    ", depth)
	if len(f.SubFields) > 0 {
		sb.WriteString(fmt.Sprintf("%s* `%s` (type: `%s`)\n", indent, f.Name, f.Type))
		for _, sf := range f.SubFields {
			sb.WriteString(formatField(sf, depth+1))
		}
	} else {
		sb.WriteString(fmt.Sprintf("%s* `%s` (type: `%s`)\n", indent, f.Name, f.Type))
	}
	return sb.String()
}

func formatReturn(r ast.ReturnSpec) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("* `return %d`", r.StatusCode))
	if r.Payload != "" {
		sb.WriteString(fmt.Sprintf(" **%s**", r.Payload))
	}
	if r.Condition != nil {
		keyword := strings.ToUpper(r.Condition.Type)
		sb.WriteString(fmt.Sprintf(" | **%s** `%s`", keyword, r.Condition.Expression))
	}
	sb.WriteString("\n")
	return sb.String()
}
