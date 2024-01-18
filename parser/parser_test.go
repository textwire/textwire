package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/token"
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

func parseStatements(t *testing.T, inp string, stmtCount int, inserts map[string]*ast.InsertStatement) []ast.Statement {
	l := lexer.New(inp)
	p := New(l, inserts)

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

func testFloatLiteral(t *testing.T, exp ast.Expression, value float64) bool {
	float, ok := exp.(*ast.FloatLiteral)

	if !ok {
		t.Errorf("exp is not a FloatLiteral, got %T", exp)
		return false
	}

	if float.Value != value {
		t.Errorf("float.Value is not %f, got %f", value, float.Value)
		return false
	}

	if float.String() != fmt.Sprintf("%g", value) {
		t.Errorf("float.String() is not %f, got %s", value, float.String())
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
	case float64:
		return testFloatLiteral(t, exp, v)
	case string:
		return testStringLiteral(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case nil:
		return testNilLiteral(t, exp)
	}

	t.Errorf("type of exp not handled. got=%T", exp)

	return false
}

func TestParseIdentifier(t *testing.T) {
	stmts := parseStatements(t, "{{ myName }}", 1, nil)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	if !testIdentifier(t, stmt.Expression, "myName") {
		return
	}
}

func TestParseIntegerLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ 234 }}", 1, nil)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	if !testIntegerLiteral(t, stmt.Expression, 234) {
		return
	}
}

func TestParseFloatLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ 2.34149 }}", 1, nil)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	if !testFloatLiteral(t, stmt.Expression, 2.34149) {
		return
	}
}

func TestParseNilLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ nil }}", 1, nil)

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
		stmts := parseStatements(t, tt.input, 1, nil)

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
	inp := `{{ "Serhii" + " Anna" }}`

	stmts := parseStatements(t, inp, 1, nil)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExpression)

	if !ok {
		t.Fatalf("stmt is not an InfixExpression, got %T", stmt.Expression)
	}

	if exp.Left.TokenLiteral() != "Serhii" {
		t.Fatalf("exp.Left is not %s, got %s", "Serhii", exp.Left.String())
	}

	if exp.Operator != "+" {
		t.Fatalf("exp.Operator is not %s, got %s", "+", exp.Operator)
	}

	if exp.Right.TokenLiteral() != " Anna" {
		t.Fatalf("exp.Right is not %s, got %s", " Anna", exp.Right.String())
	}
}

func TestGroupedExpression(t *testing.T) {
	test := "{{ (5 + 5) * 2 }}"

	stmts := parseStatements(t, test, 1, nil)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExpression)

	if !ok {
		t.Fatalf("stmt is not an InfixExpression, got %T", stmt.Expression)
	}

	if !testIntegerLiteral(t, exp.Right, 2) {
		return
	}

	if exp.Operator != "*" {
		t.Fatalf("exp.Operator is not %s, got %s", "*", exp.Operator)
	}

	infix, ok := exp.Left.(*ast.InfixExpression)

	if !ok {
		t.Fatalf("exp.Left is not an InfixExpression, got %T", exp.Left)
	}

	if !testIntegerLiteral(t, infix.Left, 5) {
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
		{`{{ "me" + "her" }}`, "me", "+", "her"},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.input, 1, nil)

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

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectBoolean bool
	}{
		{"{{ true }}", true},
		{"{{ false }}", false},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.input, 1, nil)

		stmt, ok := stmts[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
		}

		if !testBooleanLiteral(t, stmt.Expression, tt.expectBoolean) {
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
		{"{{ !true }}", "!", true},
		{"{{ !false }}", "!", false},
		{`{{ !"" }}`, "!", ""},
		{`{{ !0 }}`, "!", 0},
		{`{{ -0 }}`, "-", 0},
		{`{{ -0.0 }}`, "-", 0.0},
		{`{{ !0.0 }}`, "!", 0.0},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.input, 1, nil)

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
		p := New(l, nil)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input      string
		errMessage string
	}{
		{"{{ 5 + }}", ERR_EXPECTED_EXPRESSION},
		{"{{ }}", ERR_EMPTY_BRACKETS},
		{"{{ true ? 100 }}", fmt.Sprintf(ERR_WRONG_NEXT_TOKEN, token.TokenString(token.COLON), token.TokenString(token.RBRACES))},
		{"{{ ) }}", fmt.Sprintf(ERR_NO_PREFIX_PARSE_FUNC, token.TokenString(token.RPAREN))},
		{"{{ 5 }", fmt.Sprintf(ERR_ILLEGAL_TOKEN, "}")},
		{`{{ reserve "title" }}`, fmt.Sprintf(ERR_INSERT_NOT_DEFINED, "title")},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		p.ParseProgram()

		if len(p.Errors()) == 0 {
			t.Errorf("no errors found in input %q", tt.input)
			return
		}

		err := p.Errors()[0]

		if err.Error() != tt.errMessage {
			t.Errorf("expected error message %q, got %q", tt.errMessage, err.Error())
		}
	}
}

func TestTernaryExpression(t *testing.T) {
	inp := `{{ true ? 100 : "Some string" }}`

	stmts := parseStatements(t, inp, 1, nil)

	stmt, ok := stmts[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.TernaryExpression)

	if !ok {
		t.Fatalf("stmt is not a TernaryExpression, got %T", stmt.Expression)
	}

	if !testBooleanLiteral(t, exp.Condition, true) {
		return
	}

	if !testIntegerLiteral(t, exp.Consequence, 100) {
		return
	}

	if !testStringLiteral(t, exp.Alternative, "Some string") {
		return
	}
}

func TestIfStatement(t *testing.T) {
	inp := `{{ if true }}1{{ end }}`

	stmts := parseStatements(t, inp, 1, nil)

	ifStmt, ok := stmts[0].(*ast.IfStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an IfStatement, got %T", stmts[0])
	}

	if !testBooleanLiteral(t, ifStmt.Condition, true) {
		return
	}

	if len(ifStmt.Consequence.Statements) != 1 {
		t.Errorf("ifStmt.Consequence.Statements does not contain 1 statement, got %d", len(ifStmt.Consequence.Statements))
	}

	consequence, ok := ifStmt.Consequence.Statements[0].(*ast.HTMLStatement)

	if !ok {
		t.Fatalf("ifStmt.Consequence.Statements[0] is not an HTMLStatement, got %T", ifStmt.Consequence.Statements[0])
	}

	if consequence.String() != "1" {
		t.Errorf("consequence.String() is not %s, got %s", "1", consequence.String())
	}

	if ifStmt.Alternative != nil {
		t.Errorf("ifStmt.Alternative is not nil, got %T", ifStmt.Alternative)
	}

	if len(ifStmt.Alternatives) != 0 {
		t.Errorf("ifStmt.Alternatives is not empty, got %d", len(ifStmt.Alternatives))
	}
}

func TestIfElseStatement(t *testing.T) {
	inp := `{{ if true }}1{{ else }}2{{ end }}`

	stmts := parseStatements(t, inp, 1, nil)

	ifStmt, ok := stmts[0].(*ast.IfStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an IfStatement, got %T", stmts[0])
	}

	if ifStmt.Alternative == nil {
		t.Errorf("ifStmt.Alternative is nil")
	}

	if len(ifStmt.Alternative.Statements) != 1 {
		t.Errorf("ifStmt.Alternative.Statements does not contain 1 statement, got %d", len(ifStmt.Alternative.Statements))
	}

	alternative, ok := ifStmt.Alternative.Statements[0].(*ast.HTMLStatement)

	if !ok {
		t.Fatalf("ifStmt.Alternative.Statements[0] is not an HTMLStatement, got %T", ifStmt.Alternative.Statements[0])
	}

	if alternative.String() != "2" {
		t.Errorf("alternative.String() is not %s, got %s", "2", alternative.String())
	}

	if len(ifStmt.Alternatives) != 0 {
		t.Errorf("ifStmt.Alternatives is not empty, got %d", len(ifStmt.Alternatives))
	}
}

func TestIfElseIfStatement(t *testing.T) {
	inp := `{{ if true }}1{{ else if false }}2{{ end }}`

	stmts := parseStatements(t, inp, 1, nil)

	ifStmt, ok := stmts[0].(*ast.IfStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an IfStatement, got %T", stmts[0])
	}

	if ifStmt.Alternative != nil {
		t.Errorf("ifStmt.Alternative is not nil, got %T", ifStmt.Alternative)
	}

	if len(ifStmt.Alternatives) != 1 {
		t.Errorf("ifStmt.Alternatives does not contain 1 statement, got %d", len(ifStmt.Alternatives))
	}

	alternative := ifStmt.Alternatives[0]

	if !testBooleanLiteral(t, alternative.Condition, false) {
		return
	}

	if len(alternative.Consequence.Statements) != 1 {
		t.Errorf("alternative.Consequence.Statements does not contain 1 statement, got %d", len(alternative.Consequence.Statements))
	}

	consequence, ok := alternative.Consequence.Statements[0].(*ast.HTMLStatement)

	if !ok {
		t.Fatalf("alternative.Consequence.Statements[0] is not an HTMLStatement, got %T", alternative.Consequence.Statements[0])
	}

	if consequence.String() != "2" {
		t.Errorf("consequence.String() is not %s, got %s", "2", consequence.String())
	}
}

func TestIfElseIfElseStatement(t *testing.T) {
	inp := `{{ if true }}1{{ else if false }}2{{ else }}3{{ end }}`

	stmts := parseStatements(t, inp, 1, nil)

	ifStmt, ok := stmts[0].(*ast.IfStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an IfStatement, got %T", stmts[0])
	}

	if ifStmt.Alternative == nil {
		t.Errorf("ifStmt.Alternative is nil")
	}

	if len(ifStmt.Alternative.Statements) != 1 {
		t.Errorf("ifStmt.Alternative.Statements does not contain 1 statement, got %d", len(ifStmt.Alternative.Statements))
	}

	alternative, ok := ifStmt.Alternative.Statements[0].(*ast.HTMLStatement)

	if !ok {
		t.Fatalf("ifStmt.Alternative.Statements[0] is not an HTMLStatement, got %T", ifStmt.Alternative.Statements[0])
	}

	if alternative.String() != "3" {
		t.Errorf("alternative.String() is not %s, got %s", "3", alternative.String())
	}

	if len(ifStmt.Alternatives) != 1 {
		t.Errorf("ifStmt.Alternatives does not contain 1 statement, got %d", len(ifStmt.Alternatives))
	}

	elseIfAlternative := ifStmt.Alternatives[0]

	if !testBooleanLiteral(t, elseIfAlternative.Condition, false) {
		return
	}

	if len(elseIfAlternative.Consequence.Statements) != 1 {
		t.Errorf("alternative.Consequence.Statements does not contain 1 statement, got %d", len(elseIfAlternative.Consequence.Statements))
	}

	consequence, ok := elseIfAlternative.Consequence.Statements[0].(*ast.HTMLStatement)

	if !ok {
		t.Fatalf("alternative.Consequence.Statements[0] is not an HTMLStatement, got %T", elseIfAlternative.Consequence.Statements[0])
	}

	if consequence.String() != "2" {
		t.Errorf("consequence.String() is not %s, got %s", "2", consequence.String())
	}
}

func TestDefineStatement(t *testing.T) {
	tests := []struct {
		input    string
		varName  string
		varValue interface{}
	}{
		{`{{ var name = "Anna" }}`, "name", "Anna"},
		{`{{ var myAge = 34 }}`, "myAge", 34},
		{`{{ name := "Anna" }}`, "name", "Anna"},
		{`{{ myAge := 34 }}`, "myAge", 34},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.input, 1, nil)

		stmt, ok := stmts[0].(*ast.DefineStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not a DeclStatement, got %T", stmts[0])
		}

		if stmt.Name.Value != tt.varName {
			t.Errorf("stmt.Name.Value is not %s, got %s", tt.varName, stmt.Name.Value)
		}

		if !testLiteralExpression(t, stmt.Value, tt.varValue) {
			return
		}

		if stmt.String() != tt.input {
			t.Errorf("stmt.String() is not %s, got %s", tt.input, stmt.String())
		}
	}
}

func TestParseLayoutStatement(t *testing.T) {
	inp := `{{ layout "main" }}`

	stmts := parseStatements(t, inp, 1, nil)

	stmt, ok := stmts[0].(*ast.LayoutStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not a LayoutStatement, got %T", stmts[0])
	}

	if stmt.Path.Value != "main" {
		t.Errorf("stmt.Path.Value is not 'main', got %s", stmt.Path.Value)
	}

	if stmt.Program != nil {
		t.Errorf("stmt.Program is not nil, got %T", stmt.Program)
	}
}

func TestParseReserveStatement(t *testing.T) {
	inp := `{{ reserve "content" }}`

	stmts := parseStatements(t, inp, 1, map[string]*ast.InsertStatement{
		"content": {
			Name: &ast.StringLiteral{Value: "content"},
			Block: &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.HTMLStatement{
						Token: token.Token{Type: token.HTML, Literal: "<h1>Some content</h1>"},
					},
				},
			},
		},
	})

	stmt, ok := stmts[0].(*ast.ReserveStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not a ReserveStatement, got %T", stmts[0])
	}

	if stmt.Name.Value != "content" {
		t.Errorf("stmt.Name.Value is not 'content', got %s", stmt.Name.Value)
	}
}

func TestInsertStatement(t *testing.T) {
	t.Run("Insert with block", func(tt *testing.T) {
		inp := `{{ insert "content" }}<h1>Some content</h1>{{ end }}`

		stmts := parseStatements(t, inp, 1, nil)

		stmt, ok := stmts[0].(*ast.InsertStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not a InsertStatement, got %T", stmts[0])
		}

		if stmt.Name.Value != "content" {
			t.Errorf("stmt.Name.Value is not 'content', got %s", stmt.Name.Value)
		}

		if stmt.Block.String() != "<h1>Some content</h1>" {
			t.Errorf("stmt.Block.String() is not '<h1>Some content</h1>', got %s", stmt.Block.String())
		}
	})

	t.Run("Insert with argument", func(tt *testing.T) {
		inp := `{{ insert "content", "Some content" }}`

		stmts := parseStatements(t, inp, 1, nil)

		stmt, ok := stmts[0].(*ast.InsertStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not a InsertStatement, got %T", stmts[0])
		}

		if stmt.Name.Value != "content" {
			t.Errorf("stmt.Name.Value is not 'content', got %s", stmt.Name.Value)
		}

		if stmt.Block != nil {
			t.Errorf("stmt.Block is not nil, got %T", stmt.Block)
		}

		if stmt.Arguments[0].String() != `"Some content"` {
			t.Errorf("stmt.Arguments[0].String() is not 'Some content', got %s", stmt.Arguments[0].String())
		}
	})
}
