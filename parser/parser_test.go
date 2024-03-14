package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/fail"
	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/token"
)

func checkParserErrors(t *testing.T, p *Parser) {
	if !p.HasErrors() {
		return
	}

	t.Errorf("parser has %d errors", len(p.Errors()))

	for _, msg := range p.Errors() {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func parseStatements(t *testing.T, inp string, stmtCount int, inserts map[string]*ast.InsertStmt) []ast.Statement {
	l := lexer.New(inp)
	p := New(l, "")
	prog := p.ParseProgram()
	err := prog.ApplyInserts(inserts, "")

	if err != nil {
		t.Fatalf("error applying inserts: %s", err)
	}

	checkParserErrors(t, p)

	if len(prog.Statements) != stmtCount {
		t.Fatalf("prog must have %d statement, got %d", stmtCount, len(prog.Statements))
	}

	return prog.Statements
}

func testInfixExp(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	infix, ok := exp.(*ast.InfixExp)

	if !ok {
		t.Errorf("exp is not an InfixExp, got %T", exp)
		return false
	}

	if !testLiteralExpression(t, infix.Left, left) {
		return false
	}

	if infix.Operator != operator {
		t.Errorf("infix.Operator is not %s, got %s", operator, infix.Operator)
		return false
	}

	if !testLiteralExpression(t, infix.Right, right) {
		return false
	}

	return true
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
	boolean, ok := exp.(*ast.BooleanLiteral)

	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if boolean.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, boolean.TokenLiteral())
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

func testConsequence(t *testing.T, stmt ast.Statement, condition interface{}, consequence string) bool {
	ifStmt, ok := stmt.(*ast.IfStmt)

	if !ok {
		t.Errorf("stmt is not an IfStmt, got %T", stmt)
		return false
	}

	if !testLiteralExpression(t, ifStmt.Condition, condition) {
		return false
	}

	if ifStmt.Consequence.String() != consequence {
		t.Errorf("ifStmt.Consequence.String() is not %q, got %q",
			consequence, ifStmt.Consequence.String())
		return false
	}

	return true
}

func testAlternative(t *testing.T, alt *ast.BlockStmt, altValue string) bool {
	if alt == nil {
		t.Errorf("alternative is nil")
		return false
	}

	if len(alt.Statements) != 1 {
		t.Errorf("alternative.Statements does not contain 1 statement, got %d",
			len(alt.Statements))

		return false
	}

	if alt.String() != altValue {
		t.Errorf("alternative.String() is not %q, got %q", alt.String(), altValue)
		return false
	}

	return true
}

func TestParseIdentifier(t *testing.T) {
	stmts := parseStatements(t, "{{ myName }}", 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	if !testIdentifier(t, stmt.Expression, "myName") {
		return
	}
}

func TestParseIntegerLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ 234 }}", 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	if !testIntegerLiteral(t, stmt.Expression, 234) {
		return
	}
}

func TestParseFloatLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ 2.34149 }}", 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	if !testFloatLiteral(t, stmt.Expression, 2.34149) {
		return
	}
}

func TestParseNilLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ nil }}", 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	testNilLiteral(t, stmt.Expression)
}

func TestParseStringLiteral(t *testing.T) {
	tests := []struct {
		inp    string
		expect string
	}{
		{`{{ "Hello World" }}`, "Hello World"},
		{`{{ "Serhii \"Cho\"" }}`, `Serhii "Cho"`},
		{`{{ 'Hello World' }}`, "Hello World"},
		{`{{ "" }}`, ""},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)

		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		if !testStringLiteral(t, stmt.Expression, tt.expect) {
			return
		}
	}
}

func TestStringConcatenation(t *testing.T) {
	inp := `{{ "Serhii" + " Anna" }}`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExp)

	if !ok {
		t.Fatalf("stmt is not an InfixExp, got %T", stmt.Expression)
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
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExp)

	if !ok {
		t.Fatalf("stmt is not an InfixExp, got %T", stmt.Expression)
	}

	if !testIntegerLiteral(t, exp.Right, 2) {
		return
	}

	if exp.Operator != "*" {
		t.Fatalf("exp.Operator is not %s, got %s", "*", exp.Operator)
	}

	infix, ok := exp.Left.(*ast.InfixExp)

	if !ok {
		t.Fatalf("exp.Left is not an InfixExp, got %T", exp.Left)
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

func TestInfixExp(t *testing.T) {
	tests := []struct {
		inp      string
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
		{`{{ 14 == 14 }}`, 14, "==", 14},
		{`{{ 10 != 1 }}`, 10, "!=", 1},
		{`{{ 19 > 31 }}`, 19, ">", 31},
		{`{{ 20 < 11 }}`, 20, "<", 11},
		{`{{ 19 >= 31 }}`, 19, ">=", 31},
		{`{{ 20 <= 11 }}`, 20, "<=", 11},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)

		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		testInfixExp(t, stmt.Expression, tt.left, tt.operator, tt.right)
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		inp           string
		expectBoolean bool
	}{
		{"{{ true }}", true},
		{"{{ false }}", false},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)

		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		if !testBooleanLiteral(t, stmt.Expression, tt.expectBoolean) {
			return
		}
	}
}

func TestPrefixExp(t *testing.T) {
	tests := []struct {
		inp      string
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
		stmts := parseStatements(t, tt.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)

		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExp)

		if !ok {
			t.Fatalf("stmt is not a PrefixExp, got %T", stmt.Expression)
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
		inp      string
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
		l := lexer.New(tt.inp)
		p := New(l, "")

		prog := p.ParseProgram()

		checkParserErrors(t, p)

		actual := prog.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		inp string
		err *fail.Error
	}{
		{
			"{{ 5 + }}",
			fail.New(1, "", "parser", fail.ErrExpectedExpression),
		},
		{
			"{{ }}",
			fail.New(1, "", "parser", fail.ErrEmptyBrackets),
		},
		{
			"{{ true ? 100 }}",
			fail.New(1, "", "parser", fail.ErrWrongNextToken,
				token.String(token.COLON),
				token.String(token.RBRACES)),
		},
		{
			"{{ ) }}",
			fail.New(1, "", "parser", fail.ErrNoPrefixParseFunc,
				token.String(token.RPAREN)),
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.inp)
		p := New(l, "")

		p.ParseProgram()

		if len(p.Errors()) == 0 {
			t.Errorf("no errors found in input %q", tt.inp)
			return
		}

		err := p.Errors()[0]

		if err.String() != tt.err.String() {
			t.Errorf("expected error message %q, got %q", tt.err, err.String())
		}
	}
}

func TestTernaryExp(t *testing.T) {
	inp := `{{ true ? 100 : "Some string" }}`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.TernaryExp)

	if !ok {
		t.Fatalf("stmt is not a TernaryExp, got %T", stmt.Expression)
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

func TestParseIfStmt(t *testing.T) {
	inp := `@if(true)1@end`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.IfStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	if !testConsequence(t, stmt, true, "1") {
		return
	}

	if stmt.Alternative != nil {
		t.Errorf("ifStmt.Alternative is not nil, got %T", stmt.Alternative)
	}

	if len(stmt.Alternatives) != 0 {
		t.Errorf("ifStmt.Alternatives is not empty, got %d", len(stmt.Alternatives))
	}
}

func TestParseIfElseStatement(t *testing.T) {
	inp := `@if(true)1@else2@end`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.IfStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	if !testConsequence(t, stmt, true, "1") {
		return
	}

	if !testAlternative(t, stmt.Alternative, "2") {
		return
	}
}

func TestParseNestedIfElseStatement(t *testing.T) {
	inp := `
		@if(true)
			@if(false)
				James
			@elseif(false)
				John
			@else
				@if(true){{ "Marry" }}@end
			@end
		@else
			@if(true)Anna@end
		@end
	`

	stmts := parseStatements(t, inp, 3, nil)

	if _, ok := stmts[0].(*ast.HTMLStmt); !ok {
		t.Fatalf("stmts[0] is not an HTMLStmt, got %T", stmts[0])
	}

	if _, ok := stmts[2].(*ast.HTMLStmt); !ok {
		t.Fatalf("stmts[2] is not an HTMLStmt, got %T", stmts[0])
	}

	ifStmt, isNotIfStmt := stmts[1].(*ast.IfStmt)

	if !isNotIfStmt {
		t.Fatalf("stmts[1] is not an IfStmt, got %T", stmts[0])
	}

	if len(ifStmt.Consequence.Statements) != 3 {
		t.Fatalf("ifStmt.Consequence.Statements does not contain 3 statement, got %d",
			len(ifStmt.Consequence.Statements))
	}
}

func TestParseIfElseIfStmt(t *testing.T) {
	inp := `@if(true)1@elseif(false)2@end`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.IfStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	if !testConsequence(t, stmt, true, "1") {
		return
	}

	if stmt.Alternative != nil {
		t.Errorf("ifStmt.Alternative is not nil, got %T", stmt.Alternative)
	}

	if len(stmt.Alternatives) != 1 {
		t.Errorf("ifStmt.Alternatives does not contain 1 statement, got %d",
			len(stmt.Alternatives))
	}

	alternative := stmt.Alternatives[0]

	if !testBooleanLiteral(t, alternative.Condition, false) {
		return
	}

	if len(alternative.Consequence.Statements) != 1 {
		t.Errorf("alternative.Consequence.Statements does not contain 1 statement, got %d",
			len(alternative.Consequence.Statements))
	}

	consequence, ok := alternative.Consequence.Statements[0].(*ast.HTMLStmt)

	if !ok {
		t.Fatalf("alternative.Consequence.Statements[0] is not an HTMLStmt, got %T",
			alternative.Consequence.Statements[0])
	}

	if consequence.String() != "2" {
		t.Errorf("consequence.String() is not %s, got %s", "2", consequence.String())
	}
}

func TestParseIfElseIfElseStatement(t *testing.T) {
	inp := `@if(true)1@elseif(false)2@else3@end`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.IfStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	if !testConsequence(t, stmt, true, "1") {
		return
	}

	if !testAlternative(t, stmt.Alternative, "3") {
		return
	}

	if len(stmt.Alternatives) != 1 {
		t.Errorf("ifStmt.Alternatives does not contain 1 statement, got %d",
			len(stmt.Alternatives))
	}

	elseIfAlternative := stmt.Alternatives[0]

	if !testBooleanLiteral(t, elseIfAlternative.Condition, false) {
		return
	}

	if len(elseIfAlternative.Consequence.Statements) != 1 {
		t.Errorf("alternative.Consequence.Statements does not contain 1 statement, got %d",
			len(elseIfAlternative.Consequence.Statements))
	}

	consequence, ok := elseIfAlternative.Consequence.Statements[0].(*ast.HTMLStmt)

	if !ok {
		t.Fatalf("alternative.Consequence.Statements[0] is not an HTMLStmt, got %T",
			elseIfAlternative.Consequence.Statements[0])
	}

	if consequence.String() != "2" {
		t.Errorf("consequence.String() is not %s, got %s", "2", consequence.String())
	}
}

func TestParseAssignStmt(t *testing.T) {
	tests := []struct {
		inp      string
		varName  string
		varValue interface{}
		str      string
	}{
		{`{{ name = "Anna" }}`, "name", "Anna", `name = "Anna"`},
		{`{{ myAge = 34 }}`, "myAge", 34, `myAge = 34`},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.AssignStmt)

		if !ok {
			t.Fatalf("stmts[0] is not a AssignStmt, got %T", stmts[0])
		}

		if stmt.Name.Value != tt.varName {
			t.Errorf("stmt.Name.Value is not %s, got %s", tt.varName, stmt.Name.Value)
		}

		if !testLiteralExpression(t, stmt.Value, tt.varValue) {
			return
		}

		if stmt.String() != tt.str {
			t.Errorf("stmt.String() is not %s, got %s", tt.inp, stmt.String())
		}
	}
}

func TestParseUseStmt(t *testing.T) {
	inp := `@use("main")`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.UseStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a UseStmt, got %T", stmts[0])
	}

	if stmt.Name.Value != "main" {
		t.Errorf("stmt.Path.Value is not 'main', got %s", stmt.Name.Value)
	}

	if stmt.Program != nil {
		t.Errorf("stmt.Program is not nil, got %T", stmt.Program)
	}

	if stmt.String() != inp {
		t.Errorf("stmt.String() is not %s, got %s", inp, stmt.String())
	}
}

func TestParseReserveStmt(t *testing.T) {
	inp := `<div>@reserve("content")</div>`

	stmts := parseStatements(t, inp, 3, map[string]*ast.InsertStmt{
		"content": {
			Name: &ast.StringLiteral{Value: "content"},
			Block: &ast.BlockStmt{
				Statements: []ast.Statement{
					&ast.HTMLStmt{
						Token: token.Token{
							Type:    token.HTML,
							Literal: "<h1>Some content</h1>",
						},
					},
				},
			},
		},
	})

	stmt, ok := stmts[1].(*ast.ReserveStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ReserveStmt, got %T", stmts[0])
	}

	if stmt.Name.Value != "content" {
		t.Errorf("stmt.Name.Value is not 'content', got %s", stmt.Name.Value)
	}

	if stmt.String() == inp {
		t.Errorf("stmt.String() is not %s, got %s", inp, stmt.String())
	}
}

func TestInsertStmt(t *testing.T) {
	t.Run("Insert with block", func(tt *testing.T) {
		inp := `@insert("content")<h1>Some content</h1>@end`

		stmts := parseStatements(t, inp, 1, nil)
		stmt, ok := stmts[0].(*ast.InsertStmt)

		if !ok {
			t.Fatalf("stmts[0] is not a InsertStmt, got %T", stmts[0])
		}

		if stmt.Name.Value != "content" {
			t.Errorf("stmt.Name.Value is not 'content', got %s", stmt.Name.Value)
		}

		if stmt.Block.String() != "<h1>Some content</h1>" {
			t.Errorf("stmt.Block.String() is not '<h1>Some content</h1>', got %s",
				stmt.Block.String())
		}
	})

	t.Run("Insert with argument", func(tt *testing.T) {
		inp := `@insert("content", "Some content")`

		stmts := parseStatements(t, inp, 1, nil)
		stmt, ok := stmts[0].(*ast.InsertStmt)

		if !ok {
			t.Fatalf("stmts[0] is not a InsertStmt, got %T", stmts[0])
		}

		if stmt.Name.Value != "content" {
			t.Errorf("stmt.Name.Value is not 'content', got %s", stmt.Name.Value)
		}

		if stmt.Block != nil {
			t.Errorf("stmt.Block is not nil, got %T", stmt.Block)
		}

		if stmt.Argument.String() != `"Some content"` {
			t.Errorf("stmt.Argument.String() is not 'Some content', got %s", stmt.Argument.String())
		}
	})
}

func TestParseArray(t *testing.T) {
	inp := `{{ [11, 234] }}`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	arr, ok := stmt.Expression.(*ast.ArrayLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not a ArrayLiteral, got %T", stmt.Expression)
	}

	if len(arr.Elements) != 2 {
		t.Fatalf("len(arr.Elements) is not 2, got %d", len(arr.Elements))
	}

	if !testIntegerLiteral(t, arr.Elements[0], 11) {
		return
	}

	if !testIntegerLiteral(t, arr.Elements[1], 234) {
		return
	}

	if arr.String() != "[11, 234]" {
		t.Errorf("arr.String() is not '[11, 234]', got %s", arr.String())
	}
}

func TestParseIndexExp(t *testing.T) {
	inp := `{{ arr[1 + 2][2] }}`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	indexExp, ok := stmt.Expression.(*ast.IndexExp)

	if !ok {
		t.Fatalf("stmt.Expression is not a IndexExp, got %T", stmt.Expression)
	}

	if indexExp.String() != "((arr[(1 + 2)])[2])" {
		t.Errorf("indexExp.String() is not '(arr[(1 + 2)])', got %s",
			indexExp.String())
	}
}

func TestParsePostfixExp(t *testing.T) {
	tests := []struct {
		inp      string
		ident    string
		operator string
		str      string
	}{
		{`{{ i++ }}`, "i", "++", "(i++)"},
		{`{{ num-- }}`, "num", "--", "(num--)"},
	}

	for _, tt := range tests {
		stmts := parseStatements(t, tt.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)

		if !ok {
			t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
		}

		postfix, ok := stmt.Expression.(*ast.PostfixExp)

		if !ok {
			t.Fatalf("stmt.Expression is not a PostfixExp, got %T", stmt.Expression)
		}

		if !testIdentifier(t, postfix.Left, tt.ident) {
			return
		}

		if postfix.Operator != tt.operator {
			t.Errorf("postfix.Operator is not '%s', got %s", tt.operator,
				postfix.Operator)
		}

		if postfix.String() != tt.str {
			t.Errorf("postfix.String() is not '%s', got %s", tt.str, postfix.String())
		}
	}
}

func TestParseTwoStatements(t *testing.T) {
	inp := `{{ name = "Anna"; name }}`

	stmts := parseStatements(t, inp, 2, nil)

	if !testIdentifier(t, stmts[0].(*ast.AssignStmt).Name, "name") {
		return
	}

	if !testStringLiteral(t, stmts[0].(*ast.AssignStmt).Value, "Anna") {
		return
	}

	if !testIdentifier(t, stmts[1].(*ast.ExpressionStmt).Expression, "name") {
		return
	}

	if stmts[0].String() != `name = "Anna"` {
		t.Errorf("stmts[0].String() is not '{{ name = \"Anna\" }}', got %s",
			stmts[0].String())
	}

	if stmts[1].String() != `name` {
		t.Errorf("stmts[1].String() is not '{{ name }}', got %s", stmts[1].String())
	}
}

func TestParseTwoExpression(t *testing.T) {
	inp := `{{ 1; 2 }}`

	stmts := parseStatements(t, inp, 2, nil)

	if !testIntegerLiteral(t, stmts[0].(*ast.ExpressionStmt).Expression, 1) {
		return
	}

	if !testIntegerLiteral(t, stmts[1].(*ast.ExpressionStmt).Expression, 2) {
		return
	}
}

func TestParseCallExp(t *testing.T) {
	inp := `{{ "Serhii Cho".split(" ") }}`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	callExp, ok := stmt.Expression.(*ast.CallExp)

	if !ok {
		t.Fatalf("stmt.Expression is not a CallExp, got %T", stmt.Expression)
	}

	if !testStringLiteral(t, callExp.Receiver, "Serhii Cho") {
		return
	}

	if !testIdentifier(t, callExp.Function, "split") {
		return
	}

	if len(callExp.Arguments) != 1 {
		t.Fatalf("len(callExp.Arguments) is not 1, got %d", len(callExp.Arguments))
	}

	if !testStringLiteral(t, callExp.Arguments[0], " ") {
		return
	}
}

func TestParseCallExpWithEmptyString(t *testing.T) {
	inp := `{{ "".len() }}`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	callExp, ok := stmt.Expression.(*ast.CallExp)

	if !ok {
		t.Fatalf("stmt.Expression is not a CallExp, got %T", stmt.Expression)
	}

	if !testStringLiteral(t, callExp.Receiver, "") {
		return
	}

	if !testIdentifier(t, callExp.Function, "len") {
		return
	}
}

func TestParseForStmt(t *testing.T) {
	inp := `@for(i = 0; i < 10; i++){{ i }}@end`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ForStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ForStmt, got %T", stmts[0])
	}

	if stmt.Init.String() != `i = 0` {
		t.Errorf("stmt.Init.String() is not 'i = 0', got %s", stmt.Init.String())
	}

	if stmt.Condition.String() != `(i < 10)` {
		t.Errorf("stmt.Condition.String() is not '(i < 10)', got %s",
			stmt.Condition.String())
	}

	if stmt.Post.String() != `(i++)` {
		t.Errorf("stmt.Post.String() is not '(i++)', got %s", stmt.Post.String())
	}

	if stmt.Block.String() != `{{ i }}` {
		t.Errorf("stmt.Block.String() is not '{{ i }}', got %s", stmt.Block.String())
	}

	if stmt.Alternative != nil {
		t.Errorf("stmt.Alternative is not nil, got %T", stmt.Alternative)
	}
}

func TestParseForElseStatement(t *testing.T) {
	inp := `@for(i = 0; i < 0; i++){{ i }}@elseEmpty@end`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ForStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ForStmt, got %T", stmts[0])
	}

	if stmt.Alternative == nil {
		t.Fatalf("stmt.Alternative is nil")
	}

	if stmt.Alternative.String() != "Empty" {
		t.Errorf("stmt.Alternative.String() is not 'Empty', got %s",
			stmt.Alternative.String())
	}
}

func TestParseInfiniteForStmt(t *testing.T) {
	inp := `@for(;;)1@end`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ForStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ForStmt, got %T", stmts[0])
	}

	if stmt.Init != nil {
		t.Errorf("stmt.Init is not nil, got %s", stmt.Init.String())
	}

	if stmt.Condition != nil {
		t.Errorf("stmt.Condition is not nil, got %s", stmt.Condition.String())
	}

	if stmt.Post != nil {
		t.Errorf("stmt.Post is not nil, got %s", stmt.Post.String())
	}

	if stmt.Block.String() != "1" {
		t.Errorf("stmt.Block.String() is not '1', got %s", stmt.Block.String())
	}
}

func TestParseEachStmt(t *testing.T) {
	inp := `@each(name in ["anna", "serhii"]){{ name }}@end`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.EachStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a EachStmt, got %T", stmts[0])
	}

	if stmt.Var.String() != `name` {
		t.Errorf("stmt.Var.String() is not 'name', got %s", stmt.Var.String())
	}

	if stmt.Array.String() != `["anna", "serhii"]` {
		t.Errorf(`stmt.Array.String() is not '["anna", "serhii"]', got %s`,
			stmt.Array.String())
	}

	if stmt.Block.String() != `{{ name }}` {
		t.Errorf("stmt.Block.String() is not '{{ name }}', got %s", stmt.Block.String())
	}

	if stmt.Alternative != nil {
		t.Errorf("stmt.Alternative is not nil, got %T", stmt.Alternative)
	}
}

func TestParseEachElseStatement(t *testing.T) {
	inp := `@each(v in []){{ v }}@elseTest@end`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.EachStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a EachStmt, got %T", stmts[0])
	}

	if stmt.Alternative.String() != "Test" {
		t.Errorf("stmt.Alternative.String() is not 'Test', got %s",
			stmt.Alternative.String())
	}
}

func TestParseObjectStatement(t *testing.T) {
	inp := `{{ {"father": {name: "John"}} }}`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	obj, ok := stmt.Expression.(*ast.ObjectLiteral)

	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	if len(obj.Pairs) != 1 {
		t.Fatalf("len(obj.Pairs) is not 1, got %d", len(obj.Pairs))
	}

	if obj.String() != `{"father": {"name": "John"}}` {
		t.Fatalf(`obj.String() is not '{"father": {"name": "John" }}', got %s`,
			obj.String())
	}

	nested, ok := obj.Pairs["father"].(*ast.ObjectLiteral)

	if !ok {
		t.Fatalf("obj.Pairs['father'] is not a ObjectLiteral, got %T",
			obj.Pairs["father"])
	}

	testStringLiteral(t, nested.Pairs["name"], "John")
}

func TestParseObjectWithShorthandPropertyNotation(t *testing.T) {
	inp := `{{ { name, age } }}`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	obj, ok := stmt.Expression.(*ast.ObjectLiteral)

	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	if len(obj.Pairs) != 2 {
		t.Fatalf("len(obj.Pairs) is not 2, got %d", len(obj.Pairs))
	}

	if obj.String() != `{"name": name, "age": age}` {
		t.Fatalf(`obj.String() is not '{"name": name, "age": age}', got %s`,
			obj.String())
	}
}

func TestParseDotExp(t *testing.T) {
	inp := `{{ person.father.name }}`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	dotExp, ok := stmt.Expression.(*ast.DotExp)

	if !ok {
		t.Fatalf("stmt.Expression is not a DotExp, got %T", stmt.Expression)
	}

	if dotExp.String() != "((person.father).name)" {
		t.Fatalf("dotExp.String() is not '((person.father).name)', got %s",
			dotExp.String())
	}

	if !testIdentifier(t, dotExp.Key, "name") {
		return
	}

	if dotExp.Left == nil {
		t.Fatalf("dotExp.Left is nil")
	}

	dotExp, ok = dotExp.Left.(*ast.DotExp)

	if dotExp == nil {
		t.Fatalf("dotExp is nil")
		return
	}

	if !ok {
		t.Fatalf("dotExp.Left is not a DotExp, got %T", dotExp.Left)
	}

	if !testIdentifier(t, dotExp.Key, "father") {
		return
	}

	if !testIdentifier(t, dotExp.Left, "person") {
		return
	}
}

func TestParseBreakDirective(t *testing.T) {
	inp := `@break`

	stmts := parseStatements(t, inp, 1, nil)
	_, ok := stmts[0].(*ast.BreakStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a BreakStmt, got %T", stmts[0])
	}
}

func TestParseContinueDirective(t *testing.T) {
	inp := `@continue`

	stmts := parseStatements(t, inp, 1, nil)
	_, ok := stmts[0].(*ast.ContinueStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ContinueStmt, got %T", stmts[0])
	}
}

func TestParseBreakIfDirective(t *testing.T) {
	inp := `@breakIf(true)`

	stmts := parseStatements(t, inp, 1, nil)
	breakStmt, ok := stmts[0].(*ast.BreakIfStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a BreakIfStmt, got %T", stmts[0])
	}

	testBooleanLiteral(t, breakStmt.Condition, true)
}

func TestParseContinueIfDirective(t *testing.T) {
	inp := `@continueIf(false)`
	stmts := parseStatements(t, inp, 1, nil)

	contStmt, ok := stmts[0].(*ast.ContinueIfStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ContinueIfStmt, got %T", stmts[0])
	}

	testBooleanLiteral(t, contStmt.Condition, false)
}

func TestParseComponentDirective(t *testing.T) {
	inp := `<ul>@component("components/book-card", card)</ul>`
	stmts := parseStatements(t, inp, 3, nil)

	compStmt, ok := stmts[1].(*ast.ComponentStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ComponentStmt, got %T", stmts[0])
	}

	testStringLiteral(t, compStmt.Name, "components/book-card")

	if len(compStmt.Arguments) != 1 {
		t.Fatalf("len(compStmt.Arguments) is not 1, got %d", len(compStmt.Arguments))
	}

	testIdentifier(t, compStmt.Arguments[0], "card")

	if compStmt.Block != nil {
		t.Fatalf("compStmt.Block is not nil, got %T", compStmt.Block)
	}
}
