![Butter](butter.png)

**Butter** is an intent specification language for AI agents. A `.butter` file describes what software is intended to do, what must happen, what must not happen, which constraints apply, and what outcomes are expected. Butter intentionally leaves implementation technology choices to the implementation context supplied alongside the specification. Compile a spec to a structured prompt, then provide the prompt and the desired implementation context to an AI agent.

---

## Table of Contents

- [Design Philosophy](#design-philosophy)
- [AI Workflow](#ai-workflow)
- [Language Specification](#language-specification)
  - [Keywords](#keywords)
  - [Parameter Syntax](#parameter-syntax)
  - [Action Syntax](#action-syntax)
  - [Response Syntax](#response-syntax)
  - [Return Syntax](#return-syntax)
- [Example](#example)
- [Installation](#installation)
  - [From Source](#from-source)
  - [Install Script](#install-script)
- [Usage](#usage)
- [Compiler Architecture](#compiler-architecture)
- [VS Code Extension](#vs-code-extension)

---

## Design Philosophy

Butter is an **intent specification language for AI agents**. It communicates what software is intended to do, what must happen, what must not happen, which constraints apply, and what outcomes are expected. It does not attempt to describe every implementation detail.

> **Butter specifies what must be true; the implementation context specifies how to make it true.**

Butter stays intentionally small. A new keyword or construct is only useful when the existing language cannot express an important form of intent clearly. Prefer expressing intent with the constructs that already exist.

### Butter and implementation context

Butter communicates application intent:

- business requirements and application behavior
- business rules, constraints, and invariants
- validation requirements
- workflows and ordered actions
- expected outcomes
- API and interface contracts
- boundaries that are part of the application's intent

The implementation context communicates technology and project decisions:

- programming language and framework
- database, libraries, and infrastructure
- deployment environment
- architecture choices and coding conventions
- existing project structure
- technology-specific documentation

Technology choices belong in the implementation context or pre-prompt supplied to the AI agent alongside the Butter specification. The same Butter specification can therefore be implemented with different technology stacks without changing the application intent.

### Core Principles

- **Intent over implementation** — Butter describes what the application must do and what must remain true. It does not prescribe a programming language, framework, database, library, or UI technology.

- **Small and focused** — Existing constructs express behavior, rules, validation, workflows, constraints, interfaces, and outcomes. New syntax is not added merely to make a detail more explicit.

- **Explicit constraints** — `rules` and `enforce` express application-level boundaries. They should not be used for technology choices such as “Use Laravel” or “Use React.”

- **Ordered behavior** — Actions describe a required sequence of behavior. The implementation chooses how to realize that sequence in the selected technology context.

- **Implementation independence** — A Butter specification remains stable when the implementation stack changes. Supply the stack and project conventions separately when asking an agent to implement the specification.

- **Clear, not guaranteed** — A structured specification makes intent easier to communicate and review. It does not guarantee a correct implementation or a particular level of accuracy.

### Implementation context

A Butter specification describes the application's intent. When asking an AI agent to implement it, provide technology and project decisions separately, for example:

```text
Implementation context
- Programming language: PHP
- Framework: Laravel
- UI framework: Vue
- Database: PostgreSQL
- Follow the existing project conventions and architecture
```

The same specification can be paired with a different context:

```text
Implementation context
- Programming language: Python
- Framework: Django
- UI framework: React
- Database: PostgreSQL
- Follow the existing project conventions and architecture
```

The Butter file does not change. Only the implementation context changes.

---

## Language Specification

Butter uses strict 2-space indentation per nesting depth. Comments use `#` or `//`.

```text
# Root Application Specification
app <name>
├── description <string>
├── rules
│   └── <string>
├── feature <name>
└── endpoint <name> <route>

# Subsystem Feature Block Specification
feature <name>
├── description <string>
├── version <string>
├── params
│   └── <name> <type>
└── actions
    └── <string>
        └── enforce <string>

# Synchronous HTTP Endpoint Transport Block Specification
endpoint <name> <route>
├── description <string>
├── version <string>
├── method <string>
├── params
│   └── <name> <type>
├── responses
│   └── <name>
│       └── <name> <type>
├── actions
│   └── <string>
│       └── enforce <string>
└── returns
    └── <code> <ResponseName|string>
```

### Keywords

| Keyword       | Context       | Semantic Purpose |
| :---          | :---          | :--- |
| `app`         | Top-level     | Defines the namespace or structural root of the configuration |
| `description` | Top/Block     | Provides context or documentation string metadata |
| `version`     | Top/Block     | Declares the version identifier for the application, feature, or endpoint |
| `feature`     | Block-level   | Declares a sub-system module or discrete capability |
| `endpoint`    | Block-level   | Declares a synchronous HTTP transport contract with a route: `endpoint Name "route"` |
| `rules`       | App block     | App-wide application intent, business constraints, and invariants: `rules` followed by quoted strings |
| `method`      | Endpoint      | The HTTP verb, quoted: `method "POST"` |
| `params`      | Block-level   | A dedicated container block for parameter definitions |
| `actions`     | Block-level   | A dedicated container block for execution steps |
| `enforce`     | Action        | A constraint string directly below its parent action |
| `responses`   | Block-level   | A dedicated container block for response schema definitions |
| `returns`     | Endpoint      | Maps status codes to response payloads: `returns` followed by `200 ResponseName` or `500 "error"` |

### Parameter Syntax

Parameters are defined as `<name> <type>` lines under a `params` block. `<name>` is a snake_case identifier and `<type>` is one of `string`, `integer`, `double`, `boolean`, `enum[...]`, or `array[...]`:

```butter
params
  name string
  todo_id integer
  completed boolean
  priority enum["low", "medium", "high"]
  members_id array[integer]
```

### Action Syntax

Actions are bare quoted strings under an `actions` block, with optional `enforce <string>` grandchildren:

```butter
actions
  "Validate input is not empty"
  "Sanitize input"
    enforce "Reject empty strings"
  "Process payment"
```

### Response Syntax

Responses define reusable payload schemas under a `responses` block. Each response has a PascalCase header and nested snake_case `<name> <type>` field lines:

```butter
responses
  OrderSuccess
    order_id string
    total_amount double
    line_items array[string]
```

### Return Syntax

The `returns` block maps HTTP status codes (`100`–`599`) to either a referenced response name or a quoted string literal:

```butter
returns
  201 OrderSuccess
  500 "Internal server error"
```

---

## Example

Save the following as `demo.butter`:

```butter
# Global application declaration
app OrderProcessor
  description "Handles high-throughput retail checkout workflows safely"
  version "2.1.0"

  rules
    "Customers may only access orders they own"
    "An order cannot be charged for more than its calculated total"
    "Operations that would create duplicate orders must be rejected"

feature ProcessPayment
  description "Processes financial transactions through multiple payment gateways"
  version "1.0.0"

  params
    order_id string
    amount double
    payment_method enum["credit_card", "crypto", "bank_transfer"]

  actions
    "Validate routing balance metrics"
    "Apply cryptocurrency transaction surcharge"
    "Flag transaction for manual risk mitigation review"
    "Maintain continuous transaction ledger heartbeat"
```

Compile it:

```bash
butter compile demo.butter
```

A longer working example with multiple features and an endpoint is available in [`todo.butter`](specs/todo.butter) and [`test-endpoint.butter`](specs/test-endpoint.butter). Each feature's actions run as sequential execution steps, one after another.

Output (`demo.prompt.md`):

```markdown
# [SYSTEM SPEC] OrderProcessor
> **Version:** 2.1.0
> **Description:** Handles high-throughput retail checkout workflows safely

### Rules
**CRITICAL:** The following rules MUST be respected throughout the entire implementation:

* Customers may only access orders they own
* An order cannot be charged for more than its calculated total
* Operations that would create duplicate orders must be rejected

## Feature: ProcessPayment
**Version:** 1.0.0
Processes financial transactions through multiple payment gateways

### Params
* `order_id` (string)
* `amount` (double)
* `payment_method` (enum["credit_card", "crypto", "bank_transfer"])

### Execution Sequence
**CRITICAL:** Execute the following steps strictly in order. Do not proceed to the next step until the current one is complete.

1. **Validate routing balance metrics**
2. **Apply cryptocurrency transaction surcharge**
3. **Flag transaction for manual risk mitigation review**
4. **Maintain continuous transaction ledger heartbeat**
```

---

## Installation

### From Source

Requires [Go](https://go.dev/dl/) 1.21+.

```bash
git clone https://github.com/lebohang0824/butter.git butter
cd butter
go build -o butter main.go
```

Then install it with the install script (see below).

### Install Script

**Linux / macOS:**

```bash
chmod +x install.sh
./install.sh          # install compiler + VS Code extension
./install.sh update   # rebuild and reinstall both
./install.sh binary   # compiler only
./install.sh extension # VS Code extension only
```

**Windows (PowerShell):**

```powershell
.\install.ps1                # install compiler + VS Code extension
.\install.ps1 -Command update
.\install.ps1 -Command binary
.\install.ps1 -Command extension
```

---

## Usage

```text
butter compile [input file] [flags]
butter fmt    [input file] [flags]
```

### `butter compile`

| Flag | Shorthand | Description |
| :--- | :--- | :--- |
| `--output` | `-o` | Custom output path (defaults to `<input>.prompt.md`) |
| `--format` | `-f` | Output format (default: `prompt`). Run `butter compile --help` to see all registered formats |
| `--check` | | Validate syntax and semantics without generating output |

```bash
butter compile demo.butter
butter compile demo.butter -f json
butter compile demo.butter -f yaml
butter compile demo.butter -f json -o result.json
butter compile --check demo.butter
butter --version
```

### `butter fmt`

Formats a `.butter` file according to standard conventions — normalizes indentation, removes blank lines after section keywords, and adds blank lines before `params`, `actions`, `responses`, `returns`, `rules`, and between top-level `feature`/`endpoint` blocks.

| Flag | Description |
| :--- | :--- |
| `--check` | Check formatting without modifying |

```bash
butter fmt demo.butter
butter fmt --check demo.butter
```

Only `.butter` files are accepted as input. Use `--check` to validate syntax and semantics without writing an output file — useful for editor integration and CI pipelines.

---

### Output Extensions

Butter's output layer is fully pluggable. The built-in JSON, YAML, and prompt serializers implement a simple three-method `Extension` interface. Anyone can write a new extension — for TOML, XML, Protobuf, Markdown, or anything else — and plug it in with a single import.

To write an extension, implement the `output.Extension` interface and call `output.Register()`:

```go
package toml

import "butter/pkg/output"

func init() { output.Register(tomlExt{}) }

type tomlExt struct{}
func (tomlExt) Name() string          { return "toml" }
func (tomlExt) FileExtension() string { return ".toml" }
func (tomlExt) Serialize(spec *ast.AppSpec) ([]byte, error) {
    // your serialization logic
}
```

Then add a blank import in `cmd/root.go` and rebuild. The extension appears automatically in `--format` help text and error messages.

Built-in extensions are documented in the [VS Code Extension](docs/extension.html) page.

## AI Workflow

Butter's true value emerges when you feed the compiled output to an AI agent. Here's the workflow:

1. **Write a `.butter` spec** — Declare your features, their typed parameters, and the sequential actions that implement each feature.

2. **Compile it** — `butter compile spec.butter` produces `spec.prompt.md`.

3. **Feed the prompt and implementation context to an AI agent** — Include the compiled prompt with a short instruction: *"Implement the application intent in this specification. Each feature's actions are sequential execution steps — run them one after another in the listed order. Respect all types, rules, enforce constraints, and interface contracts."*

4. **Review the implementation** — The structured spec makes intent easier to communicate and check. The agent chooses how to realize it in the supplied implementation context, and the result should be reviewed and tested like any other implementation.

### Example

```text
Using this specification and the implementation context below, build the application. Each feature's
actions are sequential execution steps — implement them strictly one after the other
in the listed order. Respect the application rules, params, enforce expressions,
and interface contracts.

Implementation context:
- Use the selected language, framework, database, and project conventions

[paste compiled spec.prompt.md here]
```

The spec defines *what* to build. The AI figures out *how*. That's the division of labour.

## Compiler Architecture

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
       │ (Abstract Syntax Tree)
       ▼
 ┌───────────┐
 │ Semantic  │ <--- Checks: duplicate names, valid types, valid
 │  Analysis │       HTTP methods/status codes, response refs, enums
 └─────┬─────┘
       │ (Validated AST)
       ▼
  ┌──────────────┐
  │ Output       │ <--- Pluggable Extension Registry
  │ Extension    │       (json, yaml, prompt, + custom)
  │  Registry    │
  └──────┬───────┘
         │
         ▼
  [ .prompt.md / .json / .yaml / custom file ]
```

### Semantic Analysis

After parsing, a dedicated semantic analysis pass validates the AST against the following rules:

| Check | Severity | Description |
| :--- | :--- | :--- |
| Duplicate feature names | Error | Two features with the same name (includes first-definition line) |
| Duplicate parameter names | Error | Two params with the same name within a feature |
| Unknown parameter/field type | Error | Type is not `string`, `integer`, `double`, `boolean`, `enum[...]`, or `array[...]` |
| Duplicate enum value | Error | An `enum[...]` list contains the same value twice |
| Invalid HTTP method | Error | `method` is not one of POST/GET/PUT/DELETE/PATCH |
| Invalid status code | Error | `returns` code outside the range 100–599 |
| Undefined response ref | Error | `returns` references a `responses` name that isn't declared |
| Missing route/method | Error | Endpoint lacks a required `route` or `method` |

Errors block output generation.

### Lexical Analysis (The Off-side Rule)

Because Butter uses whitespace indentation to mark boundaries, the lexer reads files sequentially while maintaining a **LIFO Indentation Stack** tracking current space depth levels:

- When a newline occurs, the lexer scans consecutive leading whitespace characters.
- If the space-count exceeds the value on top of the stack, it pushes the new count and emits an implicit `INDENT` token.
- If the space-count is less than the top of the stack, it pops elements, emitting a `DEDENT` token for each, until a matching level is found. Any mismatch throws a syntax error.

### Abstract Syntax Tree (AST)

The parser constructs a typed AST graph mapped directly to Go structures:

- **AppSpec** — Root node: app name, description, version, rules, and features/endpoints
- **FeatureSpec** — Named feature with optional description, version, params, and actions
- **EndpointSpec** — Named endpoint with route, method, params, responses, actions, and returns
- **ParamSpec** — Parameter with a name and type
- **ActionSpec** — Action statement with optional enforce(s)

---

## VS Code Extension

A VS Code extension providing syntax highlighting, indentation support, and language configuration is included in the `butter-extension/` directory.

**Features:**
- Full TextMate grammar with named capture highlighting for `app`, `feature`, `endpoint`, and parameter identifiers
- **Butter Docs Colors** theme — a VS Code color theme that matches the docs color scheme (amber keywords, green strings, blue functions, purple params). Select "Butter Docs Colors" from the theme picker.
- On-save formatting — automatically applies `butter fmt` every time a file is saved, no configuration needed
- On-save linting — validates syntax via `butter compile --check` after formatting and surfaces errors with red squiggly underlines
- `Butter: Lint current file` command in the command palette
- `Butter: Format current file` command in the command palette
- Auto-indentation for `feature`, `endpoint`, `params`, `actions`, `responses`, `returns`, and `rules` blocks
- Configurable compiler path (`butter.compilerPath`)
- Comment toggle with `#` and `//`
- Auto-closing pairs for `"` and `[]`
- Document file icon for `.butter` files

Install via the install script (`./install.sh extension`) or manually with:

```bash
code --install-extension butter-extension.vsix
```

Or open the `butter-extension/` directory in VS Code and press F5.

---

## License

Butter is open source software. See the project repository for license information.
