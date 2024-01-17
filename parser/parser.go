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
	ERR_ILLEGAL_TOKEN        = "illegal token '%s' found"
)

var precedences = map[token.TokenType]int{
	token.QUESTION: TERNARY,
	token.EQ:       EQ,
	token.NOT_EQ:   EQ,
	token.LTHAN:    LESS_GREATER,
	token.GTHAN:    LESS_GREATER,
	token.PERIOD:   SUM,
	token.ADD:      SUM,
	token.SUB:      SUM,
	token.DIV:      PRODUCT,
	token.MOD:      PRODUCT,
	token.MUL:      PRODUCT,
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
	p.registerPrefix(token.SUB, p.parsePrefixExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	// Infix operators
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.QUESTION, p.parseTernaryExpression)
	p.registerInfix(token.ADD, p.parseInfixExpression)
	p.registerInfix(token.SUB, p.parseInfixExpression)
	p.registerInfix(token.MUL, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{}
	prog.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if p.curTokenIs(token.ILLEGAL) {
			p.newError(ERR_ILLEGAL_TOKEN, p.curToken.Literal)
			return nil
		}

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

func (p *Parser) parseEmbeddedCode() ast.Statement {
	p.nextToken() // skip "{{"

	if p.curTokenIs(token.RBRACES) {
		p.newError(ERR_EMPTY_BRACKETS)
		return nil
	}

	switch p.curToken.Type {
	case token.IF:
		return p.parseIfStatement()
	case token.VAR:
		return p.parseVarStatement()
	case token.LAYOUT:
		return p.parseLayoutStatement()
	case token.RESERVE:
		return p.parseReserveStatement()
	case token.INSERT:
		return p.parseInsertStatement()
	case token.IDENT:
		if p.peekTokenIs(token.DEFINE) {
			return p.parseDefineStatement()
		}
	}

	return p.parseExpressionStatement()
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
		token.TokenString(tok),
		token.TokenString(p.peekToken.Type),
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

func (p *Parser) parseDefineStatement() ast.Statement {
	ident := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	stmt := &ast.DefineStatement{
		Token: p.curToken,
		Name:  ident,
	}

	if !p.expectPeek(token.DEFINE) { // move to ":="
		return nil
	}

	p.nextToken() // skip ":="

	if p.curTokenIs(token.RBRACES) {
		p.newError(ERR_EXPECTED_EXPRESSION)
		return nil
	}

	stmt.Value = p.parseExpression(SUM)

	return stmt
}

func (p *Parser) parseLayoutStatement() ast.Statement {
	stmt := &ast.LayoutStatement{Token: p.curToken}

	if !p.expectPeek(token.STR) { // move to string
		return nil
	}

	stmt.Path = &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return stmt
}

func (p *Parser) parseReserveStatement() ast.Statement {
	stmt := &ast.ReserveStatement{Token: p.curToken}

	if !p.expectPeek(token.STR) { // move to string
		return nil
	}

	stmt.Name = &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return stmt
}

func (p *Parser) parseInsertStatement() ast.Statement {
	stmt := &ast.InsertStatement{Token: p.curToken}

	if !p.expectPeek(token.STR) { // move to string
		return nil
	}

	stmt.Name = &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return stmt
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

func (p *Parser) parseVarStatement() ast.Statement {
	stmt := &ast.DefineStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) { // move to identifier
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) { // move to "="
		return nil
	}

	p.nextToken() // skip "="

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
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
		isOpening := p.curTokenIs(token.LBRACES)    // "{{"
		isPeekEnd := p.peekTokenIs(token.END)       // "end
		isPeekElse := p.peekTokenIs(token.ELSE)     // "else"
		isPeekElseIf := p.peekTokenIs(token.ELSEIF) // "else if"

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
		p.newError(ERR_NO_PREFIX_PARSE_FUNC, token.TokenString(p.curToken.Type))
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
