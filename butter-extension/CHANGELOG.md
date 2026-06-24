# Change Log

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
