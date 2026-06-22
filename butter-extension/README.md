# Butter Configuration Language Support

![Butter Logo](butter.png)

A VS Code extension providing syntax highlighting, indentation support, and language configuration for the **Butter** DSL — an indentation-aware specification language that compiles to JSON.

## Features

- **Syntax Highlighting** — Full TextMate grammar for keywords, strings, conditionals, comments, and boolean constants
- **Auto-Indentation** — Smart indent/dedent for `feature`, `params`, `actions`, and `param` blocks
- **Comment Support** — `#` line comments with toggle support
- **Auto-Closing Pairs** — Automatic `"` and `[]` pair completion

## Usage

Install the extension and open any `.butter` file. The language mode is automatically detected.

### Example

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

  actions
    action "Validate routing balance metrics"
    action "Flag for review" | if "Amount > 10000"
    action "Bypass fraud check" | unless "Amount > 50"
```

## Compiler

The Butter compiler is a standalone Go CLI tool. See the [Butter repository](https://github.com/butter-io/butter) for instructions on building and using the compiler.

## Release Notes

See [CHANGELOG.md](CHANGELOG.md) for version history.
