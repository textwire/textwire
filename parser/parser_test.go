package parser

import (
	"strconv"
	"testing"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/lexer"
)

func TestParseIdentifier(t *testing.T) {
	stmts := parseStatements(t, "{{ myName }}", 1)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	if !checkIdentifier(t, stmt.Expression, "myName") {
		return
	}
}

func TestParseIntegerLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ 234 }}", 1)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	if !checkIntegerLiteral(t, stmt.Expression, 234) {
		return
	}
}

func parseStatements(t *testing.T, inp string, stmtCount int) []ast.Statement {
	l := lexer.New(inp)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != stmtCount {
		t.Fatalf("program must have %d statement, got %d", stmtCount, len(program.Statements))
	}

	return program.Statements
}

func checkIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integer, ok := exp.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("exp is not an IntegerLiteral, got %T", exp)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value is not %d, got %d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != strconv.FormatInt(value, 10) {
		t.Errorf("integer.TokenLiteral() is not %d, got %s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func checkIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("exp is not an Identifier, got %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value is not %s, got %s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() is not %s, got %s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
