package parser

import (
	"fmt"
	"strconv"

	"butter/pkg/ast"
	"butter/pkg/lexer"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) skipNewlines() {
	for p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
		p.nextToken()
	}
}

func (p *Parser) Parse() (*ast.AppSpec, error) {
	appSpec := &ast.AppSpec{}

	seenApp := false
	for p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
			p.nextToken()
			continue
		}

		if p.curToken.Type == lexer.TokenIndent || p.curToken.Type == lexer.TokenDedent {
			p.nextToken()
			continue
		}

		if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "app" {
			if seenApp {
				return nil, fmt.Errorf("line %d: only one 'app' keyword is allowed in a spec file", p.curToken.Line)
			}
			seenApp = true
			p.nextToken()
			if p.curToken.Type != lexer.TokenIdentifier {
				return nil, fmt.Errorf("line %d: expected an application name after 'app'", p.curToken.Line)
			}
			appSpec.App = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "description" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenString {
				return nil, fmt.Errorf("line %d: expected a quoted string for description", p.curToken.Line)
			}
			appSpec.Description = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "version" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenString {
				return nil, fmt.Errorf("line %d: expected quoted version string for the application", p.curToken.Line)
			}
			appSpec.Version = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "rules" {
			rules, err := p.parseRules()
			if err != nil {
				return nil, err
			}
			appSpec.Rules = append(appSpec.Rules, rules...)
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "feature" {
			feat, err := p.parseFeature()
			if err != nil {
				return nil, err
			}
			appSpec.Features = append(appSpec.Features, *feat)
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "endpoint" {
			ep, err := p.parseEndpoint()
			if err != nil {
				return nil, err
			}
			appSpec.Endpoints = append(appSpec.Endpoints, *ep)
		} else {
			return nil, fmt.Errorf("line %d: unexpected '%s' at the top level — expected 'app', 'description', 'version', 'rules', 'feature', or 'endpoint'", p.curToken.Line, p.curToken.Value)
		}
	}

	return appSpec, nil
}

func (p *Parser) parseRules() ([]ast.RuleSpec, error) {
	var rules []ast.RuleSpec
	p.nextToken()
	p.skipNewlines()
	if p.curToken.Type != lexer.TokenIndent {
		return nil, fmt.Errorf("line %d: expected an indented block under 'rules'", p.curToken.Line)
	}
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
			p.nextToken()
			continue
		}
		if p.curToken.Type != lexer.TokenString {
			return nil, fmt.Errorf("line %d: expected a quoted-string rule inside this block, got '%s'", p.curToken.Line, p.curToken.Value)
		}
		rules = append(rules, ast.RuleSpec{Statement: p.curToken.Value, Line: p.curToken.Line})
		p.nextToken()
		if p.curToken.Type == lexer.TokenNewline {
			p.nextToken()
		}
	}
	p.nextToken()
	return rules, nil
}

func (p *Parser) parseFeature() (*ast.FeatureSpec, error) {
	p.nextToken()
	if p.curToken.Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("line %d: expected a feature name after 'feature'", p.curToken.Line)
	}

	feat := &ast.FeatureSpec{Name: p.curToken.Value, Line: p.curToken.Line}
	p.nextToken()

	if p.curToken.Type != lexer.TokenNewline {
		return nil, fmt.Errorf("line %d: expected a newline after the feature name", p.curToken.Line)
	}
	p.skipNewlines()

	if p.curToken.Type != lexer.TokenIndent {
		return nil, fmt.Errorf("line %d: expected an indented block under this feature", p.curToken.Line)
	}
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
			p.nextToken()
			continue
		}

		if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "description" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenString {
				return nil, fmt.Errorf("line %d: expected quoted string for feature description", p.curToken.Line)
			}
			feat.Description = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "version" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenString {
				return nil, fmt.Errorf("line %d: expected quoted version string for the feature", p.curToken.Line)
			}
			feat.Version = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "params" {
			params, err := p.parseParams()
			if err != nil {
				return nil, err
			}
			feat.Params = append(feat.Params, params...)
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "actions" {
			actions, err := p.parseActions()
			if err != nil {
				return nil, err
			}
			feat.Actions = append(feat.Actions, actions...)
		} else {
			return nil, fmt.Errorf("line %d: unexpected '%s' inside feature — expected 'description', 'version', 'params', or 'actions'", p.curToken.Line, p.curToken.Value)
		}
	}

	if p.curToken.Type == lexer.TokenDedent {
		p.nextToken()
	}

	return feat, nil
}

func (p *Parser) parseEndpoint() (*ast.EndpointSpec, error) {
	p.nextToken()
	if p.curToken.Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("line %d: expected an endpoint name after 'endpoint'", p.curToken.Line)
	}

	ep := &ast.EndpointSpec{Name: p.curToken.Value, Line: p.curToken.Line}
	p.nextToken()

	if p.curToken.Type != lexer.TokenString {
		return nil, fmt.Errorf("line %d: expected a quoted route string, e.g. endpoint %s \"/route\"", p.curToken.Line, ep.Name)
	}
	ep.Route = p.curToken.Value
	p.nextToken()

	if p.curToken.Type != lexer.TokenNewline {
		return nil, fmt.Errorf("line %d: expected a newline after the endpoint header", p.curToken.Line)
	}
	p.skipNewlines()

	if p.curToken.Type != lexer.TokenIndent {
		return nil, fmt.Errorf("line %d: expected an indented block under this endpoint", p.curToken.Line)
	}
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
			p.nextToken()
			continue
		}

		if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "description" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenString {
				return nil, fmt.Errorf("line %d: expected quoted string for endpoint description", p.curToken.Line)
			}
			ep.Description = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "version" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenString {
				return nil, fmt.Errorf("line %d: expected quoted version string for the endpoint", p.curToken.Line)
			}
			ep.Version = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "method" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenIdentifier {
				return nil, fmt.Errorf("line %d: expected an HTTP method after 'method'", p.curToken.Line)
			}
			ep.Method = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "params" {
			params, err := p.parseParams()
			if err != nil {
				return nil, err
			}
			ep.Params = append(ep.Params, params...)
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "responses" {
			responses, err := p.parseResponses()
			if err != nil {
				return nil, err
			}
			ep.Responses = append(ep.Responses, responses...)
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "actions" {
			actions, err := p.parseActions()
			if err != nil {
				return nil, err
			}
			ep.Actions = append(ep.Actions, actions...)
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "returns" {
			returns, err := p.parseReturns()
			if err != nil {
				return nil, err
			}
			ep.Returns = append(ep.Returns, returns...)
		} else {
			return nil, fmt.Errorf("line %d: unexpected '%s' inside endpoint — expected 'description', 'version', 'method', 'params', 'responses', 'actions', or 'returns'", p.curToken.Line, p.curToken.Value)
		}
	}

	if p.curToken.Type == lexer.TokenDedent {
		p.nextToken()
	}

	return ep, nil
}

func (p *Parser) parseParams() ([]ast.ParamSpec, error) {
	var params []ast.ParamSpec
	p.nextToken()
	p.skipNewlines()
	if p.curToken.Type != lexer.TokenIndent {
		return nil, fmt.Errorf("line %d: expected an indented block under 'params'", p.curToken.Line)
	}
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
			p.nextToken()
			continue
		}
		if p.curToken.Type != lexer.TokenIdentifier {
			return nil, fmt.Errorf("line %d: expected a parameter name inside this block, got '%s'", p.curToken.Line, p.curToken.Value)
		}
		param := &ast.ParamSpec{Name: p.curToken.Value, Type: "string", Line: p.curToken.Line}
		p.nextToken()
		if p.curToken.Type == lexer.TokenIdentifier {
			param.Type = p.curToken.Value
			p.nextToken()
		}
		params = append(params, *param)
		if p.curToken.Type == lexer.TokenNewline {
			p.nextToken()
		}
	}
	p.nextToken()
	return params, nil
}

func (p *Parser) parseActions() ([]ast.ActionSpec, error) {
	var actions []ast.ActionSpec
	p.nextToken()
	p.skipNewlines()
	if p.curToken.Type != lexer.TokenIndent {
		return nil, fmt.Errorf("line %d: expected an indented block under 'actions'", p.curToken.Line)
	}
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
			p.nextToken()
			continue
		}
		if p.curToken.Type != lexer.TokenString {
			return nil, fmt.Errorf("line %d: expected an action statement string, got '%s'", p.curToken.Line, p.curToken.Value)
		}
		action := &ast.ActionSpec{Statement: p.curToken.Value, Line: p.curToken.Line}
		p.nextToken()

		for p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
			p.nextToken()
		}

		if p.curToken.Type == lexer.TokenIndent {
			p.nextToken()
			for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
				if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
					p.nextToken()
					continue
				}
				if p.curToken.Type != lexer.TokenIdentifier || p.curToken.Value != "enforce" {
					return nil, fmt.Errorf("line %d: expected 'enforce' under this action, got '%s'", p.curToken.Line, p.curToken.Value)
				}
				p.nextToken()
				if p.curToken.Type != lexer.TokenString {
					return nil, fmt.Errorf("line %d: expected a quoted string after 'enforce'", p.curToken.Line)
				}
				action.Enforce = append(action.Enforce, ast.EnforceSpec{Expression: p.curToken.Value, Line: p.curToken.Line})
				p.nextToken()
			}
			p.nextToken()
		}

		actions = append(actions, *action)
	}
	p.nextToken()
	return actions, nil
}

func (p *Parser) parseResponses() ([]ast.ResponseSpec, error) {
	var responses []ast.ResponseSpec
	p.nextToken()
	p.skipNewlines()
	if p.curToken.Type != lexer.TokenIndent {
		return nil, fmt.Errorf("line %d: expected an indented block under 'responses'", p.curToken.Line)
	}
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
			p.nextToken()
			continue
		}

		if p.curToken.Type == lexer.TokenIndent {
			p.nextToken()
			if len(responses) == 0 {
				return nil, fmt.Errorf("line %d: response fields found without a response header", p.curToken.Line)
			}
			idx := len(responses) - 1
			for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
				if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
					p.nextToken()
					continue
				}
				if p.curToken.Type != lexer.TokenIdentifier {
					return nil, fmt.Errorf("line %d: expected a response field name, got '%s'", p.curToken.Line, p.curToken.Value)
				}
				field := &ast.FieldSpec{Name: p.curToken.Value, Type: "string", Line: p.curToken.Line}
				p.nextToken()
				if p.curToken.Type == lexer.TokenIdentifier {
					field.Type = p.curToken.Value
					p.nextToken()
				}
				responses[idx].Fields = append(responses[idx].Fields, *field)
				if p.curToken.Type == lexer.TokenNewline {
					p.nextToken()
				}
			}
			p.nextToken()
			continue
		}

		if p.curToken.Type != lexer.TokenIdentifier {
			return nil, fmt.Errorf("line %d: expected a response name inside this block, got '%s'", p.curToken.Line, p.curToken.Value)
		}
		resp := &ast.ResponseSpec{Name: p.curToken.Value, Line: p.curToken.Line}
		p.nextToken()
		if p.curToken.Type == lexer.TokenIdentifier {
			p.nextToken()
		}
		responses = append(responses, *resp)
		if p.curToken.Type == lexer.TokenNewline {
			p.nextToken()
		}
	}
	p.nextToken()
	return responses, nil
}

func (p *Parser) parseReturns() ([]ast.ReturnSpec, error) {
	var returns []ast.ReturnSpec
	p.nextToken()
	p.skipNewlines()
	if p.curToken.Type != lexer.TokenIndent {
		return nil, fmt.Errorf("line %d: expected an indented block under 'returns'", p.curToken.Line)
	}
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
			p.nextToken()
			continue
		}
		if p.curToken.Type != lexer.TokenIdentifier {
			return nil, fmt.Errorf("line %d: expected an HTTP status code inside this block, got '%s'", p.curToken.Line, p.curToken.Value)
		}
		statusStr := p.curToken.Value
		statusCode, err := strconv.Atoi(statusStr)
		if err != nil {
			return nil, fmt.Errorf("line %d: expected an HTTP status code integer, got '%s'", p.curToken.Line, statusStr)
		}
		ret := &ast.ReturnSpec{StatusCode: statusCode, Line: p.curToken.Line}
		p.nextToken()
		if p.curToken.Type == lexer.TokenString {
			ret.Payload = p.curToken.Value
			ret.PayloadIsString = true
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier {
			ret.Payload = p.curToken.Value
			p.nextToken()
		}
		returns = append(returns, *ret)
		if p.curToken.Type == lexer.TokenNewline {
			p.nextToken()
		}
	}
	p.nextToken()
	return returns, nil
}
