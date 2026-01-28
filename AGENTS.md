## Description
Textwire is a templating language for Go. It is designed to easily inject variables from Go code into a template file or just a regular string. It uses directives like `@if()`, `@for()`, `@each()`, and expressions like `{{ x + y }}` and `{{ "Hello, World!".lower() }}`.

### Globals
You can pass globals to `Configure` or `NewTemplate` function like this:
```go
tpl, err = textwire.NewTemplate(&config.Config{
    GlobalData: map[string]any{
        "env": "development",
    },
})
```
In template files you can use globals like this `{{ globals.env }}`.

## Project Overview
- **Module**: `github.com/textwire/textwire/v2`
- **Template Extensions**: `.tw`
- **Architecture**: Modular design with separate packages for lexer, parser, evaluator, AST, token, and object handling
- **LSP Support**: Full Language Server Protocol implementation for editor integration
- **Precedence**: All the precedence defined in `./parser/parser.go:20-70`

## Development Environment & Commands
### Available Commands
- `make test` - Run all tests
- `make fmt` - Format Go code
- `make lint` - Run linter
- `make shell` - Start the REPL
- `make check` - Combines commands like: test, fmt, lint

### Single Test Execution
Use these patterns for running specific tests:
```bash
# Run specific test function
go test -run TestSpecificFunction ./parser/

# Run tests for specific package with verbose output
go test -v ./parser/

# Run all tests in parser package
go test ./parser/...

# Run tests with coverage
go test -cover ./parser/
```

## Project Architecture
### Core Packages
- `lexer/` - Tokenizes input strings into tokens
- `parser/` - Builds AST from token stream
- `evaluator/` - Evaluates AST nodes to produce output
- `token/` - Token definitions and utilities
- `ast/` - Abstract Syntax Tree node definitions
- `object/` - Runtime object system

### Support Packages
- `config/` - Configuration management
- `fail/` - Structured error handling
- `utils/` - Utility functions
- `lsp/` - Language Server Protocol implementation

### Entry Points
- `textwire.go` - Main public API
- `textwire/example/main.go` - Example web server usage
- `repl/repl.go` - Interactive REPL

## Code Style & Conventions
### Naming Conventions
- **Exported**: PascalCase (e.g., `NewIdentifier`, `ParseProgram`)
- **Unexported**: camelCase (e.g., `parseExpression`, `checkErrors`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `ErrEmptyBraces`, `LOWEST`)
- **Interfaces**: Simple, descriptive names (e.g., `Node`, `Expression`)
- **Token Names**: Simple, unprefixed (e.g., `AND`, `OR`, `NOT` - not `LOGICAL_AND`)

### Constructor Patterns
- Use `NewX` pattern for all constructors
- Example: `NewIdentifier(tok token.Token, val string) *Identifier`
- Always accept token as first parameter for AST nodes
- Return pointer types for AST nodes

### AST Node Patterns
- Embed `BaseNode` in all AST nodes
- Implement `Node` interface methods: `Tok()`, `String()`, `Line()`, `Position()`, `SetEndPosition()`
- Use `expressionNode()` and `statementNode()` marker methods
- Example structure:
```go
type Identifier struct {
    BaseNode
    Value string
}

func NewIdentifier(tok token.Token, val string) *Identifier {
    return &Identifier{
        BaseNode: NewBaseNode(tok),
        Value:    val,
    }
}
```

### Import Organization
- Group imports: standard library, third-party, local packages
- Local packages use `github.com/textwire/textwire/v2/` prefix
- Keep imports sorted and remove unused imports

### Formatting Requirements
- ALWAYS run `make check` before committing changes
- For non-Go files, respect `.prettierrc` configuration (tabWidth: 4, singleQuote: true)

## Error Handling Patterns
### Fail Package Usage
Import the fail package for structured error handling:
```go
import "github.com/textwire/textwire/v2/fail"
```

### Error Constants
Use predefined error constants from the fail package:
- Parser errors: `ErrEmptyBraces`, `ErrWrongNextToken`, `ErrExpectedExpression`
- Evaluator errors: `ErrUnknownNodeType`, `ErrIdentifierNotFound`, `ErrTypeMismatch`
- Function errors: `ErrNoFuncForThisType`, `ErrFuncRequiresOneArg`

### Error Creation
Create errors using the fail package:
```go
// Basic error
fail.New(line, filepath, component, message, args...)

// From existing Go error
fail.FromError(err, line, filepath, component).Error()

// Example in parser
p.newError(p.curToken.ErrorLine(), fail.ErrWrongNextToken, 
    token.String(token.LPAREN), token.String(p.peekToken.Type))
```

### Error Testing
Always test error conditions:
- Test that appropriate errors are generated for invalid input
- Test error messages match expected constants
- Use `checkParserErrors(t, p)` helper in parser tests
- Test both success and failure paths

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
- Example pattern:
```go
func parseStatements(t *testing.T, inp string, opts parseOpts) []ast.Statement {
    l := lexer.New(inp)
    p := New(l, "")
    prog := p.Parse()
    
    if opts.checkErrors {
        checkParserErrors(t, p)
    }
    
    return prog.Statements
}
```

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
- Follow semantic versioning (currently v2.x)
- Update CHANGELOG.md for significant changes
- Use appropriate emoji categories in changelog
- Maintain compatibility within major version

## Specific Implementation Notes
### Parser Precedence
Follow standard operator precedence when adding new operators:
- Logical operators: `&&` higher than `||`
- Comparison operators: `==`, `!=`, `<`, `>`, `<=`, `>=`
- Arithmetic operators: `*`, `/`, `%` higher than `+`, `-`
- Use appropriate precedence constants in parser.go

### Token Naming
Keep token names simple and unprefixed:
- `AND` (not `LOGICAL_AND`)
- `OR` (not `LOGICAL_OR`) 
- `NOT` (not `LOGICAL_NOT`)
- Follow existing patterns in token/token.go

### Error Message Style
Write clear, descriptive error messages:
- Include context about what went wrong
- Provide suggestions for fixing the issue
- Use consistent terminology
- Follow existing message patterns in fail/fail.go
