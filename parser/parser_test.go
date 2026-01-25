package parser

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/textwire/textwire/v2/ast"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/lexer"
	"github.com/textwire/textwire/v2/token"
	"github.com/textwire/textwire/v2/utils"
)

type parseOpts struct {
	stmtCount   int
	inserts     map[string]*ast.InsertStmt
	checkErrors bool
}

var defaultParseOpts = parseOpts{
	stmtCount:   1,
	inserts:     nil,
	checkErrors: true,
}

func parseStatements(t *testing.T, inp string, opts parseOpts) []ast.Statement {
	l := lexer.New(inp)
	p := New(l, "")
	prog := p.ParseProgram()
	err := prog.ApplyInserts(opts.inserts, "")

	if err != nil {
		t.Fatalf("error applying inserts: %s", err)
	}

	if opts.checkErrors {
		checkParserErrors(t, p)
	}

	if len(prog.Statements) != opts.stmtCount {
		t.Fatalf("prog must have %d statement, got %d for input: %q", opts.stmtCount, len(prog.Statements), inp)
	}

	return prog.Statements
}

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

func testPosition(t *testing.T, actual, expect token.Position) {
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

	if float.String() != utils.FloatToStr(value) {
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
	expect any,
) bool {
	switch v := expect.(type) {
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

func testToken(t *testing.T, tok ast.Node, expect token.TokenType) {
	if tok.Tok().Type != expect {
		t.Errorf("Token type is not %q, got %q", token.String(expect), token.String(tok.Tok().Type))
	}
}

func TestParseIdentifier(t *testing.T) {
	stmts := parseStatements(t, "{{ myName }}", defaultParseOpts)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.IDENT)
	testToken(t, stmt.Expression, token.IDENT)

	if !testIdentifier(t, stmt.Expression, "myName") {
		return
	}
}

func TestParseExpressionStatement(t *testing.T) {
	stmts := parseStatements(t, "{{ 3 / 2 }}", defaultParseOpts)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.INT)

	testPosition(t, stmt.Position(), token.Position{
		StartCol: 3,
		EndCol:   7,
	})
}

func TestParseIntegerLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ 234 }}", defaultParseOpts)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.INT)

	if !testIntegerLiteral(t, stmt.Expression, 234) {
		return
	}

	testPosition(t, stmt.Expression.Position(), token.Position{
		StartCol: 3,
		EndCol:   5,
	})
}

func TestParseFloatLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ 2.34149 }}", defaultParseOpts)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.FLOAT)

	if !testFloatLiteral(t, stmt.Expression, 2.34149) {
		return
	}

	testPosition(t, stmt.Expression.Position(), token.Position{
		StartCol: 3,
		EndCol:   9,
	})
}

func TestParseNilLiteral(t *testing.T) {
	stmts := parseStatements(t, "{{ nil }}", defaultParseOpts)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.NIL)
	testNilLiteral(t, stmt.Expression)

	testPosition(t, stmt.Expression.Position(), token.Position{
		StartCol: 3,
		EndCol:   5,
	})
}

func TestParseStringLiteral(t *testing.T) {
	cases := []struct {
		inp      string
		expect   string
		startCol uint
		endCol   uint
	}{
		{`{{ "Hello World" }}`, "Hello World", 3, 15},
		{`{{ "Serhii \"Cho\"" }}`, `Serhii "Cho"`, 3, 18},
		{`{{ 'Hello World' }}`, "Hello World", 3, 15},
		{`{{ "" }}`, "", 3, 4},
	}

	for _, tc := range cases {
		stmts := parseStatements(t, tc.inp, defaultParseOpts)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)

		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		testToken(t, stmt, token.STR)

		if !testStringLiteral(t, stmt.Expression, tc.expect) {
			return
		}

		testPosition(t, stmt.Expression.Position(), token.Position{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})
	}
}

func TestStringConcatenation(t *testing.T) {
	inp := `{{ "Serhii" + " Anna" }}`

	stmts := parseStatements(t, inp, defaultParseOpts)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExp)

	if !ok {
		t.Fatalf("stmt is not an InfixExp, got %T", stmt.Expression)
	}

	testToken(t, stmt, token.STR)

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
func TestExpression(t *testing.T) {
	test := "{{ 5 + 2 }}"

	stmts := parseStatements(t, test, defaultParseOpts)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExp)
	if !ok {
		t.Fatalf("stmt is not an InfixExp, got %T", stmt.Expression)
	}

	testToken(t, stmt, token.INT)

	testPosition(t, exp.Position(), token.Position{
		StartCol: 3,
		EndCol:   7,
	})

	if !testIntegerLiteral(t, exp.Right, 2) {
		return
	}

	if exp.Operator != "+" {
		t.Fatalf("exp.Operator is not %s, got %s", "+", exp.Operator)
	}

	if !testIntegerLiteral(t, exp.Left, 5) {
		return
	}
}

func TestGroupedExpression(t *testing.T) {
	test := "{{ (5 + 5) * 2 }}"

	stmts := parseStatements(t, test, defaultParseOpts)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExp)
	if !ok {
		t.Fatalf("stmt is not an InfixExp, got %T", stmt.Expression)
	}

	testToken(t, stmt, token.LPAREN)

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
	cases := []struct {
		inp      string
		left     any
		operator string
		right    any
		endCol   uint
		expTok   token.TokenType
	}{
		{"{{ 5 + 8 }}", 5, "+", 8, 7, token.INT},
		{"{{ 10 - 2 }}", 10, "-", 2, 8, token.INT},
		{"{{ 2 * 2 }}", 2, "*", 2, 7, token.INT},
		{"{{ 44 / 4 }}", 44, "/", 4, 8, token.INT},
		{"{{ 5 % 4 }}", 5, "%", 4, 7, token.INT},
		{`{{ "me" + "her" }}`, "me", "+", "her", 14, token.STR},
		{`{{ 14 == 14 }}`, 14, "==", 14, 10, token.INT},
		{`{{ 10 != 1 }}`, 10, "!=", 1, 9, token.INT},
		{`{{ 19 > 31 }}`, 19, ">", 31, 9, token.INT},
		{`{{ 20 < 11 }}`, 20, "<", 11, 9, token.INT},
		{`{{ 19 >= 31 }}`, 19, ">=", 31, 10, token.INT},
		{`{{ 20 <= 11 }}`, 20, "<=", 11, 10, token.INT},
		{`{{ true && true }}`, true, "&&", true, 14, token.TRUE},
		{`{{ false || false }}`, false, "||", false, 16, token.FALSE},
	}

	for _, tc := range cases {
		stmts := parseStatements(t, tc.inp, defaultParseOpts)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		testToken(t, stmt, tc.expTok)

		testPosition(t, stmt.Expression.Position(), token.Position{
			StartCol: 3,
			EndCol:   tc.endCol,
		})

		testInfixExp(t, stmt.Expression, tc.left, tc.operator, tc.right)
	}
}

func TestBooleanExpression(t *testing.T) {
	cases := []struct {
		inp           string
		expectBoolean bool
		startCol      uint
		endCol        uint
		expTok        token.TokenType
	}{
		{"{{ true }}", true, 3, 6, token.TRUE},
		{"{{ false }}", false, 3, 7, token.FALSE},
	}

	for _, tc := range cases {
		stmts := parseStatements(t, tc.inp, defaultParseOpts)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		testToken(t, stmt, tc.expTok)

		if !testBooleanLiteral(t, stmt.Expression, tc.expectBoolean) {
			return
		}

		testPosition(t, stmt.Expression.Position(), token.Position{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})
	}
}

func TestPrefixExp(t *testing.T) {
	cases := []struct {
		inp      string
		operator string
		value    any
		endCol   uint
		expTok   token.TokenType
	}{
		{"{{ -5 }}", "-", 5, 4, token.SUB},
		{"{{ -10 }}", "-", 10, 5, token.SUB},
		{"{{ !true }}", "!", true, 7, token.NOT},
		{"{{ !false }}", "!", false, 8, token.NOT},
		{`{{ !"" }}`, "!", "", 5, token.NOT},
		{`{{ !0 }}`, "!", 0, 4, token.NOT},
		{`{{ -0 }}`, "-", 0, 4, token.SUB},
		{`{{ -0.0 }}`, "-", 0.0, 6, token.SUB},
		{`{{ !0.0 }}`, "!", 0.0, 6, token.NOT},
	}

	for _, tc := range cases {
		stmts := parseStatements(t, tc.inp, defaultParseOpts)
		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExp)
		if !ok {
			t.Fatalf("stmt is not a PrefixExp, got %T", stmt.Expression)
		}

		testToken(t, stmt, tc.expTok)

		testPosition(t, exp.Position(), token.Position{
			StartCol: 3,
			EndCol:   tc.endCol,
		})

		if exp.Operator != tc.operator {
			t.Fatalf("exp.Operator is not %s, got %s", tc.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tc.value) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
	}{
		{
			inp:    "{{ 1 * 2 }}",
			expect: "{{ (1 * 2) }}",
		},
		{
			inp:    "<h2>{{ -2 + 3 }}</h2>",
			expect: "<h2>{{ ((-2) + 3) }}</h2>",
		},
		{
			inp:    "{{ a + b + c }}",
			expect: "{{ ((a + b) + c) }}",
		},
		{
			inp:    "{{ a + b / c }}",
			expect: "{{ (a + (b / c)) }}",
		},
		{
			inp:    "{{ -2.float() }}",
			expect: "{{ ((-2).float()) }}",
		},
		{
			inp:    "{{ -5.0.int() }}",
			expect: "{{ ((-5.0).int()) }}",
		},
		{
			inp:    "{{ -obj.test }}",
			expect: "{{ ((-obj).test) }}",
		},
		{
			inp:    "{{ true && true || false }}",
			expect: "{{ ((true && true) || false) }}",
		},
		{
			inp:    "{{ true ? 1 : 0 }}",
			expect: "{{ (true ? 1 : 0) }}",
		},
		{
			inp:    "{{ true && false ? 1 : 0 }}",
			expect: "{{ ((true && false) ? 1 : 0) }}",
		},
		{
			inp:    "{{ true && false || 1 ? 1 : 0 }}",
			expect: "{{ (((true && false) || 1) ? 1 : 0) }}",
		},
		{
			inp:    "{{ -2.float() && -2.0.int() ? 1 : 0 }}",
			expect: "{{ ((((-2).float()) && ((-2.0).int())) ? 1 : 0) }}",
		},
	}

	for _, tc := range cases {
		l := lexer.New(tc.inp)
		p := New(l, "")

		prog := p.ParseProgram()

		checkParserErrors(t, p)

		actual := prog.String()

		if actual != tc.expect {
			t.Errorf("expect=%q, got=%q", tc.expect, actual)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	cases := []struct {
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
		{
			"@use(1)",
			fail.New(1, "", "parser", fail.ErrUseStmtFirstArgStr, token.String(token.INT)),
		},
	}

	for _, tc := range cases {
		l := lexer.New(tc.inp)
		p := New(l, "")

		p.ParseProgram()

		if len(p.Errors()) == 0 {
			t.Errorf("no errors found in input %q", tc.inp)
			return
		}

		err := p.Errors()[0]
		if err.String() != tc.err.String() {
			t.Errorf("expect error message %q, got %q", tc.err, err.String())
		}
	}
}

func TestTernaryExp(t *testing.T) {
	inp := `{{ true ? 100 : "Some string" }}`

	stmts := parseStatements(t, inp, defaultParseOpts)
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.TernaryExp)
	if !ok {
		t.Fatalf("stmt is not a TernaryExp, got %T", stmt.Expression)
	}

	testToken(t, stmt, token.TRUE)

	testPosition(t, exp.Position(), token.Position{
		StartCol: 3,
		EndCol:   28,
	})

	testBooleanLiteral(t, exp.Condition, true)
	testIntegerLiteral(t, exp.Consequence, 100)
	testStringLiteral(t, exp.Alternative, "Some string")
}

func TestParseIfStmt(t *testing.T) {
	inp := `@if(true)1@end`

	stmts := parseStatements(t, inp, defaultParseOpts)
	stmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.IF)

	testPosition(t, stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   13,
	})

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

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.IF)
	testToken(t, stmt.Consequence, token.HTML)
	testToken(t, stmt.Alternative, token.HTML)

	testPosition(t, stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   19,
	})

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
        @end`

	stmts := parseStatements(t, inp, parseOpts{stmtCount: 2, checkErrors: true})

	if _, ok := stmts[0].(*ast.HTMLStmt); !ok {
		t.Fatalf("stmts[0] is not an HTMLStmt, got %T", stmts[0])
	}

	ifStmt, isNotIfStmt := stmts[1].(*ast.IfStmt)
	if !isNotIfStmt {
		t.Fatalf("stmts[1] is not an IfStmt, got %T", stmts[1])
	}

	if len(ifStmt.Consequence.Statements) != 3 {
		t.Fatalf("ifStmt.Consequence.Statements does not contain 3 statement, got %d",
			len(ifStmt.Consequence.Statements))
	}

	testToken(t, ifStmt, token.IF)
	testToken(t, ifStmt.Consequence, token.HTML)
	testToken(t, ifStmt.Alternative, token.HTML)

	testPosition(t, ifStmt.Position(), token.Position{
		StartLine: 1,
		EndLine:   11,
		StartCol:  8,
		EndCol:    11,
	})

	testPosition(t, ifStmt.Consequence.Position(), token.Position{
		StartLine: 1,
		EndLine:   9,
		StartCol:  17,
		EndCol:    7,
	})

	testPosition(t, ifStmt.Alternative.Position(), token.Position{
		StartLine: 9,
		EndLine:   11,
		StartCol:  13,
		EndCol:    7,
	})
}

func TestParseIfElseIfStmt(t *testing.T) {
	inp := `@if(true)first@elseif(false)second@end`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	if !testConsequence(t, stmt, true, "first") {
		return
	}

	if stmt.Alternative != nil {
		t.Errorf("ifStmt.Alternative is not nil, got %T", stmt.Alternative)
	}

	if len(stmt.Alternatives) != 1 {
		t.Errorf("ifStmt.Alternatives does not contain 1 statement, got %d",
			len(stmt.Alternatives))
	}

	alt := stmt.Alternatives[0]
	if elseIfStmt, ok := alt.(*ast.ElseIfStmt); ok {
		if !testBooleanLiteral(t, elseIfStmt.Condition, false) {
			return
		}

		if len(elseIfStmt.Consequence.Statements) != 1 {
			t.Errorf("elseIfStmt.Consequence.Statements does not contain 1 statement, got %d",
				len(elseIfStmt.Consequence.Statements))
		}

		cons, ok := elseIfStmt.Consequence.Statements[0].(*ast.HTMLStmt)

		if !ok {
			t.Fatalf("elseIfStmt.Consequence.Statements[0] is not an HTMLStmt, got %T",
				elseIfStmt.Consequence.Statements[0])
		}

		if cons.String() != "second" {
			t.Errorf("cons.String() is not %q, got %q", "second", cons.String())
		}
		return
	}

	t.Errorf("stmt.Alternatives[0] is not an ElseIfStmt, got %T", alt)
}

func TestParseIfElseIfElseStatement(t *testing.T) {
	inp := `@if(true)1@elseif(false)2@else3@end`

	stmts := parseStatements(t, inp, defaultParseOpts)

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

	if elseIfAlt, ok := stmt.Alternatives[0].(*ast.ElseIfStmt); ok {
		if !testBooleanLiteral(t, elseIfAlt.Condition, false) {
			return
		}

		if len(elseIfAlt.Consequence.Statements) != 1 {
			t.Errorf("alternative.Consequence.Statements does not contain 1 statement, got %d",
				len(elseIfAlt.Consequence.Statements))
		}

		consequence, ok := elseIfAlt.Consequence.Statements[0].(*ast.HTMLStmt)

		if !ok {
			t.Fatalf("alternative.Consequence.Statements[0] is not an HTMLStmt, got %T",
				elseIfAlt.Consequence.Statements[0])
		}

		if consequence.String() != "2" {
			t.Errorf("consequence.String() is not %s, got %s", "2", consequence.String())
		}
	}
}

func TestParseAssignStmt(t *testing.T) {
	cases := []struct {
		inp      string
		varName  string
		varValue any
		str      string
		startCol uint
		endCol   uint
	}{
		{
			inp:      `{{ name = "Anna" }}`,
			varName:  "name",
			varValue: "Anna",
			str:      `name = "Anna"`,
			startCol: 3,
			endCol:   15,
		},
		{
			inp:      `{{ myAge = 34 }}`,
			varName:  "myAge",
			varValue: 34,
			str:      `myAge = 34`,
			startCol: 3,
			endCol:   12,
		},
	}

	for _, tc := range cases {
		stmts := parseStatements(t, tc.inp, defaultParseOpts)

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

		testPosition(t, stmt.Position(), token.Position{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})

		if stmt.String() != tc.str {
			t.Errorf("stmt.String() is not %s, got %s", tc.inp, stmt.String())
		}
	}
}

func TestParseUseStmt(t *testing.T) {
	inp := `@use("main")`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.UseStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a UseStmt, got %T", stmts[0])
	}

	if stmt.Name.Value != "main" {
		t.Errorf("stmt.Path.Value is not 'main', got %s", stmt.Name.Value)
	}

	testToken(t, stmt, token.USE)

	testPosition(t, stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   11,
	})

	if stmt.Program != nil {
		t.Errorf("stmt.Program is not nil, got %T", stmt.Program)
	}

	if stmt.String() != inp {
		t.Errorf("stmt.String() is not %s, got %s", inp, stmt.String())
	}
}

func TestParseReserveStmt(t *testing.T) {
	inp := `<div>@reserve("content")</div>`

	opts := parseOpts{
		stmtCount:   3,
		checkErrors: true,
		inserts: map[string]*ast.InsertStmt{
			"content": {
				Name: &ast.StringLiteral{Value: "content"},
				Block: &ast.BlockStmt{
					Statements: []ast.Statement{
						&ast.HTMLStmt{
							BaseNode: ast.BaseNode{
								Token: token.Token{
									Type:    token.HTML,
									Literal: "<h1>Some content</h1>",
								},
							},
						},
					},
				},
			},
		},
	}

	stmts := parseStatements(t, inp, opts)

	stmt, ok := stmts[1].(*ast.ReserveStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ReserveStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.RESERVE)

	testPosition(t, stmt.Position(), token.Position{
		StartCol: 5,
		EndCol:   23,
	})

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

		stmts := parseStatements(t, inp, parseOpts{stmtCount: 3, checkErrors: true})

		stmt, ok := stmts[1].(*ast.InsertStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a InsertStmt, got %T", stmts[0])
		}

		testToken(t, stmt, token.INSERT)

		testPosition(t, stmt.Position(), token.Position{
			StartCol: 4,
			EndCol:   46,
		})

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

		stmts := parseStatements(t, inp, parseOpts{stmtCount: 3, checkErrors: true})

		stmt, ok := stmts[1].(*ast.InsertStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a InsertStmt, got %T", stmts[0])
		}

		testPosition(t, stmt.Position(), token.Position{
			StartCol: 4,
			EndCol:   37,
		})

		testToken(t, stmt, token.INSERT)

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

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	arr, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not a ArrayLiteral, got %T", stmt.Expression)
	}

	testToken(t, arr, token.LBRACKET)

	testPosition(t, arr.Position(), token.Position{
		StartCol: 3,
		EndCol:   12,
	})

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

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.IndexExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a IndexExp, got %T", stmt.Expression)
	}

	testToken(t, exp, token.LBRACKET)

	// testing the last index [2]
	testPosition(t, exp.Position(), token.Position{
		StartCol: 13,
		EndCol:   15,
	})

	if exp.String() != "((arr[(1 + 2)])[2])" {
		t.Errorf("indexExp.String() is not '(arr[(1 + 2)])', got %s",
			exp.String())
	}
}

func TestParsePostfixExp(t *testing.T) {
	cases := []struct {
		inp      string
		ident    string
		operator string
		str      string
		expTok   token.TokenType
	}{
		{`{{ i++ }}`, "i", "++", "(i++)", token.INC},
		{`{{ num-- }}`, "num", "--", "(num--)", token.DEC},
	}

	for _, tc := range cases {
		stmts := parseStatements(t, tc.inp, defaultParseOpts)

		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
		}

		postfix, ok := stmt.Expression.(*ast.PostfixExp)
		if !ok {
			t.Fatalf("stmt.Expression is not a PostfixExp, got %T", stmt.Expression)
		}

		testToken(t, postfix, tc.expTok)

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

	stmts := parseStatements(t, inp, parseOpts{stmtCount: 2, checkErrors: true})

	if !testIdentifier(t, stmts[0].(*ast.AssignStmt).Name, "name") {
		return
	}

	testToken(t, stmts[0], token.IDENT)
	if !testStringLiteral(t, stmts[0].(*ast.AssignStmt).Value, "Anna") {
		return
	}

	testToken(t, stmts[0], token.IDENT)
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

	stmts := parseStatements(t, inp, parseOpts{stmtCount: 2, checkErrors: true})
	if !testIntegerLiteral(t, stmts[0].(*ast.ExpressionStmt).Expression, 1) {
		return
	}

	if !testIntegerLiteral(t, stmts[1].(*ast.ExpressionStmt).Expression, 2) {
		return
	}
}

func TestParseCallExp(t *testing.T) {
	inp := `{{ "Serhii Cho".split(" ") }}`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a CallExp, got %T", stmt.Expression)
	}

	testToken(t, exp, token.IDENT)

	testPosition(t, exp.Position(), token.Position{
		StartCol: 16,
		EndCol:   25,
	})

	if !testStringLiteral(t, exp.Receiver, "Serhii Cho") {
		return
	}

	if !testIdentifier(t, exp.Function, "split") {
		return
	}

	if len(exp.Arguments) != 1 {
		t.Fatalf("len(callExp.Arguments) is not 1, got %d", len(exp.Arguments))
	}

	if !testStringLiteral(t, exp.Arguments[0], " ") {
		return
	}
}

func TestParseCallExpWithExpressionList(t *testing.T) {
	inp := `{{ "nice".replace("n", "") }}`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a CallExp, got %T", stmt.Expression)
	}

	testToken(t, exp, token.IDENT)

	testPosition(t, exp.Position(), token.Position{
		StartCol: 10,
		EndCol:   25,
	})

	if len(exp.Arguments) != 2 {
		t.Fatalf("len(callExp.Arguments) is not 2, got %d", len(exp.Arguments))
	}
}

func TestParseCallExpWithEmptyString(t *testing.T) {
	inp := `{{ "".len() }}`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	callExp, ok := stmt.Expression.(*ast.CallExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a CallExp, got %T", stmt.Expression)
	}

	testToken(t, callExp, token.IDENT)

	if !testStringLiteral(t, callExp.Receiver, "") {
		return
	}

	if !testIdentifier(t, callExp.Function, "len") {
		return
	}
}

func TestParseForStmt(t *testing.T) {
	inp := `@for(i = 0; i < 10; i++)
        {{ i }}
    @end`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ForStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ForStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.FOR)

	testPosition(t, stmt.Pos, token.Position{
		EndLine: 2,
		EndCol:  7,
	})

	testPosition(t, stmt.Block.Pos, token.Position{
		EndLine:  2,
		StartCol: 24,
		EndCol:   3,
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

	actual := strings.Trim(stmt.Block.String(), " \n\t")
	if actual != "{{ i }}" {
		t.Errorf("actual is not '%q', got %q", "{{ i }}", actual)
	}

	if stmt.Alternative != nil {
		t.Errorf("stmt.Alternative is not nil, got %T", stmt.Alternative)
	}
}

func TestParseForElseStatement(t *testing.T) {
	inp := `@for(i = 0; i < 0; i++){{ i }}@elseEmpty@end`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ForStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ForStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.FOR)

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

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ForStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ForStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.FOR)

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
	inp := "@each(name in ['anna', 'serhii']){{ name }}@end"

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.EachStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a EachStmt, got %T", stmts[0])
	}

	testPosition(t, stmt.Pos, token.Position{EndCol: 46})
	testToken(t, stmt, token.EACH)

	testPosition(t, stmt.Block.Pos, token.Position{
		StartCol: 33,
		EndCol:   42,
	})

	if stmt.Var.String() != `name` {
		t.Errorf("stmt.Var.String() is not 'name', got %s", stmt.Var.String())
	}

	if stmt.Array.String() != `["anna", "serhii"]` {
		t.Errorf(`stmt.Array.String() is not '["anna", "serhii"]', got %s`,
			stmt.Array.String())
	}

	actual := stmt.Block.String()
	if actual != "{{ name }}" {
		t.Errorf("actual is not %q, got %q", "{{ name }}", actual)
	}

	if stmt.Alternative != nil {
		t.Errorf("stmt.Alternative is not nil, got %T", stmt.Alternative)
	}
}

func TestParseStmtCanHaveEmptyBody(t *testing.T) {
	cases := []struct {
		inp       string
		endColPos uint
		tok       token.TokenType
	}{
		{"@each(name in ['anna', 'serhii'])@end", 36, token.EACH},
		{"@for(i = 0; i < 10; i++)@end", 27, token.FOR},
		{"@if(true)@end", 12, token.IF},
		{"@insert('content')@end", 21, token.INSERT},
		{"@component('user')@slot('footer')@end@end", 40, token.COMPONENT},
	}

	for _, tc := range cases {
		stmts := parseStatements(t, tc.inp, defaultParseOpts)

		stmt, ok := stmts[0].(ast.NodeWithStatements)
		if !ok {
			t.Fatalf("stmts[0] is not a EachStmt, got %T", stmts[0])
		}

		testToken(t, stmt, tc.tok)
		testPosition(t, stmt.Position(), token.Position{EndCol: tc.endColPos})

		actual := len(stmt.Stmts())
		if actual != 0 {
			t.Errorf("len(stmt.Stmts()) has to be empty, got %d", actual)
		}
	}
}

func TestParseEachElseStatement(t *testing.T) {
	inp := `@each(v in []){{ v }}@elseTest@end`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.EachStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a EachStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.EACH)

	if stmt.Alternative.String() != "Test" {
		t.Errorf("stmt.Alternative.String() is not 'Test', got %s",
			stmt.Alternative.String())
	}
}

func TestParseObjectStatement(t *testing.T) {
	inp := `{{ {"father": {name: "John"},} }}`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	obj, ok := stmt.Expression.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	testPosition(t, obj.Position(), token.Position{
		StartCol: 3,
		EndCol:   29,
	})

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
	testToken(t, obj, token.LBRACE)
}

func TestParseObjectWithShorthandPropertyNotation(t *testing.T) {
	inp := `{{ { name, age } }}`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	obj, ok := stmt.Expression.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	testToken(t, obj, token.LBRACE)

	testPosition(t, obj.Position(), token.Position{
		StartCol: 3,
		EndCol:   15,
	})

	if len(obj.Pairs) != 2 {
		t.Fatalf("len(obj.Pairs) is not 2, got %d", len(obj.Pairs))
	}
}

func TestParseHTMLStmt(t *testing.T) {
	inp := "<div><span>Hello</span></div>"

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.HTMLStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a HTMLStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.HTML)

	testPosition(t, stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   28,
	})
}

func TestParseDotExp(t *testing.T) {
	inp := `{{ person.father.name }}`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.DotExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a DotExp, got %T", stmt.Expression)
	}

	testToken(t, exp, token.DOT)

	// position of the last dot between "father" and "name"
	testPosition(t, exp.Position(), token.Position{
		StartCol: 16,
		EndCol:   16,
	})

	if exp.String() != "((person.father).name)" {
		t.Fatalf("dotExp.String() is not '((person.father).name)', got %s",
			exp.String())
	}

	if !testIdentifier(t, exp.Key, "name") {
		return
	}

	if exp.Left == nil {
		t.Fatalf("dotExp.Left is nil")
	}

	exp, ok = exp.Left.(*ast.DotExp)
	if exp == nil {
		t.Fatalf("dotExp is nil")
		return
	}

	if !ok {
		t.Fatalf("dotExp.Left is not a DotExp, got %T", exp.Left)
	}

	if !testIdentifier(t, exp.Key, "father") {
		return
	}

	if !testIdentifier(t, exp.Left, "person") {
		return
	}
}

func TestParseBreakDirective(t *testing.T) {
	inp := `@break`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.BreakStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a BreakStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.BREAK)
}

func TestParseContinueDirective(t *testing.T) {
	inp := `@continue`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ContinueStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ContinueStmt, got %T", stmts[0])
	}

	testToken(t, stmt, token.CONTINUE)
}

func TestParseBreakIfDirective(t *testing.T) {
	inp := `@breakIf(true)`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.BreakIfStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a BreakIfStmt, got %T", stmts[0])
	}

	testPosition(t, stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   13,
	})

	testToken(t, stmt, token.BREAK_IF)
	testBooleanLiteral(t, stmt.Condition, true)

	expect := "@breakIf(true)"

	if stmt.String() != expect {
		t.Fatalf("breakStmt.String() is not '%s', got %s", expect, stmt.String())
	}
}

func TestParseContinueIfDirective(t *testing.T) {
	inp := "@continueIf(false)"
	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ContinueIfStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ContinueIfStmt, got %T", stmts[0])
	}

	testPosition(t, stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   17,
	})

	testToken(t, stmt, token.CONTINUE_IF)
	testBooleanLiteral(t, stmt.Condition, false)

	expect := "@continueIf(false)"

	if stmt.String() != expect {
		t.Fatalf("stmt.String() is not '%s', got %s", expect, stmt.String())
	}
}

func TestParseComponentDirective(t *testing.T) {
	t.Run("@component without slots", func(t *testing.T) {
		inp := "<ul>@component('components/book-card', { c: card })</ul>"
		stmts := parseStatements(t, inp, parseOpts{stmtCount: 3, checkErrors: true})

		stmt, ok := stmts[1].(*ast.ComponentStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a ComponentStmt, got %T", stmts[1])
		}

		testPosition(t, stmt.Position(), token.Position{
			StartCol: 4,
			EndCol:   50,
		})

		testToken(t, stmt, token.COMPONENT)
		testStringLiteral(t, stmt.Name, "components/book-card")

		if len(stmt.Argument.Pairs) != 1 {
			t.Fatalf("len(stmt.Arguments) is not 1, got %d", len(stmt.Argument.Pairs))
		}

		testIdentifier(t, stmt.Argument.Pairs["c"], "card")

		if len(stmt.Slots) != 0 {
			t.Fatalf("len(stmt.Slots) is not empty, got '%d' slots", len(stmt.Slots))
		}

		expect := `@component("components/book-card", {"c": card})`
		if stmt.String() != expect {
			t.Fatalf(`stmt.String() is not '%s', got %s`, expect, stmt.String())
		}
	})

	t.Run("@component with 1 slot", func(t *testing.T) {
		inp := `<ul>
			@component("components/book-card", { c: card })
				@slot("header")<h1>Header</h1>@end
				@slot("footer")<footer>Footer</footer>@end
			@end
		</ul>`

		stmts := parseStatements(t, inp, parseOpts{stmtCount: 3, checkErrors: true})

		stmt, ok := stmts[1].(*ast.ComponentStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a ComponentStmt, got %T", stmts[1])
		}

		testPosition(t, stmt.Position(), token.Position{
			StartLine: 1,
			EndLine:   4,
			StartCol:  3, // because tabs before @component
			EndCol:    6, // because tabs before @end
		})

		if len(stmt.Slots) != 2 {
			t.Fatalf("len(stmt.Slots) is not 2, got %d", len(stmt.Slots))
		}

		testToken(t, stmt, token.COMPONENT)
		testStringLiteral(t, stmt.Slots[0].Name, "header")
		testStringLiteral(t, stmt.Slots[1].Name, "footer")

		expect := "@slot(\"header\")\n<h1>Header</h1>\n@end"
		if stmt.Slots[0].String() != expect {
			t.Fatalf("stmt.Slots[0].String() is not '%q', got %q", expect,
				stmt.Slots[0].String())
		}

		expect = "@slot(\"footer\")\n<footer>Footer</footer>\n@end"
		if stmt.Slots[1].String() != expect {
			t.Fatalf("stmt.Slots[0].String() is not '%q', got %q", expect,
				stmt.Slots[1].String())
		}
	})

	t.Run("@component with whitespace at the end", func(t *testing.T) {
		inp := "@component('some')\n <b>Book</b>"
		stmts := parseStatements(t, inp, parseOpts{stmtCount: 2, checkErrors: true})

		stmt, ok := stmts[0].(*ast.ComponentStmt)
		if !ok {
			t.Fatalf("stmts[0] is not a ComponentStmt, got %T", stmts[1])
		}

		testToken(t, stmt, token.COMPONENT)
		testStringLiteral(t, stmt.Name, "some")

		expect := "@component(\"some\")"
		if stmt.String() != expect {
			t.Fatalf("stmt.String() is not `%s`, got `%s`", expect, stmt.String())
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
		stmts := parseStatements(t, inp, parseOpts{stmtCount: 3, checkErrors: true})

		stmt, ok := stmts[1].(*ast.SlotStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a SlotStmt, got %T", stmts[1])
		}

		testPosition(t, stmt.Position(), token.Position{
			StartCol: 4,
			EndCol:   18,
		})

		testToken(t, stmt, token.SLOT)
		testStringLiteral(t, stmt.Name, "header")

		expect := "@slot(\"header\")"

		if stmt.String() != expect {
			t.Fatalf("stmt.String() is not `%s`, got `%s`", expect, stmt.String())
		}
	})

	t.Run("default slot without end", func(t *testing.T) {
		t.Skip()
		inp := `<header>@slot</header>`
		stmts := parseStatements(t, inp, parseOpts{stmtCount: 3, checkErrors: true})

		stmt, ok := stmts[1].(*ast.SlotStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a SlotStmt, got %T", stmts[1])
		}

		testToken(t, stmt, token.SLOT)
		testNilLiteral(t, stmt.Name)
		testPosition(t, stmt.Position(), token.Position{
			StartCol: 8,
			EndCol:   12,
		})

		if stmt.String() != "@slot" {
			t.Fatalf("slot.String() is not @slot, got `%s`", stmt.String())
		}
	})
}

func TestParseDumpStmt(t *testing.T) {
	inp := `@dump("test", 1 + 2, false)`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.DumpStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an DumpStmt, got %T", stmts[0])
	}

	if len(stmt.Arguments) != 3 {
		t.Fatalf("len(stmt.Arguments) is not 3, got %d", len(stmt.Arguments))
	}

	testPosition(t, stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   26,
	})

	testToken(t, stmt, token.DUMP)
	testStringLiteral(t, stmt.Arguments[0], "test")
	testInfixExp(t, stmt.Arguments[1], 1, "+", 2)
	testBooleanLiteral(t, stmt.Arguments[2], false)
}

func TestParseBodyAsIllegalNode(t *testing.T) {
	inp := "@if(false)@dump(@end"

	stmts := parseStatements(t, inp, parseOpts{stmtCount: 1, checkErrors: false})
	stmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Errorf("stmts[0] is not an IfStmt, got %T", stmt)
	}

	dump, ok := stmt.Consequence.Statements[0].(*ast.DumpStmt)
	if !ok {
		t.Errorf("stmt.Consequence.Statements[0] is not an DumpStmt, got %T", dump)
	}

	_, ok = dump.Arguments[0].(*ast.IllegalNode)
	if !ok {
		t.Errorf("dump.Arguments[0] is not an IllegalNode, got %T", dump.Arguments[0])
	}
}

func TestParseIllegalNode(t *testing.T) {
	cases := []struct {
		inp       string
		stmtCount int
	}{
		{"@if(false", 1},
		{"@if(loop. {{ 'nice' }}@end", 1},
		{"@if {{ 'nice' }}@end", 1},
		{"@if( {{ 'nice' }}@end", 1},
		{"@each( {{ 'nice' }}@end", 1},
		{"@each() {{ 'nice' }}@end", 1},
		{"@each(loop. {{ 'nice' }}@end", 1},
		{"@each(nice in []{{ 'nice' }}@end", 1},
		{"@each(nice in {{ 'nice' }}@end", 1},
		{"@for( {{ 'nice' }}@end", 1},
		{"@for() {{ 'nice' }}@end", 1},
		{"@for(i {{ 'nice' }}@end", 1},
		{"@for(i = 0; i < []; i++{{ 'nice' }}@end", 1},
		{"@for(i = 0; i < [] {{ 'nice' }}@end", 1},
		{"@component('~user'", 1},
		{"@component('", 1},
		{"@component", 1},
		{"@insert('nice", 1},
		{"@insert('nice'", 1},
		{"@insert('nice'@end", 1},
		{"@insert('nice' {{ 'nice' }}@end", 1},
		{`@if(loop.
            {{ loop.first }}
            Iteration number is {{ loop.iter }}
        @end`, 1},
	}

	for _, tc := range cases {
		stmts := parseStatements(t, tc.inp, parseOpts{
			stmtCount:   tc.stmtCount,
			checkErrors: false,
		})

		_, ok := stmts[0].(*ast.IllegalNode)
		if !ok {
			t.Errorf("stmts[0] is not an IllegalNode, got %T for %s", stmts[0], tc.inp)
		}
	}
}
