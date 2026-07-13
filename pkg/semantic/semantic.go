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
	a.checkDefaultTypes()
	a.checkEnumDefaults()
	a.checkRequiredDefault()
	a.checkDuplicateEndpoints()
	a.checkDuplicateEndpointParams()
	a.checkEndpointResponseRefs()
	a.checkEndpointRequiredDefault()
	a.checkEndpointMissingRoute()
	a.checkEndpointMissingMethod()
	a.checkEndpointDefaultTypes()
	a.checkEndpointEnumDefaults()
	a.checkDuplicateListeners()
	a.checkDuplicateListenerParams()
	a.checkListenerMissingTopic()
	a.checkListenerReturnStates()
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
				a.addError(ret.Line, "undefined response %q referenced in return statement of endpoint %q", ret.Payload, ep.Name)
			}
		}
	}
}

func (a *Analyzer) checkEndpointRequiredDefault() {
	for _, ep := range a.app.Endpoints {
		for _, p := range ep.Params {
			if p.Required && p.Default != nil {
				a.addWarning(p.Line, "parameter %q in endpoint %q is required and has a default value — the default is redundant", p.Name, ep.Name)
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

func (a *Analyzer) checkEndpointDefaultTypes() {
	for _, ep := range a.app.Endpoints {
		for _, p := range ep.Params {
			if p.Default == nil {
				continue
			}
			defaultStr := fmt.Sprintf("%v", p.Default)

			switch p.Type {
			case "int":
				if _, err := strconv.Atoi(defaultStr); err != nil {
					a.addError(p.Line, "parameter %q in endpoint %q has type int but default value %q is not an integer", p.Name, ep.Name, defaultStr)
				}
			case "float":
				if _, err := strconv.ParseFloat(defaultStr, 64); err != nil {
					a.addError(p.Line, "parameter %q in endpoint %q has type float but default value %q is not a number", p.Name, ep.Name, defaultStr)
				}
			case "bool":
				if defaultStr != "true" && defaultStr != "false" {
					a.addError(p.Line, "parameter %q in endpoint %q has type bool but default value %q is not true or false", p.Name, ep.Name, defaultStr)
				}
			}
		}
	}
}

func (a *Analyzer) checkEndpointEnumDefaults() {
	for _, ep := range a.app.Endpoints {
		for _, p := range ep.Params {
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
				a.addError(p.Line, "default value %q for parameter %q in endpoint %q is not in the enum list %v", defaultStr, p.Name, ep.Name, values)
			}
		}
	}
}

func (a *Analyzer) checkDuplicateListeners() {
	seen := make(map[string]int)
	for _, l := range a.app.Listeners {
		if prevLine, ok := seen[l.Name]; ok {
			a.addError(l.Line, "duplicate listener %q (first defined at line %d)", l.Name, prevLine)
		} else {
			seen[l.Name] = l.Line
		}
	}
}

func (a *Analyzer) checkDuplicateListenerParams() {
	for _, l := range a.app.Listeners {
		seen := make(map[string]int)
		for _, p := range l.Params {
			if prevLine, ok := seen[p.Name]; ok {
				a.addError(p.Line, "duplicate parameter %q in listener %q (first defined at line %d)", p.Name, l.Name, prevLine)
			} else {
				seen[p.Name] = p.Line
			}
		}
	}
}

func (a *Analyzer) checkListenerMissingTopic() {
	for _, l := range a.app.Listeners {
		if l.Topic == "" {
			a.addError(l.Line, "listener %q is missing required 'topic'", l.Name)
		}
	}
}

func (a *Analyzer) checkListenerReturnStates() {
	validStates := map[string]bool{
		"ack":   true,
		"nack":  true,
		"retry": true,
		"dlq":   true,
	}
	for _, l := range a.app.Listeners {
		for _, ret := range l.Returns {
			if !validStates[ret.State] {
				a.addError(ret.Line, "invalid message state %q in listener %q — expected 'ack', 'nack', 'retry', or 'dlq'", ret.State, l.Name)
			}
		}
	}
}
