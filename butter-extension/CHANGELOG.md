# Change Log

## [1.6.0] - 2026-06-29

- Semantic analysis pass: catches duplicate feature/param names, undefined param references in conditions, type-default mismatches, enum defaults, and redundant required+default combos
- `butter compile --check` now validates semantics in addition to syntax
- Source location tracking in AST for precise error reporting

## [1.5.0] - 2026-06-26

- YAML output support via `--format yaml` or `-f yaml` flag
- `butter compile --check` now reports semantic errors separately from syntax errors
- Documentation site at `docs/index.html` with full language guide and CLI reference
- `compiler.butter` self-documenting spec describes all four compilation stages

## [1.4.0] - 2026-06-25

- Added `bool|boolean|length` to keyword syntax highlighting
- `length` param field for exact length constraints (mutually exclusive with `validate`)

## [1.3.0] - 2026-06-25

- `product` keyword accepted as alias for `app` at the top level
- Numeric literals (e.g. `default 50`) now parse correctly
- Parser guards against infinite loops on unexpected token types in parameter fields
- Added `todo.butter` example using `product` with four features, enum types, and integer defaults
- `validate` rule format now checked at parse time — must be a valid numeric comparison (operator + number)
- `validate` only allowed on `int` or `float` parameter types

## [1.2.0] - 2026-06-24

- `butter fmt` command — automatically formats `.butter` files according to standard conventions
- Format on save — the extension automatically applies `butter fmt` every time a file is saved
- `Butter: Format current file` command in the command palette

## [1.1.0] - 2026-06-23

- Named capture highlighting for `app`, `feature`, and `param` identifiers
- `version` keyword support in syntax highlighting
- Document file icon (`.svg`) for `.butter` files
- On-save linting — validates syntax via `butter compile --check` and surfaces errors with red squiggly underlines
- `Butter: Lint current file` command in the command palette
- Configurable compiler path (`butter.compilerPath`)

## [1.0.0] - 2026-06-22

- Initial release
- Syntax highlighting for `.butter` files
- Indentation rules for `feature`, `params`, `actions`, and `param` blocks
- Comment toggling with `#`
- Auto-closing pairs for `"` and `[]`
