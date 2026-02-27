package parser

import (
	"strconv"

	"github.com/textwire/textwire/v3/pkg/file"
	"github.com/textwire/textwire/v3/pkg/token"

	"slices"

	"github.com/textwire/textwire/v3/pkg/ast"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/lexer"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	TERNARY       // a ? b : c
	LOGICAL_OR    // ||
	LOGICAL_AND   // &&
	EQ            // ==
	LESS_GREATER  // > or <
	SUM           // +
	PRODUCT       // *
	PREFIX        // -X or !X
	INDEX         // array[index]
	POSTFIX       // X++ or X--
	MEMBER_ACCESS // <expr>.<ident>
	CALL          // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.QUESTION: TERNARY,
	token.OR:       LOGICAL_OR,
	token.AND:      LOGICAL_AND,
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
	token.LBRACKET: INDEX,
	token.INC:      POSTFIX,
	token.DEC:      POSTFIX,
	token.DOT:      MEMBER_ACCESS,
	token.LPAREN:   CALL,
}

type Parser struct {
	l      *lexer.Lexer
	errors []*fail.Error

	file *file.SourceFile

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	// _useStmt is used to reference the use statement in program.
	// We need it because the final program object must have a field UseStmt.
	// After parsing a program we link this pointer to program.UseStmt.
	_useStmt *ast.UseStmt

	components []*ast.ComponentStmt
	inserts    map[string]*ast.InsertStmt
	reserves   map[string]*ast.ReserveStmt
}

func New(lexer *lexer.Lexer, f *file.SourceFile) *Parser {
	if f == nil {
		f = file.New("", "", "", nil)
	}

	p := &Parser{
		l:          lexer,
		file:       f,
		errors:     []*fail.Error{},
		components: []*ast.ComponentStmt{},
		inserts:    map[string]*ast.InsertStmt{},
		reserves:   map[string]*ast.ReserveStmt{},
	}

	p.nextToken() // fill curToken
	p.nextToken() // fill peekToken

	// Prefix operators
	p.prefixParseFns = map[token.TokenType]prefixParseFn{}

	p.registerPrefix(token.IDENT, p.identifier)
	p.registerPrefix(token.INT, p.integerLiteral)
	p.registerPrefix(token.FLOAT, p.floatLiteral)
	p.registerPrefix(token.STR, p.stringLiteral)
	p.registerPrefix(token.NIL, p.nilLiteral)
	p.registerPrefix(token.TRUE, p.booleanLiteral)
	p.registerPrefix(token.FALSE, p.booleanLiteral)
	p.registerPrefix(token.SUB, p.prefixExp)
	p.registerPrefix(token.NOT, p.prefixExp)
	p.registerPrefix(token.LPAREN, p.groupedExpression)
	p.registerPrefix(token.LBRACKET, p.arrayLiteral)
	p.registerPrefix(token.LBRACE, p.objectLiteral)

	// Infix operators
	p.infixParseFns = map[token.TokenType]infixParseFn{}
	p.registerInfix(token.ADD, p.infixExp)
	p.registerInfix(token.SUB, p.infixExp)
	p.registerInfix(token.MUL, p.infixExp)
	p.registerInfix(token.DIV, p.infixExp)
	p.registerInfix(token.MOD, p.infixExp)

	p.registerInfix(token.EQ, p.infixExp)
	p.registerInfix(token.NOT_EQ, p.infixExp)
	p.registerInfix(token.LTHAN, p.infixExp)
	p.registerInfix(token.GTHAN, p.infixExp)
	p.registerInfix(token.LTHAN_EQ, p.infixExp)
	p.registerInfix(token.GTHAN_EQ, p.infixExp)
	p.registerInfix(token.AND, p.infixExp)
	p.registerInfix(token.OR, p.infixExp)

	p.registerInfix(token.QUESTION, p.ternaryExp)
	p.registerInfix(token.LBRACKET, p.indexExp)
	p.registerInfix(token.INC, p.postfixExp)
	p.registerInfix(token.DEC, p.postfixExp)
	p.registerInfix(token.DOT, p.dotExp)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := ast.NewProgram(p.curToken)
	prog.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.statement()

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
	prog.UseStmt = p._useStmt
	prog.Reserves = p.reserves
	prog.AbsPath = p.file.Abs

	return prog
}

func (p *Parser) Errors() []*fail.Error {
	return p.errors
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}

func (p *Parser) statement() ast.Statement {
	switch p.curToken.Type {
	case token.HTML:
		return p.htmlStmt()
	case token.LBRACES, token.SEMI:
		return p.embeddedCode()
	case token.IF:
		return p.ifStmt()
	case token.FOR:
		return p.forStmt()
	case token.EACH:
		return p.eachStmt()
	case token.USE:
		return p.useStmt()
	case token.RESERVE:
		return p.reserveStmt()
	case token.INSERT:
		return p.insertStmt()
	case token.BREAK_IF:
		return p.breakifStmt()
	case token.CONTINUE_IF:
		return p.continueifStmt()
	case token.COMPONENT:
		return p.componentStmt()
	case token.SLOT:
		return p.slotStmt()
	case token.DUMP:
		return p.dumpStmt()
	case token.BREAK:
		return ast.NewBreakStmt(p.curToken)
	case token.CONTINUE:
		return ast.NewContinueStmt(p.curToken)
	case token.SLOT_IF:
		p.newError(p.curToken.ErrorLine(), fail.ErrSlotifPosition)
		return nil
	default:
		return p.illegalNode()
	}
}

func (p *Parser) embeddedCode() ast.Statement {
	p.nextToken() // skip "{{" or ";" or "("

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.ErrorLine(), fail.ErrEmptyBraces)
		return nil
	}

	if p.curToken.Type == token.IDENT && p.peekTokenIs(token.ASSIGN) {
		return p.assignStmt()
	}

	return p.expressionStmt()
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) curTokenIs(tokens ...token.TokenType) bool {
	if len(tokens) == 1 {
		return tokens[0] == p.curToken.Type
	}
	return slices.Contains(tokens, p.curToken.Type)
}

func (p *Parser) peekTokenIs(tokens ...token.TokenType) bool {
	if len(tokens) == 1 {
		return tokens[0] == p.peekToken.Type
	}
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
	newErr := fail.New(line, p.file.Abs, "parser", msg, args...)
	p.errors = append(p.errors, newErr)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) identifier() ast.Expression {
	ident := ast.NewIdentifier(p.curToken, p.curToken.Literal)
	if p.peekTokenIs(token.LPAREN) {
		return p.globalCallExp(ident)
	}

	return ident
}

func (p *Parser) integerLiteral() ast.Expression {
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

func (p *Parser) floatLiteral() ast.Expression {
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

func (p *Parser) nilLiteral() ast.Expression {
	return ast.NewNilLiteral(p.curToken)
}

func (p *Parser) stringLiteral() ast.Expression {
	return ast.NewStringLiteral(p.curToken, p.curToken.Literal)
}

func (p *Parser) booleanLiteral() ast.Expression {
	return ast.NewBooleanLiteral(p.curToken, p.curTokenIs(token.TRUE))
}

func (p *Parser) arrayLiteral() ast.Expression {
	arr := ast.NewArrayLiteral(p.curToken)
	arr.Elements = p.expressionList(token.RBRACKET)
	arr.SetEndPosition(p.curToken.Pos)

	return arr
}

func (p *Parser) objectLiteral() ast.Expression {
	obj := ast.NewObjectLiteral(p.curToken)

	obj.Pairs = map[string]ast.Expression{}

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

			obj.Pairs[key] = p.expression(LOWEST)
		} else {
			obj.Pairs[key] = p.expression(LOWEST)
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

func (p *Parser) htmlStmt() ast.Statement {
	return ast.NewHTMLStmt(p.curToken)
}

func (p *Parser) assignStmt() ast.Statement {
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

	stmt.Right = p.expression(LOWEST)
	stmt.SetEndPosition(p.curToken.Pos)

	if p.peekTokenIs(token.RBRACES) {
		p.nextToken() // move to '}}'
	}

	return stmt
}

func (p *Parser) useStmt() ast.Statement {
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

	if p._useStmt != nil {
		p.newError(p.curToken.ErrorLine(), fail.ErrOnlyOneUseDir)
	}

	p._useStmt = stmt

	return stmt
}

func (p *Parser) breakifStmt() ast.Statement {
	stmt := ast.NewBreakIfStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	stmt.Condition = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) continueifStmt() ast.Statement {
	stmt := ast.NewContinueIfStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	stmt.Condition = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) componentStmt() ast.Statement {
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

		obj, ok := p.expression(LOWEST).(*ast.ObjectLiteral)
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

	if p.peekTokenIs(token.SLOT, token.SLOT_IF) {
		p.nextToken() // skip whitespace
		if illegal := p.assignSlotsToComp(stmt); illegal != nil {
			return illegal
		}
	}

	p.components = append(p.components, stmt)

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) assignSlotsToComp(stmt *ast.ComponentStmt) ast.Statement {
	slots := p.slots(stmt.Name.Value)
	stmt.Slots = make([]ast.SlotStatement, len(slots))
	for i := range slots {
		slot, ok := slots[i].(ast.SlotStatement)
		if !ok {
			return slots[i]
		}

		stmt.Slots[i] = slot
	}

	return nil
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

// slotStmt parses an external slot statement inside a component file.
// Slots inside a @component are parsed by other function.
func (p *Parser) slotStmt() ast.Statement {
	tok := p.curToken // "@slot"

	// Handle default @slot without name
	if !p.peekTokenIs(token.LPAREN) {
		name := ast.NewStringLiteral(p.curToken, "")
		stmt := ast.NewSlotStmt(tok, name, p.file.Name, false)
		stmt.SetIsDefault(true)
		return stmt
	}

	p.nextToken() // skip "@slot"
	p.nextToken() // skip "("

	// Handle named @slot with name
	name := ast.NewStringLiteral(p.curToken, p.curToken.Literal)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	stmt := ast.NewSlotStmt(tok, name, p.file.Name, false)
	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) dumpStmt() ast.Statement {
	tok := p.curToken // "@dump"

	var args []ast.Expression

	if !p.expectPeek(token.LPAREN) { // move to "("
		return ast.NewIllegalNode(tok)
	}

	args = p.expressionList(token.RPAREN)

	stmt := ast.NewDumpStmt(tok, args)
	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

// slots parses local slots inside of @component directive's body
func (p *Parser) slots(compName string) []ast.Statement {
	var slots []ast.Statement

	for p.curTokenIs(token.SLOT, token.SLOT_IF) {
		slotName := ast.NewStringLiteral(p.curToken, "")

		switch p.curToken.Type {
		case token.SLOT:
			slots = append(slots, p.localSlotStmt(slotName, compName))
		case token.SLOT_IF:
			slots = append(slots, p.slotifStmt(slotName, compName))
		default:
			panic("Unknown slot token when parsing component slots")
		}

		for p.curTokenIs(token.HTML) {
			p.nextToken() // skip whitespace
		}
	}

	return slots
}

func (p *Parser) localSlotStmt(name *ast.StringLiteral, compName string) ast.Statement {
	stmt := ast.NewSlotStmt(p.curToken, name, compName, true)
	stmt.SetIsDefault(!p.peekTokenIs(token.LPAREN))

	// When slot has a name @slot('name')
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken() // move to "(" from '@slot'
		p.nextToken() // skip "("

		name.Token = p.curToken
		name.Value = p.curToken.Literal

		if !p.expectPeek(token.RPAREN) { // move to ")"
			return p.illegalNode() // create an error
		}

		p.nextToken() // skip ")"
	} else {
		p.nextToken() // skip "@slot"
	}

	hasEmptyBlock := p.curTokenIs(token.END)

	if hasEmptyBlock {
		p.nextToken() // skip "@end"
		stmt.SetEndPosition(p.curToken.Pos)
	} else {
		stmt.SetBlock(p.blockStmt())
		stmt.SetEndPosition(p.curToken.Pos)
		p.nextToken() // skip block statement
		p.nextToken() // skip "@end"
	}

	return stmt
}

func (p *Parser) slotifStmt(name *ast.StringLiteral, compName string) ast.Statement {
	stmt := ast.NewSlotifStmt(p.curToken, name, compName)

	if !p.expectPeek(token.LPAREN) { // move from "@slotif" to "("
		p.illegalNode()
	}

	p.nextToken() // skip "("

	stmt.Condition = p.expression(LOWEST)

	// When slot has name
	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // move to ","
		p.nextToken() // skip ","

		name.Token = p.curToken
		name.Value = p.curToken.Literal
	} else {
		stmt.SetIsDefault(true)
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	p.nextToken() // skip ")"

	stmt.SetBlock(p.blockStmt())
	stmt.SetEndPosition(p.curToken.Pos)
	p.nextToken() // skip block statement
	p.nextToken() // skip "@end"

	return stmt
}

func (p *Parser) reserveStmt() ast.Statement {
	stmt := ast.NewReserveStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	stmt.Name = ast.NewStringLiteral(p.curToken, p.curToken.Literal)

	// Handle when has second argument (fallback value) after comma
	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // move to "," from string
		p.nextToken() // move to expression from ","
		stmt.Fallback = p.expression(LOWEST)
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	// Check for duplicate reserve statements
	if _, ok := p.reserves[stmt.Name.Value]; ok {
		p.newError(stmt.Token.ErrorLine(), fail.ErrDuplicateReserves, stmt.Name.Value, p.file.Abs)
		return nil
	}

	p.reserves[stmt.Name.Value] = stmt

	return stmt
}

func (p *Parser) insertStmt() ast.Statement {
	stmt := ast.NewInsertStmt(p.curToken, p.file.Abs)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	stmt.Name = ast.NewStringLiteral(p.curToken, p.curToken.Literal)

	if ok := p.checkDuplicateInserts(stmt); ok {
		return nil
	}

	// Handle inline @insert without block
	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip insert name
		p.nextToken() // skip ","
		stmt.Argument = p.expression(LOWEST)

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
	stmt.Block = p.blockStmt()

	// skip block and move to @end
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

func (p *Parser) indexExp(left ast.Expression) ast.Expression {
	exp := ast.NewIndexExp(p.curToken, left)

	p.nextToken() // skip "["

	exp.Index = p.expression(LOWEST)

	if !p.expectPeek(token.RBRACKET) { // move to "]"
		return p.illegalNode()
	}

	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) postfixExp(left ast.Expression) ast.Expression {
	return ast.NewPostfixExp(p.curToken, left, p.curToken.Literal)
}

func (p *Parser) dotExp(left ast.Expression) ast.Expression {
	exp := ast.NewDotExp(p.curToken, left)

	if !p.expectPeek(token.IDENT) { // skip "." and move to identifier
		return p.illegalNode()
	}

	if p.peekTokenIs(token.LPAREN) {
		return p.callExp(left)
	}

	exp.Key = ast.NewIdentifier(p.curToken, p.curToken.Literal)

	return exp
}

func (p *Parser) globalCallExp(ident *ast.Identifier) ast.Expression {
	exp := ast.NewGlobalCallExp(p.curToken, ident)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	exp.Arguments = p.expressionList(token.RPAREN)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) callExp(receiver ast.Expression) ast.Expression {
	ident := ast.NewIdentifier(p.curToken, p.curToken.Literal)
	exp := ast.NewCallExp(p.curToken, receiver, ident)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	exp.Arguments = p.expressionList(token.RPAREN)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) infixExp(left ast.Expression) ast.Expression {
	exp := ast.NewInfixExp(*left.Tok(), left, p.curToken.Literal)

	precedence := precedences[p.curToken.Type]

	p.nextToken() // skip operator

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.ErrorLine(), fail.ErrExpectedExpression)
		return nil
	}

	exp.Right = p.expression(precedence)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) ternaryExp(left ast.Expression) ast.Expression {
	exp := ast.NewTernaryExp(*left.Tok(), left)

	p.nextToken() // skip "?"

	exp.IfBlock = p.expression(TERNARY)

	if !p.expectPeek(token.COLON) { // move to ":"
		return p.illegalNode()
	}

	p.nextToken() // skip ":"

	exp.ElseBlock = p.expression(LOWEST)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) ifStmt() ast.Statement {
	stmt := ast.NewIfStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip "("

	stmt.Condition = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip ")"

	stmt.IfBlock = p.blockStmt()
	if stmt.IfBlock == nil {
		stmt.SetEndPosition(p.curToken.Pos)
		return stmt
	}

	for p.peekTokenIs(token.ELSE_IF) {
		elseifStmt := p.elseifStmt()
		stmt.ElseifStmts = append(stmt.ElseifStmts, elseifStmt)
	}

	if p.peekTokenIs(token.ELSE) {
		stmt.ElseBlock = p.elseBlock()
		if stmt.ElseBlock == nil {
			return nil
		}
	}

	if !p.expectPeek(token.END) { // move to "@end"
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) elseifStmt() ast.Statement {
	if !p.expectPeek(token.ELSE_IF) { // move to "@elseif"
		return p.illegalNode()
	}

	stmt := ast.NewElseIfStmt(p.curToken)

	p.nextToken() // skip "@elseif"
	p.nextToken() // skip "("

	stmt.Condition = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	p.nextToken() // skip ")"

	stmt.Block = p.blockStmt()
	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) elseBlock() *ast.BlockStmt {
	p.nextToken() // move to "@else"
	p.nextToken() // skip "@else"

	stmt := p.blockStmt()

	if p.peekTokenIs(token.ELSE_IF) {
		p.newError(p.peekToken.ErrorLine(), fail.ErrElseifCannotFollowElse)
		return nil
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) forStmt() ast.Statement {
	stmt := ast.NewForStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.END)
	}

	// Parse Init
	if !p.peekTokenIs(token.SEMI) {
		stmt.Init = p.embeddedCode()
	}

	if !p.expectPeek(token.SEMI) { // move to ";"
		return p.illegalNodeUntil(token.END)
	}

	// Parse Condition
	if !p.peekTokenIs(token.SEMI) {
		p.nextToken() // skip ";"
		stmt.Condition = p.expression(LOWEST)
	}

	if !p.expectPeek(token.SEMI) { // move to ";"
		return p.illegalNodeUntil(token.END)
	}

	// Parse Post statement
	if !p.peekTokenIs(token.RPAREN) {
		stmt.Post = p.embeddedCode()
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip ")"

	stmt.Block = p.blockStmt()
	if stmt.Block == nil {
		stmt.SetEndPosition(p.curToken.Pos)
		return stmt
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // move to "@else"
		p.nextToken() // skip "@else"
		stmt.ElseBlock = p.blockStmt()
	}

	if !p.expectPeek(token.END) { // move to "@end"
		return p.illegalNodeUntil(token.END)
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) eachStmt() ast.Statement {
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

	stmt.Array = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip ")"

	stmt.Block = p.blockStmt()
	if stmt.Block == nil {
		stmt.SetEndPosition(p.curToken.Pos)
		return stmt
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // move to "@else"
		p.nextToken() // skip "@else"
		stmt.ElseBlock = p.blockStmt()
	}

	if !p.expectPeek(token.END) { // move to "@end"
		return p.illegalNodeUntil(token.END)
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) blockStmt() *ast.BlockStmt {
	if p.curTokenIs(token.END) {
		return nil
	}

	stmt := ast.NewBlockStmt(p.curToken)
	stmt.SetEndPosition(p.curToken.Pos)

	for !p.curTokenIs(token.END) && !p.curTokenIs(token.EOF) {
		block := p.statement()
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

func (p *Parser) expressionStmt() ast.Statement {
	prevTok := p.curToken

	exp := p.expression(LOWEST)

	stmt := ast.NewExpressionStmt(prevTok, exp)
	stmt.SetEndPosition(p.curToken.Pos)

	if p.peekTokenIs(token.RBRACES) {
		p.nextToken() // skip "}}"
	}

	return stmt
}

func (p *Parser) expression(precedence int) ast.Expression {
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

func (p *Parser) prefixExp() ast.Expression {
	exp := ast.NewPrefixExp(p.curToken, p.curToken.Literal)

	p.nextToken() // skip prefix operator

	exp.Right = p.expression(PREFIX)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) groupedExpression() ast.Expression {
	p.nextToken() // skip "("

	exp := p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	return exp
}

func (p *Parser) expressionList(endTok token.TokenType) []ast.Expression {
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

	result = append(result, p.expression(LOWEST))

	for p.peekTokenIs(token.COMMA) && !p.curTokenIs(token.EOF) {
		p.nextToken() // move to ","

		// break when has a trailing comma
		if p.peekTokenIs(endTok) {
			break
		}

		p.nextToken() // skip ","
		result = append(result, p.expression(LOWEST))
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
