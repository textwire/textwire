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

func parseStatements(inp string, opts parseOpts) ([]ast.Statement, error) {
	l := lexer.New(inp)
	p := New(l, nil)
	prog := p.ParseProgram()

	if opts.checkErrors && p.HasErrors() {
		return nil, p.Errors()[0].Error()
	}

	if len(prog.Statements) != opts.stmtCount {
		return nil, fmt.Errorf(
			"Program must have %d statement but got %d for input %q",
			opts.stmtCount,
			len(prog.Statements),
			inp,
		)
	}

	return prog.Statements, nil
}

func testInfixExp(exp ast.Expression, left any, op string, right any) error {
	infix, ok := exp.(*ast.InfixExp)
	if !ok {
		return fmt.Errorf("Variable exp is not an InfixExp, got %T", exp)
	}

	if err := testLiteralExpression(infix.Left, left); err != nil {
		return err
	}

	if infix.Op != op {
		return fmt.Errorf("infix.Op is not %s, got %s", op, infix.Op)
	}

	if err := testLiteralExpression(infix.Right, right); err != nil {
		return err
	}

	return nil
}

func testPosition(actual, expect token.Position) error {
	if expect.StartLine != actual.StartLine {
		return fmt.Errorf("expect.StartLine is not %d, got %d", expect.StartLine, actual.StartLine)
	}

	if expect.EndLine != actual.EndLine {
		return fmt.Errorf("expect.EndLine is not %d, got %d", expect.EndLine, actual.EndLine)
	}

	if expect.StartCol != actual.StartCol {
		return fmt.Errorf("expect.StartCol is not %d, got %d", expect.StartCol, actual.StartCol)
	}

	if expect.EndCol != actual.EndCol {
		return fmt.Errorf("expect.EndCol is not %d, got %d", expect.EndCol, actual.EndCol)
	}

	return nil
}

func testIntegerLiteral(exp ast.Expression, value int64) error {
	integer, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		return fmt.Errorf("exp is not an IntegerLiteral, got %T", exp)
	}

	if integer.Value != value {
		return fmt.Errorf("integer.Value is not %d, got %d", value, integer.Value)
	}

	if integer.Tok().Literal != strconv.FormatInt(value, 10) {
		return fmt.Errorf("integer.Tok().Literal is not %d, got %s", value, integer.Tok().Literal)
	}

	return nil
}

func testFloatLiteral(exp ast.Expression, value float64) error {
	float, ok := exp.(*ast.FloatLiteral)
	if !ok {
		return fmt.Errorf("exp is not a FloatLiteral, got %T", exp)
	}

	if float.Value != value {
		return fmt.Errorf("float.Value is not %f, got %f", value, float.Value)
	}

	if float.String() != utils.FloatToStr(value) {
		return fmt.Errorf("float.String() is not %f, got %s", value, float)
	}

	return nil
}

func testNilLiteral(exp ast.Expression) error {
	nilLit, ok := exp.(*ast.NilLiteral)
	if !ok {
		return fmt.Errorf("exp is not a NilLiteral, got %T", exp)
	}

	if nilLit.Tok().Literal != "nil" {
		return fmt.Errorf("nilLit.Tok().Literal is not 'nil', got %s", nilLit.Tok().Literal)
	}

	return nil
}

func testStringLiteral(exp ast.Expression, value string) error {
	str, ok := exp.(*ast.StringLiteral)
	if !ok {
		return fmt.Errorf("exp is not a StringLiteral, got %T", exp)
	}

	if str.Value != value {
		return fmt.Errorf("str.Value is not %s, got %s", value, str.Value)
	}

	if str.Tok().Literal != value {
		return fmt.Errorf("str.Tok().Literal is not %s, got %s", value, str.Tok().Literal)
	}

	return nil
}

func testBooleanLiteral(exp ast.Expression, value bool) error {
	b, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		return fmt.Errorf("exp not *ast.Boolean, got %T", exp)
	}

	if b.Value != value {
		return fmt.Errorf("bo.Value not %t, got %t", value, b.Value)
	}

	if b.Tok().Literal != fmt.Sprintf("%t", value) {
		return fmt.Errorf("b.Tok().Literal is not %t, got %s", value, b.Tok().Literal)
	}

	return nil
}

func testIdentifier(exp ast.Expression, value string) error {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		return fmt.Errorf("exp is not an Identifier, got %T", exp)
	}

	if ident.Name != value {
		return fmt.Errorf("ident.Name is not %s, got %s", value, ident.Name)
	}

	if ident.Tok().Literal != value {
		return fmt.Errorf("ident.Tok().Literal is not %s, got %s", value, ident.Tok().Literal)
	}

	return nil
}

func testLiteralExpression(exp ast.Expression, expect any) error {
	switch v := expect.(type) {
	case int:
		return testIntegerLiteral(exp, int64(v))
	case int64:
		return testIntegerLiteral(exp, v)
	case float64:
		return testFloatLiteral(exp, v)
	case string:
		return testStringLiteral(exp, v)
	case bool:
		return testBooleanLiteral(exp, v)
	case nil:
		return testNilLiteral(exp)
	default:
		return fmt.Errorf("type of exp not handled. got %T", exp)
	}
}

func testIfBlock(stmt ast.Statement, cond any, ifBlock string) error {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok {
		return fmt.Errorf("stmt is not an IfStmt, got %T", stmt)
	}

	if err := testLiteralExpression(ifStmt.Condition, cond); err != nil {
		return err
	}

	if ifStmt.IfBlock.String() != ifBlock {
		return fmt.Errorf("ifStmt.IfBlock.String() is not %q, got %q", ifBlock, ifStmt.IfBlock)
	}

	return nil
}

func testElseBlock(elseBlock *ast.BlockStmt, elseVal string) error {
	if elseBlock == nil {
		return fmt.Errorf("elseBlock is nil")
	}

	if len(elseBlock.Statements) != 1 {
		return fmt.Errorf(
			"elseBlock.Statements does not contain 1 statement, got %d",
			len(elseBlock.Statements),
		)
	}

	if elseBlock.String() != elseVal {
		return fmt.Errorf("elseBlock.String() is not %q, got %q", elseBlock, elseVal)
	}

	return nil
}

func testToken(tok ast.Node, expect token.TokenType) error {
	if tok.Tok().Type != expect {
		return fmt.Errorf(
			"Token type is not %q, got %q",
			token.String(expect),
			token.String(tok.Tok().Type),
		)
	}
	return nil
}

func TestParseIdentifier(t *testing.T) {
	stmts, err := parseStatements("{{ myName }}", defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testToken(stmt.Expression, token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testIdentifier(stmt.Expression, "myName"); err != nil {
		t.Fatal(err)
	}
}

func TestParseExpressionStatement(t *testing.T) {
	stmts, err := parseStatements("{{ 3 / 2 }}", defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.INT); err != nil {
		t.Fatal(err)
	}

	err = testPosition(stmt.Position(), token.Position{
		StartCol: 3,
		EndCol:   7,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseIntegerLiteral(t *testing.T) {
	stmts, err := parseStatements("{{ 234 }}", defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.INT); err != nil {
		t.Fatal(err)
	}

	if err := testIntegerLiteral(stmt.Expression, 234); err != nil {
		t.Fatal(err)
	}
	err = testPosition(stmt.Expression.Position(), token.Position{
		StartCol: 3,
		EndCol:   5,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseFloatLiteral(t *testing.T) {
	stmts, err := parseStatements("{{ 2.34149 }}", defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.FLOAT); err != nil {
		t.Fatal(err)
	}

	if err := testFloatLiteral(stmt.Expression, 2.34149); err != nil {
		t.Fatal(err)
	}
	err = testPosition(stmt.Expression.Position(), token.Position{
		StartCol: 3,
		EndCol:   9,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseNilLiteral(t *testing.T) {
	stmts, err := parseStatements("{{ nil }}", defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.NIL); err != nil {
		t.Fatal(err)
	}

	if err := testNilLiteral(stmt.Expression); err != nil {
		t.Fatal(err)
	}

	err = testPosition(stmt.Expression.Position(), token.Position{
		StartCol: 3,
		EndCol:   5,
	})
	if err != nil {
		t.Fatal(err)
	}
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
		stmts, err := parseStatements(tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		if err := testToken(stmt, token.STR); err != nil {
			t.Fatal(err)
		}
		if err := testStringLiteral(stmt.Expression, tc.expect); err != nil {
			t.Fatal(err)
		}
		err = testPosition(stmt.Expression.Position(), token.Position{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestStringConcatenation(t *testing.T) {
	inp := `{{ "Serhii" + " Anna" }}`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExp)

	if !ok {
		t.Fatalf("stmt is not an InfixExp, got %T", stmt.Expression)
	}

	if err := testToken(stmt, token.STR); err != nil {
		t.Fatal(err)
	}

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

	stmts, err := parseStatements(test, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExp)
	if !ok {
		t.Fatalf("stmt is not an InfixExp, got %T", stmt.Expression)
	}

	if err := testToken(stmt, token.INT); err != nil {
		t.Fatal(err)
	}

	if err := testIntegerLiteral(exp.Right, 2); err != nil {
		t.Fatal(err)
	}
	err = testPosition(exp.Position(), token.Position{
		StartCol: 3,
		EndCol:   7,
	})
	if err != nil {
		t.Fatal(err)
	}

	if exp.Op != "+" {
		t.Fatalf("exp.Op is not %s, got %s", "+", exp.Op)
	}

	if err := testIntegerLiteral(exp.Left, 5); err != nil {
		t.Fatal(err)
	}
}

func TestGroupedExpression(t *testing.T) {
	test := "{{ (5 + 5) * 2 }}"

	stmts, err := parseStatements(test, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExp)
	if !ok {
		t.Fatalf("stmt is not an InfixExp, got %T", stmt.Expression)
	}

	if err := testToken(stmt, token.LPAREN); err != nil {
		t.Fatal(err)
	}

	if err := testIntegerLiteral(exp.Right, 2); err != nil {
		t.Fatal(err)
	}

	if exp.Op != "*" {
		t.Fatalf("exp.Op is not %s, got %s", "*", exp.Op)
	}

	infix, ok := exp.Left.(*ast.InfixExp)
	if !ok {
		t.Fatalf("exp.Left is not an InfixExp, got %T", exp.Left)
	}

	if err := testIntegerLiteral(infix.Left, 5); err != nil {
		t.Fatal(err)
	}

	if infix.Op != "+" {
		t.Fatalf("infix.Op is not %s, got %s", "+", infix.Op)
	}

	if err := testLiteralExpression(infix.Right, 5); err != nil {
		t.Fatal(err)
	}
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
		stmts, err := parseStatements(tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		if err := testToken(stmt, tc.expTok); err != nil {
			t.Fatal(err)
		}

		err = testPosition(stmt.Expression.Position(), token.Position{
			StartCol: 3,
			EndCol:   tc.endCol,
		})
		if err != nil {
			t.Fatal(err)
		}

		if err := testInfixExp(stmt.Expression, tc.left, tc.op, tc.right); err != nil {
			t.Fatal(err)
		}
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
		stmts, err := parseStatements(tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		if err := testToken(stmt, tc.expTok); err != nil {
			t.Fatal(err)
		}
		if err := testBooleanLiteral(stmt.Expression, tc.expectBoolean); err != nil {
			t.Fatal(err)
		}
		err = testPosition(stmt.Expression.Position(), token.Position{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})
		if err != nil {
			t.Fatal(err)
		}
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
		stmts, err := parseStatements(tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExp)
		if !ok {
			t.Fatalf("stmt is not a PrefixExp, got %T", stmt.Expression)
		}

		if err := testToken(stmt, tc.expTok); err != nil {
			t.Fatal(err)
		}
		err = testPosition(exp.Position(), token.Position{
			StartCol: 3,
			EndCol:   tc.endCol,
		})
		if err != nil {
			t.Fatal(err)
		}

		if exp.Op != tc.op {
			t.Fatalf("exp.Op is not %s, got %s", tc.op, exp.Op)
		}

		if err := testLiteralExpression(exp.Right, tc.value); err != nil {
			t.Fatal(err)
		}
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

		if p.HasErrors() {
			t.Fatalf("Case: %d. Parser error:  %v", tc.id, p.Errors()[0])
		}

		if prog.String() != tc.expect {
			t.Fatalf("Case: %d. Expect %s but got %s", tc.id, tc.expect, prog)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	cases := []struct {
		inp string
		err *fail.Error
	}{
		{
			inp: `{{ { "1st": "nice" }.1st }}`,
			err: fail.New(1, "", "parser", fail.ErrObjectKeyUseGet),
		},
		{
			inp: "<div>@slotif(true)No@end</div>",
			err: fail.New(1, "", "parser", fail.ErrSlotifPosition),
		},
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
			inp: "@use('')",
			err: fail.New(1, "", "parser", fail.ErrExpectedUseName),
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

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.TernaryExp)
	if !ok {
		t.Fatalf("stmt is not a TernaryExp, got %T", stmt.Expression)
	}

	if err := testToken(stmt, token.TRUE); err != nil {
		t.Fatal(err)
	}

	err = testPosition(exp.Position(), token.Position{
		StartCol: 3,
		EndCol:   28,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := testBooleanLiteral(exp.Condition, true); err != nil {
		t.Fatal(err)
	}

	if err := testIntegerLiteral(exp.IfBlock, 100); err != nil {
		t.Fatal(err)
	}

	if err := testStringLiteral(exp.ElseBlock, "Some string"); err != nil {
		t.Fatal(err)
	}
}

func TestParseIfStmt(t *testing.T) {
	t.Run("regular @if", func(t *testing.T) {
		inp := `@if(true)1@end`

		stmts, err := parseStatements(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[0].(*ast.IfStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
		}

		if err := testToken(stmt, token.IF); err != nil {
			t.Fatal(err)
		}

		err = testPosition(stmt.Position(), token.Position{
			StartCol: 0,
			EndCol:   13,
		})
		if err != nil {
			t.Fatal(err)
		}

		if err := testIfBlock(stmt, true, "1"); err != nil {
			t.Fatal(err)
		}

		if stmt.ElseBlock != nil {
			t.Fatalf("ifStmt.ElseBlock is not nil, got %T", stmt.ElseBlock)
		}

		if len(stmt.ElseifStmts) != 0 {
			t.Fatalf("ifStmt.ElseIfStmts is not empty, got %d", len(stmt.ElseifStmts))
		}
	})

	t.Run("@if with @else", func(t *testing.T) {
		inp := `@if(true)1@else2@end`

		stmts, err := parseStatements(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[0].(*ast.IfStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
		}

		if err := testToken(stmt, token.IF); err != nil {
			t.Fatal(err)
		}

		if err := testToken(stmt.IfBlock, token.HTML); err != nil {
			t.Fatal(err)
		}

		if err := testToken(stmt.ElseBlock, token.HTML); err != nil {
			t.Fatal(err)
		}

		if err := testIfBlock(stmt, true, "1"); err != nil {
			t.Fatal(err)
		}

		if err := testElseBlock(stmt.ElseBlock, "2"); err != nil {
			t.Fatal(err)
		}

		err = testPosition(stmt.Position(), token.Position{
			StartCol: 0,
			EndCol:   19,
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("nested @if with @else", func(t *testing.T) {
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

		stmts, err := parseStatements(inp, parseOpts{stmtCount: 2, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

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

		if err := testToken(ifStmt, token.IF); err != nil {
			t.Fatal(err)
		}

		if err := testToken(ifStmt.IfBlock, token.HTML); err != nil {
			t.Fatal(err)
		}

		if err := testToken(ifStmt.ElseBlock, token.HTML); err != nil {
			t.Fatal(err)
		}
	})
}

func TestParseIfElseIfStmt(t *testing.T) {
	inp := `@if(true)first@elseif(false)second@end`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	if err := testIfBlock(stmt, true, "first"); err != nil {
		t.Fatal(err)
	}

	if stmt.ElseBlock != nil {
		t.Fatalf("ifStmt.ElseBlock is not nil, got %T", stmt.ElseBlock)
	}

	if len(stmt.ElseifStmts) != 1 {
		t.Fatalf("ifStmt.ElseifStmts does not contain 1 statement, got %d", len(stmt.ElseifStmts))
	}

	elseifStmt := stmt.ElseifStmts[0]
	if elseifStmt, ok := elseifStmt.(*ast.ElseIfStmt); ok {
		if err := testBooleanLiteral(elseifStmt.Condition, false); err != nil {
			t.Fatal(err)
		}

		if len(elseifStmt.Block.Statements) != 1 {
			t.Fatalf(
				"elseifStmt.Block.Statements does not contain 1 statement, got %d",
				len(elseifStmt.Block.Statements),
			)
		}

		htmlStmt, ok := elseifStmt.Block.Statements[0].(*ast.HTMLStmt)
		if !ok {
			t.Fatalf(
				"elseifStmt.Block.Statements[0] is not an HTMLStmt, got %T",
				elseifStmt.Block.Statements[0],
			)
		}

		if htmlStmt.String() != "second" {
			t.Fatalf("htmlStmt.String() is not %q, got %q", "second", htmlStmt)
		}

		return
	}

	t.Fatalf("stmt.ElseifStmts[0] is not an ElseifStmt, got %T", elseifStmt)
}

func TestParseElseIfWithElseStatement(t *testing.T) {
	inp := `@if(true)1@elseif(false)2@else3@end`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	if err := testIfBlock(stmt, true, "1"); err != nil {
		t.Fatal(err)
	}

	if err := testElseBlock(stmt.ElseBlock, "3"); err != nil {
		t.Fatal(err)
	}

	if len(stmt.ElseifStmts) != 1 {
		t.Fatalf(
			"ifStmt.ElseifStmts does not contain 1 statement, got %d",
			len(stmt.ElseifStmts),
		)
	}

	if elseifStmt, ok := stmt.ElseifStmts[0].(*ast.ElseIfStmt); ok {
		if err := testBooleanLiteral(elseifStmt.Condition, false); err != nil {
			t.Fatal(err)
		}

		if len(elseifStmt.Block.Statements) != 1 {
			t.Fatalf(
				"elseifStmt.Block.Statements does not contain 1 statement, got %d",
				len(elseifStmt.Block.Statements),
			)
		}

		htmlStmt, ok := elseifStmt.Block.Statements[0].(*ast.HTMLStmt)
		if !ok {
			t.Fatalf(
				"elseifStmt.Block.Statements[0] is not an HTMLStmt, got %T",
				elseifStmt.Block.Statements[0],
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
		left     string
		varValue any
		str      string
		startCol uint
		endCol   uint
	}{
		{
			inp:      `{{ name = "Anna" }}`,
			left:     "name",
			varValue: "Anna",
			str:      `name = "Anna"`,
			startCol: 3,
			endCol:   15,
		},
		{
			inp:      `{{ myAge = 34 }}`,
			left:     "myAge",
			varValue: 34,
			str:      `myAge = 34`,
			startCol: 3,
			endCol:   12,
		},
	}

	for _, tc := range cases {
		stmts, err := parseStatements(tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[0].(*ast.AssignStmt)
		if !ok {
			t.Fatalf("stmts[0] is not a AssignStmt, got %T", stmts[0])
		}

		if stmt.Left.String() != tc.left {
			t.Fatalf("stmt.Left.String() is not %s, got %s", tc.left, stmt.Left.String())
		}

		if err := testLiteralExpression(stmt.Right, tc.varValue); err != nil {
			t.Fatal(err)
		}

		err = testPosition(stmt.Position(), token.Position{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})
		if err != nil {
			t.Fatal(err)
		}

		if stmt.String() != tc.str {
			t.Fatalf("stmt.String() is not %s, got %s", tc.inp, stmt)
		}
	}
}

func TestParseUseStmt(t *testing.T) {
	inp := `@use("main")`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.UseStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a UseStmt, got %T", stmts[0])
	}

	if stmt.Name.Value != "main" {
		t.Fatalf("stmt.Path.Value is not 'main', got %s", stmt.Name.Value)
	}

	if err := testToken(stmt, token.USE); err != nil {
		t.Fatal(err)
	}

	err = testPosition(stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   11,
	})
	if err != nil {
		t.Fatal(err)
	}

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

	stmts, err := parseStatements(inp, opts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[1].(*ast.ReserveStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ReserveStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.RESERVE); err != nil {
		t.Fatal(err)
	}

	err = testPosition(stmt.Position(), token.Position{
		StartCol: 5,
		EndCol:   23,
	})
	if err != nil {
		t.Fatal(err)
	}

	if stmt.Name.Value != "content" {
		t.Fatalf("stmt.Name.Value is not 'content', got %s", stmt.Name.Value)
	}

	if stmt.String() == inp {
		t.Fatalf("stmt.String() is not %s, got %s", inp, stmt)
	}
}

func TestInsertStmt(t *testing.T) {
	t.Run("@insert with block", func(t *testing.T) {
		inp := `<h1>@insert("content")<h1>Some content</h1>@end</h1>`

		stmts, err := parseStatements(inp, parseOpts{stmtCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[1].(*ast.InsertStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a InsertStmt, got %T", stmts[0])
		}

		if err := testToken(stmt, token.INSERT); err != nil {
			t.Fatal(err)
		}

		err = testPosition(stmt.Position(), token.Position{
			StartCol: 4,
			EndCol:   46,
		})
		if err != nil {
			t.Fatal(err)
		}

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

	t.Run("@insert with argument", func(t *testing.T) {
		inp := "<h1>@insert('content', 'Some content')</h1>"

		stmts, err := parseStatements(inp, parseOpts{stmtCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[1].(*ast.InsertStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a InsertStmt, got %T", stmts[0])
		}

		err = testPosition(stmt.Position(), token.Position{
			StartCol: 4,
			EndCol:   37,
		})
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(stmt, token.INSERT); err != nil {
			t.Fatal(err)
		}

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

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	arr, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not a ArrayLiteral, got %T", stmt.Expression)
	}

	if err := testToken(arr, token.LBRACKET); err != nil {
		t.Fatal(err)
	}

	err = testPosition(arr.Position(), token.Position{
		StartCol: 3,
		EndCol:   12,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(arr.Elements) != 2 {
		t.Fatalf("len(arr.Elements) is not 2, got %d", len(arr.Elements))
	}

	if err := testIntegerLiteral(arr.Elements[0], 11); err != nil {
		t.Fatal(err)
	}

	if err := testIntegerLiteral(arr.Elements[1], 234); err != nil {
		t.Fatal(err)
	}
}

func TestParseIndexExp(t *testing.T) {
	inp := `{{ arr[1 + 2][2] }}`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.IndexExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a IndexExp, got %T", stmt.Expression)
	}

	if err := testToken(exp, token.LBRACKET); err != nil {
		t.Fatal(err)
	}

	// testing the last index [2]
	err = testPosition(exp.Position(), token.Position{
		StartCol: 13,
		EndCol:   15,
	})
	if err != nil {
		t.Fatal(err)
	}

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
		stmts, err := parseStatements(tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
		}

		postfix, ok := stmt.Expression.(*ast.PostfixExp)
		if !ok {
			t.Fatalf("stmt.Expression is not a PostfixExp, got %T", stmt.Expression)
		}

		if err := testToken(postfix, tc.expTok); err != nil {
			t.Fatal(err)
		}
		if err := testIdentifier(postfix.Left, tc.ident); err != nil {
			t.Fatal(err)
		}

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

	stmts, err := parseStatements(inp, parseOpts{stmtCount: 2, checkErrors: true})
	if err != nil {
		t.Fatal(err)
	}

	if err := testIdentifier(stmts[0].(*ast.AssignStmt).Left, "name"); err != nil {
		t.Fatal(err)
	}

	if err := testToken(stmts[0], token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testIdentifier(stmts[1].(*ast.ExpressionStmt).Expression, "name"); err != nil {
		t.Fatal(err)
	}

	if err := testToken(stmts[0], token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testStringLiteral(stmts[0].(*ast.AssignStmt).Right, "Anna"); err != nil {
		t.Fatal(err)
	}

	if stmts[0].String() != `name = "Anna"` {
		t.Fatalf("stmts[0].String() is not '{{ name = \"Anna\" }}', got %s", stmts[0])
	}

	if stmts[1].String() != `name` {
		t.Fatalf("stmts[1].String() is not '{{ name }}', got %s", stmts[1])
	}
}

func TestParseTwoExpression(t *testing.T) {
	inp := `{{ 1; 2 }}`
	stmts, err := parseStatements(inp, parseOpts{stmtCount: 2, checkErrors: true})
	if err != nil {
		t.Fatal(err)
	}

	exp1 := stmts[0].(*ast.ExpressionStmt).Expression
	if err := testIntegerLiteral(exp1, 1); err != nil {
		t.Fatal(err)
	}

	exp2 := stmts[1].(*ast.ExpressionStmt).Expression
	if err := testIntegerLiteral(exp2, 2); err != nil {
		t.Fatal(err)
	}
}

func TestParseGlobalCallExp(t *testing.T) {
	inp := `{{ defined(var1, var2) }}`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.GlobalCallExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a GlobalCallExp, got %T", stmt.Expression)
	}

	if err := testToken(exp, token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testIdentifier(exp.Function, "defined"); err != nil {
		t.Fatal(err)
	}

	if err := testIdentifier(exp.Arguments[0], "var1"); err != nil {
		t.Fatal(err)
	}

	if err := testIdentifier(exp.Arguments[1], "var2"); err != nil {
		t.Fatal(err)
	}

	err = testPosition(exp.Position(), token.Position{
		StartCol: 3,
		EndCol:   21,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(exp.Arguments) != 2 {
		t.Fatalf("len(callExp.Arguments) is not 2, got %d", len(exp.Arguments))
	}
}

func TestParseCallExp(t *testing.T) {
	inp := `{{ "Serhii Cho".split(" ") }}`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a CallExp, got %T", stmt.Expression)
	}

	if err := testToken(exp, token.IDENT); err != nil {
		t.Fatal(err)
	}

	err = testPosition(exp.Position(), token.Position{
		StartCol: 16,
		EndCol:   25,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := testStringLiteral(exp.Receiver, "Serhii Cho"); err != nil {
		t.Fatal(err)
	}

	if err := testIdentifier(exp.Function, "split"); err != nil {
		t.Fatal(err)
	}

	if len(exp.Arguments) != 1 {
		t.Fatalf("len(callExp.Arguments) is not 1, got %d", len(exp.Arguments))
	}

	if err := testStringLiteral(exp.Arguments[0], " "); err != nil {
		t.Fatal(err)
	}
}

func TestParseCallExpWithExpressionList(t *testing.T) {
	inp := `{{ "nice".replace("n", "") }}`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a CallExp, got %T", stmt.Expression)
	}

	if err := testToken(exp, token.IDENT); err != nil {
		t.Fatal(err)
	}

	err = testPosition(exp.Position(), token.Position{
		StartCol: 10,
		EndCol:   25,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(exp.Arguments) != 2 {
		t.Fatalf("len(callExp.Arguments) is not 2, got %d", len(exp.Arguments))
	}
}

func TestParseCallExpWithEmptyString(t *testing.T) {
	inp := `{{ "".len() }}`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	callExp, ok := stmt.Expression.(*ast.CallExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a CallExp, got %T", stmt.Expression)
	}

	if err := testToken(callExp, token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testStringLiteral(callExp.Receiver, ""); err != nil {
		t.Fatal(err)
	}

	if err := testIdentifier(callExp.Function, "len"); err != nil {
		t.Fatal(err)
	}
}

func TestParseForStmt(t *testing.T) {
	t.Run("regular @for", func(t *testing.T) {
		inp := "@for(i = 0; i < 10; i++){{ i }}@end"

		stmts, err := parseStatements(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.ForStmt)
		if !ok {
			t.Fatalf("stmts[0] is not a ForStmt, got %T", stmts[0])
		}

		if err := testToken(stmt, token.FOR); err != nil {
			t.Fatal(err)
		}

		if err = testPosition(stmt.Pos, token.Position{EndCol: 34}); err != nil {
			t.Fatal(err)
		}

		err = testPosition(stmt.Block.Pos, token.Position{StartCol: 24, EndCol: 30})
		if err != nil {
			t.Fatal(err)
		}

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
	})

	t.Run("@for with @else block", func(t *testing.T) {
		inp := `@for(i = 0; i < 0; i++){{ i }}@elseEmpty@end`

		stmts, err := parseStatements(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.ForStmt)
		if !ok {
			t.Fatalf("stmts[0] is not a ForStmt, got %T", stmts[0])
		}

		if err := testToken(stmt, token.FOR); err != nil {
			t.Fatal(err)
		}

		if stmt.ElseBlock == nil {
			t.Fatalf("stmt.ElseBlock is nil")
		}

		if stmt.ElseBlock.String() != "Empty" {
			t.Fatalf("stmt.ElseBlock.String() is not 'Empty', got %s", stmt.ElseBlock)
		}
	})

	t.Run("infinite @for", func(t *testing.T) {
		inp := `@for(;;)1@end`

		stmts, err := parseStatements(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.ForStmt)
		if !ok {
			t.Fatalf("stmts[0] is not a ForStmt, got %T", stmts[0])
		}

		if err := testToken(stmt, token.FOR); err != nil {
			t.Fatal(err)
		}

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
	})
}

func TestParseEachStmt(t *testing.T) {
	t.Run("regular @each", func(t *testing.T) {
		inp := "@each(name in ['anna', 'serhii']){{ name }}@end"

		stmts, err := parseStatements(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.EachStmt)
		if !ok {
			t.Fatalf("stmts[0] is not a EachStmt, got %T", stmts[0])
		}

		err = testPosition(stmt.Pos, token.Position{EndCol: 46})
		if err != nil {
			t.Fatal(err)
		}
		if err := testToken(stmt, token.EACH); err != nil {
			t.Fatal(err)
		}

		err = testPosition(stmt.Block.Pos, token.Position{
			StartCol: 33,
			EndCol:   42,
		})
		if err != nil {
			t.Fatal(err)
		}

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
	})

	t.Run("@each with @else", func(t *testing.T) {
		inp := `@each(v in []){{ v }}@elseTest@end`

		stmts, err := parseStatements(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.EachStmt)
		if !ok {
			t.Fatalf("stmts[0] is not a EachStmt, got %T", stmts[0])
		}

		if err := testToken(stmt, token.EACH); err != nil {
			t.Fatal(err)
		}

		if stmt.ElseBlock.String() != "Test" {
			t.Fatalf("stmt.ElseBlock.String() is not 'Test', got %s", stmt.ElseBlock)
		}
	})
}

func TestParseEmptyBlock(t *testing.T) {
	cases := []struct {
		id        int
		inp       string
		endColPos uint
		tok       token.TokenType
	}{
		{10, "@if(x)@end", 9, token.IF},
		{11, "@if(x)    @end", 13, token.IF},
		{20, "@if(x)1@else@end", 15, token.IF},
		{30, "@if(x)@else@end", 14, token.IF},
		{40, "@if(x)@else1@end", 15, token.IF},
		{41, "@if(x) @else @end", 16, token.IF},
		{50, "@each(x in a)@end", 16, token.EACH},
		{51, "@each(x in a)  @end", 18, token.EACH},
		{60, "@each(x in a)@else@end", 21, token.EACH},
		{70, "@each(x in a)1@else@end", 22, token.EACH},
		{80, "@each(x in a)@else1@end", 22, token.EACH},
		{81, "@each(x in a) @else @end", 23, token.EACH},
		{90, "@for(i = 0; i < x; i++)@end", 26, token.FOR},
		{100, "@for(i = 0; i < x; i++)@else@end", 31, token.FOR},
		{110, "@for(i = 0; i < x; i++)1@else@end", 32, token.FOR},
		{120, "@for(i = 0; i < x; i++)@else1@end", 32, token.FOR},
		{121, "@for(i = 0; i < x; i++) @else @end", 33, token.FOR},
		{130, "@insert('x')@end", 15, token.INSERT},
		{141, "@component('x')@slot@end@end", 27, token.COMPONENT},
		{142, "@component('x') @slot@end @end", 29, token.COMPONENT},
		{143, "@component('x') @slot  @end @end", 31, token.COMPONENT},
		{150, "@component('x')@slot('x')@end@end", 32, token.COMPONENT},
		{160, "@component('x')@slotif(x, 'x')@end@end", 37, token.COMPONENT},
		{170, "@component('x')@slotif(x)@end@end", 32, token.COMPONENT},
		{180, "@component('x')@slotif(x)@end@slotif(x, 'x')@end@end", 51, token.COMPONENT},
		{190, "@component('x')@slot@end@slot('x')@end@end", 41, token.COMPONENT},
		{200, "@component('x')@slot@end@slot('x')@end@slotif(x)@end@end", 55, token.COMPONENT},
		{
			201,
			"@component('x') @slot @end @slot('x') @end @slotif(x) @end @end",
			62,
			token.COMPONENT,
		},
		{210, "@if(x)1@elseif(y)@end", 20, token.IF},
		{220, "@if(x)@elseif(y)@end", 19, token.IF},
		{230, "@if(x)@elseif(y)1@else@end", 25, token.IF},
		{240, "@if(x)@elseif(y)@elseif(z)@end", 29, token.IF},
		{250, "@if(x)@elseif(y)@elseif(z)1@else@end", 35, token.IF},
		{260, "@if(x) @elseif(y) @elseif(z) @else @end", 38, token.IF},
	}

	for _, tc := range cases {
		stmts, err := parseStatements(tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}

		stmt, ok := stmts[0].(ast.NodeWithStatements)
		if !ok {
			t.Fatalf("Case: %d. stmts[0] is not a NodeWithStatements, got %T", tc.id, stmts[0])
		}

		if err := testToken(stmt, tc.tok); err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}

		if err = testPosition(stmt.Position(), token.Position{EndCol: tc.endColPos}); err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}

		if len(stmt.Stmts()) != 0 {
			t.Fatalf("len(stmt.Stmts()) has to be empty, got %d", len(stmt.Stmts()))
		}
	}
}

func TestParseObjectStatement(t *testing.T) {
	inp := `{{ {"father": {name: "John"},} }}`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	obj, ok := stmt.Expression.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	err = testPosition(obj.Position(), token.Position{
		StartCol: 3,
		EndCol:   29,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(obj.Pairs) != 1 {
		t.Fatalf("len(obj.Pairs) is not 1, got %d", len(obj.Pairs))
	}

	if obj.String() != `{"father": {"name": "John"}}` {
		t.Fatalf(`obj.String() is not '{"father": {"name": "John" }}', got %s`, obj)
	}

	nested, ok := obj.Pairs["father"].(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("obj.Pairs['father'] is not a ObjectLiteral, got %T", obj.Pairs["father"])
	}

	if err := testStringLiteral(nested.Pairs["name"], "John"); err != nil {
		t.Fatal(err)
	}

	if err := testToken(obj, token.LBRACE); err != nil {
		t.Fatal(err)
	}
}

func TestParseObjectWithShorthandKeyNotation(t *testing.T) {
	inp := `{{ { name, age } }}`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	obj, ok := stmt.Expression.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	if err := testToken(obj, token.LBRACE); err != nil {
		t.Fatal(err)
	}

	err = testPosition(obj.Position(), token.Position{
		StartCol: 3,
		EndCol:   15,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(obj.Pairs) != 2 {
		t.Fatalf("len(obj.Pairs) is not 2, got %d", len(obj.Pairs))
	}
}

func TestParseHTMLStmt(t *testing.T) {
	inp := "<div><span>Hello</span></div>"

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.HTMLStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a HTMLStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.HTML); err != nil {
		t.Fatal(err)
	}

	err = testPosition(stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   28,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseDotExp(t *testing.T) {
	inp := `{{ person.father.name }}`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.DotExp)
	if !ok {
		t.Fatalf("stmt.Expression is not a DotExp, got %T", stmt.Expression)
	}

	if err := testToken(exp, token.DOT); err != nil {
		t.Fatal(err)
	}

	// position of the last dot between "father" and "name"
	err = testPosition(exp.Position(), token.Position{
		StartCol: 16,
		EndCol:   16,
	})
	if err != nil {
		t.Fatal(err)
	}

	if exp.String() != "((person.father).name)" {
		t.Fatalf("dotExp.String() is not '((person.father).name)', got %s", exp)
	}

	if err := testIdentifier(exp.Key, "name"); err != nil {
		t.Fatal(err)
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

	if err := testIdentifier(exp.Key, "father"); err != nil {
		t.Fatal(err)
	}

	if err := testIdentifier(exp.Left, "person"); err != nil {
		t.Fatal(err)
	}
}

func TestParseBreakStmt(t *testing.T) {
	inp := `@break`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.BreakStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a BreakStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.BREAK); err != nil {
		t.Fatal(err)
	}
}

func TestParseContinueStmt(t *testing.T) {
	inp := `@continue`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ContinueStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ContinueStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.CONTINUE); err != nil {
		t.Fatal(err)
	}
}

func TestParseBreakifStmt(t *testing.T) {
	inp := `@breakif(true)`

	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.BreakifStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a BreakIfStmt, got %T", stmts[0])
	}

	err = testPosition(stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   13,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(stmt, token.BREAKIF); err != nil {
		t.Fatal(err)
	}

	if err := testBooleanLiteral(stmt.Condition, true); err != nil {
		t.Fatal(err)
	}

	expect := "@breakif(true)"

	if stmt.String() != expect {
		t.Fatalf("breakStmt.String() is not '%s', got %s", expect, stmt)
	}
}

func TestParseContinueifStmt(t *testing.T) {
	inp := "@continueif(false)"
	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ContinueifStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ContinueIfStmt, got %T", stmts[0])
	}

	err = testPosition(stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   17,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(stmt, token.CONTINUEIF); err != nil {
		t.Fatal(err)
	}

	if err := testBooleanLiteral(stmt.Condition, false); err != nil {
		t.Fatal(err)
	}

	expect := "@continueif(false)"

	if stmt.String() != expect {
		t.Fatalf("stmt.String() is not '%s', got %s", expect, stmt)
	}
}

func TestParseComponentStmt(t *testing.T) {
	t.Run("@component without slots", func(t *testing.T) {
		inp := "<ul>@component('components/book-card', { c: card })</ul>"
		stmts, err := parseStatements(inp, parseOpts{stmtCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[1].(*ast.ComponentStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a ComponentStmt, got %T", stmts[1])
		}

		err = testPosition(stmt.Position(), token.Position{
			StartCol: 4,
			EndCol:   50,
		})

		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(stmt, token.COMPONENT); err != nil {
			t.Fatal(err)
		}

		if err := testStringLiteral(stmt.Name, "components/book-card"); err != nil {
			t.Fatal(err)
		}

		if len(stmt.Argument.Pairs) != 1 {
			t.Fatalf("len(stmt.Arguments) is not 1, got %d", len(stmt.Argument.Pairs))
		}

		if err := testIdentifier(stmt.Argument.Pairs["c"], "card"); err != nil {
			t.Fatal(err)
		}

		if len(stmt.Slots) != 0 {
			t.Fatalf("len(stmt.Slots) is not empty, got '%d' slots", len(stmt.Slots))
		}

		expect := `@component("components/book-card", {"c": card})`
		if stmt.String() != expect {
			t.Fatalf(`stmt.String() is not '%s', got %s`, expect, stmt)
		}
	})

	t.Run("@component with default slot", func(t *testing.T) {
		inp := `@component("components/book-card")@slot<h1>Header</h1>@end@end`

		stmts, err := parseStatements(inp, parseOpts{stmtCount: 1, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.ComponentStmt)
		if !ok {
			t.Fatalf("stmts[0] is not a ComponentStmt, got %T", stmts[1])
		}

		if len(stmt.Slots) != 1 {
			t.Fatalf("len(stmt.Slots) is not 1, got %d", len(stmt.Slots))
		}

		if err := testToken(stmt, token.COMPONENT); err != nil {
			t.Fatal(err)
		}

		name := stmt.Slots[0].Name().Value
		if name != "" {
			t.Fatalf("name must be empty string, got: %s", name)
		}

		expect := `@slot<h1>Header</h1>@end`
		if stmt.Slots[0].String() != expect {
			t.Fatalf("stmt.Slots[0].String() is not '%q', got %q", expect, stmt.Slots[0])
		}
	})

	t.Run("@component with 2 slots", func(t *testing.T) {
		inp := `<ul>
			@component("components/book-card", { c: card })
				@slot("header")<h1>Header</h1>@end
				@slot("footer")<footer>Footer</footer>@end
			@end
		</ul>`

		stmts, err := parseStatements(inp, parseOpts{stmtCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[1].(*ast.ComponentStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a ComponentStmt, got %T", stmts[1])
		}

		err = testPosition(stmt.Position(), token.Position{
			StartLine: 1,
			EndLine:   4,
			StartCol:  3, // because tabs before @component
			EndCol:    6, // because tabs before @end
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(stmt.Slots) != 2 {
			t.Fatalf("len(stmt.Slots) is not 2, got %d", len(stmt.Slots))
		}

		if err := testToken(stmt, token.COMPONENT); err != nil {
			t.Fatal(err)
		}
		if err := testStringLiteral(stmt.Slots[0].Name(), "header"); err != nil {
			t.Fatal(err)
		}
		if err := testStringLiteral(stmt.Slots[1].Name(), "footer"); err != nil {
			t.Fatal(err)
		}

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
		stmts, err := parseStatements(inp, parseOpts{stmtCount: 2, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.ComponentStmt)
		if !ok {
			t.Fatalf("stmts[0] is not a ComponentStmt, got %T", stmts[1])
		}

		if err := testToken(stmt, token.COMPONENT); err != nil {
			t.Fatal(err)
		}
		if err := testStringLiteral(stmt.Name, "some"); err != nil {
			t.Fatal(err)
		}

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

func TestParseSlotStmt(t *testing.T) {
	t.Run("named slot", func(t *testing.T) {
		inp := "<h2>@slot('header')</h2>"
		stmts, err := parseStatements(inp, parseOpts{stmtCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[1].(*ast.SlotStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a SlotStmt, got %T", stmts[1])
		}

		err = testPosition(stmt.Position(), token.Position{
			StartCol: 4,
			EndCol:   18,
		})
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(stmt, token.SLOT); err != nil {
			t.Fatal(err)
		}
		if err := testStringLiteral(stmt.Name(), "header"); err != nil {
			t.Fatal(err)
		}

		expect := `@slot("header")`
		if stmt.String() != expect {
			t.Fatalf("stmt.String() is not `%s`, got `%s`", expect, stmt)
		}
	})

	t.Run("default slot without end", func(t *testing.T) {
		t.Skip()
		inp := `<header>@slot</header>`
		stmts, err := parseStatements(inp, parseOpts{stmtCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[1].(*ast.SlotStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a SlotStmt, got %T", stmts[1])
		}

		if err := testToken(stmt, token.SLOT); err != nil {
			t.Fatal(err)
		}
		if err := testNilLiteral(stmt.Name()); err != nil {
			t.Fatal(err)
		}
		err = testPosition(stmt.Position(), token.Position{
			StartCol: 8,
			EndCol:   12,
		})
		if err != nil {
			t.Fatal(err)
		}

		if stmt.String() != "@slot" {
			t.Fatalf("slot.String() is not @slot, got `%s`", stmt)
		}
	})
}

func TestParseSlotifStmt(t *testing.T) {
	t.Run("default slotif", func(t *testing.T) {
		inp := `@component('test')@slotif(true)Test@end@end`
		stmts, err := parseStatements(inp, parseOpts{stmtCount: 1, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		comp, ok := stmts[0].(*ast.ComponentStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a ComponentStmt, got %T", stmts[0])
		}

		if len(comp.Slots) > 1 {
			t.Fatalf("len(comp.Slots) must be 1, got %d", len(comp.Slots))
		}

		err = testPosition(comp.Position(), token.Position{EndCol: 42})
		if err != nil {
			t.Fatal(err)
		}
		if err := testToken(comp, token.COMPONENT); err != nil {
			t.Fatal(err)
		}

		slot, ok := comp.Slots[0].(*ast.SlotifStmt)
		if !ok {
			t.Fatalf("comp.Slots[0] is not a SlotifStmt, got %T", stmts[0])
		}

		if err := testBooleanLiteral(slot.Condition, true); err != nil {
			t.Fatal(err)
		}

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
		stmts, err := parseStatements(inp, parseOpts{stmtCount: 1, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		comp, ok := stmts[0].(*ast.ComponentStmt)
		if !ok {
			t.Fatalf("stmts[1] is not a ComponentStmt, got %T", stmts[0])
		}

		err = testPosition(comp.Position(), token.Position{EndCol: 52})
		if err != nil {
			t.Fatal(err)
		}
		if err := testToken(comp, token.COMPONENT); err != nil {
			t.Fatal(err)
		}

		if len(comp.Slots) > 1 {
			t.Fatalf("len(comp.Slots) must be 1, got %d", len(comp.Slots))
		}

		slot, ok := comp.Slots[0].(*ast.SlotifStmt)
		if !ok {
			t.Fatalf("comp.Slots[0] is not a SlotifStmt, got %T", stmts[0])
		}

		if err := testBooleanLiteral(slot.Condition, false); err != nil {
			t.Fatal(err)
		}

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
	stmts, err := parseStatements(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.DumpStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an DumpStmt, got %T", stmts[0])
	}

	if len(stmt.Arguments) != 3 {
		t.Fatalf("len(stmt.Arguments) is not 3, got %d", len(stmt.Arguments))
	}

	err = testPosition(stmt.Position(), token.Position{
		StartCol: 0,
		EndCol:   26,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(stmt, token.DUMP); err != nil {
		t.Fatal(err)
	}

	if err := testStringLiteral(stmt.Arguments[0], "test"); err != nil {
		t.Fatal(err)
	}

	if err := testInfixExp(stmt.Arguments[1], 1, "+", 2); err != nil {
		t.Fatal(err)
	}

	if err := testBooleanLiteral(stmt.Arguments[2], false); err != nil {
		t.Fatal(err)
	}
}

func TestParseBlockAsIllegalNode(t *testing.T) {
	inp := "@if(false)@dump(@end"

	stmts, err := parseStatements(inp, parseOpts{stmtCount: 1, checkErrors: false})
	if err != nil {
		t.Fatal(err)
	}
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
		id        int
		inp       string
		stmtCount int
	}{
		{10, "@if(false", 1},
		{20, "@if  (loop. {{ 'nice' }}@end", 1},
		{30, "@if {{ 'nice' }}@end", 1},
		{40, "@if( {{ 'nice' }}@end", 1},
		{50, "@each( {{ 'nice' }}@end", 1},
		{60, "@each() {{ 'nice' }}@end", 1},
		{70, "@each (loop. {{ 'nice' }}@end", 1},
		{80, "@each(nice in []{{ 'nice' }}@end", 1},
		{90, "@each(nice in {{ 'nice' }}@end", 1},
		{100, "@for( {{ 'nice' }}@end", 1},
		{110, "@for() {{ 'nice' }}@end", 1},
		{120, "@for(i {{ 'nice' }}@end", 1},
		{130, "@for(i = 0; i < []; i++{{ 'nice' }}@end", 1},
		{140, "@for(i = 0; i < [] {{ 'nice' }}@end", 1},
		{150, "@component('~user'", 1},
		{160, "@component   ('", 1},
		{170, "@component", 1},
		{180, "@insert('nice", 1},
		{190, "@insert ('nice'", 1},
		{200, "@insert('nice'@end", 1},
		{210, "@insert    ('nice' {{ 'nice' }}@end", 1},
		{220, `@if(loop.
            {{ loop.first }}
            Iteration number is {{ loop.iter }}
        @end`, 1},
	}

	for _, tc := range cases {
		stmts, err := parseStatements(tc.inp, parseOpts{
			stmtCount:   tc.stmtCount,
			checkErrors: false,
		})
		if err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}

		_, ok := stmts[0].(*ast.IllegalNode)
		if !ok {
			t.Fatalf(
				"Case: %d. stmts[0] is not an IllegalNode, got %T for %s",
				tc.id,
				stmts[0],
				tc.inp,
			)
		}
	}
}
