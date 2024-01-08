package parser

import (
	"errors"
	"fmt"
	"strconv"

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
	errors []error

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		l:      lexer,
		errors: []error{},
	}

	p.nextToken() // fill curToken
	p.nextToken() // fill peekToken

	// Prefix operators
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STR, p.parseStringLiteral)
	p.registerPrefix(token.NIL, p.parseNilLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	// Infix operators
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)

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

		p.nextToken() // skip "}}"
	}

	return prog
}

func (p *Parser) Errors() []error {
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

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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

func (p *Parser) expectPeek(tok token.TokenType) bool {
	if p.peekTokenIs(tok) {
		p.nextToken()
		return true
	}

	msg := fmt.Sprintf("expected next token to be %d, got %d instead", tok, p.peekToken.Type)

	p.errors = append(p.errors, errors.New(msg))

	return false
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

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)

	if err != nil {
		p.errors = append(p.errors, errors.New("could not parse "+p.curToken.Literal+" as integer"))
		return nil
	}

	return &ast.IntegerLiteral{
		Token: p.curToken,
		Value: val,
	}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseNilLiteral() ast.Expression {
	return &ast.NilLiteral{Token: p.curToken}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parseHTMLStatement() *ast.HTMLStatement {
	return &ast.HTMLStatement{Token: p.curToken}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	p.nextToken() // skip operator

	exp.Right = p.parseExpression(SUM)

	return exp
}

func (p *Parser) parseEmbeddedCode() ast.Statement {
	p.nextToken() // skip "{{"
	return p.parseExpressionStatement()
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	expr := p.parseExpression(LOWEST)
	result := &ast.ExpressionStatement{Token: p.curToken, Expression: expr}

	if p.peekTokenIs(token.RBRACES) {
		p.nextToken() // skip "}}"
	}

	return result
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.errors = append(p.errors, errors.New("no prefix parse function for "+p.curToken.Literal))
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.RBRACES) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		p.nextToken() // skip operator

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken() // skip operator

	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // skip "("

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return nil
	}

	return exp
}
