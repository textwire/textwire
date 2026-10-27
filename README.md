<p align="center"><img src="/textwire/textwire/raw/branch/master/.forgejo/assets/textwire-logo.png" width="200" height="200" /></p>

# Textwire. A templating language for Go

Textwire is a simple yet powerful templating language for Go. It is designed to easily inject variables from Go code into a template file or just a regular string.

### [Read Official Documentation](https://textwire.serhiicho.com)

## Installation

```bash
go get codeberg.org/textwire/textwire/v5
```

## Neovim and VSCode Support

If you use [Neovim](https://neovim.io/) or [VSCode](https://code.visualstudio.com/) as your primary editor, you can install the [Neovim Plugin](https://codeberg.org/textwire/textwire.nvim) or [VSCode Extension](https://marketplace.visualstudio.com/items?itemName=SerhiiCho.textwire) to get syntax highlighting and other features for Textwire.

## Globals Usage

You can pass globals to `Configure` or `NewTemplate` function like this:

```go
tpl, err = textwire.NewTemplate(&config.Config{
    GlobalData: map[string]any{
        "env": "development",
    },
})
```

In template files you can use globals like this `{{ global.env }}`.

## Project Overview

- **Module**: `codeberg.org/textwire/textwire/v5`
- **Template Extensions**: `.tw`
- **Architecture**: Modular design with separate packages for lexer, parser, evaluator, AST, token, and object handling
- **LSP Support**: Full Language Server Protocol implementation for editor integration
- **Precedence**: All the precedence defined in `./pkg/parser/parser.go:20-70`

## Development Environment & Commands

### Available Commands

All commands run inside a container. Pull the image once before using them:

```bash
make pull
```

- `make pull` - Pull the checks container image
- `make test` - Run all tests
- `make fmt` - Format Go code
- `make lint` - Run linter
- `make repl` - Start the REPL
- `make check` - Combines commands like: test, fmt, lint
- `make` - Defaults to `make check`

Commands use `podman-compose` by default. To use Docker Compose instead, override `OCI`:

```bash
make check OCI='docker compose'
```

## Project Architecture

### Core Packages
- `pkg/lexer/` - Tokenizes input strings into tokens
- `pkg/parser/` - Builds AST from token stream
- `pkg/evaluator/` - Evaluates AST nodes to produce output
- `pkg/token/` - Token definitions and utilities
- `pkg/ast/` - Abstract Syntax Tree node definitions
- `pkg/value/` - Runtime object system

### Support Packages

- `config/` - Configuration management
- `pkg/fail/` - Structured error handling
- `pkg/utils/` - Utility functions
- `pkg/lsp/` - Language Server Protocol implementation

### Entry Points

- `textwire.go` - Main public API
- `pkg/repl/repl.go` - Interactive REPL

## Code Style & Conventions

### Import Organization

- Group imports: standard library, third-party, local packages
- Local packages use `codeberg.org/textwire/textwire/v5/` prefix
- Keep imports sorted and remove unused imports

### Formatting Requirements

- ALWAYS run `make check` before committing changes
- For non-Go files, respect `.prettierrc` configuration (tabWidth: 4, singleQuote: true)

## Testing Guidelines

### Test Coverage Requirements

- Add tests for ALL new functionality
- Maintain high test coverage across all packages
- Test both success and failure scenarios
- Include edge cases and error conditions

### Test Organization

- Co-locate test files with source files (`parser_test.go` with `parser.go`)
- Use `textwire/testdata/` for template test data
- Organize test data into `good/` and `bad/` subdirectories
- Use descriptive test names that explain what is being tested

### Test Patterns

- Use helper functions for common test setup
- Use option structs for configurable test behavior

### Test Data

- Store template examples in `textwire/testdata/`
- Use `.tw` extension
- Include both valid and invalid examples for testing
- Name test files descriptively (e.g., `if-statements.tw`)

## LSP Integration Guidelines

### Metadata Updates

When adding new features, update LSP metadata in `lsp/metadata/en/`:
- Add completion items for new directives
- Update function signatures

### Completion Items

- Add completions for new directives (`@if`, `@for`, etc.)
- Include function parameter hints
- Provide snippet completions for complex structures
- Update context-aware completions

## Development Workflow

### Pre-Commit Checklist

- [ ] Run `make check` before commit
- [ ] Update LSP metadata if adding new features
- [ ] Add tests for new functionality

### Versioning

- Follow semantic versioning (currently v2.x.x)
- Update CHANGELOG.md for significant changes
- Use appropriate emoji categories in changelog
- Maintain compatibility within major version

## Specific Implementation Notes

### Error Message Style

Write clear, descriptive error messages:
- Include context about what went wrong
- Provide suggestions for fixing the issue
- Use consistent terminology
- Follow existing message patterns in fail/fail.go

## License

The Textwire project is licensed under the [MIT License](https://codeberg.org/textwire/textwire/src/branch/master/LICENSE)

## Thanks To

Thanks to Codeberg and Forgejo for hosting Textwire.

<p align="center">
<img src="/textwire/textwire/raw/branch/master/.forgejo/assets/codeberg-logo.png" width="50" height="50" />
<img src="/textwire/textwire/raw/branch/master/.forgejo/assets/textwire-logo.png" width="50" height="50" />
<img src="/textwire/textwire/raw/branch/master/.forgejo/assets/forgejo-logo.png" width="50" height="50" />
</p>
