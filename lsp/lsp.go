package lsp

import (
	"github.com/textwire/textwire/v2/token"

	"github.com/textwire/textwire/v2/ast"
	"github.com/textwire/textwire/v2/lexer"
	"github.com/textwire/textwire/v2/parser"
)

// IsInLoop checks if given position of the cursor is inside of a loop
func IsInLoop(doc, filePath string, line, char uint) bool {
	l := lexer.New(doc)
	p := parser.New(l, filePath)
	program := p.ParseProgram()

	if program == nil {
		return false
	}

	for _, stmt := range program.Statements {
		isEachLoop := stmt.Tok().Type == token.EACH
		isForLoop := stmt.Tok().Type == token.FOR

		if !isEachLoop && !isForLoop {
			continue
		}

		loopStmt := stmt.(ast.LoopStmt)
		pos := loopStmt.LoopBodyBlock().Pos

		if pos.StartLine > line || pos.EndLine < line {
			continue
		}

		if pos.StartCol > char || pos.EndCol < char {
			continue
		}

		return true
	}

	return false
}
