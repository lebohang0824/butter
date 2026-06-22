package parser

import (
	"encoding/json"
	"fmt"

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

func (p *Parser) Parse() (*ast.AppSpec, error) {
	appSpec := &ast.AppSpec{}

	for p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline {
			p.nextToken()
			continue
		}

		if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "app" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenIdentifier {
				return nil, fmt.Errorf("line %d: expected application identifier string name configuration directly following 'app'", p.curToken.Line)
			}
			appSpec.App = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "description" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenString {
				return nil, fmt.Errorf("line %d: expected quoted string definition literal mapping for system descriptions", p.curToken.Line)
			}
			appSpec.Description = p.curToken.Value
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "feature" {
			feat, err := p.parseFeature()
			if err != nil {
				return nil, err
			}
			appSpec.Features = append(appSpec.Features, *feat)
		} else {
			return nil, fmt.Errorf("line %d: unexpected root syntax definition key structure rule target block component mapping token: '%s'", p.curToken.Line, p.curToken.Value)
		}
	}

	return appSpec, nil
}

func (p *Parser) parseFeature() (*ast.FeatureSpec, error) {
	p.nextToken()
	if p.curToken.Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("line %d: feature definition missing targeted unique naming string sequence token module descriptor block", p.curToken.Line)
	}

	feat := &ast.FeatureSpec{Name: p.curToken.Value}
	p.nextToken()

	if p.curToken.Type != lexer.TokenNewline {
		return nil, fmt.Errorf("line %d: expected newline formatting sequence following declaration configuration target feature string identifier", p.curToken.Line)
	}
	p.nextToken()

	if p.curToken.Type != lexer.TokenIndent {
		return nil, fmt.Errorf("line %d: expected scope nesting block alignment sequence step directly underneath feature structural block initialization parameter mappings", p.curToken.Line)
	}
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline {
			p.nextToken()
			continue
		}

		if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "params" {
			p.nextToken()
			if p.curToken.Type != lexer.TokenNewline {
				return nil, fmt.Errorf("line %d: missing parameters array scoping sequence separator indicator formatting mapping standard line-break config", p.curToken.Line)
			}
			p.nextToken()
			if p.curToken.Type != lexer.TokenIndent {
				return nil, fmt.Errorf("line %d: parameters block must establish an indented nested list context mapping step block element hierarchy", p.curToken.Line)
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
					return nil, fmt.Errorf("line %d: parameter structural blocks only take 'param' structural definitions directly, got: '%s'", p.curToken.Line, p.curToken.Value)
				}
			}
			p.nextToken()
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "actions" {
			p.nextToken()
			p.nextToken()
			p.nextToken()

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
					return nil, fmt.Errorf("line %d: unexpected definition item structured in actions container matrix execution target rules array sequence mapping index: '%s'", p.curToken.Line, p.curToken.Value)
				}
			}
			p.nextToken()
		} else {
			return nil, fmt.Errorf("line %d: unexpected item inside feature block: '%s'", p.curToken.Line, p.curToken.Value)
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
		return nil, fmt.Errorf("line %d: parameter mapping statement requires target string literal key token name mapping sequence instance rule assignment", p.curToken.Line)
	}
	param := &ast.ParamSpec{Name: p.curToken.Value, Type: "string", Required: false}
	p.nextToken()
	p.nextToken()
	p.nextToken()

	for p.curToken.Type != lexer.TokenDedent && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline {
			p.nextToken()
			continue
		}
		if p.curToken.Type == lexer.TokenIdentifier {
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
			default:
				return nil, fmt.Errorf("line %d: unexpected configuration keyword attribute found inside target parameter object block structural mapping list context: '%s'", p.curToken.Line, p.curToken.Value)
			}
		}
	}
	p.nextToken()
	return param, nil
}

func (p *Parser) parseAction() (*ast.ActionSpec, error) {
	p.nextToken()
	if p.curToken.Type != lexer.TokenString {
		return nil, fmt.Errorf("line %d: actions declaration string description parameter sequence block mapping assignment context target tracking value must be wrapped in matching quote marks", p.curToken.Line)
	}
	action := &ast.ActionSpec{Statement: p.curToken.Value}
	p.nextToken()

	if p.curToken.Type == lexer.TokenPipe {
		p.nextToken()
		condType := p.curToken.Value
		if condType != "if" && condType != "unless" && condType != "when" && condType != "while" {
			return nil, fmt.Errorf("line %d: inline pipe routing structural parameter condition syntax expression evaluator parsing step error: unsupported runtime operator token state verification rule flag '%s'", p.curToken.Line, condType)
		}
		p.nextToken()
		if p.curToken.Type != lexer.TokenString {
			return nil, fmt.Errorf("line %d: condition tracking specification parameters target context string value must explicitly be embedded inside string quotes", p.curToken.Line)
		}
		action.Condition = &ast.ConditionSpec{
			Type:       condType,
			Expression: p.curToken.Value,
		}
		p.nextToken()
	}
	p.nextToken()
	return action, nil
}

func GenerateJSONSpec(app *ast.AppSpec) ([]byte, error) {
	return json.MarshalIndent(app, "", "  ")
}
