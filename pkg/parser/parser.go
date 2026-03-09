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

	p.next() // fill curToken
	p.next() // fill peekToken

	// Prefix operators
	p.prefixParseFns = map[token.TokenType]prefixParseFn{}

	p.registerPrefix(token.IDENT, p.ident)
	p.registerPrefix(token.INT, p.intLit)
	p.registerPrefix(token.FLOAT, p.floatLit)
	p.registerPrefix(token.STR, p.strLit)
	p.registerPrefix(token.NIL, p.nilLit)
	p.registerPrefix(token.TRUE, p.boolLit)
	p.registerPrefix(token.FALSE, p.boolLit)
	p.registerPrefix(token.SUB, p.prefixExp)
	p.registerPrefix(token.NOT, p.prefixExp)
	p.registerPrefix(token.LPAREN, p.groupedExp)
	p.registerPrefix(token.LBRACKET, p.arrLit)
	p.registerPrefix(token.LBRACE, p.objLit)

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
			p.next()
			continue
		}

		if stmt == nil {
			p.next() // skip to next token
			continue
		}

		prog.Statements = append(prog.Statements, stmt)

		p.next() // skip to next token
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
	case token.TEXT:
		return p.textStmt()
	case token.LBRACES, token.SEMI:
		return p.embeddedCode()
	case token.IF:
		return p.ifStmt()
	case token.FOR:
		return p._for()
	case token.EACH:
		return p.eachStmt()
	case token.USE:
		return p.useStmt()
	case token.RESERVE:
		return p.reserveStmt()
	case token.INSERT:
		return p.insertStmt()
	case token.BREAKIF:
		return p.breakifStmt()
	case token.CONTINUEIF:
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
	case token.SLOTIF:
		p.newError(p.curToken.ErrorLine(), fail.ErrSlotifPosition)
		return nil
	default:
		return p.illegalNode()
	}
}

func (p *Parser) embeddedCode() ast.Statement {
	p.next() // skip "{{" or ";" or "("

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.ErrorLine(), fail.ErrEmptyBraces)
		return nil
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

func (p *Parser) peekIs(tokens ...token.TokenType) bool {
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
	if p.peekIs(tok) {
		p.next()
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

func (p *Parser) next() {
	p.curToken = p.peekToken
	p.peekToken = p.l.Next()
}

func (p *Parser) ident() ast.Expression {
	ident := ast.NewIdent(p.curToken, p.curToken.Lit)
	if p.peekIs(token.LPAREN) {
		return p.globalCallExp(ident)
	}

	return ident
}

func (p *Parser) intLit() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Lit, 10, 64)
	if err != nil {
		p.newError(
			p.curToken.ErrorLine(),
			fail.ErrCouldNotParseAs,
			p.curToken.Lit,
			"INT",
		)
		return nil
	}

	return ast.NewIntLit(p.curToken, val)
}

func (p *Parser) floatLit() ast.Expression {
	val, err := strconv.ParseFloat(p.curToken.Lit, 64)
	if err != nil {
		p.newError(
			p.curToken.ErrorLine(),
			fail.ErrCouldNotParseAs,
			p.curToken.Lit,
			"FLOAT",
		)
		return nil
	}

	return ast.NewFloatLit(p.curToken, val)
}

func (p *Parser) nilLit() ast.Expression {
	return ast.NewNilLit(p.curToken)
}

func (p *Parser) strLit() ast.Expression {
	return ast.NewStrLit(p.curToken, p.curToken.Lit)
}

func (p *Parser) boolLit() ast.Expression {
	return ast.NewBoolLit(p.curToken, p.curTokenIs(token.TRUE))
}

func (p *Parser) arrLit() ast.Expression {
	arr := ast.NewArrLit(p.curToken)
	arr.Elements = p.expressionList(token.RBRACKET)
	arr.SetEndPosition(p.curToken.Pos)

	return arr
}

func (p *Parser) objLit() ast.Expression {
	obj := ast.NewObjLit(p.curToken)

	obj.Pairs = map[string]ast.Expression{}

	p.next() // skip "{"

	if p.curTokenIs(token.RBRACE) {
		obj.SetEndPosition(p.curToken.Pos)
		return obj
	}

	for !p.curTokenIs(token.RBRACE) {
		key := p.curToken.Lit

		if p.peekIs(token.COLON) {
			p.next() // move to ":"
			p.next() // skip to ":"

			obj.Pairs[key] = p.expression(LOWEST)
		} else {
			obj.Pairs[key] = p.expression(LOWEST)
		}

		if p.peekIs(token.RBRACE) {
			p.next() // skip "}"
			break
		}

		if p.peekIs(token.COMMA) {
			p.next() // move to ","
			p.next() // skip ","
		}
	}

	obj.SetEndPosition(p.curToken.Pos)

	return obj
}

func (p *Parser) textStmt() ast.Statement {
	return ast.NewTextStmt(p.curToken)
}

func (p *Parser) assignStmt(left ast.Expression) ast.Statement {
	stmt := ast.NewAssignStmt(*left.Tok(), left)

	if !p.expectPeek(token.ASSIGN) { // move to "="
		return p.illegalNode()
	}

	p.next() // skip "="

	if p.curTokenIs(token.RBRACES) {
		p.newError(p.curToken.ErrorLine(), fail.ErrExpectedExpression)
		return nil
	}

	stmt.Right = p.expression(LOWEST)
	stmt.SetEndPosition(p.curToken.Pos)

	if p.peekIs(token.RBRACES) {
		p.next() // move to '}}'
	}

	return stmt
}

func (p *Parser) useStmt() ast.Statement {
	stmt := ast.NewUseStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.next() // skip "("

	if p.curToken.Type != token.STR {
		p.newError(
			p.curToken.ErrorLine(),
			fail.ErrUseStmtFirstArgStr,
			token.String(p.curToken.Type),
		)
	}

	if p.curToken.Lit == "" {
		p.newError(p.curToken.ErrorLine(), fail.ErrExpectedUseName)
	}

	stmt.Name = ast.NewStrLit(
		p.curToken,
		file.ReplacePathAlias(p.curToken.Lit, file.PathAliasUse),
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

	p.next() // skip "("

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

	p.next() // skip "("

	stmt.Condition = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) componentStmt() ast.Statement {
	stmt := ast.NewComponentStmt(p.curToken)

	if illegal := p.componentStmtHeader(stmt); illegal != nil {
		return illegal
	}

	if p.peekIs(token.SLOT, token.SLOTIF) {
		p.next() // skip whitespace
		if illegal := p.assignSlotsToComp(stmt); illegal != nil {
			return illegal
		}
	}

	p.components = append(p.components, stmt)

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) componentStmtHeader(stmt *ast.ComponentStmt) *ast.IllegalNode {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.RPAREN)
	}

	p.next() // skip "("

	if p.curToken.Lit == "" {
		p.newError(p.curToken.ErrorLine(), fail.ErrExpectedComponentName)
	}

	stmt.Name = ast.NewStrLit(
		p.curToken,
		file.ReplacePathAlias(p.curToken.Lit, file.PathAliasComp),
	)

	if p.peekIs(token.COMMA) {
		p.next() // move to ","
		p.next() // skip ","

		obj, ok := p.expression(LOWEST).(*ast.ObjLit)
		if !ok {
			p.newError(p.curToken.ErrorLine(), fail.ErrExpectedObjLit, p.curToken.Lit)
			return nil
		}

		stmt.Argument = obj
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.TEXT)
	}

	if p.peekIs(token.TEXT) && isWhitespace(p.peekToken.Lit) {
		p.next() // move to ")"
	}

	return nil
}

func (p *Parser) assignSlotsToComp(stmt *ast.ComponentStmt) ast.Statement {
	slots := p.slots(stmt.Name.Val)
	stmt.Slots = make([]ast.SlotCommand, len(slots))
	for i := range slots {
		slot, ok := slots[i].(ast.SlotCommand)
		if !ok {
			return slots[i]
		}

		stmt.Slots[i] = slot
	}

	return nil
}

// slotStmt parses an external slot statement inside a component file.
// Slots inside a @component are parsed by other function.
func (p *Parser) slotStmt() ast.Statement {
	tok := p.curToken // "@slot"

	// Handle default @slot without name
	if !p.peekIs(token.LPAREN) {
		name := ast.NewStrLit(p.curToken, "")
		stmt := ast.NewSlotStmt(tok, name, p.file.Name, false)
		stmt.SetIsDefault(true)
		return stmt
	}

	p.next() // skip "@slot"
	p.next() // skip "("

	// Handle named @slot with name
	name := ast.NewStrLit(p.curToken, p.curToken.Lit)

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

	for p.curTokenIs(token.SLOT, token.SLOTIF) {
		slotName := ast.NewStrLit(p.curToken, "")

		switch p.curToken.Type {
		case token.SLOT:
			slots = append(slots, p.localSlotStmt(slotName, compName))
		case token.SLOTIF:
			slots = append(slots, p.slotifStmt(slotName, compName))
		default:
			panic("Unknown slot token when parsing component slots")
		}

		for p.curTokenIs(token.TEXT) {
			p.next() // skip whitespace
		}
	}

	return slots
}

func (p *Parser) localSlotStmt(name *ast.StrLit, compName string) ast.Statement {
	stmt := ast.NewSlotStmt(p.curToken, name, compName, true)
	stmt.SetIsDefault(!p.peekIs(token.LPAREN))

	// When slot has a name @slot('name')
	if p.peekIs(token.LPAREN) {
		p.next() // move to "(" from '@slot'
		p.next() // skip "("

		name.Token = p.curToken
		name.Val = p.curToken.Lit

		if !p.expectPeek(token.RPAREN) { // move to ")"
			return p.illegalNode() // create an error
		}

		p.next() // skip ")"
	} else {
		p.next() // skip "@slot"
	}

	if !p.curTokenIs(token.END) {
		stmt.SetBlock(p.blockStmt())
	}

	p.next() // skip "@end"
	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) slotifStmt(name *ast.StrLit, compName string) ast.Statement {
	stmt := ast.NewSlotifStmt(p.curToken, name, compName)

	if illegal := p.slotifStmtHeader(stmt, name); illegal != nil { // skips ")"
		return illegal
	}

	// Handle empty slotif body
	if p.curTokenIs(token.END) {
		stmt.SetEndPosition(p.curToken.Pos)
		p.next() // skip "@end"
		return stmt
	}

	stmt.SetBlock(p.blockStmt())

	if !p.curTokenIs(token.END) {
		return p.illegalNodeUntil(token.END)
	}

	stmt.SetEndPosition(p.curToken.Pos)

	p.next() // skip "@end"

	return stmt
}

func (p *Parser) slotifStmtHeader(stmt *ast.SlotifStmt, name *ast.StrLit) *ast.IllegalNode {
	if !p.expectPeek(token.LPAREN) { // move from "@slotif" to "("
		p.illegalNodeUntil(token.END)
	}

	p.next() // skip "("

	stmt.Condition = p.expression(LOWEST)

	// When slot has name
	if p.peekIs(token.COMMA) {
		p.next() // move to ","
		p.next() // skip ","

		name.Token = p.curToken
		name.Val = p.curToken.Lit
	} else {
		stmt.SetIsDefault(true)
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.next() // skip ")"

	return nil
}

func (p *Parser) reserveStmt() ast.Statement {
	stmt := ast.NewReserveStmt(p.curToken)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	p.next() // skip "("

	stmt.Name = ast.NewStrLit(p.curToken, p.curToken.Lit)

	// Handle when has second argument (fallback value) after comma
	if p.peekIs(token.COMMA) {
		p.next() // move to "," from string
		p.next() // move to expression from ","
		stmt.Fallback = p.expression(LOWEST)
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	// Check for duplicate reserve statements
	if _, ok := p.reserves[stmt.Name.Val]; ok {
		p.newError(stmt.Token.ErrorLine(), fail.ErrDuplicateReserves, stmt.Name.Val, p.file.Abs)
		return nil
	}

	p.reserves[stmt.Name.Val] = stmt

	return stmt
}

func (p *Parser) insertStmt() ast.Statement {
	stmt := ast.NewInsertStmt(p.curToken, p.file.Abs)

	illegal, done := p.insertStmtHeader(stmt) // moves to ")"
	if illegal != nil {
		return illegal
	}

	if done {
		return stmt
	}

	p.next() // skip ")"

	// Handle empty block
	if p.curTokenIs(token.END) {
		stmt.SetEndPosition(p.curToken.Pos)
		return stmt
	}

	stmt.Block = p.blockStmt()

	// skip block and move to @end
	if !p.curTokenIs(token.END) {
		return p.illegalNodeUntil(token.END)
	}

	stmt.SetEndPosition(p.curToken.Pos)

	p.inserts[stmt.Name.Val] = stmt

	return stmt
}

func (p *Parser) insertStmtHeader(stmt *ast.InsertStmt) (*ast.IllegalNode, bool) {
	done := false
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode(), done
	}

	p.next() // skip "("

	stmt.Name = ast.NewStrLit(p.curToken, p.curToken.Lit)

	if ok := p.checkDuplicateInserts(stmt); ok {
		return nil, done
	}

	// Handle inline @insert without block
	if p.peekIs(token.COMMA) {
		p.next() // skip insert name
		p.next() // skip ","
		stmt.Argument = p.expression(LOWEST)

		if !p.expectPeek(token.RPAREN) { // move to ")"
			return p.illegalNodeUntil(token.RBRACE), done
		}

		stmt.SetEndPosition(p.curToken.Pos)

		p.inserts[stmt.Name.Val] = stmt
		done = true

		return nil, done
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END), done
	}

	return nil, done
}

func (p *Parser) checkDuplicateInserts(stmt *ast.InsertStmt) bool {
	if _, hasDuplicate := p.inserts[stmt.Name.Val]; hasDuplicate {
		p.newError(
			stmt.Token.ErrorLine(),
			fail.ErrDuplicateInserts,
			stmt.Name.Val,
		)

		return true
	}

	return false
}

func (p *Parser) indexExp(left ast.Expression) ast.Expression {
	exp := ast.NewIndexExp(*left.Tok(), left)

	p.next() // skip "["

	exp.Index = p.expression(LOWEST)

	if !p.expectPeek(token.RBRACKET) { // move to "]"
		return p.illegalNode()
	}

	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) incStmt(left ast.Expression) ast.Statement {
	p.next() // skip to "++"

	stmt := ast.NewIncStmt(p.curToken, left)

	if p.peekIs(token.RBRACES) {
		p.next() // skip "}}"
	}

	return stmt
}

func (p *Parser) decStmt(left ast.Expression) ast.Statement {
	p.next() // skip to "++"

	stmt := ast.NewDecStmt(p.curToken, left)

	if p.peekIs(token.RBRACES) {
		p.next() // skip "}}"
	}

	return stmt
}

func (p *Parser) dotExp(left ast.Expression) ast.Expression {
	exp := ast.NewDotExp(*left.Tok(), left)

	if p.peekIs(token.INT) {
		p.newError(p.curToken.ErrorLine(), fail.ErrObjKeyUseGet)
		return nil
	}

	if !p.expectPeek(token.IDENT) { // skip "." and move to identifier
		return p.illegalNode()
	}

	if p.peekIs(token.LPAREN) {
		return p.callExp(left)
	}

	exp.Key = ast.NewIdent(p.curToken, p.curToken.Lit)

	return exp
}

func (p *Parser) globalCallExp(ident *ast.Ident) ast.Expression {
	exp := ast.NewGlobalCallExp(p.curToken, ident)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	exp.Arguments = p.expressionList(token.RPAREN)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) callExp(receiver ast.Expression) ast.Expression {
	ident := ast.NewIdent(p.curToken, p.curToken.Lit)
	exp := ast.NewCallExp(p.curToken, receiver, ident)

	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNode()
	}

	exp.Arguments = p.expressionList(token.RPAREN)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) infixExp(left ast.Expression) ast.Expression {
	exp := ast.NewInfixExp(*left.Tok(), left, p.curToken.Lit)

	precedence := precedences[p.curToken.Type]

	p.next() // skip operator

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

	p.next() // skip "?"

	exp.IfBlock = p.expression(TERNARY)

	if !p.expectPeek(token.COLON) { // move to ":"
		return p.illegalNode()
	}

	p.next() // skip ":"

	exp.ElseBlock = p.expression(LOWEST)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) ifStmt() ast.Statement {
	stmt := ast.NewIfStmt(p.curToken)

	if illegal := p.ifStmtHeader(stmt); illegal != nil { // skips ")"
		return illegal
	}

	stmt.IfBlock = p.blockStmt()

	if p.curTokenIs(token.END) {
		stmt.SetEndPosition(p.curToken.Pos)
		return stmt
	}

	for p.curTokenIs(token.ELSEIF) {
		elseifStmt := p.elseifStmt()
		stmt.ElseifStmts = append(stmt.ElseifStmts, elseifStmt)
	}

	stmt.ElseBlock = p.elseBlock()

	if !p.curTokenIs(token.END) { // move to "@end"
		return p.illegalNode()
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) ifStmtHeader(stmt *ast.IfStmt) *ast.IllegalNode {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.END)
	}

	p.next() // skip "("

	stmt.Condition = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.next() // skip ")"

	return nil
}

func (p *Parser) elseifStmt() ast.Statement {
	stmt := ast.NewElseIfStmt(p.curToken)

	p.next() // skip "@elseif"
	p.next() // skip "("

	stmt.Condition = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	p.next() // skip ")"

	stmt.Block = p.blockStmt()
	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) elseBlock() *ast.BlockStmt {
	if p.curTokenIs(token.ELSE) {
		p.next() // skip "@else"
	}

	if p.curTokenIs(token.ELSE, token.END) {
		return nil
	}

	stmt := p.blockStmt()

	if p.peekIs(token.ELSEIF) {
		p.newError(p.peekToken.ErrorLine(), fail.ErrElseifCannotFollowElse)
		return nil
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) _for() ast.Statement {
	stmt := ast.NewForStmt(p.curToken)

	if illegal := p.forStmtHeader(stmt); illegal != nil { // skips ")"
		return illegal
	}

	stmt.Block = p.blockStmt()

	if p.curTokenIs(token.END) {
		stmt.SetEndPosition(p.curToken.Pos)
		return stmt
	}

	stmt.ElseBlock = p.elseBlock()

	if !p.curTokenIs(token.END) {
		return p.illegalNodeUntil(token.END)
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) forStmtHeader(stmt *ast.ForStmt) *ast.IllegalNode {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.END)
	}

	// Parse Init
	if !p.peekIs(token.SEMI) {
		stmt.Init = p.embeddedCode()
	}

	if !p.expectPeek(token.SEMI) { // move to ";"
		return p.illegalNodeUntil(token.END)
	}

	// Parse Condition
	if !p.peekIs(token.SEMI) {
		p.next() // skip ";"
		stmt.Condition = p.expression(LOWEST)
	}

	if !p.expectPeek(token.SEMI) { // move to ";"
		return p.illegalNodeUntil(token.END)
	}

	// Parse Post statement
	if !p.peekIs(token.RPAREN) {
		stmt.Post = p.statement()
	}

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.next() // skip ")"

	return nil
}

func (p *Parser) eachStmt() ast.Statement {
	stmt := ast.NewEachStmt(p.curToken)

	if illegal := p.eachStmtHeader(stmt); illegal != nil { // skips ")"
		return illegal
	}

	stmt.Block = p.blockStmt()

	if p.curTokenIs(token.END) {
		stmt.SetEndPosition(p.curToken.Pos)
		return stmt
	}

	stmt.ElseBlock = p.elseBlock()

	if !p.curTokenIs(token.END) {
		return p.illegalNodeUntil(token.END)
	}

	stmt.SetEndPosition(p.curToken.Pos)

	return stmt
}

func (p *Parser) eachStmtHeader(stmt *ast.EachStmt) *ast.IllegalNode {
	if !p.expectPeek(token.LPAREN) { // move to "("
		return p.illegalNodeUntil(token.END)
	}

	p.next() // skip "("

	stmt.Var = ast.NewIdent(p.curToken, p.curToken.Lit)

	if !p.expectPeek(token.IN) { // move to "in"
		return p.illegalNodeUntil(token.END)
	}

	p.next() // skip "in"

	stmt.Arr = p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNodeUntil(token.END)
	}

	p.next() // skip ")"

	return nil
}

func (p *Parser) blockStmt() *ast.BlockStmt {
	if p.curTokenIs(token.ELSE, token.ELSEIF, token.END) {
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

		if p.peekIs(token.ELSE, token.ELSEIF, token.END) {
			p.next() // skip statement
			break
		}

		p.next() // skip statement
	}

	return stmt
}

func (p *Parser) expressionStmt() ast.Statement {
	prevTok := p.curToken

	exp := p.expression(LOWEST)

	switch p.peekToken.Type {
	case token.ASSIGN:
		return p.assignStmt(exp)
	case token.INC:
		return p.incStmt(exp)
	case token.DEC:
		return p.decStmt(exp)
	}

	stmt := ast.NewExpressionStmt(prevTok, exp)
	stmt.SetEndPosition(p.curToken.Pos)

	if p.peekIs(token.RBRACES) {
		p.next() // skip "}}"
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

	for !p.peekIs(token.RBRACES, token.SEMI, token.RPAREN) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		p.next()

		leftExp = infix(leftExp)
	}

	if leftExp != nil {
		leftExp.SetEndPosition(p.curToken.Pos)
	}

	return leftExp
}

func (p *Parser) prefixExp() ast.Expression {
	exp := ast.NewPrefixExp(p.curToken, p.curToken.Lit)

	p.next() // skip prefix operator

	exp.Right = p.expression(PREFIX)
	exp.SetEndPosition(p.curToken.Pos)

	return exp
}

func (p *Parser) groupedExp() ast.Expression {
	p.next() // skip "("

	exp := p.expression(LOWEST)

	if !p.expectPeek(token.RPAREN) { // move to ")"
		return p.illegalNode()
	}

	return exp
}

func (p *Parser) expressionList(endTok token.TokenType) []ast.Expression {
	var expressions []ast.Expression

	if p.peekIs(endTok) {
		p.next() // skip endTok token
		return expressions
	}

	if p.peekIs(token.END) {
		expressions = append(expressions, p.illegalNode())
		return expressions
	}

	p.next() // move to first expression

	expressions = append(expressions, p.expression(LOWEST))

	for p.peekIs(token.COMMA) && !p.curTokenIs(token.EOF) {
		p.next() // move to ","

		// break when has a trailing comma
		if p.peekIs(endTok) {
			break
		}

		p.next() // skip ","
		expressions = append(expressions, p.expression(LOWEST))
	}

	if !p.expectPeek(endTok) { // move to endTok
		expressions = append(expressions, ast.NewIllegalNode(p.curToken))
		return expressions
	}

	return expressions
}

func (p *Parser) illegalNode() *ast.IllegalNode {
	p.newError(
		p.curToken.ErrorLine(),
		fail.ErrIllegalToken,
		p.curToken.Lit,
	)

	return ast.NewIllegalNode(p.curToken)
}

func (p *Parser) illegalNodeUntil(tok token.TokenType) *ast.IllegalNode {
	for !p.curTokenIs(tok) && !p.curTokenIs(token.EOF) {
		p.next()
	}

	return p.illegalNode()
}
