package semantic

import (
	"fmt"
	"regexp"
	"strings"

	"butter/pkg/ast"
)

var enumTypeRe = regexp.MustCompile(`^enum\[(.+)\]$`)

var validTypes = map[string]bool{
	"string":  true,
	"integer": true,
	"double":  true,
	"boolean": true,
	"enum":    true,
	"array":   true,
}

var validMethods = map[string]bool{
	"POST":   true,
	"GET":    true,
	"PUT":    true,
	"DELETE": true,
	"PATCH":  true,
}

type Analyzer struct {
	app   *ast.AppSpec
	diags []Diagnostic
}

func Analyze(app *ast.AppSpec) []Diagnostic {
	a := &Analyzer{app: app}
	a.checkDuplicateFeatures()
	a.checkDuplicateParams()
	a.checkParamTypes()
	a.checkEnumValues()
	a.checkDuplicateEndpoints()
	a.checkDuplicateEndpointParams()
	a.checkEndpointResponseRefs()
	a.checkEndpointMissingRoute()
	a.checkEndpointMissingMethod()
	a.checkEndpointMethod()
	a.checkEndpointStatusCodes()
	a.checkEndpointParamTypes()
	a.checkEndpointEnumValues()
	a.checkEndpointResponseFieldTypes()
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

func (a *Analyzer) checkParamTypes() {
	for _, f := range a.app.Features {
		for _, p := range f.Params {
			base := p.Type
			if idx := strings.Index(base, "["); idx >= 0 {
				base = base[:idx]
			}
			if !validTypes[base] {
				a.addError(p.Line, "parameter %q in feature %q has unknown type %q — expected string, integer, double, boolean, enum[...], or array[...]", p.Name, f.Name, p.Type)
			}
		}
	}
}

func (a *Analyzer) checkEnumValues() {
	for _, f := range a.app.Features {
		for _, p := range f.Params {
			if values := extractEnumValues(p.Type); values != nil {
				seen := make(map[string]int)
				for i, v := range values {
					if prev, ok := seen[v]; ok {
						a.addError(p.Line, "enum value %q in parameter %q of feature %q duplicates index %d", v, p.Name, f.Name, prev)
					} else {
						seen[v] = i
					}
				}
			}
		}
	}
}

func (a *Analyzer) checkDuplicateEndpoints() {
	seen := make(map[string]int)
	for _, ep := range a.app.Endpoints {
		if prevLine, ok := seen[ep.Name]; ok {
			a.addError(ep.Line, "duplicate endpoint %q (first defined at line %d)", ep.Name, prevLine)
		} else {
			seen[ep.Name] = ep.Line
		}
	}
}

func (a *Analyzer) checkDuplicateEndpointParams() {
	for _, ep := range a.app.Endpoints {
		seen := make(map[string]int)
		for _, p := range ep.Params {
			if prevLine, ok := seen[p.Name]; ok {
				a.addError(p.Line, "duplicate parameter %q in endpoint %q (first defined at line %d)", p.Name, ep.Name, prevLine)
			} else {
				seen[p.Name] = p.Line
			}
		}
	}
}

func (a *Analyzer) checkEndpointResponseRefs() {
	for _, ep := range a.app.Endpoints {
		responseNames := make(map[string]bool)
		for _, r := range ep.Responses {
			responseNames[r.Name] = true
		}
		for _, ret := range ep.Returns {
			if ret.Payload == "" || ret.PayloadIsString {
				continue
			}
			if !responseNames[ret.Payload] {
				a.addError(ret.Line, "undefined response %q referenced in returns of endpoint %q", ret.Payload, ep.Name)
			}
		}
	}
}

func (a *Analyzer) checkEndpointMissingRoute() {
	for _, ep := range a.app.Endpoints {
		if ep.Route == "" {
			a.addError(ep.Line, "endpoint %q is missing required 'route'", ep.Name)
		}
	}
}

func (a *Analyzer) checkEndpointMissingMethod() {
	for _, ep := range a.app.Endpoints {
		if ep.Method == "" {
			a.addError(ep.Line, "endpoint %q is missing required 'method'", ep.Name)
		}
	}
}

func (a *Analyzer) checkEndpointMethod() {
	for _, ep := range a.app.Endpoints {
		if ep.Method != "" && !validMethods[ep.Method] {
			a.addError(ep.Line, "endpoint %q has invalid method %q — expected POST, GET, PUT, DELETE, or PATCH", ep.Name, ep.Method)
		}
	}
}

func (a *Analyzer) checkEndpointStatusCodes() {
	for _, ep := range a.app.Endpoints {
		for _, ret := range ep.Returns {
			if ret.StatusCode < 100 || ret.StatusCode > 599 {
				a.addError(ret.Line, "endpoint %q has invalid HTTP status code %d — must be in the range 100-599", ep.Name, ret.StatusCode)
			}
		}
	}
}

func (a *Analyzer) checkEndpointParamTypes() {
	for _, ep := range a.app.Endpoints {
		for _, p := range ep.Params {
			base := p.Type
			if idx := strings.Index(base, "["); idx >= 0 {
				base = base[:idx]
			}
			if !validTypes[base] {
				a.addError(p.Line, "parameter %q in endpoint %q has unknown type %q — expected string, integer, double, boolean, enum[...], or array[...]", p.Name, ep.Name, p.Type)
			}
		}
	}
}

func (a *Analyzer) checkEndpointEnumValues() {
	for _, ep := range a.app.Endpoints {
		for _, p := range ep.Params {
			if values := extractEnumValues(p.Type); values != nil {
				seen := make(map[string]int)
				for i, v := range values {
					if prev, ok := seen[v]; ok {
						a.addError(p.Line, "enum value %q in parameter %q of endpoint %q duplicates index %d", v, p.Name, ep.Name, prev)
					} else {
						seen[v] = i
					}
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

func (a *Analyzer) checkEndpointResponseFieldTypes() {
	for _, ep := range a.app.Endpoints {
		for _, r := range ep.Responses {
			for _, f := range r.Fields {
				base := f.Type
				if idx := strings.Index(base, "["); idx >= 0 {
					base = base[:idx]
				}
				if !validTypes[base] {
					a.addError(f.Line, "field %q in response %q of endpoint %q has unknown type %q", f.Name, r.Name, ep.Name, f.Type)
				}
			}
		}
	}
}
