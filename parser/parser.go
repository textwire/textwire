package parser

import (
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

	// Error messages
	ERR_EMPTY_BRACKETS       = "bracket statement must contain an expression '{{ <expression> }}'"
	ERR_WRONG_NEXT_TOKEN     = "expected next token to be %s, got %s instead"
	ERR_EXPECTED_EXPRESSION  = "expected expression, got '}}'"
	ERR_COULD_NOT_PARSE_AS   = "could not parse %s as %s"
	ERR_NO_PREFIX_PARSE_FUNC = "no prefix parse function for %s"
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
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STR, p.parseStringLiteral)
	p.registerPrefix(token.NIL, p.parseNilLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	// Infix operators
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.QUESTION, p.parseTernaryExpression)
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

	msg := fmt.Sprintf(
		ERR_WRONG_NEXT_TOKEN,
		token.TypeName(tok),
		token.TypeName(p.peekToken.Type),
	)

	p.newError(msg)

	return false
}

func (p *Parser) newError(msg string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Errorf(msg, args...))
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
		p.newError(ERR_COULD_NOT_PARSE_AS, p.curToken.Literal, "INT")
		return nil
	}

	return &ast.IntegerLiteral{
		Token: p.curToken,
		Value: val,
	}
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	val, err := strconv.ParseFloat(p.curToken.Literal, 10)

	if err != nil {
		p.newError(ERR_COULD_NOT_PARSE_AS, p.curToken.Literal, "FLOAT")
		return nil
	}

	return &ast.FloatLiteral{
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

	if p.curTokenIs(token.RBRACES) {
		p.newError(ERR_EXPECTED_EXPRESSION)
		return nil
	}

	exp.Right = p.parseExpression(SUM)

	return exp
}

func (p *Parser) parseTernaryExpression(left ast.Expression) ast.Expression {
	exp := &ast.TernaryExpression{
		Token:     p.curToken,
		Condition: left,
	}

	p.nextToken() // skip "?"

	exp.Consequence = p.parseExpression(TERNARY)

	if !p.expectPeek(token.COLON) { // move to ":"
		return nil
	}

	p.nextToken() // skip ":"

	exp.Alternative = p.parseExpression(LOWEST)

	return exp
}

func (p *Parser) parseEmbeddedCode() ast.Statement {
	p.nextToken() // skip "{{"

	if p.curTokenIs(token.RBRACES) {
		p.newError(ERR_EMPTY_BRACKETS)
		return nil
	}

	switch p.curToken.Type {
	case token.IF:
		return p.parseIfStatement()
	}

	return p.parseExpressionStatement()
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}

	p.nextToken() // skip "if" or "else if"

	stmt.Condition = p.parseExpression(LOWEST)

	p.nextToken() // skip "}}"

	stmt.Consequence = p.parseBlockStatement()

	for p.peekTokenIs(token.ELSEIF) {
		p.nextToken() // skip "{{"
		p.nextToken() // skip "else if"

		stmt.Alternatives = append(stmt.Alternatives, &ast.ElseIfStatement{
			Token:       p.curToken,
			Condition:   p.parseExpression(LOWEST),
			Consequence: p.parseBlockStatement(),
		})
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // skip "{{"
		p.nextToken() // skip "else"
		p.nextToken() // skip "}}"

		stmt.Alternative = p.parseBlockStatement()

		if p.peekTokenIs(token.ELSEIF) {
			p.newError("ELSEIF statement cannot follow ELSE statement")
			return nil
		}
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	stmt := &ast.BlockStatement{Token: p.curToken}

	for {
		isOpening := p.curTokenIs(token.LBRACES)
		isPeekEnd := p.peekTokenIs(token.END)
		isPeekElse := p.peekTokenIs(token.ELSE)
		isPeekElseIf := p.peekTokenIs(token.ELSEIF)

		if isOpening && (isPeekEnd || isPeekElse || isPeekElseIf) {
			break
		}

		block := p.parseStatement()

		if block != nil {
			stmt.Statements = append(stmt.Statements, block)
		}

		p.nextToken() // skip "}}"
	}

	return stmt
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
		p.newError(ERR_NO_PREFIX_PARSE_FUNC, token.TypeName(p.curToken.Type))
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
