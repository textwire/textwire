package parser

import (
	"strconv"

	"github.com/textwire/textwire/v2/token"

	"slices"

	"github.com/textwire/textwire/v2/ast"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/lexer"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	TERNARY       // a ? b : c
	EQ            // ==
	LESS_GREATER  // > or <
	SUM           // +
	PRODUCT       // *
	MEMBER_ACCESS // <expr>.<ident>
	PREFIX        // -X or !X
	CALL          // myFunction(X)
	INDEX         // array[index]
	POSTFIX       // X++ or X--
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
	token.DOT:      MEMBER_ACCESS,
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
	prog := ast.NewProgram(p.curToken)
	prog.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		// if the end of the {{ expression }}
		if p.curTokenIs(token.RBRACES) {
			prog.Statements = append(prog.Statements, stmt)
			p.nextToken()
			continue
		}

		if stmt == nil {
			p.nextToken() // skip to next token
			continue
		}

		prog.Statements = append(prog.Statements, stmt)

		p.nextToken() // skip to next token
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
	case token.SLOT:
		return p.parseSlotStmt()
	case token.DUMP:
		return p.parseDumpStmt()
	case token.BREAK:
		return ast.NewBreakStmt(p.curToken)
	case token.CONTINUE:
		return ast.NewContinueStmt(p.curToken)
	default:
		return p.illegalNode()
	}
}

func (p *Parser) parseEmbeddedCode() ast.Statement {
	p.nextToken() // skip "{{" or ";" or "("

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.ErrorLine(), fail.ErrEmptyBraces)
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
	return slices.Contains(tokens, p.peekToken.Type)
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
		p.peekToken.ErrorLine(),
		fail.ErrWrongNextToken,
		token.String(tok),
		token.String(p.peekToken.Type),
	)

	return false
}

func (p *Parser) newError(line uint, msg string, args ...any) {
	newErr := fail.New(line, p.filepath, "parser", msg, args...)
	p.errors = append(p.errors, newErr)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.NewIdentifier(p.curToken, p.curToken.Literal)
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.newError(
			p.curToken.ErrorLine(),
			fail.ErrCouldNotParseAs,
			p.curToken.Literal,
			"INT",
		)
		return nil
	}

	return ast.NewIntegerLiteral(p.curToken, val)
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.newError(
			p.curToken.ErrorLine(),
			fail.ErrCouldNotParseAs,
			p.curToken.Literal,
			"FLOAT",
		)
		return nil
	}

	return ast.NewFloatLiteral(p.curToken, val)
}

func (p *Parser) parseNilLiteral() ast.Expression {
	return ast.NewNilLiteral(p.curToken)
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return ast.NewStringLiteral(p.curToken, p.curToken.Literal)
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return ast.NewBooleanLiteral(p.curToken, p.curTokenIs(token.TRUE))
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := ast.NewArrayLiteral(p.curToken)
	arr.Elements = p.parseExpressionList(token.RBRACKET)
	arr.SetEndPosition(p.curToken.Pos)

	return arr
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	obj := ast.NewObjectLiteral(p.curToken)

	obj.Pairs = make(map[string]ast.Expression)

	p.nextToken() // skip "{"

	if p.curTokenIs(token.RBRACE) {
		obj.SetEndPosition(p.curToken.Pos)
		return obj
	}

	for !p.curTokenIs(token.RBRACE) {
		key := p.curToken.Literal

		if p.peekTokenIs(token.COLON) {
			p.nextToken() // move to ":"
			p.nextToken() // skip to ":"

			obj.Pairs[key] = p.parseExpression(LOWEST)
		} else {
			obj.Pairs[key] = p.parseExpression(LOWEST)
		}

		if p.peekTokenIs(token.RBRACE) {
			p.nextToken() // skip "}"
			break
		}

		if p.peekTokenIs(token.COMMA) {
			p.nextToken() // move to ","
			p.nextToken() // skip ","
		}
	}

	obj.SetEndPosition(p.curToken.Pos)

	return obj
}

func (p *Parser) parseHTMLStmt() ast.Statement {
	return ast.NewHTMLStmt(p.curToken)
}

func (p *Parser) parseAssignStmt() ast.Statement {
	ident := ast.NewIdentifier(p.curToken, p.curToken.Literal)

	stmt := ast.NewAssignStmt(p.curToken, ident)

	if !p.expectPeek(token.ASSIGN) { // move to "="
		return p.illegalNode()
	}

	p.nextToken() // skip "="

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.ErrorLine(), fail.ErrExpectedExpression)
		return nil
	}

	stmt.Value = p.parseExpression(SUM)
	stmt.SetEndPosition(p.curToken.Pos)

	if p.peekTokenIs(token.RBRACES) {
		p.nextToken() // move to '}}'
	}

	return stmt
}

func (p *Parser) parseUseStmt() ast.Statement {
	stmt := ast.NewUseStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	if p.curToken.Type != token.STR {
		p.newError(
			p.curToken.ErrorLine(),
			fail.ErrUseStmtFirstArgStr,
			token.String(p.curToken.Type),
		)
	}

	stmt.Name = ast.NewStringLiteral(
		p.curToken,
		p.parseAliasPathShortcut("layouts"),
	)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	p.useStmt = stmt

	return stmt
}

func (p *Parser) parseBreakIfStmt() ast.Statement {
	stmt := ast.NewBreakIfStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) parseContinueIfStmt() ast.Statement {
	stmt := ast.NewContinueIfStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) parseComponentStmt() ast.Statement {
	stmt := ast.NewComponentStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.RPAREN)
	}

	p.nextToken() // skip "("

	stmt.Name = ast.NewStringLiteral(
		p.curToken,
		p.parseAliasPathShortcut("components"),
	)

	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // move to ","
		p.nextToken() // skip ","

		obj, ok := p.parseExpression(LOWEST).(*ast.ObjectLiteral)

		if !ok {
			p.newError(p.curToken.ErrorLine(), fail.ErrExpectedObjectLiteral, p.curToken.Literal)
			return nil
		}

		stmt.Argument = obj
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.HTML)
	}

	if p.peekTokenIs(token.HTML) && isWhitespace(p.peekToken.Literal) {
		p.nextToken() // move to ")"
	}

	if p.peekTokenIs(token.SLOT) {
		p.nextToken() // skip whitespace
		stmt.Slots = p.parseSlots()
	}

	p.components = append(p.components, stmt)

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) parseAliasPathShortcut(shortenTo string) string {
	name := p.curToken.Literal

	if name == "" {
		p.newError(p.curToken.ErrorLine(), fail.ErrExpectedComponentName)
		return ""
	}

	if name[0] == '~' {
		name = shortenTo + "/" + name[1:]
	}

	return name
}

// parseSlotStmt parses a slot statement inside a component file.
// Slots inside a component are parsed by other function
func (p *Parser) parseSlotStmt() ast.Statement {
	tok := p.curToken // "@slot"

	if !p.peekTokenIs(token.LPAREN) {
		return ast.NewSlotStmt(tok, ast.NewStringLiteral(p.curToken, ""))
	}

	p.nextToken() // skip "@slot"
	p.nextToken() // skip "("

	slotName := ast.NewStringLiteral(p.curToken, p.curToken.Literal)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	stmt := ast.NewSlotStmt(tok, slotName)
	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) parseDumpStmt() ast.Statement {
	tok := p.curToken // "@dump"

	var args []ast.Expression

	if !p.expectPeek(token.LPAREN) { // move to "("
		return ast.NewIllegalNode(tok)
	}

	args = p.parseExpressionList(token.RPAREN)

	stmt := ast.NewDumpStmt(tok, args)
	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) parseSlots() []*ast.SlotStmt {
	var slots []*ast.SlotStmt

	for p.curTokenIs(token.SLOT) {
		slotName := ast.NewStringLiteral(p.curToken, "")

		tok := p.curToken

		if p.peekTokenIs(token.LPAREN) {
			p.nextToken() // move to "("
			p.nextToken() // skip "("

			slotName.Token = p.curToken
			slotName.Value = p.curToken.Literal

			if !p.expectPeek(token.RPAREN) { // move to ")"
				p.illegalNode() // create an error
			}

			p.nextToken() // skip ")"
		}

		stmt := ast.NewSlotStmt(tok, slotName)
		hasEmptyBody := p.curTokenIs(token.END)

		if hasEmptyBody {
			p.nextToken() // skip "@end"
			stmt.SetEndPosition(p.curToken.Pos)
		} else {
			stmt.Body = p.parseBlockStmt()
			stmt.SetEndPosition(p.curToken.Pos)
			p.nextToken() // skip block statement
			p.nextToken() // skip "@end"
		}

		slots = append(slots, stmt)

		for p.curTokenIs(token.HTML) {
			p.nextToken() // skip whitespace
		}
	}

	return slots
}

func (p *Parser) parseReserveStmt() ast.Statement {
	stmt := ast.NewReserveStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	stmt.Name = ast.NewStringLiteral(p.curToken, p.curToken.Literal)

	if !p.expectPeek(token.RPAREN) { // skip string token
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	p.reserves[stmt.Name.Value] = stmt

	return stmt
}

func (p *Parser) parseInsertStmt() ast.Statement {
	stmt := ast.NewInsertStmt(p.curToken, p.filepath)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	stmt.Name = ast.NewStringLiteral(p.curToken, p.curToken.Literal)

	if hasDuplicates := p.checkDuplicateInserts(stmt); hasDuplicates {
		return nil
	}

	// Handle inline @insert without body
	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip insert name
		p.nextToken() // skip ","
		stmt.Argument = p.parseExpression(LOWEST)

		if !p.expectPeek(token.RPAREN) { // move to ")"
			return p.illegalNodeUntil(token.RBRACE)
		}

		stmt.SetEndPosition(p.curToken.Pos)

		p.inserts[stmt.Name.Value] = stmt

		return stmt
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	if p.peekTokenIs(token.END) {
		p.nextToken() // move to "@end"
		stmt.SetEndPosition(p.curToken.Pos)
		return stmt
	}

	p.nextToken() // skip ")"
	stmt.Block = p.parseBlockStmt()

	// skip body block and move to @end
	if !p.expectPeek(token.END) {
		return p.illegalNodeUntil(token.END)
	}

	stmt.SetEndPosition(p.curToken.Pos)

	p.inserts[stmt.Name.Value] = stmt

	return stmt
}

func (p *Parser) checkDuplicateInserts(stmt *ast.InsertStmt) bool {
	if _, hasDuplicate := p.inserts[stmt.Name.Value]; hasDuplicate {
		p.newError(
			stmt.Token.ErrorLine(),
			fail.ErrDuplicateInserts,
			stmt.Name.Value,
		)

		return true
	}

	return false
}

func (p *Parser) parseIndexExp(left ast.Expression) ast.Expression {
	exp := ast.NewIndexExp(p.curToken, left)

	p.nextToken() // skip "["

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) { // move to "]"
		return p.illegalNode()
	}

	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) parsePostfixExp(left ast.Expression) ast.Expression {
	return ast.NewPostfixExp(p.curToken, left, p.curToken.Literal)
}

func (p *Parser) parseDotExp(left ast.Expression) ast.Expression {
	exp := ast.NewDotExp(p.curToken, left)

	if !p.expectPeek(token.IDENT) { // skip "." and move to identifier
		return p.illegalNode()
	}

	if p.peekTokenIs(token.LPAREN) {
		return p.parseCallExp(left)
	}

	exp.Key = p.parseIdentifier()

	return exp
}

func (p *Parser) parseCallExp(receiver ast.Expression) ast.Expression {
	ident := ast.NewIdentifier(p.curToken, p.curToken.Literal)
	exp := ast.NewCallExp(p.curToken, receiver, ident)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	exp.Arguments = p.parseExpressionList(token.RPAREN)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) parseInfixExp(left ast.Expression) ast.Expression {
	exp := ast.NewInfixExp(*left.Tok(), left, p.curToken.Literal)

	p.nextToken() // skip operator

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.ErrorLine(), fail.ErrExpectedExpression)
		return nil
	}

	exp.Right = p.parseExpression(SUM)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) parseTernaryExp(left ast.Expression) ast.Expression {
	exp := ast.NewTernaryExp(*left.Tok(), left)

	p.nextToken() // skip "?"

	exp.Consequence = p.parseExpression(TERNARY)

	if !p.expectPeek(token.COLON) { // move to ":"
		return p.illegalNode()
	}

	p.nextToken() // skip ":"

	exp.Alternative = p.parseExpression(LOWEST)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) parseIfStmt() ast.Statement {
	stmt := ast.NewIfStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip "("

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip ")"

	stmt.Consequence = p.parseBlockStmt()
	if stmt.Consequence == nil {
		stmt.SetEndPosition(p.curToken.Pos)
		return stmt
	}

	for p.peekTokenIs(token.ELSE_IF) {
		alt := p.parseElseIfStmt()
		stmt.Alternatives = append(stmt.Alternatives, alt)
	}

	if p.peekTokenIs(token.ELSE) {
		stmt.Alternative = p.parseAlternativeBlock()

		if stmt.Alternative == nil {
			return nil
		}
	}

	if !p.expectPeek(token.END) { // move to "@end"
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) parseElseIfStmt() ast.Statement {
	if !p.expectPeek(token.ELSE_IF) { // move to "@elseif"
		return p.illegalNode()
	}

	stmt := ast.NewElseIfStmt(p.curToken)

	p.nextToken() // skip "@elseif"
	p.nextToken() // skip "("

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	p.nextToken() // skip ")"

	stmt.Consequence = p.parseBlockStmt()
	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) parseAlternativeBlock() *ast.BlockStmt {
	p.nextToken() // move to "@else"
	p.nextToken() // skip "@else"

	stmt := p.parseBlockStmt()

	if p.peekTokenIs(token.ELSE_IF) {
		p.newError(p.peekToken.ErrorLine(), fail.ErrElseifCannotFollowElse)
		return nil
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) parseForStmt() ast.Statement {
	stmt := ast.NewForStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.END)
	}

	// Parse Init
	if !p.peekTokenIs(token.SEMI) {
		stmt.Init = p.parseEmbeddedCode()
	}

	if !p.expectPeek(token.SEMI) { // move to ";"
		return p.illegalNodeUntil(token.END)
	}

	// Parse Condition
	if !p.peekTokenIs(token.SEMI) {
		p.nextToken() // skip ";"
		stmt.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.SEMI) { // move to ";"
		return p.illegalNodeUntil(token.END)
	}

	// Parse Post statement
	if !p.peekTokenIs(token.RPAREN) {
		stmt.Post = p.parseEmbeddedCode()
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip ")"

	stmt.Block = p.parseBlockStmt()
	if stmt.Block == nil {
		stmt.SetEndPosition(p.curToken.Pos)
		return stmt
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // move to "@else"
		p.nextToken() // skip "@else"
		stmt.Alternative = p.parseBlockStmt()
	}

	if !p.expectPeek(token.END) { // move to "@end"
		return p.illegalNodeUntil(token.END)
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) parseEachStmt() ast.Statement {
	stmt := ast.NewEachStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip "("

	stmt.Var = ast.NewIdentifier(p.curToken, p.curToken.Literal)

	if !p.expectPeek(token.IN) { // move to "in"
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip "in"

	stmt.Array = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip ")"

	stmt.Block = p.parseBlockStmt()
	if stmt.Block == nil {
		stmt.SetEndPosition(p.curToken.Pos)
		return stmt
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // move to "@else"
		p.nextToken() // skip "@else"
		stmt.Alternative = p.parseBlockStmt()
	}

	if !p.expectPeek(token.END) { // move to "@end"
		return p.illegalNodeUntil(token.END)
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	if p.curTokenIs(token.END) {
		return nil
	}

	stmt := ast.NewBlockStmt(p.curToken)
	stmt.SetEndPosition(p.curToken.Pos)

	for !p.curTokenIs(token.END) && !p.curTokenIs(token.EOF) {
		block := p.parseStatement()
		stmt.SetEndPosition(p.curToken.Pos)

		if block != nil {
			stmt.Statements = append(stmt.Statements, block)
		}

		if p.peekTokenIs(token.ELSE, token.ELSE_IF, token.END) {
			break
		}

		p.nextToken() // skip statement
	}

	return stmt
}

func (p *Parser) parseExpressionStmt() ast.Statement {
	prevTok := p.curToken

	exp := p.parseExpression(LOWEST)

	stmt := ast.NewExpressionStmt(prevTok, exp)
	stmt.SetEndPosition(p.curToken.Pos)

	if p.peekTokenIs(token.RBRACES) {
		p.nextToken() // skip "}}"
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.newError(
			p.curToken.ErrorLine(),
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
	exp := ast.NewPrefixExp(p.curToken, p.curToken.Literal)

	p.nextToken() // skip prefix operator

	exp.Right = p.parseExpression(PREFIX)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // skip "("

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	return exp
}

func (p *Parser) parseExpressionList(endTok token.TokenType) []ast.Expression {
	var result []ast.Expression

	if p.peekTokenIs(endTok) {
		p.nextToken() // skip endTok token
		return result
	}

	if p.peekTokenIs(token.END) {
		result = append(result, p.illegalNode())
		return result
	}

	p.nextToken() // move to first expression

	result = append(result, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) && !p.curTokenIs(token.EOF) {
		p.nextToken() // move to ","

		// break when has a trailing comma
		if p.peekTokenIs(endTok) {
			break
		}

		p.nextToken() // skip ","
		result = append(result, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(endTok) { // move to endTok
		result = append(result, ast.NewIllegalNode(p.curToken))
		return result
	}

	return result
}

func (p *Parser) illegalNode() *ast.IllegalNode {
	p.newError(
		p.curToken.ErrorLine(),
		fail.ErrIllegalToken,
		p.curToken.Literal,
	)

	return ast.NewIllegalNode(p.curToken)
}

func (p *Parser) illegalNodeUntil(tok token.TokenType) *ast.IllegalNode {
	for !p.curTokenIs(tok) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	return p.illegalNode()
}
