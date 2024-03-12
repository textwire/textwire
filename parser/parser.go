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
	INDEX        // array[index] or obj.key
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
	token.ADD:      SUM,
	token.SUB:      SUM,
	token.DIV:      PRODUCT,
	token.MOD:      PRODUCT,
	token.MUL:      PRODUCT,
	token.LPAREN:   CALL,
	token.DOT:      INDEX,
	token.LBRACKET: INDEX,
	token.INC:      POSTFIX,
	token.DEC:      POSTFIX,
}

type Parser struct {
	l        *lexer.Lexer
	errors   []*fail.Error
	filepath string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	useStmt    *ast.UseStmt
	components []*ast.ComponentStmt
	inserts    map[string]*ast.InsertStmt
	reserves   map[string]*ast.ReserveStmt
}

func New(lexer *lexer.Lexer, filepath string) *Parser {
	p := &Parser{
		l:          lexer,
		filepath:   filepath,
		errors:     []*fail.Error{},
		components: []*ast.ComponentStmt{},
		inserts:    map[string]*ast.InsertStmt{},
		reserves:   map[string]*ast.ReserveStmt{},
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
	p.registerPrefix(token.SUB, p.parsePrefixExp)
	p.registerPrefix(token.NOT, p.parsePrefixExp)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseObjectLiteral)

	// Infix operators
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.ADD, p.parseInfixExp)
	p.registerInfix(token.SUB, p.parseInfixExp)
	p.registerInfix(token.MUL, p.parseInfixExp)
	p.registerInfix(token.DIV, p.parseInfixExp)
	p.registerInfix(token.MOD, p.parseInfixExp)

	p.registerInfix(token.EQ, p.parseInfixExp)
	p.registerInfix(token.NOT_EQ, p.parseInfixExp)
	p.registerInfix(token.LTHAN, p.parseInfixExp)
	p.registerInfix(token.GTHAN, p.parseInfixExp)
	p.registerInfix(token.LTHAN_EQ, p.parseInfixExp)
	p.registerInfix(token.GTHAN_EQ, p.parseInfixExp)

	p.registerInfix(token.QUESTION, p.parseTernaryExp)
	p.registerInfix(token.LBRACKET, p.parseIndexExp)
	p.registerInfix(token.INC, p.parsePostfixExp)
	p.registerInfix(token.DEC, p.parsePostfixExp)
	p.registerInfix(token.DOT, p.parseDotExp)

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

	prog.Components = p.components
	prog.Inserts = p.inserts
	prog.UseStmt = p.useStmt
	prog.Reserves = p.reserves

	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.HTML:
		return p.parseHTMLStmt()
	case token.LBRACES:
		return p.parseEmbeddedCode()
	case token.SEMI:
		return p.parseEmbeddedCode()
	case token.IF:
		return p.parseIfStmt()
	case token.FOR:
		return p.parseForStmt()
	case token.EACH:
		return p.parseEachStmt()
	case token.USE:
		return p.parseUseStmt()
	case token.RESERVE:
		return p.parseReserveStmt()
	case token.INSERT:
		return p.parseInsertStmt()
	case token.BREAK_IF:
		return p.parseBreakIfStmt()
	case token.CONTINUE_IF:
		return p.parseContinueIfStmt()
	case token.COMPONENT:
		return p.parseComponentStmt()
	case token.BREAK:
		return &ast.BreakStmt{Token: p.curToken}
	case token.CONTINUE:
		return &ast.ContinueStmt{Token: p.curToken}
	default:
		return nil
	}
}

func (p *Parser) parseEmbeddedCode() ast.Statement {
	p.nextToken() // skip "{{" or ";" or "("

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.Line, fail.ErrEmptyBrackets)
		return nil
	}

	if p.curToken.Type == token.IDENT && p.peekTokenIs(token.ASSIGN) {
		return p.parseAssignStmt()
	}

	return p.parseExpressionStmt()
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
	newErr := fail.New(line, p.filepath, "parser", msg, args...)
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
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)

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

func (p *Parser) parseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.curToken} // "{"

	obj.Pairs = make(map[string]ast.Expression)

	p.nextToken() // skip "{"

	if p.curTokenIs(token.RBRACE) {
		return obj
	}

	for !p.peekTokenIs(token.RBRACE) {
		key := p.curToken.Literal

		if !p.expectPeek(token.COLON) { // move to ":"
			return nil
		}

		p.nextToken() // skip ":"

		obj.Pairs[key] = p.parseExpression(LOWEST)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken() // skip value
			p.nextToken() // skip ","
		}
	}

	if !p.expectPeek(token.RBRACE) { // move to "}"
		return nil
	}

	return obj
}

func (p *Parser) parseHTMLStmt() *ast.HTMLStmt {
	return &ast.HTMLStmt{Token: p.curToken}
}

func (p *Parser) parseAssignStmt() ast.Statement {
	ident := &ast.Identifier{
		Token: p.curToken, // identifier
		Value: p.curToken.Literal,
	}

	stmt := &ast.AssignStmt{
		Token: p.curToken, // identifier
		Name:  ident,
	}

	if !p.expectPeek(token.ASSIGN) { // move to "="
		return nil
	}

	p.nextToken() // skip "="

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.Line, fail.ErrExpectedExpression)
		return nil
	}

	stmt.Value = p.parseExpression(SUM)

	return stmt
}

func (p *Parser) parseUseStmt() ast.Statement {
	stmt := &ast.UseStmt{
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

	p.useStmt = stmt

	return stmt
}

func (p *Parser) parseBreakIfStmt() ast.Statement {
	stmt := &ast.BreakIfStmt{
		Token: p.curToken, // "@breakIf"
	}

	if !p.expectPeek(token.LPAREN) { // move to "("
		return nil
	}

	p.nextToken() // skip "("

	stmt.Condition = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseContinueIfStmt() ast.Statement {
	stmt := &ast.ContinueIfStmt{
		Token: p.curToken, // "@continueIf"
	}

	if !p.expectPeek(token.LPAREN) { // move to "("
		return nil
	}

	p.nextToken() // skip "("

	stmt.Condition = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseComponentStmt() ast.Statement {
	stmt := &ast.ComponentStmt{
		Token: p.curToken, // "@component"
	}

	if !p.expectPeek(token.LPAREN) { // move to "("
		return nil
	}

	p.nextToken() // skip "("

	stmt.Name = &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip component name
		stmt.Arguments = p.parseExpressionList(token.RPAREN)
	}

	p.components = append(p.components, stmt)

	return stmt
}

func (p *Parser) parseReserveStmt() ast.Statement {
	stmt := &ast.ReserveStmt{
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

	p.reserves[stmt.Name.Value] = stmt

	return stmt
}

func (p *Parser) parseInsertStmt() ast.Statement {
	stmt := &ast.InsertStmt{
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

		p.inserts[stmt.Name.Value] = stmt

		return stmt
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return nil
	}

	p.nextToken() // skip ")"

	stmt.Block = p.parseBlockStmt()

	p.inserts[stmt.Name.Value] = stmt

	return stmt
}

func (p *Parser) parseIndexExp(left ast.Expression) ast.Expression {
	exp := &ast.IndexExp{
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

func (p *Parser) parsePostfixExp(left ast.Expression) ast.Expression {
	return &ast.PostfixExp{
		Token:    p.curToken,         // identifier
		Operator: p.curToken.Literal, // "++" or "--"
		Left:     left,
	}
}

func (p *Parser) parseDotExp(left ast.Expression) ast.Expression {
	exp := &ast.DotExp{
		Token: p.curToken, // "."
		Left:  left,
	}

	if !p.expectPeek(token.IDENT) { // skip "." and move to identifier
		return nil
	}

	if p.peekTokenIs(token.LPAREN) {
		return p.parseCallExp(left)
	}

	exp.Key = p.parseIdentifier()

	return exp
}

func (p *Parser) parseCallExp(receiver ast.Expression) ast.Expression {
	ident, ok := p.parseIdentifier().(*ast.Identifier)

	if !ok {
		p.newError(p.curToken.Line, fail.ErrExpectedIdentifier, p.curToken.Literal)
		return nil
	}

	exp := &ast.CallExp{
		Token:    p.curToken, // identifier
		Receiver: receiver,
		Function: ident,
	}

	if !p.expectPeek(token.LPAREN) { // move to "("
		return nil
	}

	exp.Arguments = p.parseExpressionList(token.RPAREN)

	return exp
}

func (p *Parser) parseInfixExp(left ast.Expression) ast.Expression {
	exp := &ast.InfixExp{
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

func (p *Parser) parseTernaryExp(left ast.Expression) ast.Expression {
	exp := &ast.TernaryExp{
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

func (p *Parser) parseIfStmt() *ast.IfStmt {
	stmt := &ast.IfStmt{Token: p.curToken} // "@if"

	if !p.expectPeek(token.LPAREN) { // move to "("
		return nil
	}

	p.nextToken() // skip "("

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return nil
	}

	p.nextToken() // skip ")"

	stmt.Consequence = p.parseBlockStmt()

	for p.peekTokenIs(token.ELSEIF) {
		alt := p.parseElseIfStmt()

		if alt == nil {
			return nil
		}

		stmt.Alternatives = append(stmt.Alternatives, alt)
	}

	if p.peekTokenIs(token.ELSE) {
		stmt.Alternative = p.parseAlternativeBlock()

		if stmt.Alternative == nil {
			return nil
		}
	}

	if !p.expectPeek(token.END) { // move to "@end"
		return nil
	}

	return stmt
}

func (p *Parser) parseElseIfStmt() *ast.ElseIfStmt {
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

	return &ast.ElseIfStmt{
		Token:       p.curToken,
		Condition:   condition,
		Consequence: p.parseBlockStmt(),
	}
}

func (p *Parser) parseAlternativeBlock() *ast.BlockStmt {
	p.nextToken() // move to "@else"
	p.nextToken() // skip "@else"

	alt := p.parseBlockStmt()

	if p.peekTokenIs(token.ELSEIF) {
		p.newError(p.peekToken.Line, fail.ErrElseifCannotFollowElse)
		return nil
	}

	return alt
}

func (p *Parser) parseForStmt() *ast.ForStmt {
	stmt := &ast.ForStmt{Token: p.curToken} // "@for"

	if !p.expectPeek(token.LPAREN) { // move to "("
		return nil
	}

	// Parse Init
	if !p.peekTokenIs(token.SEMI) {
		stmt.Init = p.parseEmbeddedCode()
	}

	if !p.expectPeek(token.SEMI) { // move to ";"
		return nil
	}

	// Parse Condition
	if !p.peekTokenIs(token.SEMI) {
		p.nextToken() // skip ";"
		stmt.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.SEMI) { // move to ";"
		return nil
	}

	// Parse Post statement
	if !p.peekTokenIs(token.RPAREN) {
		stmt.Post = p.parseEmbeddedCode()
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return nil
	}

	p.nextToken() // skip ")"

	stmt.Block = p.parseBlockStmt()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // skip "@else"
		stmt.Alternative = p.parseBlockStmt()
	}

	if !p.expectPeek(token.END) { // move to "@end"
		return nil
	}

	return stmt
}

func (p *Parser) parseEachStmt() *ast.EachStmt {
	stmt := &ast.EachStmt{Token: p.curToken} // "@each"

	if !p.expectPeek(token.LPAREN) { // move to "("
		return nil
	}

	p.nextToken() // skip "("

	stmt.Var = &ast.Identifier{
		Token: p.curToken, // identifier
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.IN) { // move to "in"
		return nil
	}

	p.nextToken() // skip "in"

	stmt.Array = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return nil
	}

	p.nextToken() // skip ")"

	stmt.Block = p.parseBlockStmt()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // skip "@else"
		stmt.Alternative = p.parseBlockStmt()
	}

	if !p.expectPeek(token.END) { // move to "@end"
		return nil
	}

	return stmt
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	stmt := &ast.BlockStmt{Token: p.curToken}

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

func (p *Parser) parseExpressionStmt() ast.Statement {
	exp := p.parseExpression(LOWEST)

	result := &ast.ExpressionStmt{Token: p.curToken, Expression: exp}

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

func (p *Parser) parsePrefixExp() ast.Expression {
	exp := &ast.PrefixExp{
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
