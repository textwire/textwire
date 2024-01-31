package parser

import (
	"strconv"

	"github.com/textwire/textwire/token"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/fail"
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
	POSTFIX      // X++ or X--
)

var precedences = map[token.TokenType]int{
	token.QUESTION: TERNARY,
	token.EQ:       EQ,
	token.NOT_EQ:   EQ,
	token.LTHAN:    LESS_GREATER,
	token.GTHAN:    LESS_GREATER,
	token.LTHAN_EQ: LESS_GREATER,
	token.GTHAN_EQ: LESS_GREATER,
	token.PERIOD:   SUM,
	token.ADD:      SUM,
	token.SUB:      SUM,
	token.DIV:      PRODUCT,
	token.MOD:      PRODUCT,
	token.MUL:      PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
	token.INC:      POSTFIX,
	token.DEC:      POSTFIX,
}

type Parser struct {
	l      *lexer.Lexer
	errors []*fail.Error

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		l:      lexer,
		errors: []*fail.Error{},
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
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	// Infix operators
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.ADD, p.parseInfixExpression)
	p.registerInfix(token.SUB, p.parseInfixExpression)
	p.registerInfix(token.MUL, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)

	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LTHAN, p.parseInfixExpression)
	p.registerInfix(token.GTHAN, p.parseInfixExpression)
	p.registerInfix(token.LTHAN_EQ, p.parseInfixExpression)
	p.registerInfix(token.GTHAN_EQ, p.parseInfixExpression)

	p.registerInfix(token.QUESTION, p.parseTernaryExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.INC, p.parsePostfixExpression)
	p.registerInfix(token.DEC, p.parsePostfixExpression)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{Token: p.curToken}
	prog.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if p.curTokenIs(token.ILLEGAL) {
			p.newError(
				p.curToken.Line,
				fail.ErrIllegalToken,
				p.curToken.Literal,
			)
			return nil
		}

		if stmt == nil {
			p.nextToken()
			continue
		}

		prog.Statements = append(prog.Statements, stmt)

		p.nextToken() // skip "}}"
	}

	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.HTML:
		return p.parseHTMLStatement()
	case token.LBRACES:
		return p.parseEmbeddedCode()
	case token.SEMI:
		return p.parseEmbeddedCode()
	case token.IF:
		return p.parseIfStatement()
	case token.USE:
		return p.parseUseStatement()
	case token.RESERVE:
		return p.parseReserveStatement()
	case token.INSERT:
		return p.parseInsertStatement()
	default:
		return nil
	}
}

func (p *Parser) parseEmbeddedCode() ast.Statement {
	p.nextToken() // skip "{{" or ";"

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.Line, fail.ErrEmptyBrackets)
		return nil
	}

	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
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

func (p *Parser) peekTokenIs(tokens ...token.TokenType) bool {
	for _, tok := range tokens {
		if p.peekToken.Type == tok {
			return true
		}
	}

	return false
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

	p.newError(
		p.peekToken.Line,
		fail.ErrWrongNextToken,
		token.String(tok),
		token.String(p.peekToken.Type),
	)

	return false
}

func (p *Parser) newError(line uint, msg string, args ...interface{}) {
	newErr := fail.New(line, "parser", msg, args...)
	p.errors = append(p.errors, newErr)
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
		p.newError(
			p.curToken.Line,
			fail.ErrCouldNotParseAs,
			p.curToken.Literal,
			"INT",
		)
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
		p.newError(
			p.curToken.Line,
			fail.ErrCouldNotParseAs,
			p.curToken.Literal,
			"FLOAT",
		)
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
		Token: p.curToken, // "true" or "false"
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.curToken} // "["
	arr.Elements = p.parseExpressionList(token.RBRACKET)
	return arr
}

func (p *Parser) parseHTMLStatement() *ast.HTMLStatement {
	return &ast.HTMLStatement{Token: p.curToken}
}

func (p *Parser) parseDefineStatement() ast.Statement {
	ident := &ast.Identifier{
		Token: p.curToken, // identifier
		Value: p.curToken.Literal,
	}

	stmt := &ast.DefineStatement{
		Token: p.curToken, // identifier
		Name:  ident,
	}

	if !p.expectPeek(token.DEFINE) { // move to ":="
		return nil
	}

	p.nextToken() // skip ":="

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.Line, fail.ErrExpectedExpression)
		return nil
	}

	stmt.Value = p.parseExpression(SUM)

	return stmt
}

func (p *Parser) parseUseStatement() ast.Statement {
	stmt := &ast.UseStatement{
		Token: p.curToken, // "@use"
	}

	if !p.expectPeek(token.LPAREN) { // move to "("
		return nil
	}

	p.nextToken() // skip "("

	stmt.Name = &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return stmt
}

func (p *Parser) parseReserveStatement() ast.Statement {
	stmt := &ast.ReserveStatement{
		Token: p.curToken, // "@reserve"
	}

	if !p.expectPeek(token.LPAREN) { // move to "("
		return nil
	}

	p.nextToken() // skip "("

	stmt.Name = &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return stmt
}

func (p *Parser) parseInsertStatement() ast.Statement {
	stmt := &ast.InsertStatement{
		Token: p.curToken, // "@insert"
	}

	if !p.expectPeek(token.LPAREN) { // move to "("
		return nil
	}

	p.nextToken() // skip "("

	stmt.Name = &ast.StringLiteral{
		Token: p.curToken, // The name of the insert statement
		Value: p.curToken.Literal,
	}

	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip insert name
		p.nextToken() // skip ","
		stmt.Argument = p.parseExpression(LOWEST)
		return stmt
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return nil
	}

	p.nextToken() // skip ")"

	stmt.Block = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.curToken, // "["
		Left:  left,
	}

	p.nextToken() // skip "["

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) { // move to "]"
		return nil
	}

	return exp
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{
		Token:    p.curToken,         // identifier
		Operator: p.curToken.Literal, // "++" or "--"
		Left:     left,
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken, // operator
		Operator: p.curToken.Literal,
		Left:     left,
	}

	p.nextToken() // skip operator

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.Line, fail.ErrExpectedExpression)
		return nil
	}

	exp.Right = p.parseExpression(SUM)

	return exp
}

func (p *Parser) parseTernaryExpression(left ast.Expression) ast.Expression {
	exp := &ast.TernaryExpression{
		Token:     p.curToken, // "?"
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
	stmt := &ast.DefineStatement{
		Token: p.curToken, // "var"
	}

	if !p.expectPeek(token.IDENT) { // move to identifier
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken, // identifier
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
	stmt := &ast.IfStatement{
		Token: p.curToken, // "@if"
	}

	if !p.expectPeek(token.LPAREN) { // move to "("
		return nil
	}

	p.nextToken() // skip "("

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return nil
	}

	p.nextToken() // skip ")"

	stmt.Consequence = p.parseBlockStatement()

	for p.peekTokenIs(token.ELSEIF) {
		if !p.expectPeek(token.ELSEIF) { // move to "@elseif"
			return nil
		}

		p.nextToken() // skip "@elseif"
		p.nextToken() // skip "("

		condition := p.parseExpression(LOWEST)

		if !p.expectPeek(token.RPAREN) { // move to ")"
			return nil
		}

		p.nextToken() // skip ")"

		stmt.Alternatives = append(stmt.Alternatives, &ast.ElseIfStatement{
			Token:       p.curToken,
			Condition:   condition,
			Consequence: p.parseBlockStatement(),
		})
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // move to "@else"
		p.nextToken() // skip "@else"

		stmt.Alternative = p.parseBlockStatement()

		if p.peekTokenIs(token.ELSEIF) {
			p.newError(p.peekToken.Line, fail.ErrElseifCannotFollowElse)
			return nil
		}
	}

	if !p.expectPeek(token.END) { // move to "@end"
		return nil
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	stmt := &ast.BlockStatement{Token: p.curToken}

	for {
		block := p.parseStatement()

		if block != nil {
			stmt.Statements = append(stmt.Statements, block)
		}

		if p.peekTokenIs(token.ELSE, token.ELSEIF, token.END) {
			break
		}

		p.nextToken() // skip statement
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	exp := p.parseExpression(LOWEST)

	result := &ast.ExpressionStatement{Token: p.curToken, Expression: exp}

	if p.peekTokenIs(token.RBRACES) {
		p.nextToken() // skip "}}"
	}

	return result
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.newError(
			p.curToken.Line,
			fail.ErrNoPrefixParseFunc,
			token.String(p.curToken.Type),
		)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.RBRACES, token.SEMI, token.RPAREN) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken, // prefix operator
		Operator: p.curToken.Literal,
	}

	p.nextToken() // skip prefix operator

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

func (p *Parser) parseExpressionList(endTok token.TokenType) []ast.Expression {
	var result []ast.Expression

	if p.peekTokenIs(endTok) {
		p.nextToken() // skip endTok token
		return result
	}

	p.nextToken() // skip ","

	result = append(result, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip ","
		p.nextToken() // skip expression
		result = append(result, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(endTok) { // move to endTok
		return nil
	}

	return result
}
