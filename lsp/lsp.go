package lsp

import (
	"fmt"

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

		loopStmt := stmt.(*ast.EachStmt)

		if loopStmt.Block.Pos.StartLine > line || loopStmt.Block.Pos.EndLine < line {
			fmt.Printf("continue in first")
			continue
		}

		fmt.Printf("------->StartCol %d\n", loopStmt.Block.Pos.StartCol)
		fmt.Printf("------->EndCol %d\n", loopStmt.Block.Pos.EndCol)
		if loopStmt.Block.Pos.StartCol >= char || loopStmt.Block.Pos.EndCol <= char {
			fmt.Printf("continue in second")
			continue
		}

		return true
	}

	return false
}
