package semantic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"butter/pkg/ast"
)

var enumTypeRe = regexp.MustCompile(`^enum\[(.+)\]$`)

type Analyzer struct {
	app   *ast.AppSpec
	diags []Diagnostic
}

func Analyze(app *ast.AppSpec) []Diagnostic {
	a := &Analyzer{app: app}
	a.checkDuplicateFeatures()
	a.checkDuplicateParams()
	a.checkConditionRefs()
	a.checkDefaultTypes()
	a.checkEnumDefaults()
	a.checkRequiredDefault()
	return a.diags
}

func (a *Analyzer) addError(line int, format string, args ...interface{}) {
	a.diags = append(a.diags, Diagnostic{
		Line:     line,
		Severity: SemError,
		Message:  fmt.Sprintf(format, args...),
	})
}

func (a *Analyzer) addWarning(line int, format string, args ...interface{}) {
	a.diags = append(a.diags, Diagnostic{
		Line:     line,
		Severity: SemWarning,
		Message:  fmt.Sprintf(format, args...),
	})
}

func (a *Analyzer) checkDuplicateFeatures() {
	seen := make(map[string]int)
	for _, f := range a.app.Features {
		if prevLine, ok := seen[f.Name]; ok {
			a.addError(f.Line, "duplicate feature %q (first defined at line %d)", f.Name, prevLine)
		} else {
			seen[f.Name] = f.Line
		}
	}
}

func (a *Analyzer) checkDuplicateParams() {
	for _, f := range a.app.Features {
		seen := make(map[string]int)
		for _, p := range f.Params {
			if prevLine, ok := seen[p.Name]; ok {
				a.addError(p.Line, "duplicate parameter %q in feature %q (first defined at line %d)", p.Name, f.Name, prevLine)
			} else {
				seen[p.Name] = p.Line
			}
		}
	}
}

func (a *Analyzer) checkDefaultTypes() {
	for _, f := range a.app.Features {
		for _, p := range f.Params {
			if p.Default == nil {
				continue
			}
			defaultStr := fmt.Sprintf("%v", p.Default)

			switch p.Type {
			case "int":
				if _, err := strconv.Atoi(defaultStr); err != nil {
					a.addError(p.Line, "parameter %q in feature %q has type int but default value %q is not an integer", p.Name, f.Name, defaultStr)
				}
			case "float":
				if _, err := strconv.ParseFloat(defaultStr, 64); err != nil {
					a.addError(p.Line, "parameter %q in feature %q has type float but default value %q is not a number", p.Name, f.Name, defaultStr)
				}
			case "bool":
				if defaultStr != "true" && defaultStr != "false" {
					a.addError(p.Line, "parameter %q in feature %q has type bool but default value %q is not true or false", p.Name, f.Name, defaultStr)
				}
			}
		}
	}
}

func extractEnumValues(typeStr string) []string {
	matches := enumTypeRe.FindStringSubmatch(typeStr)
	if matches == nil {
		return nil
	}
	parts := strings.Split(matches[1], ",")
	values := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, `"`)
		values = append(values, p)
	}
	return values
}

func (a *Analyzer) checkEnumDefaults() {
	for _, f := range a.app.Features {
		for _, p := range f.Params {
			values := extractEnumValues(p.Type)
			if values == nil || p.Default == nil {
				continue
			}
			defaultStr := fmt.Sprintf("%v", p.Default)
			found := false
			for _, v := range values {
				if v == defaultStr {
					found = true
					break
				}
			}
			if !found {
				a.addError(p.Line, "default value %q for parameter %q in feature %q is not in the enum list %v", defaultStr, p.Name, f.Name, values)
			}
		}
	}
}

func (a *Analyzer) checkRequiredDefault() {
	for _, f := range a.app.Features {
		for _, p := range f.Params {
			if p.Required && p.Default != nil {
				a.addWarning(p.Line, "parameter %q in feature %q is required and has a default value — the default is redundant", p.Name, f.Name)
			}
		}
	}
}
