# Butter Development & Architecture Documentation

This document serves as the comprehensive development, architecture, and implementation blueprint for **Butter**, an indentation-aware domain-specific language (DSL) compiler. 

`butter` is a high-performance command-line interface (CLI) tool written in **Go** using the **Cobra** framework. It compiles clean, human-readable specification files (`.butter`) into structured, production-ready, pretty-printed JSON configurations (`.json`).

---

## 1. Executive Summary & Design Philosophy

Modern application architectures often require declarative schemas to define features, validation rules, application parameters, or workflows. While JSON and YAML are industry standards for data exchange, they can become verbose, deeply nested, error-prone, and visually exhausting for human engineers to write and maintain from scratch.

**Butter** solves this by providing an elegant, minimalistic language interface heavily inspired by Python’s significant indentation. It strips away syntactic noise—such as trailing commas, brackets, curly braces, and redundant tags—allowing system architects and developers to declare application specifications cleanly.

### Core Objectives
* **Zero Dependencies for Compilation:** The lexing, parsing, and semantic validation engines are entirely hand-written in native Go, ensuring extreme runtime speed, predictable compilation pathways, and zero supply-chain vulnerabilities.
* **Significant Indentation:** Structural scope is driven entirely by whitespaces (spaces or tabs), optimizing readability.
* **Rich Conditionals:** Built-in keywords natively capture diverse run-time semantic states (`if`, `unless`, `when`, `while`).
* **Developer Ergonomics:** Paired with an easily distributable VS Code extension framework for rich syntax highlighting and structural auto-indentation.

---

## 2. Language Specification & Grammar

The Butter grammar is defined cleanly by key blocks, nested structural declarations, line-breaks, and indentation steps.

### 2.1 Keyword Dictionary

| Keyword | Context | Semantic Purpose |
| :--- | :--- | :--- |
| `app` | Top-level | Defines the namespace or structural root of the configuration. |
| `description` | Top/Block-level | Provides context or documentation string metadata. |
| `feature` | Block-level | Declares a sub-system module, API endpoint, or discrete capability. |
| `params` | Block-level | A dedicated container block specifying input definitions. |
| `param` | Item-level | Declares a discrete parameter variable name. |
| `type` | Parameter field| Dictates data constraints (e.g., `string`, `int`, `enum[...]`). |
| `required` | Parameter field| Boolean validation rule (`true` or `false`). |
| `default` | Parameter field| Explicit fallback value if the parameter is omitted. |
| `actions` | Block-level | A dedicated container block specifying execution routines. |
| `action` | Item-level | Declares a logical execution string or mutation step. |

### 2.2 Semantic Conditionals
Butter expands standard evaluation logic beyond a simple `if` condition, offering native semantic blocks that map perfectly to backend execution engines:

* **`if`**: The action executes **only if** the target predicate expression evaluates to `true`.
* **`unless`**: The action executes **except when** the predicate expression evaluates to `true` (an elegant structural inversion shortcut for `if not`).
* **`when`**: Reactive or event-driven hook. Indicates the action triggers **asynchronously upon** an external event or state shift.
* **`while`**: Active polling or operational state persistence. The action requires this state condition to remain continuously active throughout execution.

### 2.3 Syntactic Layout Example (`demo.butter`)

```butter
# Global application declaration
app OrderProcessor
description "Handles high-throughput retail checkout workflows safely"

feature ProcessPayment
  params
    param OrderID
      type string
      required true
    param Amount
      type float
      required true
    param PaymentMethod
      type enum["CreditCard", "Crypto", "BankTransfer"]
      default "CreditCard"
    param AccountNotes
      default "Standard processing sequence" # Implicit type inference: string
      
  actions
    action "Validate routing balance metrics"
    action "Apply cryptocurrency transaction surcharge" | when "PaymentMethod is set to Crypto"
    action "Flag transaction for manual risk mitigation review" | if "Amount > 10000"
    action "Bypass fraud detection ledger verification" | unless "Amount > 50"
    action "Maintain continuous transaction ledger heartbeat" | while "Gateway Connection is unstable"
```

---

## 3. Compiler Pipeline Architecture

The implementation splits the compilation phases cleanly across decoupled steps, isolating string scanning from tree parsing.

```
[ .butter file ]
       │
       ▼
 ┌───────────┐
 │   Lexer   │ <--- Tracks Indentation Stack & emits INDENT/DEDENT/NEWLINE
 └─────┬─────┘
       │ (Stream of Tokens)
       ▼
 ┌───────────┐
 │  Parser   │ <--- Stateful Recursive Descent State Machine
 └─────┬─────┘
       │ (Abstract Syntax Tree Structure)
       ▼
 ┌───────────┐
 │JSON Engine│ <--- Go json.MarshalIndent Serialization Block
 └─────┬─────┘
       │
       ▼
 [ .json file ]
```

### 3.1 Lexical Analysis (The Off-side Rule)
Because Butter uses whitespace indentation to mark boundaries, the Lexer reads files sequentially while maintaining a **LIFO Indentation Stack** tracking current space depth levels. 
* When a newline occurs, the lexer scans consecutive leading whitespace characters.
* If the space-count exceeds the value on top of the stack, the Lexer pushes the new count onto the stack and emits an implicit `INDENT` token.
* If the space-count is less than the top of the stack, it pops elements off the stack, emitting a `DEDENT` token for each element popped until a matching level is safely located. Any mismatch throws a syntax error immediately.

### 3.2 Abstract Syntax Tree (AST) Model
The parser constructs a strict root AST graph mapped instantly to Go structures for seamless native serialization.

---

## 4. Step-by-Step Implementation Guide in Go

### 4.1 Project Blueprint & Cobra Setup

Execute these commands within your clean shell interface to scaffold the project structure correctly:

```bash
mkdir -p butter/cmd butter/pkg/ast butter/pkg/lexer butter/pkg/parser
cd butter
go mod init butter
go get github.com/spf13/cobra@latest
```

### 4.2 Main Entrypoint & Command Structuring

Create `main.go` inside the root module directory to wire execution routines into the terminal interface:

```go
// main.go
package main

import (
	"butter/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

Create the implementation files inside the `cmd/` directory to manage flag parsers and platform I/O pipelines.

```go
// cmd/root.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"butter/pkg/lexer"
	"butter/pkg/parser"

	"github.com/spf13/cobra"
)

var outputFile string

var rootCmd = &cobra.Command{
	Use:   "butter",
	Short: "Butter is a high-performance, indentation-aware specification compiler.",
	Long:  `A clean compiler framework that turns minimalist indentation-based .butter specifications into beautifully formatted JSON structures.`,
}

var compileCmd = &cobra.Command{
	Use:   "compile [input file]",
	Short: "Compile a .butter specification file down to pretty JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]
		if !strings.HasSuffix(inputFile, ".butter") {
			return fmt.Errorf("invalid file context: source files must end explicitly with the '.butter' extension")
		}

		// Read application source file directly to string
		content, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read source file: %w", err)
		}

		// Execute Lexical Scanning and AST Construction
		l := lexer.NewLexer(string(content))
		p := parser.NewParser(l)
		appAST, err := p.Parse()
		if err != nil {
			return fmt.Errorf("compilation syntax compilation error:\n%w", err)
		}

		// Convert Abstract Syntax Tree map graph to formatted JSON string
		jsonOutput, err := parser.GenerateJSONSpec(appAST)
		if err != nil {
			return fmt.Errorf("json packaging generation failed: %w", err)
		}

		// Determine target path resolution layout
		if outputFile == "" {
			outputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".json"
		}

		if err := os.WriteFile(outputFile, jsonOutput, 0644); err != nil {
			return fmt.Errorf("failed to write compiled asset to target destination disk: %w", err)
		}

		fmt.Printf("Successfully compiled: %s ──> %s\n", inputFile, outputFile)
		return nil
	},
}

func init() {
	compileCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Specify custom path for output file destination (defaults to input name + .json)")
	rootCmd.AddCommand(compileCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
```

### 4.3 Definition of Core Types (AST Graph Nodes)

Define your layout structures securely within the nested domain module layer.

```go
// pkg/ast/types.go
package ast

type AppSpec struct {
	App         string        `json:"app"`
	Description string        `json:"description,omitempty"`
	Features    []FeatureSpec `json:"features"`
}

type FeatureSpec struct {
	Name    string       `json:"name"`
	Params  []ParamSpec  `json:"params,omitempty"`
	Actions []ActionSpec `json:"actions,omitempty"`
}

type ParamSpec struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Required bool        `json:"required"`
	Default  interface{} `json:"default,omitempty"`
}

type ActionSpec struct {
	Statement string         `json:"statement"`
	Condition *ConditionSpec `json:"condition,omitempty"`
}

type ConditionSpec struct {
	Type       string `json:"type"`       // if, unless, when, while
	Expression string `json:"expression"` // Conditional evaluation string
}
```

### 4.4 Hand-Written Indentation-Aware Lexer Implementation

```go
// pkg/lexer/lexer.go
package lexer

import (
	"fmt"
	"unicode"
)

type TokenType string

const (
	TokenError      TokenType = "ERROR"
	TokenEOF        TokenType = "EOF"
	TokenIdentifier TokenType = "IDENTIFIER"
	TokenString     TokenType = "STRING"
	TokenIndent     TokenType = "INDENT"
	TokenDedent     TokenType = "DEDENT"
	TokenNewline    TokenType = "NEWLINE"
	TokenPipe       TokenType = "PIPE"
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
}

type Lexer struct {
	input       string
	pos         int
	line        int
	indentStack []int
	pendingToks []Token
	isLineStart bool
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:       input,
		line:        1,
		indentStack: []int{0}, // Standard structural base context begins at layer zero
		isLineStart: true,
	}
}

func (l *Lexer) NextToken() Token {
	if len(l.pendingToks) > 0 {
		tok := l.pendingToks[0]
		l.pendingToks = l.pendingToks[1:]
		return tok
	}

	l.skipWhitespaceAndComments()

	if l.pos >= len(l.input) {
		if len(l.indentStack) > 1 {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			return Token{Type: TokenDedent, Line: l.line}
		}
		return Token{Type: TokenEOF, Line: l.line}
	}

	if l.isLineStart {
		l.isLineStart = false
		indent := l.consumeIndentation()
		currentIndent := l.indentStack[len(l.indentStack)-1]

		if indent > currentIndent {
			l.indentStack = append(l.indentStack, indent)
			return Token{Type: TokenIndent, Line: l.line}
		}

		if indent < currentIndent {
			for len(l.indentStack) > 0 && l.indentStack[len(l.indentStack)-1] > indent {
				l.indentStack = l.indentStack[:len(l.indentStack)-1]
				l.pendingToks = append(l.pendingToks, Token{Type: TokenDedent, Line: l.line})
			}
			if l.indentStack[len(l.indentStack)-1] != indent {
				return Token{Type: TokenError, Value: "Indentation compilation alignment tracking error", Line: l.line}
			}
			if len(l.pendingToks) > 0 {
				tok := l.pendingToks[0]
				l.pendingToks = l.pendingToks[1:]
				return tok
			}
		}
	}

	ch := l.input[l.pos]

	if ch == '\n' {
		l.line++
		l.pos++
		l.isLineStart = true
		return Token{Type: TokenNewline, Line: l.line - 1}
	}

	if ch == '|' {
		l.pos++
		return Token{Type: TokenPipe, Value: "|", Line: l.line}
	}

	if ch == '"' {
		return l.readString()
	}

	if isIdentifierStart(ch) {
		return l.readIdentifier()
	}

	l.pos++
	return Token{Type: TokenError, Value: fmt.Sprintf("Unexpected standalone literal token instance: '%c'", ch), Line: l.line}
}

func (l *Lexer) consumeIndentation() int {
	count := 0
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' {
			count++
			l.pos++
		} else if ch == '\t' {
			count += 4 // Normalize tabs cleanly to base standard index parameters
			l.pos++
		} else {
			break
		}
	}
	return count
}

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\r' || ch == '\t' {
			l.pos++
		} else if ch == '#' {
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		} else {
			break
		}
	}
}

func (l *Lexer) readString() Token {
	l.pos++ // Skip quote opening token
	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		l.pos++
	}
	val := l.input[start:l.pos]
	l.pos++ // Skip close quote boundary safely
	return Token{Type: TokenString, Value: val, Line: l.line}
}

func (l *Lexer) readIdentifier() Token {
	start := l.pos
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if isIdentifierPart(ch) {
			l.pos++
		} else {
			break
		}
	}
	return Token{Type: TokenIdentifier, Value: l.input[start:l.pos], Line: l.line}
}

func isIdentifierStart(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isIdentifierPart(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' || ch == '-' || ch == '.' || ch == '[' || ch == ']' || ch == ',' || ch == '"'
}
```

### 4.5 Recursive Descent Stateful Parser Implementation

```go
// pkg/parser/parser.go
package parser

import (
	"encoding/json"
	"fmt"
	"strings"

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
	p.nextToken() // consume 'feature' keyword
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
	p.nextToken() // consume INDENT

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
			p.nextToken() // consume INDENT

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
			p.nextToken() // consume DEDENT
		} else if p.curToken.Type == lexer.TokenIdentifier && p.curToken.Value == "actions" {
			p.nextToken()
			p.nextToken() // consume Newline
			p.nextToken() // consume INDENT

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
			p.nextToken() // consume DEDENT
		} else {
			return nil, fmt.Errorf("line %d: unexpected item inside feature block: '%s'", p.curToken.Line, p.curToken.Value)
		}
	}

	if p.curToken.Type == lexer.TokenDedent {
		p.nextToken() // consume DEDENT
	}

	return feat, nil
}

func (p *Parser) parseParam() (*ast.ParamSpec, error) {
	p.nextToken() // consume 'param' keyword
	if p.curToken.Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("line %d: parameter mapping statement requires target string literal key token name mapping sequence instance rule assignment", p.curToken.Line)
	}
	param := &ast.ParamSpec{Name: p.curToken.Value, Type: "string", Required: false} // Implicit default fallback string fallback configuration standard
	p.nextToken()
	p.nextToken() // consume Newline
	p.nextToken() // consume INDENT

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
	p.nextToken() // consume DEDENT
	return param, nil
}

func (p *Parser) parseAction() (*ast.ActionSpec, error) {
	p.nextToken() // consume 'action' keyword
	if p.curToken.Type != lexer.TokenString {
		return nil, fmt.Errorf("line %d: actions declaration string description parameter sequence block mapping assignment context target tracking value must be wrapped in matching quote marks", p.curToken.Line)
	}
	action := &ast.ActionSpec{Statement: p.curToken.Value}
	p.nextToken()

	if p.curToken.Type == lexer.TokenPipe {
		p.nextToken() // consume '|'
		condType := p.curToken.Value
		if condType != "if" && condType != "unless" && condType != "when" && condType != "while" {
			return nil, fmt.Errorf("line %d: inline pipe routing structural parameter condition syntax expression evaluator parsing step error: unsupported runtime operator token state verification rule flag '%s'", p.curToken.Line, condType)
		}
		p.nextToken() // consume conditional indicator key keyword token
		if p.curToken.Type != lexer.TokenString {
			return nil, fmt.Errorf("line %d: condition tracking specification parameters target context string value must explicitly be embedded inside string quotes", p.curToken.Line)
		}
		action.Condition = &ast.ConditionSpec{
			Type:       condType,
			Expression: p.curToken.Value,
		}
		p.nextToken()
	}
	p.nextToken() // consume Newline
	return action, nil
}

func GenerateJSONSpec(app *ast.AppSpec) ([]byte, error) {
	return json.MarshalIndent(app, "", "  ")
}
```

---

## 5. VS Code Extension Blueprint

To ensure rich syntactic evaluation and seamless configuration workflow ergonomics, use the workspace structure configuration maps below to build your custom IDE system plugin.

### 5.1 Extension Directory Layout
```text
butter-extension/
├── package.json
├── language-configuration.json
└── syntaxes/
    └── butter.tmLanguage.json
```

### 5.2 `package.json`
```json
{
  "name": "butter-extension",
  "displayName": "Butter Configuration Language Support",
  "description": "Syntax structural highlighting and auto-indent configuration mapping engines designed natively for the Butter DSL file ecosystem.",
  "version": "1.0.0",
  "publisher": "butter-io",
  "engines": {
    "vscode": "^1.85.0"
  },
  "categories": [
    "Programming Languages"
  ],
  "contributes": {
    "languages": [{
      "id": "butter",
      "aliases": ["Butter", "butter"],
      "extensions": [".butter"],
      "configuration": "./language-configuration.json"
    }],
    "grammars": [{
      "language": "butter",
      "scopeName": "source.butter",
      "path": "./syntaxes/butter.tmLanguage.json"
    }]
  }
}
```

### 5.3 `language-configuration.json`
```json
{
  "comments": {
    "lineComment": "#"
  },
  "brackets": [
    ["[", "]"]
  ],
  "autoClosingPairs": [
    {"open": "\"", "close": "\""},
    {"open": "[", "close": "]"}
  ],
  "indentationRules": {
    "increaseIndentPattern": "^\\s*(feature|params|actions|param\\s+\\w+)\\b.*$",
    "decreaseIndentPattern": "^\\s*$"
  }
}
```

### 5.4 `syntaxes/butter.tmLanguage.json` (TextMate Architecture Definition)
```json
{
  "$schema": "https://raw.githubusercontent.com/martinring/tmlanguage/master/tmlanguage.json",
  "name": "Butter DSL Language Configuration Map",
  "scopeName": "source.butter",
  "patterns": [
    { "include": "#comments" },
    { "include": "#keywords" },
    { "include": "#conditionals" },
    { "include": "#strings" },
    { "include": "#constants" }
  ],
  "repository": {
    "comments": {
      "match": "#.*$",
      "name": "comment.line.number-sign.butter"
    },
    "keywords": {
      "match": "\\b(app|description|feature|params|param|type|required|default|actions|action)\\b",
      "name": "keyword.control.butter"
    },
    "conditionals": {
      "match": "\\b(if|unless|when|while)\\b",
      "name": "keyword.control.conditional.butter"
    },
    "strings": {
      "name": "string.quoted.double.butter",
      "begin": "\"",
      "end": "\""
    },
    "constants": {
      "match": "\\b(true|false)\\b",
      "name": "constant.language.boolean.butter"
    }
  }
}
```

---

## 6. Verification and Deployment Checklist

### Building the Compiler Binary
Run the following compilation script command sequences in your local shell to bundle the production executable binary cleanly:
```bash
go build -o butter main.go
```

### Comprehensive Testing Workflow Scenario
Create a temporary validation source instance verification file locally matching the specification logic pattern layout structure rules perfectly:

```bash
cat << 'EOF' > test.butter
app AutomatedGatekeeper
description "Defines network gateway configuration maps programmatically"

feature InterceptPayload
  params
    param SourceIP
      type string
      required true
    param ThreatVectorScore
      type float
      required true
    param MitigationAction
      type enum["Drop", "Quarantine", "Pass"]
      default "Quarantine"
  actions
    action "Log full system payload frame data packet maps directly to disk metadata pools"
    action "Route traffic straight into specialized network analysis pools" | if "ThreatVectorScore > 7.5"
    action "Bypass deep validation routing checking layer parameters entirely" | unless "SourceIP == 127.0.0.1"
    action "Halt secondary egress interfaces securely to safely isolate hardware zones" | when "MitigationAction == Drop"
    action "Keep auditing diagnostic internal systems infrastructure pipelines active" | while "ThreatVectorScore >= 5.0"
EOF
```

Execute your newly built local compiler assembly artifact targeting the test environment parameters securely:
```bash
./butter compile test.butter --output output.json
```

Verify your formatted schema results matrix graph output artifact directly to confirm system structural correctness mapping targets:
```bash
cat output.json
```
