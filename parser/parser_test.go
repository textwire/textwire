package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/textwire/textwire/v2/ast"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/lexer"
	"github.com/textwire/textwire/v2/token"
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

func testInfixExp(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
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

	if integer.Tok().Literal != strconv.FormatInt(value, 10) {
		t.Errorf("integer.Tok().Literal is not %d, got %s", value, integer.Tok().Literal)
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

	if nilLit.Tok().Literal != "nil" {
		t.Errorf("nilLit.Tok().Literal is not 'nil', got %s", nilLit.Tok().Literal)
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

	if str.Tok().Literal != value {
		t.Errorf("str.Tok().Literal is not %s, got %s", value, str.Tok().Literal)
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

	if boolean.Tok().Literal != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, boolean.Tok().Literal)
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

	if ident.Tok().Literal != value {
		t.Errorf("ident.Tok().Literal is not %s, got %s", value, ident.Tok().Literal)
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected any,
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

func testConsequence(t *testing.T, stmt ast.Statement, condition any, consequence string) bool {
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

	for _, tc := range tests {
		stmts := parseStatements(t, tc.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)

		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		if !testStringLiteral(t, stmt.Expression, tc.expect) {
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

	if exp.Left.Tok().Literal != "Serhii" {
		t.Fatalf("exp.Left is not %s, got %s", "Serhii", exp.Left.String())
	}

	if exp.Operator != "+" {
		t.Fatalf("exp.Operator is not %s, got %s", "+", exp.Operator)
	}

	if exp.Right.Tok().Literal != " Anna" {
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
		left     any
		operator string
		right    any
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

	for _, tc := range tests {
		stmts := parseStatements(t, tc.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		testInfixExp(t, stmt.Expression, tc.left, tc.operator, tc.right)
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

	for _, tc := range tests {
		stmts := parseStatements(t, tc.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		if !testBooleanLiteral(t, stmt.Expression, tc.expectBoolean) {
			return
		}
	}
}

func TestPrefixExp(t *testing.T) {
	tests := []struct {
		inp      string
		operator string
		value    any
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

	for _, tc := range tests {
		stmts := parseStatements(t, tc.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExp)
		if !ok {
			t.Fatalf("stmt is not a PrefixExp, got %T", stmt.Expression)
		}

		if exp.Operator != tc.operator {
			t.Fatalf("exp.Operator is not %s, got %s", tc.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tc.value) {
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
		{
			"{{ -2.float() }}",
			"{{ ((-2).float()) }}",
		},
		{
			"{{ -obj.test }}",
			"{{ ((-obj).test) }}",
		},
	}

	for _, tc := range tests {
		l := lexer.New(tc.inp)
		p := New(l, "")

		prog := p.ParseProgram()

		checkParserErrors(t, p)

		actual := prog.String()

		if actual != tc.expected {
			t.Errorf("expected=%q, got=%q", tc.expected, actual)
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
			fail.New(1, "", "parser", fail.ErrEmptyBraces),
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
		{
			"@component('')",
			fail.New(1, "", "parser", fail.ErrExpectedComponentName),
		},
	}

	for _, tc := range tests {
		l := lexer.New(tc.inp)
		p := New(l, "")

		p.ParseProgram()

		if len(p.Errors()) == 0 {
			t.Errorf("no errors found in input %q", tc.inp)
			return
		}

		err := p.Errors()[0]
		if err.String() != tc.err.String() {
			t.Errorf("expected error message %q, got %q", tc.err, err.String())
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

	testBooleanLiteral(t, exp.Condition, true)
	testIntegerLiteral(t, exp.Consequence, 100)
	testStringLiteral(t, exp.Alternative, "Some string")
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
		varValue any
		str      string
	}{
		{`{{ name = "Anna" }}`, "name", "Anna", `name = "Anna"`},
		{`{{ myAge = 34 }}`, "myAge", 34, `myAge = 34`},
	}

	for _, tc := range tests {
		stmts := parseStatements(t, tc.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.AssignStmt)

		if !ok {
			t.Fatalf("stmts[0] is not a AssignStmt, got %T", stmts[0])
		}

		if stmt.Name.Value != tc.varName {
			t.Errorf("stmt.Name.Value is not %s, got %s", tc.varName, stmt.Name.Value)
		}

		if !testLiteralExpression(t, stmt.Value, tc.varValue) {
			return
		}

		if stmt.String() != tc.str {
			t.Errorf("stmt.String() is not %s, got %s", tc.inp, stmt.String())
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
	t.Run("Insert with block", func(t *testing.T) {
		inp := `<h1>@insert("content")<h1>Some content</h1>@end</h1>`

		stmts := parseStatements(t, inp, 3, nil)
		stmt, ok := stmts[1].(*ast.InsertStmt)

		if !ok {
			t.Fatalf("stmts[1] is not a InsertStmt, got %T", stmts[0])
		}

		if stmt.Name.Value != "content" {
			t.Errorf("stmt.Name.Value is not 'content', got %s", stmt.Name.Value)
		}

		if stmt.Block.String() != "<h1>Some content</h1>" {
			t.Errorf("stmt.Block.String() is not '<h1>Some content</h1>', got %s",
				stmt.Block.String())
		}
	})

	t.Run("Insert with argument", func(t *testing.T) {
		inp := `<h1>@insert("content", "Some content")</h1>`

		stmts := parseStatements(t, inp, 3, nil)
		stmt, ok := stmts[1].(*ast.InsertStmt)

		if !ok {
			t.Fatalf("stmts[1] is not a InsertStmt, got %T", stmts[0])
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
	inp := `{{ [11, 234,] }}`

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

	for _, tc := range tests {
		stmts := parseStatements(t, tc.inp, 1, nil)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)

		if !ok {
			t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
		}

		postfix, ok := stmt.Expression.(*ast.PostfixExp)

		if !ok {
			t.Fatalf("stmt.Expression is not a PostfixExp, got %T", stmt.Expression)
		}

		if !testIdentifier(t, postfix.Left, tc.ident) {
			return
		}

		if postfix.Operator != tc.operator {
			t.Errorf("postfix.Operator is not '%s', got %s", tc.operator,
				postfix.Operator)
		}

		if postfix.String() != tc.str {
			t.Errorf("postfix.String() is not '%s', got %s", tc.str, postfix.String())
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
	inp := "@for(i = 0; i < 10; i++)\n{{ i }}\n@end"

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.ForStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ForStmt, got %T", stmts[0])
	}

	checkPosition(t, stmt.Pos, token.Position{
		StartLine: 0,
		EndLine:   2,
		StartCol:  0,
		EndCol:    3,
	})

	checkPosition(t, stmt.Block.Pos, token.Position{
		StartLine: 0,
		EndLine:   1,
		StartCol:  24,
		EndCol:    6,
	})

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

	if stmt.Block.String() != "\n{{ i }}\n" {
		t.Errorf("stmt.Block.String() is not '%q', got %q", "\n{{ i }}\n", stmt.Block.String())
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

func checkPosition(t *testing.T, actual, expect token.Position) {
	if expect.StartLine != actual.StartLine {
		t.Errorf("expect.StartLine is not %d, got %d", expect.StartLine,
			actual.StartLine)
	}

	if expect.EndLine != actual.EndLine {
		t.Errorf("expect.EndLine is not %d, got %d", expect.EndLine,
			actual.EndLine)
	}

	if expect.StartCol != actual.StartCol {
		t.Errorf("expect.StartCol is not %d, got %d", expect.StartCol,
			actual.StartCol)
	}

	if expect.EndCol != actual.EndCol {
		t.Errorf("expect.EndCol is not %d, got %d", expect.EndCol,
			actual.EndCol)
	}
}

func TestParseEachStmt(t *testing.T) {
	inp := "@each(name in ['anna', 'serhii'])\n{{ name }}\n@end"

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.EachStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a EachStmt, got %T", stmts[0])
	}

	checkPosition(t, stmt.Pos, token.Position{
		StartLine: 0,
		EndLine:   2,
		StartCol:  0,
		EndCol:    3,
	})

	checkPosition(t, stmt.Block.Pos, token.Position{
		StartLine: 0,
		EndLine:   1,
		StartCol:  33,
		EndCol:    9,
	})

	if stmt.Var.String() != `name` {
		t.Errorf("stmt.Var.String() is not 'name', got %s", stmt.Var.String())
	}

	if stmt.Array.String() != `["anna", "serhii"]` {
		t.Errorf(`stmt.Array.String() is not '["anna", "serhii"]', got %s`,
			stmt.Array.String())
	}

	if stmt.Block.String() != "\n{{ name }}\n" {
		t.Errorf("stmt.Block.String() is not %q, got %q", "\n{{ name }}\n", stmt.Block.String())
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
	inp := `{{ {"father": {name: "John"},} }}`

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

	expect := "@breakIf(true)"

	if breakStmt.String() != expect {
		t.Fatalf("breakStmt.String() is not '%s', got %s", expect, breakStmt.String())
	}
}

func TestParseContinueIfDirective(t *testing.T) {
	inp := `@continueIf(false)`
	stmts := parseStatements(t, inp, 1, nil)

	contStmt, ok := stmts[0].(*ast.ContinueIfStmt)

	if !ok {
		t.Fatalf("stmts[0] is not a ContinueIfStmt, got %T", stmts[0])
	}

	testBooleanLiteral(t, contStmt.Condition, false)

	expect := "@continueIf(false)"

	if contStmt.String() != expect {
		t.Fatalf("contStmt.String() is not '%s', got %s", expect, contStmt.String())
	}
}

func TestParseComponentDirective(t *testing.T) {
	t.Run("@component without slots", func(t *testing.T) {
		inp := `<ul>@component("components/book-card", { c: card })</ul>`
		stmts := parseStatements(t, inp, 3, nil)

		compStmt, ok := stmts[1].(*ast.ComponentStmt)

		if !ok {
			t.Fatalf("stmts[1] is not a ComponentStmt, got %T", stmts[1])
		}

		testStringLiteral(t, compStmt.Name, "components/book-card")

		if len(compStmt.Argument.Pairs) != 1 {
			t.Fatalf("len(compStmt.Arguments) is not 1, got %d", len(compStmt.Argument.Pairs))
		}

		testIdentifier(t, compStmt.Argument.Pairs["c"], "card")

		if len(compStmt.Slots) != 0 {
			t.Fatalf("len(compStmt.Slots) is not empty, got '%d' slots", len(compStmt.Slots))
		}

		expect := `@component("components/book-card", {"c": card})`

		if compStmt.String() != expect {
			t.Fatalf(`compStmt.String() is not '%s', got %s`, expect, compStmt.String())
		}
	})

	t.Run("@component with 1 slot", func(t *testing.T) {
		inp := `<ul>
			@component("components/book-card", { c: card })
				@slot("header")<h1>Header</h1>@end
				@slot("footer")<footer>Footer</footer>@end
			@end
		</ul>`

		stmts := parseStatements(t, inp, 3, nil)

		compStmt, ok := stmts[1].(*ast.ComponentStmt)

		if !ok {
			t.Fatalf("stmts[1] is not a ComponentStmt, got %T", stmts[1])
		}

		if len(compStmt.Slots) != 2 {
			t.Fatalf("len(compStmt.Slots) is not 2, got %d", len(compStmt.Slots))
		}

		testStringLiteral(t, compStmt.Slots[0].Name, "header")
		testStringLiteral(t, compStmt.Slots[1].Name, "footer")

		expect := "@slot(\"header\")\n<h1>Header</h1>\n@end"

		if compStmt.Slots[0].String() != expect {
			t.Fatalf("compStmt.Slots[0].String() is not '%s', got %s", expect,
				compStmt.Slots[0].String())
		}

		expect = "@slot(\"footer\")\n<footer>Footer</footer>\n@end"

		if compStmt.Slots[1].String() != expect {
			t.Fatalf("compStmt.Slots[0].String() is not '%s', got %s", expect,
				compStmt.Slots[1].String())
		}
	})

	t.Run("@component with whitespace at the end", func(t *testing.T) {
		inp := "@component('some')\n <b>Book</b>"
		stmts := parseStatements(t, inp, 2, nil)

		compStmt, ok := stmts[0].(*ast.ComponentStmt)

		if !ok {
			t.Fatalf("stmts[0] is not a ComponentStmt, got %T", stmts[1])
		}

		testStringLiteral(t, compStmt.Name, "some")

		expect := "@component(\"some\")"

		if compStmt.String() != expect {
			t.Fatalf("compStmt.String() is not `%s`, got `%s`", expect, compStmt.String())
		}

		htmlStmt, htmlOk := stmts[1].(*ast.HTMLStmt)

		if !htmlOk {
			t.Fatalf("stmts[1] is not a HTMLStmt, got %T", stmts[1])
		}

		expect = "\n <b>Book</b>"

		if htmlStmt.String() != expect {
			t.Fatalf("htmlStmt.String() is not `%s`, got `%s`", expect, htmlStmt.String())
		}
	})
}

func TestParseSlotDirective(t *testing.T) {
	t.Run("named slot", func(t *testing.T) {
		inp := `<h2>@slot("header")</h2>`
		stmts := parseStatements(t, inp, 3, nil)

		slotStmt, ok := stmts[1].(*ast.SlotStmt)

		if !ok {
			t.Fatalf("stmts[1] is not a SlotStmt, got %T", stmts[1])
		}

		testStringLiteral(t, slotStmt.Name, "header")

		expect := "@slot(\"header\")"

		if slotStmt.String() != expect {
			t.Fatalf("slotStmt.String() is not `%s`, got `%s`", expect, slotStmt.String())
		}
	})

	t.Run("default slot without end", func(t *testing.T) {
		t.Skip()
		inp := `<header>@slot</header>`
		stmts := parseStatements(t, inp, 3, nil)

		slotStmt, ok := stmts[1].(*ast.SlotStmt)

		if !ok {
			t.Fatalf("stmts[1] is not a SlotStmt, got %T", stmts[1])
		}

		testNilLiteral(t, slotStmt.Name)

		if slotStmt.String() != "@slot" {
			t.Fatalf("slotStmt.String() is not @slot, got `%s`", slotStmt.String())
		}
	})
}

func TestParseDumpStmt(t *testing.T) {
	inp := `@dump("test", 1 + 2, false)`

	stmts := parseStatements(t, inp, 1, nil)
	stmt, ok := stmts[0].(*ast.DumpStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an DumpStmt, got %T", stmts[0])
	}

	if len(stmt.Arguments) != 3 {
		t.Fatalf("len(stmt.Arguments) is not 3, got %d", len(stmt.Arguments))
	}

	testStringLiteral(t, stmt.Arguments[0], "test")
	testInfixExp(t, stmt.Arguments[1], 1, "+", 2)
	testBooleanLiteral(t, stmt.Arguments[2], false)
}
