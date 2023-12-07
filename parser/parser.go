package parser

import (
	"go/token"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/lexer"
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	// todo: here
}

func (p *Parser) Errors() []string {
	return p.errors
}
