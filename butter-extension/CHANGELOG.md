# Change Log

## [1.13.0] - 2026-07-13

- **`endpoint` enhancement** — New top-level block for defining synchronous HTTP transport network architecture alongside features. Supports `version`, `params`, `responses`, `actions`, and `returns` sub-blocks with strict scope separation from features.
- **`response` blocks** — Define reusable response schemas with typed fields inside endpoints.
- **`return` mapping** — Map status codes to response references or string payloads with optional `if`/`unless` conditions.
- **`field` keyword** — Typed fields inside responses with optional nested sub-fields; defaults to `string` when no `type` is specified.
- **Return statement coloring** — Status codes, response references, and string payloads in `return` statements are now syntax-highlighted with dedicated TextMate patterns.
- **Type value coloring** — Values after `type` keywords now receive `support.type.butter` scope (brown).
- **Field name coloring** — Field names now use `variable.parameter.butter` scope (pink), matching param names.
- **Response ref hover** — Hovering over a response reference in a `return` statement shows the response schema declaration.
- **Response ref navigation** — Ctrl+click on a response reference in a `return` statement navigates to its `response` declaration.
- **HTML tree visualization** — Endpoints render as node cards with route/method badges, response schema cards, and return mapping cards with connection lines.
- **AI simulator** — Endpoints appear in the sidebar with route/method display, parameter forms, response schemas, and return mapping. Code generation prompt includes endpoint data.
- **Formatter** — Handles endpoint, responses, returns, response, and field keywords with proper indentation rules.

## [1.9.0] - 2026-07-02

- **IntelliSense code completion** — Context-aware autocomplete for Butter keywords, types, parameter fields, and conditionals. The completion provider analyses indentation and parent block context to suggest only the valid keywords for the current cursor position. Includes snippets for `app`, `product`, `feature`, `param`, `action`, and `enum[...]`.
- **Hover documentation** — Hover over any Butter keyword, type, or conditional to see a markdown description with syntax examples.
- **Trigger characters** — Completion triggers on `|` (pipe for conditionals), `[` (enum), and space (after keywords like `type`, `required`).
- **Updated documentation philosophy** — All README, DOCUMENTATION, docs HTML, and .butter spec files now reflect Butter's core purpose: a specification language that communicates intent to AI agents for up to 100% one-shot results.

## [1.8.0] - 2026-07-01

- `enforce` keyword added to syntax highlighting — recognized as a keyword in both TextMate grammar and Prism docs highlighter
- Updated grammar with `contributes.tokenColorCustomizations` matching docs color scheme: amber keywords, green strings, blue functions/names, purple params, slate-blue booleans, gray italic comments
- Extended `ActionSpec` struct with optional `Enforce []string` field — zero or more enforce strings can appear as indented children under an action
- Updated `compiler.butter` self-documenting spec with enforce examples and aligned to v1.8.0

## [1.7.0] - 2026-06-30

- Pluggable output extension system: JSON and YAML refactored into standalone `Extension` implementations
- New `pkg/output` package with registry pattern — third-party formats can be added via a single-file Go implementation
- Dynamic `--format` help text lists all registered extensions
- Cleaner YAML output: internal line numbers no longer leak into generated files

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
