package parser

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/textwire/textwire/v3/pkg/ast"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/lexer"
	"github.com/textwire/textwire/v3/pkg/token"
	"github.com/textwire/textwire/v3/pkg/utils"
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
	p := New(l, nil)
	prog := p.ParseProgram()

	if opts.checkErrors {
		checkParserErrors(t, p)
	}

	if len(prog.Statements) != opts.stmtCount {
		t.Fatalf(
			"Program must have %d statementbut got %d for input %q",
			opts.stmtCount,
			len(prog.Statements),
			inp,
		)
	}

	return prog.Statements
}

func checkParserErrors(t *testing.T, p *Parser) {
	if !p.HasErrors() {
		return
	}

	t.Errorf("Parser has %d errors", len(p.Errors()))

	for _, msg := range p.Errors() {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testInfixExp(t *testing.T, exp ast.Expression, left any, op string, right any) {
	infix, ok := exp.(*ast.InfixExp)
	if !ok {
		t.Fatalf("Variable exp is not an InfixExp, got %T", exp)
	}

	testLiteralExpression(t, infix.Left, left)

	if infix.Op != op {
		t.Fatalf("infix.Op is not %s, got %s", op, infix.Op)
	}

	testLiteralExpression(t, infix.Right, right)
}

func testPosition(t *testing.T, actual, expect token.Position) {
	if expect.StartLine != actual.StartLine {
		t.Fatalf("expect.StartLine is not %d, got %d", expect.StartLine, actual.StartLine)
	}

	if expect.EndLine != actual.EndLine {
		t.Fatalf("expect.EndLine is not %d, got %d", expect.EndLine, actual.EndLine)
	}

	if expect.StartCol != actual.StartCol {
		t.Fatalf("expect.StartCol is not %d, got %d", expect.StartCol, actual.StartCol)
	}

	if expect.EndCol != actual.EndCol {
		t.Fatalf("expect.EndCol is not %d, got %d", expect.EndCol, actual.EndCol)
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	integer, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp is not an IntegerLiteral, got %T", exp)
	}

	if integer.Value != value {
		t.Fatalf("integer.Value is not %d, got %d", value, integer.Value)
	}

	if integer.Tok().Literal != strconv.FormatInt(value, 10) {
		t.Fatalf("integer.Tok().Literal is not %d, got %s", value, integer.Tok().Literal)
	}
}

func testFloatLiteral(t *testing.T, exp ast.Expression, value float64) {
	float, ok := exp.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("exp is not a FloatLiteral, got %T", exp)
	}

	if float.Value != value {
		t.Fatalf("float.Value is not %f, got %f", value, float.Value)
	}

	if float.String() != utils.FloatToStr(value) {
		t.Fatalf("float.String() is not %f, got %s", value, float)
	}
}

func testNilLiteral(t *testing.T, exp ast.Expression) {
	nilLit, ok := exp.(*ast.NilLiteral)
	if !ok {
		t.Fatalf("exp is not a NilLiteral, got %T", exp)
	}

	if nilLit.Tok().Literal != "nil" {
		t.Fatalf("nilLit.Tok().Literal is not 'nil', got %s", nilLit.Tok().Literal)
	}
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) {
	str, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp is not a StringLiteral, got %T", exp)
	}

	if str.Value != value {
		t.Fatalf("str.Value is not %s, got %s", value, str.Value)
	}

	if str.Tok().Literal != value {
		t.Fatalf("str.Tok().Literal is not %s, got %s", value, str.Tok().Literal)
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) {
	b, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Boolean, got %T", exp)
	}

	if b.Value != value {
		t.Fatalf("bo.Value not %t, got %t", value, b.Value)
	}

	if b.Tok().Literal != fmt.Sprintf("%t", value) {
		t.Fatalf("b.Tok().Literal is not %t, got %s", value, b.Tok().Literal)
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not an Identifier, got %T", exp)
	}

	if ident.Name != value {
		t.Fatalf("ident.Name is not %s, got %s", value, ident.Name)
	}

	if ident.Tok().Literal != value {
		t.Fatalf("ident.Tok().Literal is not %s, got %s", value, ident.Tok().Literal)
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expect any) {
	switch v := expect.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(v))
	case int64:
		testIntegerLiteral(t, exp, v)
	case float64:
		testFloatLiteral(t, exp, v)
	case string:
		testStringLiteral(t, exp, v)
	case bool:
		testBooleanLiteral(t, exp, v)
	case nil:
		testNilLiteral(t, exp)
	default:
		t.Fatalf("type of exp not handled. got %T", exp)
	}
}

func testIfBlock(t *testing.T, stmt ast.Statement, cond any, ifBlock string) {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok {
		t.Fatalf("stmt is not an IfStmt, got %T", stmt)
	}

	testLiteralExpression(t, ifStmt.Condition, cond)

	if ifStmt.IfBlock.String() != ifBlock {
		t.Fatalf("ifStmt.IfBlock.String() is not %q, got %q", ifBlock, ifStmt.IfBlock)
	}
}

func testElseBlock(t *testing.T, elseBlock *ast.BlockStmt, elseVal string) {
	if elseBlock == nil {
		t.Fatalf("elseBlock is nil")
	}

	if len(elseBlock.Statements) != 1 {
		t.Fatalf(
			"elseBlock.Statements does not contain 1 statement, got %d",
			len(elseBlock.Statements),
		)
	}

	if elseBlock.String() != elseVal {
		t.Fatalf("elseBlock.String() is not %q, got %q", elseBlock, elseVal)
	}
}

func testToken(t *testing.T, tok ast.Node, expect token.TokenType) {
	if tok.Tok().Type != expect {
		t.Fatalf(
			"Token type is not %q, got %q",
			token.String(expect),
			token.String(tok.Tok().Type),
		)
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
	testIdentifier(t, stmt.Expression, "myName")
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
	testIntegerLiteral(t, stmt.Expression, 234)
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
	testFloatLiteral(t, stmt.Expression, 2.34149)
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
		testStringLiteral(t, stmt.Expression, tc.expect)
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
		t.Fatalf("exp.Left is not %s, got %s", "Serhii", exp.Left)
	}

	if exp.Op != "+" {
		t.Fatalf("exp.Op is not %s, got %s", "+", exp.Op)
	}

	if exp.Right.Tok().Literal != " Anna" {
		t.Fatalf("exp.Right is not %s, got %s", " Anna", exp.Right)
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
	testIntegerLiteral(t, exp.Right, 2)
	testPosition(t, exp.Position(), token.Position{
		StartCol: 3,
		EndCol:   7,
	})

	if exp.Op != "+" {
		t.Fatalf("exp.Op is not %s, got %s", "+", exp.Op)
	}

	testIntegerLiteral(t, exp.Left, 5)
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
	testIntegerLiteral(t, exp.Right, 2)

	if exp.Op != "*" {
		t.Fatalf("exp.Op is not %s, got %s", "*", exp.Op)
	}

	infix, ok := exp.Left.(*ast.InfixExp)
	if !ok {
		t.Fatalf("exp.Left is not an InfixExp, got %T", exp.Left)
	}

	testIntegerLiteral(t, infix.Left, 5)

	if infix.Op != "+" {
		t.Fatalf("infix.Op is not %s, got %s", "+", infix.Op)
	}

	testLiteralExpression(t, infix.Right, 5)
}

func TestInfixExp(t *testing.T) {
	cases := []struct {
		inp    string
		left   any
		op     string
		right  any
		endCol uint
		expTok token.TokenType
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

		testInfixExp(t, stmt.Expression, tc.left, tc.op, tc.right)
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
		testBooleanLiteral(t, stmt.Expression, tc.expectBoolean)
		testPosition(t, stmt.Expression.Position(), token.Position{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})
	}
}

func TestPrefixExp(t *testing.T) {
	cases := []struct {
		inp    string
		op     string
		value  any
		endCol uint
		expTok token.TokenType
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

		if exp.Op != tc.op {
			t.Fatalf("exp.Op is not %s, got %s", tc.op, exp.Op)
		}

		testLiteralExpression(t, exp.Right, tc.value)
	}
}

func TestOpPrecedenceParsing(t *testing.T) {
	cases := []struct {
		id     int
		inp    string
		expect string
	}{
		{
			id:     1,
			inp:    "{{ 1 * 2 }}",
			expect: "(1 * 2)",
		},
		{
			id:     2,
			inp:    "<h2>{{ -2 + 3 }}</h2>",
			expect: "<h2>((-2) + 3)</h2>",
		},
		{
			id:     3,
			inp:    "{{ a + b + c }}",
			expect: "((a + b) + c)",
		},
		{
			id:     4,
			inp:    "{{ a + b / c }}",
			expect: "(a + (b / c))",
		},
		{
			id:     5,
			inp:    "{{ -2.float() }}",
			expect: "(-(2.float()))",
		},
		{
			id:     6,
			inp:    "{{ -5.0.int() }}",
			expect: "(-(5.0.int()))",
		},
		{
			id:     7,
			inp:    "{{ -obj.test }}",
			expect: "(-(obj.test))",
		},
		{
			id:     8,
			inp:    "{{ true && true || false }}",
			expect: "((true && true) || false)",
		},
		{
			id:     9,
			inp:    "{{ true ? 1 : 0 }}",
			expect: "(true ? 1 : 0)",
		},
		{
			id:     10,
			inp:    "{{ true && false ? 1 : 0 }}",
			expect: "((true && false) ? 1 : 0)",
		},
		{
			id:     11,
			inp:    "{{ true && false || 1 ? 1 : 0 }}",
			expect: "(((true && false) || 1) ? 1 : 0)",
		},
		{
			id:     12,
			inp:    "{{ -2.float() && -2.0.int() ? 1 : 0 }}",
			expect: "(((-(2.float())) && (-(2.0.int()))) ? 1 : 0)",
		},
		{
			id:     13,
			inp:    "{{ !defined(age) || !defined(name) ? 1 : 0 }}",
			expect: "(((!(defined(age))) || (!(defined(name)))) ? 1 : 0)",
		},
		{
			id:     14,
			inp:    "{{ defined(name) }}",
			expect: "(defined(name))",
		},
		{
			id:     15,
			inp:    "{{ long = user.name.len() > 0 }}",
			expect: "long = (((user.name).len()) > 0)",
		},
		{
			id:     16,
			inp:    "{{ user && user.name == 'serhii' }}",
			expect: `(user && ((user.name) == "serhii"))`,
		},
	}

	for _, tc := range cases {
		l := lexer.New(tc.inp)
		p := New(l, nil)
		prog := p.ParseProgram()

		checkParserErrors(t, p)

		if prog.String() != tc.expect {
			t.Fatalf("Test %d. Expect %s but got %s", tc.id, tc.expect, prog)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	cases := []struct {
		inp string
		err *fail.Error
	}{
		{
			inp: "{{ 5 + }}",
			err: fail.New(1, "", "parser", fail.ErrExpectedExpression),
		},
		{
			inp: "{{ }}",
			err: fail.New(1, "", "parser", fail.ErrEmptyBraces),
		},
		{
			inp: `{{ 1 ~ 8 }}`,
			err: fail.New(1, "", "parser", fail.ErrIllegalToken, "~"),
		},
		{
			inp: "{{ true ? 100 }}",
			err: fail.New(1, "", "parser", fail.ErrWrongNextToken,
				token.String(token.COLON), token.String(token.RBRACES)),
		},
		{
			inp: "{{ ) }}",
			err: fail.New(1, "", "parser", fail.ErrNoPrefixParseFunc,
				token.String(token.RPAREN)),
		},
		{
			inp: "@component('')",
			err: fail.New(1, "", "parser", fail.ErrExpectedComponentName),
		},
		{
			inp: "@use(1)",
			err: fail.New(1, "", "parser", fail.ErrUseStmtFirstArgStr, token.String(token.INT)),
		},
	}

	for _, tc := range cases {
		l := lexer.New(tc.inp)
		p := New(l, nil)
		p.ParseProgram()

		if !p.HasErrors() {
			t.Fatalf("no errors found in input %q", tc.inp)
		}

		if err := p.Errors()[0]; err.String() != tc.err.String() {
			t.Fatalf("expect error message %q, got %q", tc.err, err)
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
	testIntegerLiteral(t, exp.IfBlock, 100)
	testStringLiteral(t, exp.ElseBlock, "Some string")
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

	testIfBlock(t, stmt, true, "1")

	if stmt.ElseBlock != nil {
		t.Fatalf("ifStmt.ElseBlock is not nil, got %T", stmt.ElseBlock)
	}

	if len(stmt.ElseIfStmts) != 0 {
		t.Fatalf("ifStmt.ElseIfStmts is not empty, got %d", len(stmt.ElseIfStmts))
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
	testToken(t, stmt.IfBlock, token.HTML)
	testToken(t, stmt.ElseBlock, token.HTML)
	testIfBlock(t, stmt, true, "1")
	testElseBlock(t, stmt.ElseBlock, "2")

	testPosition(t, stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   19,
	})
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

	if len(ifStmt.IfBlock.Statements) != 3 {
		t.Fatalf(
			"ifStmt.IfBlock.Statements does not contain 3 statement, got %d",
			len(ifStmt.IfBlock.Statements),
		)
	}

	testToken(t, ifStmt, token.IF)
	testToken(t, ifStmt.IfBlock, token.HTML)
	testToken(t, ifStmt.ElseBlock, token.HTML)

	testPosition(t, ifStmt.Position(), token.Position{
		StartLine: 1,
		EndLine:   11,
		StartCol:  8,
		EndCol:    11,
	})

	testPosition(t, ifStmt.IfBlock.Position(), token.Position{
		StartLine: 1,
		EndLine:   9,
		StartCol:  17,
		EndCol:    7,
	})

	testPosition(t, ifStmt.ElseBlock.Position(), token.Position{
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

	testIfBlock(t, stmt, true, "first")

	if stmt.ElseBlock != nil {
		t.Fatalf("ifStmt.ElseBlock is not nil, got %T", stmt.ElseBlock)
	}

	if len(stmt.ElseIfStmts) != 1 {
		t.Fatalf("ifStmt.ElseIfStmts does not contain 1 statement, got %d", len(stmt.ElseIfStmts))
	}

	elseIfStmt := stmt.ElseIfStmts[0]
	if elseIfStmt, ok := elseIfStmt.(*ast.ElseIfStmt); ok {
		testBooleanLiteral(t, elseIfStmt.Condition, false)

		if len(elseIfStmt.Block.Statements) != 1 {
			t.Fatalf(
				"elseIfStmt.Block.Statements does not contain 1 statement, got %d",
				len(elseIfStmt.Block.Statements),
			)
		}

		htmlStmt, ok := elseIfStmt.Block.Statements[0].(*ast.HTMLStmt)
		if !ok {
			t.Fatalf(
				"elseIfStmt.Block.Statements[0] is not an HTMLStmt, got %T",
				elseIfStmt.Block.Statements[0],
			)
		}

		if htmlStmt.String() != "second" {
			t.Fatalf("htmlStmt.String() is not %q, got %q", "second", htmlStmt)
		}

		return
	}

	t.Fatalf("stmt.ElseIfStmts[0] is not an ElseIfStmt, got %T", elseIfStmt)
}

func TestParseElseIfWithElseStatement(t *testing.T) {
	inp := `@if(true)1@elseif(false)2@else3@end`

	stmts := parseStatements(t, inp, defaultParseOpts)
	stmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	testIfBlock(t, stmt, true, "1")
	testElseBlock(t, stmt.ElseBlock, "3")

	if len(stmt.ElseIfStmts) != 1 {
		t.Fatalf(
			"ifStmt.ElseIfStmts does not contain 1 statement, got %d",
			len(stmt.ElseIfStmts),
		)
	}

	if elseIfStmt, ok := stmt.ElseIfStmts[0].(*ast.ElseIfStmt); ok {
		testBooleanLiteral(t, elseIfStmt.Condition, false)

		if len(elseIfStmt.Block.Statements) != 1 {
			t.Fatalf(
				"elseIfStmt.Block.Statements does not contain 1 statement, got %d",
				len(elseIfStmt.Block.Statements),
			)
		}

		htmlStmt, ok := elseIfStmt.Block.Statements[0].(*ast.HTMLStmt)
		if !ok {
			t.Fatalf(
				"elseIfStmt.Block.Statements[0] is not an HTMLStmt, got %T",
				elseIfStmt.Block.Statements[0],
			)
		}

		if htmlStmt.String() != "2" {
			t.Fatalf("htmlStmt.String() is not %s, got %s", "2", htmlStmt)
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

		if stmt.Left.Name != tc.varName {
			t.Fatalf("stmt.Left.Name is not %s, got %s", tc.varName, stmt.Left.Name)
		}

		testLiteralExpression(t, stmt.Right, tc.varValue)

		testPosition(t, stmt.Position(), token.Position{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})

		if stmt.String() != tc.str {
			t.Fatalf("stmt.String() is not %s, got %s", tc.inp, stmt)
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
		t.Fatalf("stmt.Path.Value is not 'main', got %s", stmt.Name.Value)
	}

	testToken(t, stmt, token.USE)

	testPosition(t, stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   11,
	})

	if stmt.LayoutProg != nil {
		t.Fatalf("stmt.LayoutProg is not nil, got %T", stmt.LayoutProg)
	}

	if stmt.String() != inp {
		t.Fatalf("stmt.String() is not %s, got %s", inp, stmt)
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
		t.Fatalf("stmt.Name.Value is not 'content', got %s", stmt.Name.Value)
	}

	if stmt.String() == inp {
		t.Fatalf("stmt.String() is not %s, got %s", inp, stmt)
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
			t.Fatalf("stmt.Name.Value is not 'content', got %s", stmt.Name.Value)
		}

		if stmt.Block.String() != "<h1>Some content</h1>" {
			t.Fatalf(
				"stmt.Block.String() is not '<h1>Some content</h1>', got %s",
				stmt.Block,
			)
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
			t.Fatalf("stmt.Name.Value is not 'content', got %s", stmt.Name.Value)
		}

		if stmt.Block != nil {
			t.Fatalf("stmt.Block is not nil, got %T", stmt.Block)
		}

		if stmt.Argument.String() != `"Some content"` {
			t.Fatalf(
				"stmt.Argument.String() is not 'Some content', got %s",
				stmt.Argument,
			)
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

	testIntegerLiteral(t, arr.Elements[0], 11)
	testIntegerLiteral(t, arr.Elements[1], 234)
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
		t.Fatalf("indexExp.String() is not '(arr[(1 + 2)])', got %s", exp)
	}
}

func TestParsePostfixExp(t *testing.T) {
	cases := []struct {
		inp    string
		ident  string
		op     string
		str    string
		expTok token.TokenType
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
		testIdentifier(t, postfix.Left, tc.ident)

		if postfix.Op != tc.op {
			t.Fatalf("postfix.Op is not '%s', got %s", tc.op, postfix.Op)
		}

		if postfix.String() != tc.str {
			t.Fatalf("postfix.String() is not '%s', got %s", tc.str, postfix)
		}
	}
}

func TestParseTwoStatements(t *testing.T) {
	inp := `{{ name = "Anna"; name }}`

	stmts := parseStatements(t, inp, parseOpts{stmtCount: 2, checkErrors: true})

	testIdentifier(t, stmts[0].(*ast.AssignStmt).Left, "name")
	testToken(t, stmts[0], token.IDENT)

	testIdentifier(t, stmts[1].(*ast.ExpressionStmt).Expression, "name")
	testToken(t, stmts[0], token.IDENT)

	testStringLiteral(t, stmts[0].(*ast.AssignStmt).Right, "Anna")

	if stmts[0].String() != `name = "Anna"` {
		t.Fatalf("stmts[0].String() is not '{{ name = \"Anna\" }}', got %s", stmts[0])
	}

	if stmts[1].String() != `name` {
		t.Fatalf("stmts[1].String() is not '{{ name }}', got %s", stmts[1])
	}
}

func TestParseTwoExpression(t *testing.T) {
	inp := `{{ 1; 2 }}`
	stmts := parseStatements(t, inp, parseOpts{stmtCount: 2, checkErrors: true})
	testIntegerLiteral(t, stmts[0].(*ast.ExpressionStmt).Expression, 1)
	testIntegerLiteral(t, stmts[1].(*ast.ExpressionStmt).Expression, 2)
}

func TestParseGlobalCallExp(t *testing.T) {
	inp := `{{ defined(var1, var2) }}`

	stmts := parseStatements(t, inp, defaultParseOpts)

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.GlobalCallExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a GlobalCallExp, got %T", stmt.Expression)
	}

	testToken(t, exp, token.IDENT)
	testIdentifier(t, exp.Function, "defined")
	testIdentifier(t, exp.Arguments[0], "var1")
	testIdentifier(t, exp.Arguments[1], "var2")

	testPosition(t, exp.Position(), token.Position{
		StartCol: 3,
		EndCol:   21,
	})

	if len(exp.Arguments) != 2 {
		t.Fatalf("len(callExp.Arguments) is not 2, got %d", len(exp.Arguments))
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

	testStringLiteral(t, exp.Receiver, "Serhii Cho")
	testIdentifier(t, exp.Function, "split")

	if len(exp.Arguments) != 1 {
		t.Fatalf("len(callExp.Arguments) is not 1, got %d", len(exp.Arguments))
	}

	testStringLiteral(t, exp.Arguments[0], " ")
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
	testStringLiteral(t, callExp.Receiver, "")
	testIdentifier(t, callExp.Function, "len")
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
		t.Fatalf("stmt.Init.String() is not 'i = 0', got %s", stmt.Init)
	}

	if stmt.Condition.String() != `(i < 10)` {
		t.Fatalf("stmt.Condition.String() is not '(i < 10)', got %s", stmt.Condition)
	}

	if stmt.Post.String() != `(i++)` {
		t.Fatalf("stmt.Post.String() is not '(i++)', got %s", stmt.Post)
	}

	actual := strings.Trim(stmt.Block.String(), " \n\t")
	if actual != "{{ i }}" {
		t.Fatalf("actual is not '%q', got %q", "{{ i }}", actual)
	}

	if stmt.ElseBlock != nil {
		t.Fatalf("stmt.ElseBlock is not nil, got %T", stmt.ElseBlock)
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

	if stmt.ElseBlock == nil {
		t.Fatalf("stmt.ElseBlock is nil")
	}

	if stmt.ElseBlock.String() != "Empty" {
		t.Fatalf("stmt.ElseBlock.String() is not 'Empty', got %s", stmt.ElseBlock)
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
		t.Fatalf("stmt.Init is not nil, got %s", stmt.Init)
	}

	if stmt.Condition != nil {
		t.Fatalf("stmt.Condition is not nil, got %s", stmt.Condition)
	}

	if stmt.Post != nil {
		t.Fatalf("stmt.Post is not nil, got %s", stmt.Post)
	}

	if stmt.Block.String() != "1" {
		t.Fatalf("stmt.Block.String() is not '1', got %s", stmt.Block)
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
		t.Fatalf("stmt.Var.String() is not 'name', got %s", stmt.Var)
	}

	if stmt.Array.String() != `["anna", "serhii"]` {
		t.Fatalf(`stmt.Array.String() is not '["anna", "serhii"]', got %s`, stmt.Array)
	}

	actual := stmt.Block.String()
	if actual != "{{ name }}" {
		t.Fatalf("actual is not %q, got %q", "{{ name }}", actual)
	}

	if stmt.ElseBlock != nil {
		t.Fatalf("stmt.ElseBlock is not nil, got %T", stmt.ElseBlock)
	}
}

func TestParseStmtCanHaveEmptyBlock(t *testing.T) {
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

		if len(stmt.Stmts()) != 0 {
			t.Fatalf("len(stmt.Stmts()) has to be empty, got %d", len(stmt.Stmts()))
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

	if stmt.ElseBlock.String() != "Test" {
		t.Fatalf("stmt.ElseBlock.String() is not 'Test', got %s", stmt.ElseBlock)
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
		t.Fatalf(`obj.String() is not '{"father": {"name": "John" }}', got %s`, obj)
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
		t.Fatalf("dotExp.String() is not '((person.father).name)', got %s", exp)
	}

	testIdentifier(t, exp.Key, "name")

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

	testIdentifier(t, exp.Key, "father")
	testIdentifier(t, exp.Left, "person")
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
		t.Fatalf("breakStmt.String() is not '%s', got %s", expect, stmt)
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
		t.Fatalf("stmt.String() is not '%s', got %s", expect, stmt)
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
			t.Fatalf(`stmt.String() is not '%s', got %s`, expect, stmt)
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
		testStringLiteral(t, stmt.Slots[0].Name(), "header")
		testStringLiteral(t, stmt.Slots[1].Name(), "footer")

		expect := `@slot("header")<h1>Header</h1>@end`
		if stmt.Slots[0].String() != expect {
			t.Fatalf("stmt.Slots[0].String() is not '%q', got %q", expect, stmt.Slots[0])
		}

		expect = `@slot("footer")<footer>Footer</footer>@end`
		if stmt.Slots[1].String() != expect {
			t.Fatalf("stmt.Slots[0].String() is not '%q', got %q", expect, stmt.Slots[1])
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

		expect := `@component("some")`
		if stmt.String() != expect {
			t.Fatalf("stmt.String() is not `%s`, got `%s`", expect, stmt)
		}

		htmlStmt, htmlOk := stmts[1].(*ast.HTMLStmt)
		if !htmlOk {
			t.Fatalf("stmts[1] is not a HTMLStmt, got %T", stmts[1])
		}

		expect = "\n <b>Book</b>"
		if htmlStmt.String() != expect {
			t.Fatalf("htmlStmt.String() is not `%s`, got `%s`", expect, htmlStmt)
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
		testStringLiteral(t, stmt.Name(), "header")

		expect := `@slot("header")`
		if stmt.String() != expect {
			t.Fatalf("stmt.String() is not `%s`, got `%s`", expect, stmt)
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
		testNilLiteral(t, stmt.Name())
		testPosition(t, stmt.Position(), token.Position{
			StartCol: 8,
			EndCol:   12,
		})

		if stmt.String() != "@slot" {
			t.Fatalf("slot.String() is not @slot, got `%s`", stmt)
		}
	})
}

func TestParseSlotifDirective(t *testing.T) {
	t.Run("default slotif", func(t *testing.T) {
		inp := `@component('test')@slotif(true)Test@end@end`
		stmts := parseStatements(t, inp, parseOpts{stmtCount: 1, checkErrors: true})

		comp, ok := stmts[0].(*ast.ComponentStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a ComponentStmt, got %T", stmts[0])
		}

		if len(comp.Slots) > 1 {
			t.Fatalf("len(comp.Slots) must be 1, got %d", len(comp.Slots))
		}

		testPosition(t, comp.Position(), token.Position{EndCol: 42})
		testToken(t, comp, token.COMPONENT)

		slot, ok := comp.Slots[0].(*ast.SlotifStmt)
		if !ok {
			t.Fatalf("comp.Slots[0] is not a SlotifStmt, got %T", stmts[0])
		}

		testBooleanLiteral(t, slot.Condition, true)

		body := slot.Block().String()
		if body != "Test" {
			t.Fatalf("slotif.Block().String() is not 'Test', got %s", body)
		}

		expect := "@slotif(true)Test@end"
		if slot.String() != expect {
			t.Fatalf("slotif.String() is not '%s', got %s", expect, slot)
		}
	})

	t.Run("named slotif", func(t *testing.T) {
		inp := `@component('user')@slotif(false, 'name')Test2@end@end`
		stmts := parseStatements(t, inp, parseOpts{stmtCount: 1, checkErrors: true})

		comp, ok := stmts[0].(*ast.ComponentStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a ComponentStmt, got %T", stmts[0])
		}

		testPosition(t, comp.Position(), token.Position{EndCol: 52})
		testToken(t, comp, token.COMPONENT)

		if len(comp.Slots) > 1 {
			t.Fatalf("len(comp.Slots) must be 1, got %d", len(comp.Slots))
		}

		slot, ok := comp.Slots[0].(*ast.SlotifStmt)
		if !ok {
			t.Fatalf("comp.Slots[0] is not a SlotifStmt, got %T", stmts[0])
		}

		testBooleanLiteral(t, slot.Condition, false)

		if slot.Name().Value != "name" {
			t.Fatalf("slot.Name().Value is not 'name', got %s", slot.Name())
		}

		body := slot.Block().String()
		if body != "Test2" {
			t.Fatalf("slotif.Block().String() is not 'Test2', got %s", body)
		}

		expect := `@slotif(false, "name")Test2@end`
		if slot.String() != expect {
			t.Fatalf("slotif.String() is not '%s', got %s", expect, slot)
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

func TestParseBlockAsIllegalNode(t *testing.T) {
	inp := "@if(false)@dump(@end"

	stmts := parseStatements(t, inp, parseOpts{stmtCount: 1, checkErrors: false})
	stmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmt)
	}

	dumpStmt, ok := stmt.IfBlock.Statements[0].(*ast.DumpStmt)
	if !ok {
		t.Fatalf("stmt.IfBlock.Statements[0] is not an DumpStmt, got %T", dumpStmt)
	}

	if _, ok = dumpStmt.Arguments[0].(*ast.IllegalNode); !ok {
		t.Fatalf("dump.Arguments[0] is not an IllegalNode, got %T", dumpStmt.Arguments[0])
	}
}

func TestParseIllegalNode(t *testing.T) {
	cases := []struct {
		inp       string
		stmtCount int
	}{
		{"@if(false", 1},
		{"@if  (loop. {{ 'nice' }}@end", 1},
		{"@if {{ 'nice' }}@end", 1},
		{"@if( {{ 'nice' }}@end", 1},
		{"@each( {{ 'nice' }}@end", 1},
		{"@each() {{ 'nice' }}@end", 1},
		{"@each (loop. {{ 'nice' }}@end", 1},
		{"@each(nice in []{{ 'nice' }}@end", 1},
		{"@each(nice in {{ 'nice' }}@end", 1},
		{"@for( {{ 'nice' }}@end", 1},
		{"@for() {{ 'nice' }}@end", 1},
		{"@for(i {{ 'nice' }}@end", 1},
		{"@for(i = 0; i < []; i++{{ 'nice' }}@end", 1},
		{"@for(i = 0; i < [] {{ 'nice' }}@end", 1},
		{"@component('~user'", 1},
		{"@component   ('", 1},
		{"@component", 1},
		{"@insert('nice", 1},
		{"@insert ('nice'", 1},
		{"@insert('nice'@end", 1},
		{"@insert    ('nice' {{ 'nice' }}@end", 1},
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
			t.Fatalf("stmts[0] is not an IllegalNode, got %T for %s", stmts[0], tc.inp)
		}
	}
}
