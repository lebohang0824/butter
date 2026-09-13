package parser

import (
	"testing"

	"butter/pkg/lexer"
)

func TestParseAppRulesReturns(t *testing.T) {
	src := `app Todo
  description "Todo application"
  version "1.0.0"

  rules
    "Use React on the frontend"
    "Use Nodejs on the backend"

feature CreateTodo
  description "Todo frontend form"
  version "0.0.1"

  params
    name string
    todo_id integer
    completed boolean

  actions
    "Validate name is not empty"
    "Send data to /todo"

endpoint SaveTodo "/todo"
  description "Save todo to the database"
  version "0.0.1"
  method "POST"

  params
    name string
    todo_id integer

  responses
    TodosResponse
      id integer
      name string

  actions
    "Validate name is not empty"
    "Store data to the database"

  returns
    201 TodosResponse
    500 "Internal server error"
`
	l := lexer.NewLexer(src)
	p := NewParser(l)
	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if spec.App != "Todo" {
		t.Errorf("expected app Todo, got %q", spec.App)
	}
	if len(spec.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(spec.Rules))
	}
	if spec.Rules[0].Statement != "Use React on the frontend" {
		t.Errorf("unexpected rule: %q", spec.Rules[0].Statement)
	}
	if len(spec.Features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(spec.Features))
	}
	feat := spec.Features[0]
	if len(feat.Params) != 3 {
		t.Fatalf("expected 3 params, got %d", len(feat.Params))
	}
	if feat.Params[1].Type != "integer" {
		t.Errorf("expected integer type, got %q", feat.Params[1].Type)
	}
	if feat.Params[2].Type != "boolean" {
		t.Errorf("expected boolean type, got %q", feat.Params[2].Type)
	}
	if len(spec.Endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(spec.Endpoints))
	}
	ep := spec.Endpoints[0]
	if ep.Route != "/todo" {
		t.Errorf("expected route /todo, got %q", ep.Route)
	}
	if ep.Method != "POST" {
		t.Errorf("expected method %q, got %q", "POST", ep.Method)
	}
	if len(ep.Returns) != 2 {
		t.Fatalf("expected 2 returns, got %d", len(ep.Returns))
	}
	if ep.Returns[0].StatusCode != 201 || ep.Returns[0].Payload != "TodosResponse" {
		t.Errorf("unexpected first return: %+v", ep.Returns[0])
	}
	if !ep.Returns[1].PayloadIsString || ep.Returns[1].Payload != "Internal server error" {
		t.Errorf("unexpected second return: %+v", ep.Returns[1])
	}
}

func TestParseActionEnforce(t *testing.T) {
	src := `app A
  description "demo"
  version "1.0.0"

feature F
  description "demo"
  version "1.0.0"

  actions
    "Validate name"
      enforce "Sanitize name before storing"
      enforce "Reject empty strings"
    "Persist"
`
	l := lexer.NewLexer(src)
	p := NewParser(l)
	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	actions := spec.Features[0].Actions
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(actions))
	}
	if len(actions[0].Enforce) != 2 {
		t.Fatalf("expected 2 enforces, got %d", len(actions[0].Enforce))
	}
	if actions[0].Enforce[0].Expression != "Sanitize name before storing" {
		t.Errorf("unexpected enforce: %q", actions[0].Enforce[0].Expression)
	}
	if len(actions[1].Enforce) != 0 {
		t.Errorf("expected no enforce on second action")
	}
}

func TestParseRejectsProductAlias(t *testing.T) {
	src := `product Foo
  description "nope"
`
	l := lexer.NewLexer(src)
	p := NewParser(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for 'product' alias, got nil")
	}
}

func TestParseRejectsUnquotedMethod(t *testing.T) {
	src := `app A
  description "demo"
  version "1.0.0"

endpoint E "/e"
  description "demo"
  version "1.0.0"
  method GET
`
	l := lexer.NewLexer(src)
	p := NewParser(l)
	if _, err := p.Parse(); err == nil {
		t.Fatal("expected error for unquoted method, got nil")
	}
}

func TestParseAcceptsQuotedMethod(t *testing.T) {
	src := `app A
  description "demo"
  version "1.0.0"

endpoint E "/e"
  description "demo"
  version "1.0.0"
  method "GET"
`
	l := lexer.NewLexer(src)
	p := NewParser(l)
	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error for quoted method: %v", err)
	}
	if spec.Endpoints[0].Method != "GET" {
		t.Errorf("expected method GET, got %q", spec.Endpoints[0].Method)
	}
}

func TestParseRejectsSingularRule(t *testing.T) {
	src := `app Foo
  description "demo"

  rules
    "ok"
`
	l := lexer.NewLexer(src)
	p := NewParser(l)
	if _, err := p.Parse(); err != nil {
		t.Fatalf("unexpected error for plural rules: %v", err)
	}

	src2 := `app Foo
  description "demo"

rule "singular"
`
	l2 := lexer.NewLexer(src2)
	p2 := NewParser(l2)
	if _, err := p2.Parse(); err == nil {
		t.Fatal("expected error for singular 'rule', got nil")
	}
}

func TestParseRejectsSingularReturn(t *testing.T) {
	src := `app A
  description "demo"
  version "1.0.0"

endpoint E "/e"
  description "demo"
  version "1.0.0"
  method "GET"

  returns
    200 "ok"
`
	l := lexer.NewLexer(src)
	p := NewParser(l)
	if _, err := p.Parse(); err != nil {
		t.Fatalf("unexpected error for plural returns: %v", err)
	}

	src2 := `app A
  description "demo"
  version "1.0.0"

endpoint E "/e"
  description "demo"
  version "1.0.0"
  method "GET"

  return
    200 "ok"
`
	l2 := lexer.NewLexer(src2)
	p2 := NewParser(l2)
	if _, err := p2.Parse(); err == nil {
		t.Fatal("expected error for singular 'return', got nil")
	}
}
