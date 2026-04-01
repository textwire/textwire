package parser

import (
	"fmt"
	"strconv"

	"github.com/textwire/textwire/v4/pkg/ast"
	"github.com/textwire/textwire/v4/pkg/lexer"
	"github.com/textwire/textwire/v4/pkg/position"
	"github.com/textwire/textwire/v4/pkg/token"
	"github.com/textwire/textwire/v4/pkg/utils"
)

type parseOpts struct {
	chunksCount int
	checkErrors bool
}

var defaultParseOpts = parseOpts{
	chunksCount: 1,
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
			"program must have %d chunks but got %d for input %q",
			opts.chunksCount,
			len(prog.Chunks),
			inp,
		)
	}

	return prog.Chunks, nil
}

func parseEmbedded[T ast.Segment](inp string, opts parseOpts) (T, error) {
	var zero T

	chunks, err := parseChunks(inp, opts)
	if err != nil {
		return zero, err
	}

	embedded, ok := chunks[0].(*ast.Embedded)
	if !ok {
		return zero, fmt.Errorf("chunks[0] is not an Embedded, got %T", chunks[0])
	}

	if len(embedded.Segments) != 1 {
		return zero, fmt.Errorf(
			"embedded.Segments must contain 1 segment, got %d",
			len(embedded.Segments),
		)
	}

	segment, ok := embedded.Segments[0].(T)
	if !ok {
		return zero, fmt.Errorf(
			"embedded.Segments[0] is not %T, got %T",
			zero,
			embedded.Segments[0],
		)
	}

	return segment, nil
}

func parseEmbeddedSegments(inp string, opts parseOpts) ([]ast.Segment, error) {
	chunks, err := parseChunks(inp, opts)
	if err != nil {
		return nil, err
	}

	embedded, ok := chunks[0].(*ast.Embedded)
	if !ok {
		return nil, fmt.Errorf("chunks[0] is not Embedded, got %T", chunks[0])
	}

	return embedded.Segments, nil
}

func parseDirective[T ast.Chunk](inp string, opts parseOpts) (T, error) {
	var zero T

	chunks, err := parseChunks(inp, opts)
	if err != nil {
		return zero, err
	}

	dir, ok := chunks[0].(T)
	if !ok {
		return zero, fmt.Errorf("chunks[0] is not an %T, got %T", zero, chunks[0])
	}

	return dir, nil
}

func testInfixExpr(expr ast.Expression, left any, op string, right any) error {
	infixExpr, ok := expr.(*ast.InfixExpr)
	if !ok {
		return fmt.Errorf("expr is not an InfixExpr, got %T", expr)
	}

	if err := testLiteralExpr(infixExpr.Left, left); err != nil {
		return err
	}

	if infixExpr.Op != op {
		return fmt.Errorf("infixExpr.Op is not %s, got %s", op, infixExpr.Op)
	}

	if err := testLiteralExpr(infixExpr.Right, right); err != nil {
		return err
	}

	return nil
}

func testTokPosition(actual, expect *position.Pos) error {
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
		return fmt.Errorf("str.Val is not %q, got %q", val, str.Val)
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
			"token type is not %q, got %q",
			token.String(expect),
			token.String(tok.Tok().Type),
		)
	}
	return nil
}
