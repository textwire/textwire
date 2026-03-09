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
	chunksCount int
	inserts     map[string]*ast.InsertDir
	checkErrors bool
}

var defaultParseOpts = parseOpts{
	chunksCount: 1,
	inserts:     nil,
	checkErrors: true,
}

func parseChunks(inp string, opts parseOpts) ([]ast.Chunk, error) {
	l := lexer.New(inp)
	p := New(l, nil)
	prog := p.ParseProgram()

	if opts.checkErrors && p.HasErrors() {
		return nil, p.Errors()[0].Error()
	}

	if len(prog.Chunks) != opts.chunksCount {
		return nil, fmt.Errorf(
			"Program must have %d chunks but got %d for input %q",
			opts.chunksCount,
			len(prog.Chunks),
			inp,
		)
	}

	return prog.Chunks, nil
}

func testInfixExpr(expr ast.Expression, left any, op string, right any) error {
	infix, ok := expr.(*ast.InfixExpr)
	if !ok {
		return fmt.Errorf("Variable expr is not an InfixExpr, got %T", expr)
	}

	if err := testLiteralExpr(infix.Left, left); err != nil {
		return err
	}

	if infix.Op != op {
		return fmt.Errorf("infix.Op is not %s, got %s", op, infix.Op)
	}

	if err := testLiteralExpr(infix.Right, right); err != nil {
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

func testIntExpr(expr ast.Expression, val int64) error {
	integer, ok := expr.(*ast.IntExpr)
	if !ok {
		return fmt.Errorf("expr is not an IntExpr, got %T", expr)
	}

	if integer.Val != val {
		return fmt.Errorf("integer.Val is not %d, got %d", val, integer.Val)
	}

	if integer.Tok().Lit != strconv.FormatInt(val, 10) {
		return fmt.Errorf("integer.Tok().Lit is not %d, got %s", val, integer.Tok().Lit)
	}

	return nil
}

func testFloatExpr(expr ast.Expression, val float64) error {
	float, ok := expr.(*ast.FloatExpr)
	if !ok {
		return fmt.Errorf("expr is not a FloatExpr, got %T", expr)
	}

	if float.Val != val {
		return fmt.Errorf("float.Val is not %f, got %f", val, float.Val)
	}

	if float.String() != utils.FloatToStr(val) {
		return fmt.Errorf("float.String() is not %f, got %s", val, float)
	}

	return nil
}

func testNilExpr(expr ast.Expression) error {
	nilExpr, ok := expr.(*ast.NilExpr)
	if !ok {
		return fmt.Errorf("expr is not a NilExpr, got %T", expr)
	}

	if nilExpr.Tok().Lit != "nil" {
		return fmt.Errorf("nilExpr.Tok().Lit is not 'nil', got %s", nilExpr.Tok().Lit)
	}

	return nil
}

func testStrExpr(expr ast.Expression, val string) error {
	str, ok := expr.(*ast.StrExpr)
	if !ok {
		return fmt.Errorf("expr is not a StrExpr, got %T", expr)
	}

	if str.Val != val {
		return fmt.Errorf("str.Val is not %s, got %s", val, str.Val)
	}

	if str.Tok().Lit != val {
		return fmt.Errorf("str.Tok().Lit is not %s, got %s", val, str.Tok().Lit)
	}

	return nil
}

func testBoolExpr(expr ast.Expression, val bool) error {
	boolean, ok := expr.(*ast.BoolExpr)
	if !ok {
		return fmt.Errorf("expr not *ast.BoolExpr, got %T", expr)
	}

	if boolean.Val != val {
		return fmt.Errorf("boolean.Val not %t, got %t", val, boolean.Val)
	}

	if boolean.Tok().Lit != fmt.Sprintf("%t", val) {
		return fmt.Errorf("boolean.Tok().Lit is not %t, got %s", val, boolean.Tok().Lit)
	}

	return nil
}

func testIdentExpr(expr ast.Expression, val string) error {
	ident, ok := expr.(*ast.IdentExpr)
	if !ok {
		return fmt.Errorf("expr is not an IdentExpr, got %T", expr)
	}

	if ident.Name != val {
		return fmt.Errorf("ident.Name is not %s, got %s", val, ident.Name)
	}

	if ident.Tok().Lit != val {
		return fmt.Errorf("ident.Tok().Lit is not %s, got %s", val, ident.Tok().Lit)
	}

	return nil
}

func testLiteralExpr(expr ast.Expression, expect any) error {
	switch v := expect.(type) {
	case int:
		return testIntExpr(expr, int64(v))
	case int64:
		return testIntExpr(expr, v)
	case float64:
		return testFloatExpr(expr, v)
	case string:
		return testStrExpr(expr, v)
	case bool:
		return testBoolExpr(expr, v)
	case nil:
		return testNilExpr(expr)
	default:
		return fmt.Errorf("type of expr not handled. got %T", expr)
	}
}

func testIfDir(dir ast.Chunk, cond any, ifBlock string) error {
	ifDir, ok := dir.(*ast.IfDir)
	if !ok {
		return fmt.Errorf("dir is not an IfDir, got %T", dir)
	}

	if err := testLiteralExpr(ifDir.Cond, cond); err != nil {
		return err
	}

	if ifDir.IfBlock.String() != ifBlock {
		return fmt.Errorf("ifDir.IfBlock.String() is not %q, got %q", ifBlock, ifDir.IfBlock)
	}

	return nil
}

func testBlock(block *ast.Block, val string) error {
	if block == nil {
		return fmt.Errorf("block is nil")
	}

	if len(block.Chunks) != 1 {
		return fmt.Errorf(
			"block.Chunks must contain 1 chunk, got %d",
			len(block.Chunks),
		)
	}

	if block.String() != val {
		return fmt.Errorf("block.String() is not %q, got %q", block, val)
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
	chunks, err := parseChunks("{{ myName }}", defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	chunk, ok := chunks[0].(*ast.Embedded)
	if !ok {
		t.Fatalf("chunks[0] is not an Embedded, got %T", chunks[0])
	}

	if len(chunk.Elements) != 1 {
		t.Fatalf("chunk.Statements must contain 1 statement, got %d", len(chunk.Elements))
	}

	stmt, ok := chunk.Elements[0].(*ast.IdentExpr)
	if !ok {
		t.Fatalf("chunk.Statements[0] is not an IdentExpr, got %T", chunk.Elements[0])
	}

	if err := testToken(ident, token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testToken(ident.Expression, token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(ident.Expression, "myName"); err != nil {
		t.Fatal(err)
	}
}

func TestParseExpressionStatement(t *testing.T) {
	stmts, err := parseChunks("{{ 3 / 2 }}", defaultParseOpts)
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

func TestParseIntLit(t *testing.T) {
	stmts, err := parseChunks("{{ 234 }}", defaultParseOpts)
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

	if err := testIntExpr(stmt.Expression, 234); err != nil {
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

func TestParseFloatLit(t *testing.T) {
	stmts, err := parseChunks("{{ 2.34149 }}", defaultParseOpts)
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

	if err := testFloatExpr(stmt.Expression, 2.34149); err != nil {
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

func TestParseNilLit(t *testing.T) {
	stmts, err := parseChunks("{{ nil }}", defaultParseOpts)
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

	if err := testNilLit(stmt.Expression); err != nil {
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

func TestParseStrLit(t *testing.T) {
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
		stmts, err := parseChunks(tc.inp, defaultParseOpts)
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
		if err := testStrLit(stmt.Expression, tc.expect); err != nil {
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

func TestStrConcatenation(t *testing.T) {
	inp := `{{ "Serhii" + " Anna" }}`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExpr)

	if !ok {
		t.Fatalf("stmt is not an InfixExp, got %T", stmt.Expression)
	}

	if err := testToken(stmt, token.STR); err != nil {
		t.Fatal(err)
	}

	if exp.Left.Tok().Lit != "Serhii" {
		t.Fatalf("exp.Left is not %s, got %s", "Serhii", exp.Left)
	}

	if exp.Op != "+" {
		t.Fatalf("exp.Op is not %s, got %s", "+", exp.Op)
	}

	if exp.Right.Tok().Lit != " Anna" {
		t.Fatalf("exp.Right is not %s, got %s", " Anna", exp.Right)
	}
}
func TestExpression(t *testing.T) {
	test := "{{ 5 + 2 }}"

	stmts, err := parseChunks(test, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExpr)
	if !ok {
		t.Fatalf("stmt is not an InfixExp, got %T", stmt.Expression)
	}

	if err := testToken(stmt, token.INT); err != nil {
		t.Fatal(err)
	}

	if err := testIntExpr(exp.Right, 2); err != nil {
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

	if err := testIntExpr(exp.Left, 5); err != nil {
		t.Fatal(err)
	}
}

func TestGroupedExpression(t *testing.T) {
	test := "{{ (5 + 5) * 2 }}"

	stmts, err := parseChunks(test, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExpr)
	if !ok {
		t.Fatalf("stmt is not an InfixExp, got %T", stmt.Expression)
	}

	if err := testToken(stmt, token.LPAREN); err != nil {
		t.Fatal(err)
	}

	if err := testIntExpr(exp.Right, 2); err != nil {
		t.Fatal(err)
	}

	if exp.Op != "*" {
		t.Fatalf("exp.Op is not %s, got %s", "*", exp.Op)
	}

	infix, ok := exp.Left.(*ast.InfixExpr)
	if !ok {
		t.Fatalf("exp.Left is not an InfixExp, got %T", exp.Left)
	}

	if err := testIntExpr(infix.Left, 5); err != nil {
		t.Fatal(err)
	}

	if infix.Op != "+" {
		t.Fatalf("infix.Op is not %s, got %s", "+", infix.Op)
	}

	if err := testLiteralExpr(infix.Right, 5); err != nil {
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
		stmts, err := parseChunks(tc.inp, defaultParseOpts)
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
		stmts, err := parseChunks(tc.inp, defaultParseOpts)
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
		if err := testBoolLit(stmt.Expression, tc.expectBoolean); err != nil {
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
		val    any
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
		stmts, err := parseChunks(tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpr)
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

		if err := testLiteralExpr(exp.Right, tc.val); err != nil {
			t.Fatal(err)
		}
	}
}

func TestOpPrecedenceParsing(t *testing.T) {
	cases := []struct {
		id     uint
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
			err: fail.New(1, "", "parser", fail.ErrObjKeyUseGet),
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

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)

	if !ok {
		t.Fatalf("stmts[0] is not an ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.TernaryExpr)
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

	if err := testBoolLit(exp.Cond, true); err != nil {
		t.Fatal(err)
	}

	if err := testIntExpr(exp.IfExpr, 100); err != nil {
		t.Fatal(err)
	}

	if err := testStrLit(exp.ElseExpr, "Some string"); err != nil {
		t.Fatal(err)
	}
}

func TestParseIfStmt(t *testing.T) {
	t.Run("regular @if", func(t *testing.T) {
		inp := `@if(true)1@end`

		stmts, err := parseChunks(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[0].(*ast.IfDir)
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

		if err := testIfDir(stmt, true, "1"); err != nil {
			t.Fatal(err)
		}

		if stmt.ElseBlock != nil {
			t.Fatalf("ifStmt.ElseBlock is not nil, got %T", stmt.ElseBlock)
		}

		if len(stmt.ElseifDirs) != 0 {
			t.Fatalf("ifStmt.ElseIfStmts is not empty, got %d", len(stmt.ElseifDirs))
		}
	})

	t.Run("@if with @else", func(t *testing.T) {
		inp := `@if(true)1@else2@end`

		stmts, err := parseChunks(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[0].(*ast.IfDir)
		if !ok {
			t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
		}

		if err := testToken(stmt, token.IF); err != nil {
			t.Fatal(err)
		}

		if err := testToken(stmt.IfBlock, token.TEXT); err != nil {
			t.Fatal(err)
		}

		if err := testToken(stmt.ElseBlock, token.TEXT); err != nil {
			t.Fatal(err)
		}

		if err := testIfDir(stmt, true, "1"); err != nil {
			t.Fatal(err)
		}

		if err := testBlock(stmt.ElseBlock, "2"); err != nil {
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

		stmts, err := parseChunks(inp, parseOpts{chunksCount: 2, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := stmts[0].(*ast.Text); !ok {
			t.Fatalf("stmts[0] is not an TextStmt, got %T", stmts[0])
		}

		ifStmt, isNotIfStmt := stmts[1].(*ast.IfDir)
		if !isNotIfStmt {
			t.Fatalf("stmts[1] is not an IfStmt, got %T", stmts[1])
		}

		if len(ifStmt.IfBlock.Chunks) != 3 {
			t.Fatalf(
				"ifStmt.IfBlock.Statements does not contain 3 statement, got %d",
				len(ifStmt.IfBlock.Chunks),
			)
		}

		if err := testToken(ifStmt, token.IF); err != nil {
			t.Fatal(err)
		}

		if err := testToken(ifStmt.IfBlock, token.TEXT); err != nil {
			t.Fatal(err)
		}

		if err := testToken(ifStmt.ElseBlock, token.TEXT); err != nil {
			t.Fatal(err)
		}
	})
}

func TestParseIfElseIfStmt(t *testing.T) {
	inp := `@if(true)first@elseif(false)second@end`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.IfDir)
	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	if err := testIfDir(stmt, true, "first"); err != nil {
		t.Fatal(err)
	}

	if stmt.ElseBlock != nil {
		t.Fatalf("ifStmt.ElseBlock is not nil, got %T", stmt.ElseBlock)
	}

	if len(stmt.ElseifDirs) != 1 {
		t.Fatalf("ifStmt.ElseifStmts does not contain 1 statement, got %d", len(stmt.ElseifDirs))
	}

	elseifStmt := stmt.ElseifDirs[0]
	if elseifStmt, ok := elseifStmt.(*ast.ElseIfDir); ok {
		if err := testBoolLit(elseifStmt.Condition, false); err != nil {
			t.Fatal(err)
		}

		if len(elseifStmt.Block.Chunks) != 1 {
			t.Fatalf(
				"elseifStmt.Block.Statements does not contain 1 statement, got %d",
				len(elseifStmt.Block.Chunks),
			)
		}

		textStmt, ok := elseifStmt.Block.Chunks[0].(*ast.TextStmt)
		if !ok {
			t.Fatalf(
				"elseifStmt.Block.Statements[0] is not an TextStmt, got %T",
				elseifStmt.Block.Chunks[0],
			)
		}

		if textStmt.String() != "second" {
			t.Fatalf("textStmt.String() is not %q, got %q", "second", textStmt)
		}

		return
	}

	t.Fatalf("stmt.ElseifStmts[0] is not an ElseifStmt, got %T", elseifStmt)
}

func TestParseElseIfWithElseStatement(t *testing.T) {
	inp := `@if(true)1@elseif(false)2@else3@end`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.IfDir)
	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmts[0])
	}

	if err := testIfDir(stmt, true, "1"); err != nil {
		t.Fatal(err)
	}

	if err := testBlock(stmt.ElseBlock, "3"); err != nil {
		t.Fatal(err)
	}

	if len(stmt.ElseifDirs) != 1 {
		t.Fatalf(
			"ifStmt.ElseifStmts does not contain 1 statement, got %d",
			len(stmt.ElseifDirs),
		)
	}

	if elseifStmt, ok := stmt.ElseifDirs[0].(*ast.ElseIfDir); ok {
		if err := testBoolLit(elseifStmt.Condition, false); err != nil {
			t.Fatal(err)
		}

		if len(elseifStmt.Block.Chunks) != 1 {
			t.Fatalf(
				"elseifStmt.Block.Statements does not contain 1 statement, got %d",
				len(elseifStmt.Block.Chunks),
			)
		}

		textStmt, ok := elseifStmt.Block.Chunks[0].(*ast.TextStmt)
		if !ok {
			t.Fatalf(
				"elseifStmt.Block.Statements[0] is not an TextStmt, got %T",
				elseifStmt.Block.Chunks[0],
			)
		}

		if textStmt.String() != "2" {
			t.Fatalf("TextStmt.String() is not %s, got %s", "2", textStmt)
		}
	}
}

func TestParseAssignStmt(t *testing.T) {
	cases := []struct {
		id       uint
		inp      string
		str      string
		startCol uint
		endCol   uint
	}{
		{
			id:       10,
			inp:      `{{ name = "Anna" }}`,
			str:      `name = "Anna"`,
			startCol: 3,
			endCol:   15,
		},
		{
			id:       20,
			inp:      `{{ myAge = 34 }}`,
			str:      `myAge = 34`,
			startCol: 3,
			endCol:   12,
		},
		{
			id:       30,
			inp:      `{{ me.age = 34 }}`,
			str:      `(me.age) = 34`,
			startCol: 3,
			endCol:   13,
		},
		{
			id:       40,
			inp:      `{{ arr[0] = 1 }}`,
			str:      `(arr[0]) = 1`,
			startCol: 3,
			endCol:   12,
		},
		{
			id:       50,
			inp:      `{{ arr[234][2][23].name.first = "Anna" }}`,
			str:      `(((((arr[234])[2])[23]).name).first) = "Anna"`,
			startCol: 3,
			endCol:   37,
		},
		{
			id:       60,
			inp:      `{{ (obj.one.two) = "test" }}`,
			str:      `((obj.one).two) = "test"`,
			startCol: 4,
			endCol:   24,
		},
	}

	for _, tc := range cases {
		stmts, err := parseChunks(tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.AssignStmt)
		if !ok {
			t.Fatalf("Case: %d. stmts[0] is not a AssignStmt, got %T", tc.id, stmts[0])
		}

		stmtStr := stmt.String()
		if stmtStr != tc.str {
			t.Fatalf("Case: %d. stmt.String() is not %s, got %s", tc.id, tc.inp, stmtStr)
		}

		err = testPosition(stmt.Position(), token.Position{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})

		if err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}
	}
}

func TestParseUseStmt(t *testing.T) {
	inp := `@use("main")`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.UseDir)
	if !ok {
		t.Fatalf("stmts[0] is not a UseStmt, got %T", stmts[0])
	}

	if stmt.Name.Val != "main" {
		t.Fatalf("stmt.Path.Val is not 'main', got %s", stmt.Name.Val)
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
		chunksCount: 3,
		checkErrors: true,
		inserts: map[string]*ast.InsertDir{
			"content": {
				Name: &ast.StrExpr{Val: "content"},
				Block: &ast.Block{
					Chunks: []ast.Statement{
						&ast.Text{
							BaseNode: ast.BaseNode{
								Token: token.Token{
									Type: token.TEXT,
									Lit:  "<h1>Some content</h1>",
								},
							},
						},
					},
				},
			},
		},
	}

	stmts, err := parseChunks(inp, opts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[1].(*ast.ReserveDir)
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

	if stmt.Name.Val != "content" {
		t.Fatalf("stmt.Name.Val is not 'content', got %s", stmt.Name.Val)
	}

	if stmt.String() == inp {
		t.Fatalf("stmt.String() is not %s, got %s", inp, stmt)
	}
}

func TestInsertStmt(t *testing.T) {
	t.Run("@insert with block", func(t *testing.T) {
		inp := `<h1>@insert("content")<h1>Some content</h1>@end</h1>`

		stmts, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[1].(*ast.InsertDir)
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

		if stmt.Name.Val != "content" {
			t.Fatalf("stmt.Name.Val is not 'content', got %s", stmt.Name.Val)
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

		stmts, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[1].(*ast.InsertDir)
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

		if stmt.Name.Val != "content" {
			t.Fatalf("stmt.Name.Val is not 'content', got %s", stmt.Name.Val)
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

func TestParseArr(t *testing.T) {
	inp := `{{ [11, 234,] }}`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	arr, ok := stmt.Expression.(*ast.ArrExpr)
	if !ok {
		t.Fatalf("stmt.Expression is not a ArrLit, got %T", stmt.Expression)
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

	if err := testIntExpr(arr.Elements[0], 11); err != nil {
		t.Fatal(err)
	}

	if err := testIntExpr(arr.Elements[1], 234); err != nil {
		t.Fatal(err)
	}
}

func TestParseIndexExp(t *testing.T) {
	inp := `{{ arr[1 + 2][2] }}`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("stmt.Expression is not a IndexExp, got %T", stmt.Expression)
	}

	if err := testToken(exp, token.IDENT); err != nil {
		t.Fatal(err)
	}

	err = testPosition(exp.Position(), token.Position{
		StartCol: 3,
		EndCol:   15,
	})

	if err != nil {
		t.Fatal(err)
	}

	if exp.String() != "((arr[(1 + 2)])[2])" {
		t.Fatalf("indexExp.String() is not '(arr[(1 + 2)])', got %s", exp)
	}
}

func TestParseIncStmt(t *testing.T) {
	inp := "{{ i++ }}"

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.IncStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a IncStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.INC); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(stmt.Left, "i"); err != nil {
		t.Fatal(err)
	}

	if stmt.String() != "(i++)" {
		t.Fatalf("stmt.String() is not '(i++)', got %s", stmt)
	}
}

func TestParseDecStmt(t *testing.T) {
	inp := "{{ i-- }}"

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.DecStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a DecStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.DEC); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(stmt.Left, "i"); err != nil {
		t.Fatal(err)
	}

	if stmt.String() != "(i--)" {
		t.Fatalf("stmt.String() is not '(i--)', got %s", stmt)
	}
}

func TestParseTwoStatements(t *testing.T) {
	inp := `{{ name = "Anna"; name }}`

	stmts, err := parseChunks(inp, parseOpts{chunksCount: 2, checkErrors: true})
	if err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(stmts[0].(*ast.AssignStmt).Left, "name"); err != nil {
		t.Fatal(err)
	}

	if err := testToken(stmts[0], token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(stmts[1].(*ast.ExpressionStmt).Expression, "name"); err != nil {
		t.Fatal(err)
	}

	if err := testToken(stmts[0], token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testStrLit(stmts[0].(*ast.AssignStmt).Right, "Anna"); err != nil {
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
	stmts, err := parseChunks(inp, parseOpts{chunksCount: 2, checkErrors: true})
	if err != nil {
		t.Fatal(err)
	}

	exp1 := stmts[0].(*ast.ExpressionStmt).Expression
	if err := testIntExpr(exp1, 1); err != nil {
		t.Fatal(err)
	}

	exp2 := stmts[1].(*ast.ExpressionStmt).Expression
	if err := testIntExpr(exp2, 2); err != nil {
		t.Fatal(err)
	}
}

func TestParseGlobalCallExp(t *testing.T) {
	inp := `{{ defined(var1, var2) }}`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.GlobalCallExpr)
	if !ok {
		t.Fatalf("stmt.Expression is not a GlobalCallExp, got %T", stmt.Expression)
	}

	if err := testToken(exp, token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(exp.Function, "defined"); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(exp.Arguments[0], "var1"); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(exp.Arguments[1], "var2"); err != nil {
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

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpr)
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

	if err := testStrLit(exp.Receiver, "Serhii Cho"); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(exp.Function, "split"); err != nil {
		t.Fatal(err)
	}

	if len(exp.Arguments) != 1 {
		t.Fatalf("len(callExp.Arguments) is not 1, got %d", len(exp.Arguments))
	}

	if err := testStrLit(exp.Arguments[0], " "); err != nil {
		t.Fatal(err)
	}
}

func TestParseCallExpWithExpressionList(t *testing.T) {
	inp := `{{ "nice".replace("n", "") }}`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpr)
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

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	callExp, ok := stmt.Expression.(*ast.CallExpr)
	if !ok {
		t.Fatalf("stmt.Expression is not a CallExp, got %T", stmt.Expression)
	}

	if err := testToken(callExp, token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testStrLit(callExp.Receiver, ""); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(callExp.Function, "len"); err != nil {
		t.Fatal(err)
	}
}

func TestParseForStmt(t *testing.T) {
	t.Run("regular @for", func(t *testing.T) {
		inp := "@for(i = 0; i < 10; i++){{ i }}@end"

		stmts, err := parseChunks(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.ForDir)
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

		stmts, err := parseChunks(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.ForDir)
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

		stmts, err := parseChunks(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.ForDir)
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

		stmts, err := parseChunks(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.EachDir)
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

		if stmt.Arr.String() != `["anna", "serhii"]` {
			t.Fatalf(`stmt.Arr.String() is not '["anna", "serhii"]', got %s`, stmt.Arr)
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

		stmts, err := parseChunks(inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.EachDir)
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
		id        uint
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
		stmts, err := parseChunks(tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}

		stmt, ok := stmts[0].(ast.NodeWithChunks)
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

func TestParseObjStmt(t *testing.T) {
	inp := `{{ {"father": {name: "John"},} }}`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	obj, ok := stmt.Expression.(*ast.ObjExpr)
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

	nested, ok := obj.Pairs["father"].(*ast.ObjExpr)
	if !ok {
		t.Fatalf("obj.Pairs['father'] is not a ObjLit, got %T", obj.Pairs["father"])
	}

	if err := testStrLit(nested.Pairs["name"], "John"); err != nil {
		t.Fatal(err)
	}

	if err := testToken(obj, token.LBRACE); err != nil {
		t.Fatal(err)
	}
}

func TestParseObjWithShorthandKeyNotation(t *testing.T) {
	inp := `{{ { name, age } }}`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	obj, ok := stmt.Expression.(*ast.ObjExpr)
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

func TestParseTextStmt(t *testing.T) {
	inp := "<div><span>Hello</span></div>"

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.Text)
	if !ok {
		t.Fatalf("stmts[0] is not a TextStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.TEXT); err != nil {
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
	inp := "{{ person.father.name }}"

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmts[0] is not a ExpressionStmt, got %T", stmts[0])
	}

	exp, ok := stmt.Expression.(*ast.DotExpr)
	if !ok {
		t.Fatalf("stmt.Expression is not a DotExp, got %T", stmt.Expression)
	}

	if err := testToken(exp, token.IDENT); err != nil {
		t.Fatal(err)
	}

	// position of the last dot between "father" and "name"
	err = testPosition(exp.Position(), token.Position{
		StartCol: 3,
		EndCol:   20,
	})
	if err != nil {
		t.Fatal(err)
	}

	if exp.String() != "((person.father).name)" {
		t.Fatalf("dotExp.String() is not '((person.father).name)', got %s", exp)
	}

	if err := testIdentExpr(exp.Key, "name"); err != nil {
		t.Fatal(err)
	}

	if exp.Left == nil {
		t.Fatalf("dotExp.Left is nil")
	}

	exp, ok = exp.Left.(*ast.DotExpr)
	if exp == nil {
		t.Fatalf("dotExp is nil")
		return
	}

	if !ok {
		t.Fatalf("dotExp.Left is not a DotExp, got %T", exp.Left)
	}

	if err := testIdentExpr(exp.Key, "father"); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(exp.Left, "person"); err != nil {
		t.Fatal(err)
	}
}

func TestParseBreakStmt(t *testing.T) {
	inp := `@break`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.BreakDir)
	if !ok {
		t.Fatalf("stmts[0] is not a BreakStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.BREAK); err != nil {
		t.Fatal(err)
	}
}

func TestParseContinueStmt(t *testing.T) {
	inp := `@continue`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ContinueDir)
	if !ok {
		t.Fatalf("stmts[0] is not a ContinueStmt, got %T", stmts[0])
	}

	if err := testToken(stmt, token.CONTINUE); err != nil {
		t.Fatal(err)
	}
}

func TestParseBreakifStmt(t *testing.T) {
	inp := `@breakif(true)`

	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.BreakifDir)
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

	if err := testBoolLit(stmt.Condition, true); err != nil {
		t.Fatal(err)
	}

	expect := "@breakif(true)"

	if stmt.String() != expect {
		t.Fatalf("breakStmt.String() is not '%s', got %s", expect, stmt)
	}
}

func TestParseContinueifStmt(t *testing.T) {
	inp := "@continueif(false)"
	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	stmt, ok := stmts[0].(*ast.ContinueifDir)
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

	if err := testBoolLit(stmt.Condition, false); err != nil {
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
		stmts, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[1].(*ast.ComponentDir)
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

		if err := testStrLit(stmt.Name, "components/book-card"); err != nil {
			t.Fatal(err)
		}

		if len(stmt.Argument.Pairs) != 1 {
			t.Fatalf("len(stmt.Arguments) is not 1, got %d", len(stmt.Argument.Pairs))
		}

		if err := testIdentExpr(stmt.Argument.Pairs["c"], "card"); err != nil {
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

		stmts, err := parseChunks(inp, parseOpts{chunksCount: 1, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.ComponentDir)
		if !ok {
			t.Fatalf("stmts[0] is not a ComponentStmt, got %T", stmts[1])
		}

		if len(stmt.Slots) != 1 {
			t.Fatalf("len(stmt.Slots) is not 1, got %d", len(stmt.Slots))
		}

		if err := testToken(stmt, token.COMPONENT); err != nil {
			t.Fatal(err)
		}

		name := stmt.Slots[0].Name().Val
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

		stmts, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := stmts[1].(*ast.ComponentDir)
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
		if err := testStrLit(stmt.Slots[0].Name(), "header"); err != nil {
			t.Fatal(err)
		}
		if err := testStrLit(stmt.Slots[1].Name(), "footer"); err != nil {
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
		stmts, err := parseChunks(inp, parseOpts{chunksCount: 2, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[0].(*ast.ComponentDir)
		if !ok {
			t.Fatalf("stmts[0] is not a ComponentStmt, got %T", stmts[1])
		}

		if err := testToken(stmt, token.COMPONENT); err != nil {
			t.Fatal(err)
		}
		if err := testStrLit(stmt.Name, "some"); err != nil {
			t.Fatal(err)
		}

		expect := `@component("some")`
		if stmt.String() != expect {
			t.Fatalf("stmt.String() is not `%s`, got `%s`", expect, stmt)
		}

		textStmt, ok := stmts[1].(*ast.Text)
		if !ok {
			t.Fatalf("stmts[1] is not a TextStmt, got %T", stmts[1])
		}

		expect = "\n <b>Book</b>"
		if textStmt.String() != expect {
			t.Fatalf("textStmt.String() is not `%s`, got `%s`", expect, textStmt)
		}
	})
}

func TestParseSlotStmt(t *testing.T) {
	t.Run("named slot", func(t *testing.T) {
		inp := "<h2>@slot('header')</h2>"
		stmts, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[1].(*ast.SlotDir)
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
		if err := testStrLit(stmt.Name(), "header"); err != nil {
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
		stmts, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		stmt, ok := stmts[1].(*ast.SlotDir)
		if !ok {
			t.Fatalf("stmts[1] is not a SlotStmt, got %T", stmts[1])
		}

		if err := testToken(stmt, token.SLOT); err != nil {
			t.Fatal(err)
		}
		if err := testNilLit(stmt.Name()); err != nil {
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
		stmts, err := parseChunks(inp, parseOpts{chunksCount: 1, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		comp, ok := stmts[0].(*ast.ComponentDir)
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

		slot, ok := comp.Slots[0].(*ast.SlotifDir)
		if !ok {
			t.Fatalf("comp.Slots[0] is not a SlotifStmt, got %T", stmts[0])
		}

		if err := testBoolLit(slot.Condition, true); err != nil {
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
		stmts, err := parseChunks(inp, parseOpts{chunksCount: 1, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		comp, ok := stmts[0].(*ast.ComponentDir)
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

		slot, ok := comp.Slots[0].(*ast.SlotifDir)
		if !ok {
			t.Fatalf("comp.Slots[0] is not a SlotifStmt, got %T", stmts[0])
		}

		if err := testBoolLit(slot.Condition, false); err != nil {
			t.Fatal(err)
		}

		if slot.Name().Val != "name" {
			t.Fatalf("slot.Name().Val is not 'name', got %s", slot.Name())
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
	stmts, err := parseChunks(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.DumpDir)
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

	if err := testStrLit(stmt.Arguments[0], "test"); err != nil {
		t.Fatal(err)
	}

	if err := testInfixExp(stmt.Arguments[1], 1, "+", 2); err != nil {
		t.Fatal(err)
	}

	if err := testBoolLit(stmt.Arguments[2], false); err != nil {
		t.Fatal(err)
	}
}

func TestParseBlockAsIllegalNode(t *testing.T) {
	inp := "@if(false)@dump(@end"

	stmts, err := parseChunks(inp, parseOpts{chunksCount: 1, checkErrors: false})
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := stmts[0].(*ast.IfDir)
	if !ok {
		t.Fatalf("stmts[0] is not an IfStmt, got %T", stmt)
	}

	dumpStmt, ok := stmt.IfBlock.Chunks[0].(*ast.DumpDir)
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
		stmts, err := parseChunks(tc.inp, parseOpts{
			chunksCount: tc.stmtCount,
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
