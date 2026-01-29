package lsp

import (
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/token"

	"github.com/textwire/textwire/v2/ast"
	"github.com/textwire/textwire/v2/lexer"
	"github.com/textwire/textwire/v2/parser"
)

// IsInLoop checks if given position of the cursor is inside of a loop
func IsInLoop(doc, filePath string, line, col uint) (bool, []*fail.Error) {
	l := lexer.New(doc)
	p := parser.New(l, filePath)
	program := p.ParseProgram()

	if program == nil {
		return false, p.Errors()
	}

	for _, stmt := range program.Stmts() {
		isEachLoop := stmt.Tok().Type == token.EACH
		isForLoop := stmt.Tok().Type == token.FOR

		if !isEachLoop && !isForLoop {
			continue
		}

		loopStmt := stmt.(ast.LoopStmt)
		pos := loopStmt.LoopBodyBlock().Pos

		if IsCursorInBody(line, col, pos) {
			return true, p.Errors()
		}
	}

	return false, p.Errors()
}

func IsCursorInBody(line, col uint, pos token.Position) bool {
	// Line outside range
	if line < pos.StartLine || line > pos.EndLine {
		return false
	}

	// For inlined loops that are written in a single line
	if line == pos.StartLine && line == pos.EndLine {
		return col >= pos.StartCol && col <= pos.EndCol
	}

	// When cursor is on the start line
	if line == pos.StartLine {
		return col >= pos.StartCol
	}

	// When cursor is on the end line
	if line == pos.EndLine {
		return col <= pos.EndCol
	}

	// In middle lines any column is valid
	return true
}
