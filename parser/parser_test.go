package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/lexer"
)

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

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case nil:
		return testNilLiteral(t, exp)
	}

	t.Errorf("type of exp not handled. got=%T", exp)

	return false
}

func TestParseIdentifier(t *testing.T) {
	stmts := parseStatements(t, "{{ myName }}", 1)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	if !testIdentifier(t, stmt.Expression, "myName") {
		return
	}
}

func TestParseIntegerLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ 234 }}", 1)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	if !testIntegerLiteral(t, stmt.Expression, 234) {
		return
	}
}

func TestParseNilLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ nil }}", 1)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	testNilLiteral(t, stmt.Expression)
}

func TestParseStringLiteral(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{`{{ "Hello World" }}`, "Hello World"},
		{`{{ "Serhii \"Cho\"" }}`, `Serhii "Cho"`},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.input, 1)

		stmt, ok := stmts[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
		}

		if !testStringLiteral(t, stmt.Expression, tt.expect) {
			return
		}
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `{{ "Serhii" + " Anna" }}`

	stmts := parseStatements(t, input, 1)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExpression)

	if !ok {
		t.Fatalf("stmt is not an InfixExpression, got %T", stmt.Expression)
	}

	if exp.Left.String() != "Serhii" {
		t.Fatalf("exp.Left is not %s, got %s", "Serhii", exp.Left.String())
	}

	if exp.Operator != "+" {
		t.Fatalf("exp.Operator is not %s, got %s", "+", exp.Operator)
	}

	if exp.Right.String() != " Anna" {
		t.Fatalf("exp.Right is not %s, got %s", " Anna", exp.Right.String())
	}
}

func TestGroupedExpression(t *testing.T) {
	test := "{{ (5 + 5) * 2 }}"

	stmts := parseStatements(t, test, 1)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExpression)

	if !ok {
		t.Fatalf("stmt is not an InfixExpression, got %T", stmt.Expression)
	}

	if !testLiteralExpression(t, exp.Right, 2) {
		return
	}

	if exp.Operator != "*" {
		t.Fatalf("exp.Operator is not %s, got %s", "*", exp.Operator)
	}

	infix, ok := exp.Left.(*ast.InfixExpression)

	if !ok {
		t.Fatalf("exp.Left is not an InfixExpression, got %T", exp.Left)
	}

	if !testLiteralExpression(t, infix.Left, 5) {
		return
	}

	if infix.Operator != "+" {
		t.Fatalf("infix.Operator is not %s, got %s", "+", infix.Operator)
	}

	if !testLiteralExpression(t, infix.Right, 5) {
		return
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"{{ 5 + 8 }}", 5, "+", 8},
		{"{{ 10 - 2 }}", 10, "-", 2},
		{"{{ 2 * 2 }}", 2, "*", 2},
		{"{{ 44 / 4 }}", 44, "/", 4},
		{"{{ 5 % 4 }}", 5, "%", 4},
		{"{{ myAge + herAge }}", "myAge", "+", "herAge"},
		{"{{ myAge - herAge }}", "myAge", "-", "herAge"},
		{"{{ myAge * herAge }}", "myAge", "*", "herAge"},
		{"{{ myAge / herAge }}", "myAge", "/", "herAge"},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.input, 1)

		stmt, ok := stmts[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)

		if !ok {
			t.Fatalf("stmt is not an InfixExpression, got %T", stmt.Expression)
		}

		if !testLiteralExpression(t, exp.Left, tt.left) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s, got %s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.right) {
			return
		}
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"{{ -5 }}", "-", 5},
		{"{{ -10 }}", "-", 10},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.input, 1)

		stmt, ok := stmts[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt is not a PrefixExpression, got %T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s, got %s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
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

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
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

func testNilLiteral(t *testing.T, exp ast.Expression) bool {
	nilLit, ok := exp.(*ast.NilLiteral)

	if !ok {
		t.Errorf("exp is not a NilLiteral, got %T", exp)
		return false
	}

	if nilLit.TokenLiteral() != "nil" {
		t.Errorf("nilLit.TokenLiteral() is not 'nil', got %s", nilLit.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	str, ok := exp.(*ast.StringLiteral)

	if !ok {
		t.Errorf("exp is not a StringLiteral, got %T", exp)
		return false
	}

	if str.Value != value {
		t.Errorf("str.Value is not %s, got %s", value, str.Value)
		return false
	}

	if str.TokenLiteral() != value {
		t.Errorf("str.TokenLiteral() is not %s, got %s", value, str.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)

	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
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

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"{{ 1 * 2 }}",
			"{{ (1 * 2) }}",
		},
		{
			"<h2>{{ -2 + 3 }}</h2>",
			"<h2>{{ ((-2) + 3) }}</h2>",
		},
		{
			"{{ a + b + c }}",
			"{{ ((a + b) + c) }}",
		},
		{
			"{{ a + b / c }}",
			"{{ (a + (b / c)) }}",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
