# Butter Configuration Language Support

![Butter Logo](https://raw.githubusercontent.com/lebohang0824/butter/main/butter-extension/butter.png)

A VS Code extension providing syntax highlighting, IntelliSense, formatting, linting, and language configuration for **Butter** — a specification language that communicates intent to AI agents. Write `.butter` specs, compile to a Markdown prompt, and feed it to AI agents who produce higher first-pass accuracy.

## Features

- **IntelliSense** — Context-aware code completion and hover documentation for all Butter keywords and types. Suggests the right keyword based on indentation and parent block context (e.g., `method` in an endpoint).
- **Syntax Highlighting** — Full TextMate grammar with named capture highlighting for `app`, `feature`, `endpoint`, and `rules`/`returns` identifiers, including bare param types (`string`, `integer`, `double`, `boolean`, `enum[...]`)
- **On-Save Formatting** — Automatically applies `butter fmt` on every save, no configuration needed
- **On-Save Linting** — Validates `.butter` syntax on save using the bundled compiler and surfaces errors with red squiggly underlines
- **Manual Lint Command** — `Butter: Lint current file` in the command palette
- **Manual Format Command** — `Butter: Format current file` in the command palette
- **Auto-Indentation** — Smart indent/dedent for `feature`, `endpoint`, `params`, `actions`, `responses`, `returns`, and `rules` blocks
- **Configurable Compiler Path** — Set the path to the `butter` binary via `butter.compilerPath`
- **Comment Support** — `#` and `//` line comments with toggle support
- **Auto-Closing Pairs** — Automatic `"` and `[]` pair completion
- **Document File Icon** — Custom icon for `.butter` files

## Usage

Install the extension and open any `.butter` file. The language mode is automatically detected.

### Example

```butter
# Global application declaration
app OrderProcessor
  description "Handles high-throughput retail checkout workflows safely"
  version "2.1.0"

feature ProcessPayment
  description "Processes financial transactions through multiple payment gateways"
  version "1.0.0"

  params
    order_id string
    amount double
    payment_method enum["credit_card", "crypto", "bank_transfer"]

  actions
    "Validate routing balance metrics"
    "Flag for review"
    "Apply cryptocurrency transaction surcharge"

endpoint SaveOrder "/orders"
  description "Stores an order in the database"
  version "1.0.0"
  method POST

  params
    order_id string
    amount double

  responses
    OrderSuccess
      id integer
      amount double

  actions
    "Validate order payload"
    "Store order in the database"

  returns
    201 OrderSuccess
    500 "Internal server error"
```

## Compiler

The Butter compiler is a standalone Go CLI tool. See the [Butter repository](https://github.com/lebohang0824/butter) for instructions on building and using the compiler.

## Release Notes

See [CHANGELOG.md](CHANGELOG.md) for version history.
