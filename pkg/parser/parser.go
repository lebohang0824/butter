package parser

import (
	"fmt"
	"regexp"
	"strconv"

	"butter/pkg/ast"
	"butter/pkg/lexer"
)

var validateRuleRe = regexp.MustCompile(`^\s*(>=?|<=?|={1,2}|!=|=<)\s*[0-9]+(\.[0-9]+)?\s*$`)

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
	for p.curToken.Type == lexer.TokenNewline {
		p.nextToken()
	}
}

func (p *Parser) Parse() (*ast.AppSpec, error) {
	appSpec := &ast.AppSpec{}

	for p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline {
			p.nextToken()
			continue
		}

		if p.curToken.Type == lexer.TokenIdentifier && (p.curToken.Value == "app" || p.curToken.Value == "product") {
			p.nextToken()
			if p.curToken.Type != lexer.TokenIdentifier {
				return nil, fmt.Errorf("line %d: expected an application name after '%s'", p.curToken.Line, p.curToken.Value)
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
			return nil, fmt.Errorf("line %d: unexpected '%s' at the top level — expected 'app' (or 'product'), 'description', 'version', 'feature', or 'endpoint'", p.curToken.Line, p.curToken.Value)
		}
	}

	return appSpec, nil
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
		if p.curToken.Type == lexer.TokenNewline {
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
			p.nextToken()
			if p.curToken.Type != lexer.TokenNewline {
				return nil, fmt.Errorf("line %d: expected a newline after 'params'", p.curToken.Line)
			}
			p.skipNewlines()
			if p.curToken.Type != lexer.TokenIndent {
				return nil, fmt.Errorf("line %d: expected an indented block under 'params'", p.curToken.Line)
			}
			p.nextToken()

			for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
				if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "param" {
					param, err := p.parseParam()
					if err != nil {
						return nil, err
					}
					feat.Params = append(feat.Params, *param)
				} else if p.curToken.Type == lexer.TokenNewline {
					p.nextToken()
				} else {
					return nil, fmt.Errorf("line %d: expected 'param' inside this block, got '%s'", p.curToken.Line, p.curToken.Value)
				}
			}
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "actions" {
			p.nextToken()
			p.skipNewlines()
			if p.curToken.Type == lexer.TokenIndent {
				p.nextToken()
			}

			for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
				if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "action" {
					action, err := p.parseAction()
					if err != nil {
						return nil, err
					}
					feat.Actions = append(feat.Actions, *action)
				} else if p.curToken.Type == lexer.TokenNewline {
					p.nextToken()
				} else {
					return nil, fmt.Errorf("line %d: expected 'action' inside this block, got '%s'", p.curToken.Line, p.curToken.Value)
				}
			}
			p.nextToken()
		} else {
			return nil, fmt.Errorf("line %d: unexpected '%s' inside feature — expected 'description', 'version', 'params', or 'actions'", p.curToken.Line, p.curToken.Value)
		}
	}

	if p.curToken.Type == lexer.TokenDedent {
		p.nextToken()
	}

	return feat, nil
}

func (p *Parser) parseParam() (*ast.ParamSpec, error) {
	p.nextToken()
	if p.curToken.Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("line %d: expected a parameter name after 'param'", p.curToken.Line)
	}
	param := &ast.ParamSpec{Name: p.curToken.Value, Type: "string", Required: false, Line: p.curToken.Line}
	p.nextToken()
	p.skipNewlines()
	if p.curToken.Type == lexer.TokenIndent {
		p.nextToken()
	}

	var validateLine int
	var lengthLine int
	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline {
			p.nextToken()
			continue
		}
		switch {
		case p.curToken.Type == lexer.TokenIdentifier:
			switch p.curToken.Value {
			case "type":
				p.nextToken()
				param.Type = p.curToken.Value
				p.nextToken()
			case "required":
				p.nextToken()
				param.Required = (p.curToken.Value == "true")
				p.nextToken()
			case "default":
				p.nextToken()
				if p.curToken.Type == lexer.TokenString {
					param.Default = p.curToken.Value
				} else {
					param.Default = p.curToken.Value
				}
				p.nextToken()
			case "validate":
				p.nextToken()
				if p.curToken.Type != lexer.TokenString {
					return nil, fmt.Errorf("line %d: validate rule must be a quoted string", p.curToken.Line)
				}
				if !validateRuleRe.MatchString(p.curToken.Value) {
					return nil, fmt.Errorf("line %d: invalid validate rule %q — must be a numeric comparison like \">0\", \">=1\", \"=<100\", \"!=5\"", p.curToken.Line, p.curToken.Value)
				}
				if validateLine == 0 {
					validateLine = p.curToken.Line
				}
				param.Validate = append(param.Validate, p.curToken.Value)
				p.nextToken()
			case "length":
				p.nextToken()
				n, err := strconv.Atoi(p.curToken.Value)
				if err != nil || n < 1 {
					return nil, fmt.Errorf("line %d: length must be a positive integer, got %q", p.curToken.Line, p.curToken.Value)
				}
				lengthLine = p.curToken.Line
				param.Length = n
				p.nextToken()
			default:
				return nil, fmt.Errorf("line %d: unexpected '%s' for this parameter — expected 'type', 'required', 'default', 'validate', or 'length'", p.curToken.Line, p.curToken.Value)
			}
		default:
			return nil, fmt.Errorf("line %d: unexpected token %s in parameter fields", p.curToken.Line, p.curToken.Type)
		}
	}
	if len(param.Validate) > 0 && param.Type != "int" && param.Type != "float" {
		return nil, fmt.Errorf("line %d: validate rules require numeric type (int or float), got %q", validateLine, param.Type)
	}
	if param.Length > 0 && len(param.Validate) > 0 {
		return nil, fmt.Errorf("line %d: length and validate cannot be used together on the same parameter", lengthLine)
	}
	if param.Length > 0 && param.Type != "int" && param.Type != "float" {
		return nil, fmt.Errorf("line %d: length requires numeric type (int or float), got %q", lengthLine, param.Type)
	}
	p.nextToken()
	return param, nil
}

func (p *Parser) parseAction() (*ast.ActionSpec, error) {
	actionLine := p.curToken.Line
	p.nextToken()
	if p.curToken.Type != lexer.TokenString {
		return nil, fmt.Errorf("line %d: action statement must be a quoted string", p.curToken.Line)
	}
	action := &ast.ActionSpec{Statement: p.curToken.Value, Line: actionLine}
	p.nextToken()

	if p.curToken.Type == lexer.TokenPipe {
		p.nextToken()
		condType := p.curToken.Value
		if condType != "if" && condType != "unless" && condType != "when" && condType != "while" {
			return nil, fmt.Errorf("line %d: unsupported condition '%s' after '|' — expected if, unless, when, or while", p.curToken.Line, condType)
		}
		condLine := p.curToken.Line
		p.nextToken()
		if p.curToken.Type != lexer.TokenString {
			return nil, fmt.Errorf("line %d: condition expression after '|' must be a quoted string", p.curToken.Line)
		}
		action.Condition = &ast.ConditionSpec{
			Type:       condType,
			Expression: p.curToken.Value,
			Line:       condLine,
		}
		p.nextToken()
	}
	p.nextToken()

	if p.curToken.Type == lexer.TokenIndent {
		p.nextToken()

		for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
			if p.curToken.Type == lexer.TokenNewline {
				p.nextToken()
				continue
			}
			if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "enforce" {
				p.nextToken()
				if p.curToken.Type != lexer.TokenString {
					return nil, fmt.Errorf("line %d: action enforce must be a quoted string", p.curToken.Line)
				}
				action.Enforce = append(action.Enforce, p.curToken.Value)
				p.nextToken()
				continue
			}
			return nil, fmt.Errorf("line %d: unexpected '%s' inside action block", p.curToken.Line, p.curToken.Value)
		}
		p.nextToken()
	}

	return action, nil
}

func (p *Parser) parseEndpoint() (*ast.EndpointSpec, error) {
	p.nextToken()
	if p.curToken.Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("line %d: expected an endpoint name after 'endpoint'", p.curToken.Line)
	}

	ep := &ast.EndpointSpec{Name: p.curToken.Value, Line: p.curToken.Line}
	p.nextToken()

	if p.curToken.Type != lexer.TokenNewline {
		return nil, fmt.Errorf("line %d: expected a newline after the endpoint name", p.curToken.Line)
	}
	p.skipNewlines()

	if p.curToken.Type != lexer.TokenIndent {
		return nil, fmt.Errorf("line %d: expected an indented block under this endpoint", p.curToken.Line)
	}
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline {
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
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "route" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenString {
				return nil, fmt.Errorf("line %d: expected quoted string for route", p.curToken.Line)
			}
			ep.Route = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "method" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenString {
				return nil, fmt.Errorf("line %d: expected quoted string for HTTP method", p.curToken.Line)
			}
			ep.Method = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "params" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenNewline {
				return nil, fmt.Errorf("line %d: expected a newline after 'params'", p.curToken.Line)
			}
			p.skipNewlines()
			if p.curToken.Type != lexer.TokenIndent {
				return nil, fmt.Errorf("line %d: expected an indented block under 'params'", p.curToken.Line)
			}
			p.nextToken()

			for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
				if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "param" {
					param, err := p.parseParam()
					if err != nil {
						return nil, err
					}
					ep.Params = append(ep.Params, *param)
				} else if p.curToken.Type == lexer.TokenNewline {
					p.nextToken()
				} else {
					return nil, fmt.Errorf("line %d: expected 'param' inside this block, got '%s'", p.curToken.Line, p.curToken.Value)
				}
			}
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "responses" {
			p.nextToken()
			p.skipNewlines()
			if p.curToken.Type != lexer.TokenIndent {
				return nil, fmt.Errorf("line %d: expected an indented block under 'responses'", p.curToken.Line)
			}
			p.nextToken()

			for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
				if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "response" {
					resp, err := p.parseResponse()
					if err != nil {
						return nil, err
					}
					ep.Responses = append(ep.Responses, *resp)
				} else if p.curToken.Type == lexer.TokenNewline {
					p.nextToken()
				} else {
					return nil, fmt.Errorf("line %d: expected 'response' inside this block, got '%s'", p.curToken.Line, p.curToken.Value)
				}
			}
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "actions" {
			p.nextToken()
			p.skipNewlines()
			if p.curToken.Type == lexer.TokenIndent {
				p.nextToken()
			}

			for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
				if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "action" {
					action, err := p.parseAction()
					if err != nil {
						return nil, err
					}
					ep.Actions = append(ep.Actions, *action)
				} else if p.curToken.Type == lexer.TokenNewline {
					p.nextToken()
				} else {
					return nil, fmt.Errorf("line %d: expected 'action' inside this block, got '%s'", p.curToken.Line, p.curToken.Value)
				}
			}
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "returns" {
			p.nextToken()
			p.skipNewlines()
			if p.curToken.Type != lexer.TokenIndent {
				return nil, fmt.Errorf("line %d: expected an indented block under 'returns'", p.curToken.Line)
			}
			p.nextToken()

			for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
				if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "return" {
					ret, err := p.parseReturn()
					if err != nil {
						return nil, err
					}
					ep.Returns = append(ep.Returns, *ret)
				} else if p.curToken.Type == lexer.TokenNewline {
					p.nextToken()
				} else {
					return nil, fmt.Errorf("line %d: expected 'return' inside this block, got '%s'", p.curToken.Line, p.curToken.Value)
				}
			}
			p.nextToken()
		} else {
			return nil, fmt.Errorf("line %d: unexpected '%s' inside endpoint — expected 'description', 'version', 'route', 'method', 'params', 'responses', 'actions', or 'returns'", p.curToken.Line, p.curToken.Value)
		}
	}

	if p.curToken.Type == lexer.TokenDedent {
		p.nextToken()
	}

	return ep, nil
}

func (p *Parser) parseResponse() (*ast.ResponseSpec, error) {
	p.nextToken()
	if p.curToken.Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("line %d: expected a response name after 'response'", p.curToken.Line)
	}

	resp := &ast.ResponseSpec{Name: p.curToken.Value, Line: p.curToken.Line}
	p.nextToken()
	p.skipNewlines()

	if p.curToken.Type != lexer.TokenIndent {
		return nil, fmt.Errorf("line %d: expected an indented block under response %q", p.curToken.Line, resp.Name)
	}
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline {
			p.nextToken()
			continue
		}
		if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "field" {
			field, err := p.parseField()
			if err != nil {
				return nil, err
			}
			resp.Fields = append(resp.Fields, *field)
		} else {
			return nil, fmt.Errorf("line %d: expected 'field' inside response block, got '%s'", p.curToken.Line, p.curToken.Value)
		}
	}
	p.nextToken()
	return resp, nil
}

func (p *Parser) parseField() (*ast.FieldSpec, error) {
	p.nextToken()
	if p.curToken.Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("line %d: expected a field name after 'field'", p.curToken.Line)
	}

	field := &ast.FieldSpec{Name: p.curToken.Value, Type: "string", Line: p.curToken.Line}
	p.nextToken()
	p.skipNewlines()

	if p.curToken.Type != lexer.TokenIndent {
		return field, nil
	}
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline {
			p.nextToken()
			continue
		}
		switch p.curToken.Value {
		case "type":
			p.nextToken()
			field.Type = p.curToken.Value
			p.nextToken()
		case "field":
			subField, err := p.parseField()
			if err != nil {
				return nil, err
			}
			field.SubFields = append(field.SubFields, *subField)
		default:
			return nil, fmt.Errorf("line %d: unexpected '%s' inside field block — expected 'type' or 'field'", p.curToken.Line, p.curToken.Value)
		}
	}
	p.nextToken()
	return field, nil
}

func (p *Parser) parseReturn() (*ast.ReturnSpec, error) {
	retLine := p.curToken.Line
	p.nextToken()

	statusStr := p.curToken.Value
	statusCode, err := strconv.Atoi(statusStr)
	if err != nil {
		return nil, fmt.Errorf("line %d: expected an HTTP status code integer after 'return', got '%s'", p.curToken.Line, statusStr)
	}
	ret := &ast.ReturnSpec{StatusCode: statusCode, Line: retLine}
	p.nextToken()

	if p.curToken.Type == lexer.TokenString {
		ret.Payload = p.curToken.Value
		ret.PayloadIsString = true
		p.nextToken()
	} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value != "if" && p.curToken.Value != "unless" {
		ret.Payload = p.curToken.Value
		p.nextToken()
	}

	if p.curToken.Type == lexer.TokenPipe {
		p.nextToken()
		condType := p.curToken.Value
		if condType != "if" && condType != "unless" {
			return nil, fmt.Errorf("line %d: unsupported condition '%s' after '|' — expected 'if' or 'unless'", p.curToken.Line, condType)
		}
		condLine := p.curToken.Line
		p.nextToken()
		if p.curToken.Type != lexer.TokenString {
			return nil, fmt.Errorf("line %d: condition expression after '|' must be a quoted string", p.curToken.Line)
		}
		ret.Condition = &ast.ConditionSpec{
			Type:       condType,
			Expression: p.curToken.Value,
			Line:       condLine,
		}
		p.nextToken()
	}
	p.nextToken()

	return ret, nil
}
