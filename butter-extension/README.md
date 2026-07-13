# Butter Configuration Language Support

![Butter Logo](https://raw.githubusercontent.com/lebohang0824/butter/main/butter-extension/butter.png)

A VS Code extension providing syntax highlighting, IntelliSense, formatting, linting, and language configuration for **Butter** — a specification language that communicates intent to AI agents. Write `.butter` specs, compile to JSON, and feed to AI agents who produce implementations matching up to 100% of expected results in a single shot.

## Features

- **IntelliSense** — Context-aware code completion and hover documentation for all Butter keywords, types, parameter fields, and conditionals. Suggests the right keyword based on indentation and parent block context (e.g., `type`/`required` inside a `param`, `action` inside `actions`, conditionals after `|`).
- **Syntax Highlighting** — Full TextMate grammar with named capture highlighting for `app`, `feature`, `endpoint`, `listener`, and `param` identifiers, including listener topics and return states
- **On-Save Formatting** — Automatically applies `butter fmt` on every save, no configuration needed
- **On-Save Linting** — Validates `.butter` syntax on save using the bundled compiler and surfaces errors with red squiggly underlines
- **Manual Lint Command** — `Butter: Lint current file` in the command palette
- **Manual Format Command** — `Butter: Format current file` in the command palette
- **Auto-Indentation** — Smart indent/dedent for `feature`, `endpoint`, `listener`, `params`, `actions`, and `param` blocks
- **Configurable Compiler Path** — Set the path to the `butter` binary via `butter.compilerPath`
- **Comment Support** — `#` line comments with toggle support
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
    param OrderID
      type string
      required true
    param Amount
      type float
      required true
    param PaymentMethod
      type enum["CreditCard", "Crypto", "BankTransfer"]
      default "CreditCard"

  actions
    action "Validate routing balance metrics"
    action "Flag for review" | if "Amount > 10000"
    action "Bypass fraud check" | unless "Amount > 50"
```

## Compiler

The Butter compiler is a standalone Go CLI tool. See the [Butter repository](https://github.com/butter-io/butter) for instructions on building and using the compiler.

## Release Notes

See [CHANGELOG.md](CHANGELOG.md) for version history.
