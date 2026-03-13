package parser

import (
	"strconv"

	"github.com/textwire/textwire/v3/pkg/file"
	"github.com/textwire/textwire/v3/pkg/position"
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
	INDEX         // arr[index]
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

	// _useDir is used to reference the use directive in the program.
	// We need it because the final program object must have a field UseDir.
	// After parsing a program we link this pointer to program.UseDir.
	_useDir *ast.UseDir

	components []*ast.ComponentDir
	inserts    map[string]*ast.InsertDir
	reserves   map[string]*ast.ReserveDir
}

func New(lexer *lexer.Lexer, f *file.SourceFile) *Parser {
	if f == nil {
		f = file.New("", "", "", nil)
	}

	p := &Parser{
		l:          lexer,
		file:       f,
		errors:     []*fail.Error{},
		components: []*ast.ComponentDir{},
		inserts:    map[string]*ast.InsertDir{},
		reserves:   map[string]*ast.ReserveDir{},
	}

	p.nextToken() // fill curToken
	p.nextToken() // fill peekToken

	// Prefix operators
	p.prefixParseFns = map[token.TokenType]prefixParseFn{}

	p.registerPrefix(token.IDENT, p.ident)
	p.registerPrefix(token.INT, p.intExpr)
	p.registerPrefix(token.FLOAT, p.floatExpr)
	p.registerPrefix(token.STR, p.strExpr)
	p.registerPrefix(token.NIL, p.nilExpr)
	p.registerPrefix(token.TRUE, p.boolExpr)
	p.registerPrefix(token.FALSE, p.boolExpr)
	p.registerPrefix(token.SUB, p.prefixExpr)
	p.registerPrefix(token.NOT, p.prefixExpr)
	p.registerPrefix(token.LPAREN, p.groupedExpr)
	p.registerPrefix(token.LBRACKET, p.arrExpr)
	p.registerPrefix(token.LBRACE, p.objExpr)

	// Infix operators
	p.infixParseFns = map[token.TokenType]infixParseFn{}
	p.registerInfix(token.ADD, p.infixExpr)
	p.registerInfix(token.SUB, p.infixExpr)
	p.registerInfix(token.MUL, p.infixExpr)
	p.registerInfix(token.DIV, p.infixExpr)
	p.registerInfix(token.MOD, p.infixExpr)

	p.registerInfix(token.EQ, p.infixExpr)
	p.registerInfix(token.NOT_EQ, p.infixExpr)
	p.registerInfix(token.LTHAN, p.infixExpr)
	p.registerInfix(token.GTHAN, p.infixExpr)
	p.registerInfix(token.LTHAN_EQ, p.infixExpr)
	p.registerInfix(token.GTHAN_EQ, p.infixExpr)
	p.registerInfix(token.AND, p.infixExpr)
	p.registerInfix(token.OR, p.infixExpr)

	p.registerInfix(token.QUESTION, p.ternaryExpr)
	p.registerInfix(token.LBRACKET, p.indexExpr)
	p.registerInfix(token.DOT, p.dotExpr)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := ast.NewProgram(p.curToken)
	prog.Chunks = []ast.Chunk{}

	for !p.curTokenIs(token.EOF) {
		chunk := p.chunk()
		if chunk == nil {
			p.nextToken() // skip to next token
			continue
		}

		prog.Chunks = append(prog.Chunks, chunk)

		p.nextToken() // skip to next token
	}

	prog.Components = p.components
	prog.Inserts = p.inserts
	prog.UseDir = p._useDir
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

func (p *Parser) chunk() ast.Chunk {
	switch p.curToken.Type {
	case token.TEXT:
		return p.text()
	case token.LBRACES:
		return p.embedded()
	case token.IF:
		return p.ifDir()
	case token.FOR:
		return p.forDir()
	case token.EACH:
		return p.eachDir()
	case token.USE:
		return p.useDir()
	case token.RESERVE:
		return p.reserveDir()
	case token.INSERT:
		return p.insertDir()
	case token.BREAKIF:
		return p.breakifDir()
	case token.CONTINUEIF:
		return p.continueifDir()
	case token.COMPONENT:
		return p.compDir()
	case token.SLOT:
		return p.slotDir()
	case token.DUMP:
		return p.dumpDir()
	case token.BREAK:
		return ast.NewBreakDir(p.curToken)
	case token.CONTINUE:
		return ast.NewContinueDir(p.curToken)
	}

	return p.illegalNode()
}

func (p *Parser) embedded() ast.Chunk {
	embedded := ast.NewEmbedded(p.curToken)

	if p.peekTokenIs(token.RBRACES) {
		pos := p.curToken.Pos
		pos.EndCol = p.peekToken.Pos.EndCol
		p.newError(pos, fail.ErrEmptyBraces)
		return nil
	}

	p.nextToken() // skip "{{"

	// Loop until we find the closing "}}" or reach the end of file
	for !p.curTokenIs(token.RBRACES, token.EOF) {
		if segment := p.parseSegment(); segment != nil {
			embedded.Segments = append(embedded.Segments, segment)
			if p.curTokenIs(token.SEMI) {
				p.nextToken() // skip ";"
			}
		}
	}

	if p.peekTokenIs(token.RBRACES) {
		p.nextToken() // skip "}}"
	}

	return embedded
}

func (p *Parser) peekIsStatement() bool {
	return p.peekTokenIs(token.ASSIGN, token.INC, token.DEC)
}

// parseSegment parses individual segment
func (p *Parser) parseSegment() ast.Segment {
	left := p.expression(LOWEST)
	if left == nil {
		p.nextToken()
		return nil
	}

	if !p.peekIsStatement() {
		p.nextToken()
		return left.(ast.Segment)
	}

	p.nextToken()
	stmt := p.statement(left)
	p.nextToken()
	return stmt.(ast.Segment)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) expectType(tok token.Token, expectType token.TokenType) bool {
	if tok.Type == expectType {
		return true
	}

	p.newError(
		tok.Pos,
		fail.ErrWrongTokenType,
		token.String(expectType),
		token.String(tok.Type),
	)

	return false
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

func (p *Parser) expectPeek2(tok token.TokenType, fromPos *position.Pos) bool {
	if p.peekTokenIs(tok) {
		p.nextToken()
		return true
	}

	fromPos.EndLine = p.peekToken.Pos.EndLine
	fromPos.EndCol = p.peekToken.Pos.EndCol

	p.newError(
		fromPos,
		fail.ErrWrongPeekToken,
		token.String(tok),
		token.String(p.peekToken.Type),
	)

	return false
}

func (p *Parser) expectPeek(tok token.TokenType) bool {
	if p.peekTokenIs(tok) {
		p.nextToken()
		return true
	}

	p.newError(
		p.peekToken.Pos,
		fail.ErrWrongPeekToken,
		token.String(tok),
		token.String(p.peekToken.Type),
	)

	return false
}

func (p *Parser) newError(pos *position.Pos, msg string, args ...any) {
	newErr := fail.New(pos, p.file.Abs, fail.OriginPars, msg, args...)
	p.errors = append(p.errors, newErr)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.Next()
}

func (p *Parser) ident() ast.Expression {
	ident := ast.NewIdentExpr(p.curToken, p.curToken.Lit)
	if p.peekTokenIs(token.LPAREN) {
		return p.globalCallExpr(ident)
	}

	return ident
}

func (p *Parser) intExpr() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Lit, 10, 64)
	if err != nil {
		p.newError(
			p.curToken.Pos,
			fail.ErrCouldNotParseAs,
			p.curToken.Lit,
			"INT",
		)
		return nil
	}

	return ast.NewIntExpr(p.curToken, val)
}

func (p *Parser) floatExpr() ast.Expression {
	val, err := strconv.ParseFloat(p.curToken.Lit, 64)
	if err != nil {
		p.newError(
			p.curToken.Pos,
			fail.ErrCouldNotParseAs,
			p.curToken.Lit,
			"FLOAT",
		)
		return nil
	}

	return ast.NewFloatExpr(p.curToken, val)
}

func (p *Parser) nilExpr() ast.Expression {
	return ast.NewNilExpr(p.curToken)
}

func (p *Parser) strExpr() ast.Expression {
	return ast.NewStrExpr(p.curToken, p.curToken.Lit)
}

func (p *Parser) boolExpr() ast.Expression {
	return ast.NewBoolExpr(p.curToken, p.curTokenIs(token.TRUE))
}

func (p *Parser) arrExpr() ast.Expression {
	arr := ast.NewArrExpr(p.curToken)
	arr.Elements = p.expressionList(token.RBRACKET)
	arr.SetEndPosition(p.curToken.Pos)

	return arr
}

func (p *Parser) objExpr() ast.Expression {
	obj := ast.NewObjExpr(p.curToken)

	obj.Pairs = map[string]ast.Expression{}

	p.nextToken() // skip "{"

	if p.curTokenIs(token.RBRACE) {
		obj.SetEndPosition(p.curToken.Pos)
		return obj
	}

	for !p.curTokenIs(token.RBRACE) {
		key := p.curToken.Lit

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

func (p *Parser) text() ast.Chunk {
	return ast.NewText(p.curToken)
}

func (p *Parser) assignStmt(left ast.Expression) ast.Statement {
	stmt := ast.NewAssignStmt(*left.Tok(), left)

	p.nextToken() // skip "="

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.Pos, fail.ErrExpectedExpression)
		return nil
	}

	stmt.Right = p.expression(LOWEST)
	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) useDir() ast.Chunk {
	dir := ast.NewUseDir(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	if p.curToken.Lit == "" {
		p.newError(p.curToken.Pos, fail.ErrStrCannotBeEmpty)
		return nil
	}

	if !p.expectType(p.curToken, token.STR) {
		return nil
	}

	dir.Name = ast.NewStrExpr(
		p.curToken,
		file.ReplacePathAlias(p.curToken.Lit, file.PathAliasUse),
	)

	if p._useDir != nil {
		p.newError(dir.Name.Pos(), fail.ErrOnlyOneUseDir)
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	dir.SetEndPosition(p.curToken.Pos)

	p._useDir = dir

	return dir
}

func (p *Parser) breakifDir() ast.Chunk {
	dir := ast.NewBreakIfDir(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	dir.Cond = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) continueifDir() ast.Chunk {
	dir := ast.NewContinueIfDir(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	dir.Cond = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) compDir() ast.Chunk {
	dir := ast.NewComponentDir(p.curToken)

	if illegal := p.compDirHeader(dir); illegal != nil {
		return illegal
	}

	if p.peekTokenIs(token.SLOT, token.SLOTIF) {
		p.nextToken() // skip whitespace
		if illegal := p.attachSlotsToComp(dir); illegal != nil {
			return illegal
		}
	}

	p.components = append(p.components, dir)

	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) compDirHeader(compDir *ast.ComponentDir) *ast.IllegalNode {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.RPAREN)
	}

	p.nextToken() // skip "("

	if p.curToken.Lit == "" {
		p.newError(p.curToken.Pos, fail.ErrStrCannotBeEmpty)
	}

	compDir.Name = ast.NewStrExpr(
		p.curToken,
		file.ReplacePathAlias(p.curToken.Lit, file.PathAliasComp),
	)

	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // move to ","
		p.nextToken() // skip ","

		obj, ok := p.expression(LOWEST).(*ast.ObjExpr)
		if !ok {
			p.newError(p.curToken.Pos, fail.ErrExpectedObjLit, p.curToken.Lit)
			return nil
		}

		compDir.Argument = obj
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.TEXT)
	}

	if p.peekTokenIs(token.TEXT) && isWhitespace(p.peekToken.Lit) {
		p.nextToken() // move to ")"
	}

	return nil
}

func (p *Parser) attachSlotsToComp(compDir *ast.ComponentDir) ast.Chunk {
	slots := p.slots(compDir.Name.Val)
	compDir.Slots = make([]ast.SlotDirective, len(slots))
	for i := range slots {
		slot, ok := slots[i].(ast.SlotDirective)
		if !ok {
			return slots[i]
		}

		compDir.Slots[i] = slot
	}
	return nil
}

// slotDir parses an external slot statement inside a component file.
// Slots inside a @component are parsed by other function.
func (p *Parser) slotDir() ast.Chunk {
	tok := p.curToken // "@slot"

	// Handle default @slot without name
	if !p.peekTokenIs(token.LPAREN) {
		name := ast.NewStrExpr(p.curToken, "")
		slotDir := ast.NewSlotDir(tok, name, p.file.Name, false)
		slotDir.SetIsDefault(true)
		return slotDir
	}

	p.nextToken() // skip "@slot"
	p.nextToken() // skip "("

	// Handle named @slot with name
	name := ast.NewStrExpr(p.curToken, p.curToken.Lit)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	dir := ast.NewSlotDir(tok, name, p.file.Name, false)
	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) dumpDir() ast.Chunk {
	tok := p.curToken // "@dump"

	var args []ast.Expression

	if !p.expectPeek(token.LPAREN) { // move to "("
		return ast.NewIllegalNode(tok)
	}

	args = p.expressionList(token.RPAREN)

	dir := ast.NewDumpDir(tok, args)
	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

// slots parses local slots inside of @component directive's body
func (p *Parser) slots(compName string) []ast.Chunk {
	var slots []ast.Chunk

	for p.curTokenIs(token.SLOT, token.SLOTIF) {
		slotName := ast.NewStrExpr(p.curToken, "")

		switch p.curToken.Type {
		case token.SLOT:
			slots = append(slots, p.localSlotDir(slotName, compName))
		case token.SLOTIF:
			slots = append(slots, p.slotifDir(slotName, compName))
		default:
			panic("Unknown slot token when parsing component slots")
		}

		for p.curTokenIs(token.TEXT) {
			p.nextToken() // skip whitespace
		}
	}

	return slots
}

func (p *Parser) localSlotDir(name *ast.StrExpr, compName string) ast.Chunk {
	dir := ast.NewSlotDir(p.curToken, name, compName, true)
	dir.SetIsDefault(!p.peekTokenIs(token.LPAREN))

	// When slot has a name @slot('name')
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken() // move to "(" from '@slot'
		p.nextToken() // skip "("

		name.Token = p.curToken
		name.Val = p.curToken.Lit

		if !p.expectPeek(token.RPAREN) { // move to ")"
			return p.illegalNode() // create an error
		}

		p.nextToken() // skip ")"
	} else {
		p.nextToken() // skip "@slot"
	}

	if !p.curTokenIs(token.END) {
		dir.SetBlock(p.block())
	}

	dir.SetEndPosition(p.curToken.Pos)

	p.nextToken() // skip "@end"

	return dir
}

func (p *Parser) slotifDir(name *ast.StrExpr, compName string) ast.Chunk {
	dir := ast.NewSlotifDir(p.curToken, name, compName)

	if illegal := p.slotifDirHeader(dir, name); illegal != nil { // skips ")"
		return illegal
	}

	// Handle empty slotif body
	if p.curTokenIs(token.END) {
		dir.SetEndPosition(p.curToken.Pos)
		p.nextToken() // skip "@end"
		return dir
	}

	dir.SetBlock(p.block())

	if !p.curTokenIs(token.END) {
		return p.illegalNodeUntil(token.END)
	}

	dir.SetEndPosition(p.curToken.Pos)

	p.nextToken() // skip "@end"

	return dir
}

func (p *Parser) slotifDirHeader(dir *ast.SlotifDir, name *ast.StrExpr) *ast.IllegalNode {
	if !p.expectPeek(token.LPAREN) { // move from "@slotif" to "("
		p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip "("

	dir.Cond = p.expression(LOWEST)

	// When slot has name
	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // move to ","
		p.nextToken() // skip ","

		name.Token = p.curToken
		name.Val = p.curToken.Lit
	} else {
		dir.SetIsDefault(true)
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip ")"

	return nil
}

func (p *Parser) reserveDir() ast.Chunk {
	dir := ast.NewReserveDir(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.nextToken() // skip "("

	dir.Name = ast.NewStrExpr(p.curToken, p.curToken.Lit)

	// Handle when has second argument (fallback value) after comma
	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // move to "," from string
		p.nextToken() // move to expression from ","
		dir.Fallback = p.expression(LOWEST)
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	dir.SetEndPosition(p.curToken.Pos)

	// Check for duplicate reserve statements
	if _, ok := p.reserves[dir.Name.Val]; ok {
		p.newError(dir.Token.Pos, fail.ErrDuplicateReserves, dir.Name.Val, p.file.Abs)
		return nil
	}

	p.reserves[dir.Name.Val] = dir

	return dir
}

func (p *Parser) insertDir() ast.Chunk {
	dir := ast.NewInsertDir(p.curToken, p.file.Abs)

	illegal, done := p.insertDirHeader(dir) // moves to ")"
	if illegal != nil {
		return illegal
	}

	if done {
		return dir
	}

	p.nextToken() // skip ")"

	// Handle empty block
	if p.curTokenIs(token.END) {
		dir.SetEndPosition(p.curToken.Pos)
		return dir
	}

	dir.Block = p.block()

	// skip block and move to @end
	if !p.curTokenIs(token.END) {
		return p.illegalNodeUntil(token.END)
	}

	dir.SetEndPosition(p.curToken.Pos)

	p.inserts[dir.Name.Val] = dir

	return dir
}

func (p *Parser) insertDirHeader(dir *ast.InsertDir) (*ast.IllegalNode, bool) {
	done := false
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode(), done
	}

	p.nextToken() // skip "("

	dir.Name = ast.NewStrExpr(p.curToken, p.curToken.Lit)

	if ok := p.checkDuplicateInserts(dir); ok {
		return nil, done
	}

	// Handle inline @insert without block
	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip insert name
		p.nextToken() // skip ","
		dir.Argument = p.expression(LOWEST)

		if !p.expectPeek(token.RPAREN) { // move to ")"
			return p.illegalNodeUntil(token.RBRACE), done
		}

		dir.SetEndPosition(p.curToken.Pos)

		p.inserts[dir.Name.Val] = dir
		done = true

		return nil, done
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END), done
	}

	return nil, done
}

func (p *Parser) checkDuplicateInserts(insertDir *ast.InsertDir) bool {
	if _, hasDuplicate := p.inserts[insertDir.Name.Val]; hasDuplicate {
		p.newError(
			insertDir.Name.Pos(),
			fail.ErrDuplicateInserts,
			insertDir.Name.Val,
		)

		return true
	}

	return false
}

func (p *Parser) indexExpr(left ast.Expression) ast.Expression {
	expr := ast.NewIndexExpr(*left.Tok(), left)

	p.nextToken() // skip "["

	expr.Index = p.expression(LOWEST)

	if !p.expectPeek(token.RBRACKET) { // move to "]"
		return p.illegalNode()
	}

	expr.SetEndPosition(p.curToken.Pos)

	return expr
}

func (p *Parser) incStmt(left ast.Expression) ast.Statement {
	return ast.NewIncStmt(p.curToken, left)
}

func (p *Parser) decStmt(left ast.Expression) ast.Statement {
	return ast.NewDecStmt(p.curToken, left)
}

func (p *Parser) dotExpr(left ast.Expression) ast.Expression {
	expr := ast.NewDotExpr(*left.Tok(), left)

	if !p.expectPeek2(token.IDENT, p.curToken.Pos) { // skip "." and move to identifier
		return nil
	}

	if p.peekTokenIs(token.LPAREN) {
		return p.callExpr(left)
	}

	expr.Key = ast.NewIdentExpr(p.curToken, p.curToken.Lit)

	return expr
}

func (p *Parser) globalCallExpr(ident *ast.IdentExpr) ast.Expression {
	expr := ast.NewGlobalCallExpr(p.curToken, ident)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	expr.Arguments = p.expressionList(token.RPAREN)
	expr.SetEndPosition(p.curToken.Pos)

	return expr
}

func (p *Parser) callExpr(receiver ast.Expression) ast.Expression {
	ident := ast.NewIdentExpr(p.curToken, p.curToken.Lit)
	expr := ast.NewCallExpr(p.curToken, receiver, ident)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	expr.Arguments = p.expressionList(token.RPAREN)
	expr.SetEndPosition(p.curToken.Pos)

	return expr
}

func (p *Parser) infixExpr(left ast.Expression) ast.Expression {
	expr := ast.NewInfixExpr(*left.Tok(), left, p.curToken.Lit)

	precedence := precedences[p.curToken.Type]

	if p.peekTokenIs(token.RBRACES) {
		pos := left.Pos()
		pos.EndCol = p.curToken.Pos.EndCol
		p.newError(pos, fail.ErrExpectedExpression)
		return nil
	}

	p.nextToken() // skip operator

	expr.Right = p.expression(precedence)
	expr.SetEndPosition(p.curToken.Pos)

	return expr
}

func (p *Parser) ternaryExpr(left ast.Expression) ast.Expression {
	expr := ast.NewTernaryExpr(*left.Tok(), left)

	p.nextToken() // skip "?"

	expr.IfExpr = p.expression(TERNARY)

	if !p.expectPeek(token.COLON) { // move to ":"
		return p.illegalNode()
	}

	p.nextToken() // skip ":"

	expr.ElseExpr = p.expression(LOWEST)
	expr.SetEndPosition(p.curToken.Pos)

	return expr
}

func (p *Parser) ifDir() ast.Chunk {
	dir := ast.NewIfDir(p.curToken)

	if ok := p.ifDirHeader(dir); !ok { // skips ")"
		return nil
	}

	dir.IfBlock = p.block()
	if p.curTokenIs(token.END) {
		dir.SetEndPosition(p.curToken.Pos)
		return dir
	}

	for p.curTokenIs(token.ELSEIF) {
		elseifDir, illegal := p.elseifDir()
		if illegal != nil {
			return illegal
		}

		dir.ElseifDirs = append(dir.ElseifDirs, elseifDir)
	}

	dir.ElseBlock = p.elseBlock()

	if !p.curTokenIs(token.END) { // move to "@end"
		return p.illegalNode()
	}

	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) ifDirHeader(ifDir *ast.IfDir) bool {
	if !p.expectPeek2(token.LPAREN, ifDir.Pos()) { // move to "("
		return false
	}

	p.nextToken() // skip "("

	ifDir.Cond = p.expression(LOWEST)

	if !p.expectPeek2(token.RPAREN, ifDir.Pos()) { // move to ")"
		return false
	}

	p.nextToken() // skip ")"

	return true
}

func (p *Parser) elseifDir() (*ast.ElseIfDir, *ast.IllegalNode) {
	dir := ast.NewElseIfDir(p.curToken)

	p.nextToken() // skip "@elseif"
	p.nextToken() // skip "("

	dir.Cond = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return nil, p.illegalNode()
	}

	p.nextToken() // skip ")"

	dir.Block = p.block()
	dir.SetEndPosition(p.curToken.Pos)

	return dir, nil
}

func (p *Parser) elseBlock() *ast.Block {
	if p.curTokenIs(token.ELSE) {
		p.nextToken() // skip "@else"
	}

	if p.curTokenIs(token.ELSE, token.END) {
		return nil
	}

	block := p.block()

	if p.peekTokenIs(token.ELSEIF) {
		p.newError(p.peekToken.Pos, fail.ErrElseifCannotFollowElse)
		return nil
	}

	block.SetEndPosition(p.curToken.Pos)

	return block
}

func (p *Parser) forDir() ast.Chunk {
	dir := ast.NewForDir(p.curToken)

	if illegal := p.forDirHeader(dir); illegal != nil { // skips ")"
		return illegal
	}

	dir.Block = p.block()

	if p.curTokenIs(token.END) {
		dir.SetEndPosition(p.curToken.Pos)
		return dir
	}

	dir.ElseBlock = p.elseBlock()

	if !p.curTokenIs(token.END) {
		return p.illegalNodeUntil(token.END)
	}

	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) forDirHeader(dir *ast.ForDir) *ast.IllegalNode {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.END)
	}

	// Parse Init
	if !p.peekTokenIs(token.SEMI) {
		p.nextToken() // move to first token of init statement
		left := p.expression(LOWEST)
		p.nextToken() // move to =/++/--
		dir.Init = p.statement(left)
	}

	if !p.expectPeek(token.SEMI) { // move to ";"
		return p.illegalNodeUntil(token.END)
	}

	// Parse Condition
	if !p.peekTokenIs(token.SEMI) {
		p.nextToken() // skip ";"
		dir.Cond = p.expression(LOWEST)
	}

	if !p.expectPeek(token.SEMI) { // move to ";"
		return p.illegalNodeUntil(token.END)
	}

	// Parse Post statement
	if !p.peekTokenIs(token.RPAREN) {
		p.nextToken() // skip ";"
		left := p.expression(LOWEST)
		p.nextToken() // move to ++/--

		if p.curTokenIs(token.RPAREN) {
			p.newError(left.Pos(), fail.ErrForLoopExpectStmt, left.String())
			return nil
		}

		dir.Post = p.statement(left)
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip ")"

	return nil
}

func (p *Parser) eachDir() ast.Chunk {
	dir := ast.NewEachDir(p.curToken)

	if illegal := p.eachDirHeader(dir); illegal != nil { // skips ")"
		return illegal
	}

	dir.Block = p.block()

	if p.curTokenIs(token.END) {
		dir.SetEndPosition(p.curToken.Pos)
		return dir
	}

	dir.ElseBlock = p.elseBlock()

	if !p.curTokenIs(token.END) {
		return p.illegalNodeUntil(token.END)
	}

	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) eachDirHeader(dir *ast.EachDir) *ast.IllegalNode {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip "("

	dir.Var = ast.NewIdentExpr(p.curToken, p.curToken.Lit)

	if !p.expectPeek(token.IN) { // move to "in"
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip "in"

	dir.Arr = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.nextToken() // skip ")"

	return nil
}

func (p *Parser) block() *ast.Block {
	if p.curTokenIs(token.ELSE, token.ELSEIF, token.END) {
		return nil
	}

	block := ast.NewBlock(p.curToken)
	block.SetEndPosition(p.curToken.Pos)

	for !p.curTokenIs(token.END) && !p.curTokenIs(token.EOF) {
		chunk := p.chunk()
		block.SetEndPosition(p.curToken.Pos)

		if chunk != nil {
			block.Chunks = append(block.Chunks, chunk)
		}

		if p.peekTokenIs(token.ELSE, token.ELSEIF, token.END) {
			p.nextToken() // skip chunk
			break
		}

		p.nextToken() // skip chunk
	}

	return block
}

func (p *Parser) statement(left ast.Expression) ast.Statement {
	switch p.curToken.Type {
	case token.ASSIGN:
		return p.assignStmt(left)
	case token.INC:
		return p.incStmt(left)
	case token.DEC:
		return p.decStmt(left)
	}

	return nil
}

func (p *Parser) expression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		return p.illegalNode()
	}

	leftExpr := prefix()

	for !p.peekTokenIs(token.RBRACES, token.SEMI, token.RPAREN) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return leftExpr
		}

		p.nextToken()

		leftExpr = infix(leftExpr)
	}

	if leftExpr != nil {
		leftExpr.SetEndPosition(p.curToken.Pos)
	}

	return leftExpr
}

func (p *Parser) prefixExpr() ast.Expression {
	expr := ast.NewPrefixExpr(p.curToken, p.curToken.Lit)

	p.nextToken() // skip prefix operator

	expr.Right = p.expression(PREFIX)
	expr.SetEndPosition(p.curToken.Pos)

	return expr
}

func (p *Parser) groupedExpr() ast.Expression {
	exprTok := p.curToken

	p.nextToken() // skip "("

	expr := p.expression(LOWEST)
	expr.SetTok(exprTok)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	return expr
}

func (p *Parser) expressionList(endTok token.TokenType) []ast.Expression {
	var exprs []ast.Expression

	if p.peekTokenIs(endTok) {
		p.nextToken() // skip endTok token
		return exprs
	}

	if p.peekTokenIs(token.END) {
		exprs = append(exprs, p.illegalNode())
		return exprs
	}

	p.nextToken() // move to first expression

	exprs = append(exprs, p.expression(LOWEST))

	for p.peekTokenIs(token.COMMA) && !p.curTokenIs(token.EOF) {
		p.nextToken() // move to ","

		// break when has a trailing comma
		if p.peekTokenIs(endTok) {
			break
		}

		p.nextToken() // skip ","
		exprs = append(exprs, p.expression(LOWEST))
	}

	if !p.expectPeek(endTok) { // move to endTok
		exprs = append(exprs, ast.NewIllegalNode(p.curToken))
		return exprs
	}

	return exprs
}

func (p *Parser) illegalNode() *ast.IllegalNode {
	p.newError(
		p.curToken.Pos,
		fail.ErrIllegalToken,
		p.curToken.Lit,
	)

	return ast.NewIllegalNode(p.curToken)
}

func (p *Parser) illegalNodeUntil(tok token.TokenType) *ast.IllegalNode {
	for !p.curTokenIs(tok) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	return p.illegalNode()
}
