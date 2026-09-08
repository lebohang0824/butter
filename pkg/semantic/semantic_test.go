package semantic

import (
	"testing"

	"butter/pkg/ast"
)

func TestAnalyzeValidEndpoint(t *testing.T) {
	spec := &ast.AppSpec{
		App: "A",
		Endpoints: []ast.EndpointSpec{
			{
				Name:    "E",
				Route:   "/e",
				Method:  "POST",
				Params:  []ast.ParamSpec{{Name: "name", Type: "string"}, {Name: "count", Type: "integer"}},
				Actions: []ast.ActionSpec{{Statement: "do it"}},
				Responses: []ast.ResponseSpec{
					{Name: "Ok", Fields: []ast.FieldSpec{{Name: "id", Type: "integer"}}},
				},
				Returns: []ast.ReturnSpec{{StatusCode: 201, Payload: "Ok"}, {StatusCode: 500, Payload: "boom", PayloadIsString: true}},
			},
		},
	}
	diags := Analyze(spec)
	for _, d := range diags {
		if d.Severity == SemError {
			t.Fatalf("unexpected error: %v", d)
		}
	}
}

func TestAnalyzeUnknownParamType(t *testing.T) {
	spec := &ast.AppSpec{
		App: "A",
		Features: []ast.FeatureSpec{
			{Name: "F", Params: []ast.ParamSpec{{Name: "x", Type: "wibble"}}},
		},
	}
	diags := Analyze(spec)
	if !hasError(diags) {
		t.Fatal("expected error for unknown param type")
	}
}

func TestAnalyzeInvalidMethod(t *testing.T) {
	spec := &ast.AppSpec{
		App: "A",
		Endpoints: []ast.EndpointSpec{
			{Name: "E", Route: "/e", Method: "FETCH"},
		},
	}
	diags := Analyze(spec)
	if !hasError(diags) {
		t.Fatal("expected error for invalid method")
	}
}

func TestAnalyzeOutOfRangeStatusCode(t *testing.T) {
	spec := &ast.AppSpec{
		App: "A",
		Endpoints: []ast.EndpointSpec{
			{Name: "E", Route: "/e", Method: "GET", Returns: []ast.ReturnSpec{{StatusCode: 999, Payload: "ok", PayloadIsString: true}}},
		},
	}
	diags := Analyze(spec)
	if !hasError(diags) {
		t.Fatal("expected error for out-of-range status code")
	}
}

func TestAnalyzeUndefinedResponseRef(t *testing.T) {
	spec := &ast.AppSpec{
		App: "A",
		Endpoints: []ast.EndpointSpec{
			{Name: "E", Route: "/e", Method: "GET", Returns: []ast.ReturnSpec{{StatusCode: 200, Payload: "Missing"}}},
		},
	}
	diags := Analyze(spec)
	if !hasError(diags) {
		t.Fatal("expected error for undefined response reference")
	}
}

func TestAnalyzeEnumDuplicate(t *testing.T) {
	spec := &ast.AppSpec{
		App: "A",
		Features: []ast.FeatureSpec{
			{Name: "F", Params: []ast.ParamSpec{{Name: "p", Type: `enum["a", "a"]`}}},
		},
	}
	diags := Analyze(spec)
	if !hasError(diags) {
		t.Fatal("expected error for duplicate enum value")
	}
}

func hasError(diags []Diagnostic) bool {
	for _, d := range diags {
		if d.Severity == SemError {
			return true
		}
	}
	return false
}
