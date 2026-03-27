package parser

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/ast"
)

// trimTextChunks removes Text nodes with whitespace only characters and trims
// all the other Text nodes.
func trimTextChunks(block *ast.Block) *ast.Block {
	trimmedBlock := ast.NewBlock(*block.Tok())

	for i := range block.Chunks {
		text, ok := block.Chunks[i].(*ast.Text)
		if !ok {
			trimmedBlock.Chunks = append(trimmedBlock.Chunks, block.Chunks[i])
			continue
		}

		content := strings.Trim(text.Token.Lit, " \n\t\r")
		if content == "" {
			continue
		}

		text.Token.Lit = content
		trimmedBlock.Chunks = append(trimmedBlock.Chunks, text)
	}

	return trimmedBlock
}
