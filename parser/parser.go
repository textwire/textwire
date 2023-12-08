package parser

import (
	"github.com/textwire/textwire/token"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/lexer"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	TERNARY      // a ? b : c
	EQ           // ==
	LESS_GREATER // > or <
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X or !X
	CALL         // myFunction(X)
	INDEX        // array[index]
)

var precedences = map[token.TokenType]int{
	token.QUESTION: TERNARY,
	token.EQ:       EQ,
	token.NOT_EQ:   EQ,
	token.LTHAN:    LESS_GREATER,
	token.GTHAN:    LESS_GREATER,
	token.PERIOD:   SUM,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.MODULO:   PRODUCT,
	token.ASTERISK: PRODUCT,
}

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		l:      lexer,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	// Prefix operators
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{}
	prog.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}

		p.nextToken()
	}

	return prog
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.HTML:
		return p.parseHTMLStatement()
	case token.LBRACES:
		return p.parseEmbeddedCode()
	default:
		return nil
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) curTokenIs(tok token.TokenType) bool {
	return p.curToken.Type == tok
}

func (p *Parser) peekTokenIs(tok token.TokenType) bool {
	return p.peekToken.Type == tok
}

func (p *Parser) peekPrecedence() int {
	result, ok := precedences[p.peekToken.Type]

	if !ok {
		return LOWEST
	}

	return result
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseHTMLStatement() *ast.HTMLStatement {
	return &ast.HTMLStatement{Token: p.curToken}
}

func (p *Parser) parseEmbeddedCode() ast.Statement {
	p.nextToken()

	return p.parseExpressionStatement()
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	expr := p.parseExpression(LOWEST)
	result := &ast.ExpressionStatement{Token: p.curToken, Expression: expr}

	if p.peekTokenIs(token.RBRACES) {
		p.nextToken()
	}

	return result
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		// TODO: Handle error
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.RBRACES) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}
