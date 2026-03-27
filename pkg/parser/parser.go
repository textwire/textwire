package parser

import (
	"strconv"

	"github.com/textwire/textwire/v4/pkg/file"
	"github.com/textwire/textwire/v4/pkg/position"
	"github.com/textwire/textwire/v4/pkg/token"

	"slices"

	"github.com/textwire/textwire/v4/pkg/ast"
	"github.com/textwire/textwire/v4/pkg/fail"
	"github.com/textwire/textwire/v4/pkg/lexer"
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

	components []*ast.CompDir
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
		components: []*ast.CompDir{},
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
	case token.PASS, token.PASSIF:
		return p.passDir()
	case token.DUMP:
		return p.dumpDir()
	case token.BREAK:
		return ast.NewBreakDir(p.curToken)
	case token.CONTINUE:
		return ast.NewContinueDir(p.curToken)
	}

	return p.illegal()
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
		if segment := p.segment(); segment != nil {
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

// segment parses individual segment
func (p *Parser) segment() ast.Segment {
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

	if stmt == nil {
		return nil
	}

	return stmt.(ast.Segment)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) expectNonEmptyNameOn(node ast.Node) bool {
	if p.curToken.Lit == "" {
		p.newError(p.curToken.Pos, fail.ErrNameCannotBeEmpty, node.Tok().Lit)
		return false
	}
	return true
}

func (p *Parser) expectType(expectType token.TokenType) bool {
	if p.curToken.Type == expectType {
		return true
	}

	p.newError(
		p.curToken.Pos,
		fail.ErrWrongTokenType,
		token.String(expectType),
		token.String(p.curToken.Type),
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

func (p *Parser) expectPeek(tok token.TokenType) bool {
	if p.peekTokenIs(tok) {
		p.nextToken()
		return true
	}

	got := p.peekToken.Lit
	if len(got) > 100 {
		got = got[:100] + "..."
	}

	p.newError(p.peekToken.Pos, fail.ErrWrongPeekToken, token.String(tok), got)

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
		p.newError(p.curToken.Pos, fail.ErrExpectExprAfter, token.String(token.ASSIGN))
		return nil
	}

	stmt.Right = p.expression(LOWEST)
	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) useDir() ast.Chunk {
	useDir := ast.NewUseDir(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegal()
	}

	p.nextToken() // skip "("

	if !p.expectType(token.STR) {
		return p.illegal()
	}

	if !p.expectNonEmptyNameOn(useDir) {
		return p.illegal()
	}

	useDir.Name = ast.NewStrExpr(
		p.curToken,
		file.ReplacePathAlias(p.curToken.Lit, file.PathAliasUse),
	)

	if p._useDir != nil {
		p.newError(useDir.Pos(), fail.ErrOnlyOneUseDir)
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegal()
	}

	useDir.SetEndPosition(p.curToken.Pos)

	p._useDir = useDir

	return useDir
}

func (p *Parser) breakifDir() ast.Chunk {
	dir := ast.NewBreakIfDir(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegal()
	}

	p.nextToken() // skip "("

	dir.Cond = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegal()
	}

	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) continueifDir() ast.Chunk {
	dir := ast.NewContinueIfDir(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegal()
	}

	p.nextToken() // skip "("

	dir.Cond = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegal()
	}

	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) compDir() ast.Chunk {
	compDir := ast.NewCompDir(p.curToken)

	if illegal := p.compDirHeader(compDir); illegal != nil { // moves to ")"
		return illegal
	}

	p.nextToken() // skip ")"

	if p.curTokenIs(token.END) {
		return p.endCompDir(compDir)
	}

	block := p.block()

	// Extract @pass from block and map them to a component
	for _, chunk := range block.AllChunks() {
		passDir, ok := chunk.(*ast.PassDir)
		if !ok {
			continue
		}
		compDir.Passes = append(compDir.Passes, passDir)
	}

	defaultPass := ast.NewPassDir(*block.Tok(), ast.NewStrExpr(*block.Tok(), ""))
	defaultPass.CompName = compDir.Name.Val
	defaultPass.Block = block
	compDir.Passes = append(compDir.Passes, defaultPass)

	return p.endCompDir(compDir)
}

func (p *Parser) endCompDir(compDir *ast.CompDir) *ast.CompDir {
	p.components = append(p.components, compDir)
	compDir.SetEndPosition(p.curToken.Pos)
	return compDir
}

func (p *Parser) compDirHeader(compDir *ast.CompDir) *ast.Illegal {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalUntil(token.RPAREN)
	}

	p.nextToken() // skip "("

	if !p.expectType(token.STR) {
		return p.illegal()
	}

	if !p.expectNonEmptyNameOn(compDir) {
		return p.illegal()
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
		return p.illegalUntil(token.TEXT)
	}

	return nil
}

// slotDir parses an @slot inside a component file.
func (p *Parser) slotDir() ast.Chunk {
	tok := p.curToken // "@slot"

	// Handle default @slot without name
	if !p.peekTokenIs(token.LPAREN) {
		name := ast.NewStrExpr(p.curToken, "")
		slotDir := ast.NewSlotDir(tok, name, p.file.Name)
		return slotDir
	}

	p.nextToken() // skip "@slot"
	p.nextToken() // skip "("

	// Handle named @slot with name
	name := ast.NewStrExpr(p.curToken, p.curToken.Lit)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalUntil(token.END)
	}

	dir := ast.NewSlotDir(tok, name, p.file.Name)
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

func (p *Parser) passDir() ast.Chunk {
	passDir := ast.NewPassDir(p.curToken, nil)
	hasCondition := p.curToken.Type == token.PASSIF

	p.nextToken() // move to "("
	p.nextToken() // skip "("

	if hasCondition {
		passDir.Cond = p.expression(LOWEST)
	}

	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // move to ","
	}

	if p.peekTokenIs(token.STR) {
		p.nextToken() // move to string
	}

	passDir.Name = ast.NewStrExpr(p.curToken, p.curToken.Lit)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegal()
	}

	p.nextToken() // skip ")"

	if !p.curTokenIs(token.END) {
		passDir.Block = p.block()
	}

	passDir.SetEndPosition(p.curToken.Pos)

	return passDir
}

func (p *Parser) reserveDir() ast.Chunk {
	reserveDir := ast.NewReserveDir(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegal()
	}

	p.nextToken() // skip "("

	if !p.expectType(token.STR) {
		return p.illegal()
	}

	if !p.expectNonEmptyNameOn(reserveDir) {
		return p.illegal()
	}

	reserveDir.Name = ast.NewStrExpr(p.curToken, p.curToken.Lit)
	reserveDir.Fallback = p.reserveDirFallback()

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegal()
	}

	reserveDir.SetEndPosition(p.curToken.Pos)

	// Check for duplicate reserve statements
	if _, ok := p.reserves[reserveDir.Name.Val]; ok {
		p.newError(
			reserveDir.Name.Pos(),
			fail.ErrDuplicateReserves,
			reserveDir.Name.Val,
			p.file.Abs,
		)
		return nil
	}

	p.reserves[reserveDir.Name.Val] = reserveDir

	return reserveDir
}

// reserveDirFallback parses the second argument (fallback value) after comma
// for @reserve statement. If there are no expression, set Fallback to Empty.
func (p *Parser) reserveDirFallback() ast.Expression {
	if !p.peekTokenIs(token.COMMA) {
		return ast.NewEmpty(p.curToken.Pos)
	}

	p.nextToken() // move to "," from string
	p.nextToken() // move to expression from ","

	return p.expression(LOWEST)
}

func (p *Parser) insertDir() ast.Chunk {
	insertDir := ast.NewInsertDir(p.curToken, p.file.Abs)

	illegal, done := p.insertDirHeader(insertDir) // moves to ")"
	if illegal != nil {
		return illegal
	}

	if done {
		return insertDir
	}

	p.nextToken() // skip ")"

	// Handle empty block
	if p.curTokenIs(token.END) {
		insertDir.SetEndPosition(p.curToken.Pos)
		return insertDir
	}

	insertDir.Block = p.block()

	// skip block and move to @end
	if !p.curTokenIs(token.END) {
		p.newError(insertDir.Pos(), fail.ErrInsertMustHaveContent, insertDir.Name.Val)
		return nil
	}

	insertDir.SetEndPosition(p.curToken.Pos)

	p.inserts[insertDir.Name.Val] = insertDir

	return insertDir
}

func (p *Parser) insertDirHeader(insertDir *ast.InsertDir) (*ast.Illegal, bool) {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegal(), false
	}

	p.nextToken() // skip "("

	if !p.expectType(token.STR) {
		return p.illegal(), false
	}

	if !p.expectNonEmptyNameOn(insertDir) {
		return p.illegal(), false
	}

	insertDir.Name = ast.NewStrExpr(p.curToken, p.curToken.Lit)

	if ok := p.hasDuplicateInserts(insertDir); ok {
		return p.illegal(), false
	}

	illegal, done := p.insertDirArgument(insertDir)

	if illegal != nil {
		return illegal, false
	}

	if done {
		return nil, true
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegal(), false
	}

	return nil, false
}

func (p *Parser) insertDirArgument(insertDir *ast.InsertDir) (*ast.Illegal, bool) {
	if !p.peekTokenIs(token.COMMA) {
		insertDir.Argument = ast.NewEmpty(p.curToken.Pos)
		return nil, false
	}

	p.nextToken() // skip insert name
	p.nextToken() // skip ","
	insertDir.Argument = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegal(), false
	}

	insertDir.SetEndPosition(p.curToken.Pos)

	p.inserts[insertDir.Name.Val] = insertDir

	return nil, true
}

func (p *Parser) hasDuplicateInserts(insertDir *ast.InsertDir) bool {
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
		return p.illegal()
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

	if !p.expectPeek(token.IDENT) { // skip "." and move to identifier
		return p.illegal()
	}

	if p.peekTokenIs(token.LPAREN) {
		return p.callExpr(left)
	}

	expr.Key = ast.NewIdentExpr(p.curToken, p.curToken.Lit)

	return expr
}

func (p *Parser) globalCallExpr(ident *ast.IdentExpr) ast.Expression {
	name := ast.GlobalFuncName(ident.Name)

	rules, exists := ast.GlobalFunctions[name]
	if !exists {
		p.newError(p.peekToken.Pos, fail.ErrGlobalFuncMissing, ident.Name)
		return nil
	}

	expr := ast.NewGlobalCallExpr(p.curToken, name)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegal()
	}

	expr.Arguments = p.expressionList(token.RPAREN)
	expr.SetEndPosition(p.curToken.Pos)
	argsLen := len(expr.Arguments)

	if argsLen < rules.Min {
		p.newError(expr.Pos(), fail.ErrGlobalFuncFewArgs, name, rules.Min, argsLen)
		return nil
	}

	if argsLen > rules.Max {
		p.newError(expr.Pos(), fail.ErrGlobalFuncLotsOfArgs, name, rules.Max, argsLen)
		return nil
	}

	return expr
}

func (p *Parser) callExpr(receiver ast.Expression) ast.Expression {
	ident := ast.NewIdentExpr(p.curToken, p.curToken.Lit)
	expr := ast.NewCallExpr(p.curToken, receiver, ident)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegal()
	}

	expr.Arguments = p.expressionList(token.RPAREN)
	expr.SetEndPosition(p.curToken.Pos)

	return expr
}

func (p *Parser) infixExpr(left ast.Expression) ast.Expression {
	opTok := p.curToken
	expr := ast.NewInfixExpr(*left.Tok(), left, opTok.Lit)

	precedence := precedences[opTok.Type]

	if p.peekTokenIs(token.RBRACES) {
		pos := left.Pos()
		pos.EndCol = opTok.Pos.EndCol
		p.newError(pos, fail.ErrExpectExprAfter, opTok.Lit)
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
		return p.illegal()
	}

	p.nextToken() // skip ":"

	expr.ElseExpr = p.expression(LOWEST)
	expr.SetEndPosition(p.curToken.Pos)

	return expr
}

func (p *Parser) ifDir() ast.Chunk {
	dir := ast.NewIfDir(p.curToken)

	if illegal := p.ifDirHeader(dir); illegal != nil { // skips ")"
		return illegal
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
		return p.illegal()
	}

	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) ifDirHeader(ifDir *ast.IfDir) *ast.Illegal {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegal()
	}

	p.nextToken() // skip "("

	ifDir.Cond = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegal()
	}

	p.nextToken() // skip ")"

	return nil
}

func (p *Parser) elseifDir() (*ast.ElseIfDir, *ast.Illegal) {
	dir := ast.NewElseIfDir(p.curToken)

	p.nextToken() // skip "@elseif"
	p.nextToken() // skip "("

	dir.Cond = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return nil, p.illegal()
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
		return p.illegalUntil(token.END)
	}

	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) forDirHeader(forDir *ast.ForDir) *ast.Illegal {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalUntil(token.END)
	}

	forDir.Init = p.parserForDirInit()

	if !p.expectPeek(token.SEMI) { // move to ";"
		return p.illegalUntil(token.END)
	}

	forDir.Cond = p.parserForDirCond()

	if !p.expectPeek(token.SEMI) { // move to ";"
		return p.illegalUntil(token.END)
	}

	post := p.parseForDirPost()
	if post == nil {
		return nil
	}

	forDir.Post = post

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalUntil(token.END)
	}

	p.nextToken() // skip ")"

	return nil
}

func (p *Parser) parserForDirInit() ast.Statement {
	if p.peekTokenIs(token.SEMI) {
		return ast.NewEmpty(p.curToken.Pos)
	}

	p.nextToken() // move to first token of init statement
	left := p.expression(LOWEST)
	p.nextToken() // move to =/++/--

	return p.statement(left)
}

func (p *Parser) parserForDirCond() ast.Expression {
	if p.peekTokenIs(token.SEMI) {
		return ast.NewEmpty(p.curToken.Pos)
	}

	p.nextToken() // skip ";"

	return p.expression(LOWEST)
}

func (p *Parser) parseForDirPost() ast.Statement {
	if p.peekTokenIs(token.RPAREN) {
		return ast.NewEmpty(p.curToken.Pos)
	}

	p.nextToken() // skip ";"
	left := p.expression(LOWEST)
	p.nextToken() // move to ++/--

	if p.curTokenIs(token.RPAREN) {
		p.newError(left.Pos(), fail.ErrForLoopExpectStmt, left.String())
		return nil
	}

	return p.statement(left)
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
		return p.illegalUntil(token.END)
	}

	dir.SetEndPosition(p.curToken.Pos)

	return dir
}

func (p *Parser) eachDirHeader(dir *ast.EachDir) *ast.Illegal {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegal()
	}

	p.nextToken() // skip "("

	if !p.expectType(token.IDENT) {
		return p.illegal()
	}

	dir.Var = ast.NewIdentExpr(p.curToken, p.curToken.Lit)

	if !p.expectPeek(token.IN) { // move to "in"
		return p.illegal()
	}

	p.nextToken() // skip "in"

	dir.Arr = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegal()
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
		return p.illegal()
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
		return p.illegal()
	}

	return expr
}

func (p *Parser) expressionList(endTok token.TokenType) []ast.Expression {
	firstTok := p.curToken
	var exprs []ast.Expression

	if p.peekTokenIs(endTok) {
		p.nextToken() // skip endTok token
		return exprs
	}

	if p.peekTokenIs(token.END) {
		p.newError(p.peekToken.Pos, fail.ErrExpectExprAfter, firstTok.Lit)
		return nil
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

func (p *Parser) illegal() *ast.Illegal {
	p.newError(
		p.curToken.Pos,
		fail.ErrIllegalToken,
		p.curToken.Lit,
	)

	return ast.NewIllegalNode(p.curToken)
}

func (p *Parser) illegalUntil(tok token.TokenType) *ast.Illegal {
	for !p.curTokenIs(tok) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	return p.illegal()
}
