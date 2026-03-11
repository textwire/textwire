package evaluator

import (
	"html"
	"reflect"
	"strings"

	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/pkg/ast"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/value"
)

var (
	NIL      = &value.Nil{}
	TRUE     = &value.Bool{Val: true}
	FALSE    = &value.Bool{Val: false}
	BREAK    = &value.Break{}
	CONTINUE = &value.Continue{}
)

type Evaluator struct {
	customFunc     *config.Func
	usingTemplates bool
	usingUseStmt   bool

	// config can be nil when Textwire is used for simple string and
	// file evaluation. If config is not nil, it means we use templates.
	config *config.Config
}

func New(customFunc *config.Func, conf *config.Config) *Evaluator {
	return &Evaluator{
		customFunc:     customFunc,
		config:         conf,
		usingTemplates: conf != nil,
	}
}

func (e *Evaluator) Eval(node ast.Node, ctx *Context) value.Value {
	if lit := e.evalIntoLit(node, ctx); lit != nil {
		return lit
	}
	return e.evalIntoValue(node, ctx)
}

func (e *Evaluator) evalIntoLit(node ast.Node, ctx *Context) value.Literal {
	switch node := node.(type) {
	case *ast.Program:
		return e.program(node, ctx)
	case *ast.Embedded:
		return e.embedded(node, ctx)
	case *ast.Text:
		return &value.Str{Val: node.String()}
	case *ast.IllegalNode:
		return NIL
	case *ast.IfDir:
		return e.ifDir(node, ctx)
	case *ast.ForDir:
		return e.forDir(node, ctx)
	case *ast.EachDir:
		return e.eachDir(node, ctx)
	case *ast.AssignStmt:
		return e.assignStmt(node, ctx)
	case *ast.IncStmt:
		return e.incStmt(node, ctx)
	case *ast.DecStmt:
		return e.decStmt(node, ctx)

	// Expressions
	case *ast.IdentExpr:
		return e.identExpr(node, ctx)
	case *ast.IndexExpr:
		return e.indexExpr(node, ctx)
	case *ast.DotExpr:
		return e.dotExpr(node, ctx)
	case *ast.StrExpr:
		return e.strExpr(node)
	case *ast.BoolExpr:
		return nativeBoolToBoolObj(node.Val)
	case *ast.ObjExpr:
		return e.objExpr(node, ctx)
	case *ast.ArrExpr:
		return e.arrExpr(node, ctx)
	case *ast.PrefixExpr:
		return e.prefixExpr(node, ctx)
	case *ast.TernaryExpr:
		return e.ternaryExpr(node, ctx)
	case *ast.InfixExpr:
		return e.infixExpr(node.Op, node.Left, node.Right, ctx)
	case *ast.CallExpr:
		return e.callExpr(node, ctx)
	case *ast.GlobalCallExpr:
		return e.globalCallExpr(node, ctx)
	case *ast.IntExpr:
		return &value.Int{Val: node.Val}
	case *ast.FloatExpr:
		return &value.Float{Val: node.Val}
	case *ast.NilExpr:
		return NIL
	}
	return nil
}

func (e *Evaluator) evalIntoValue(node ast.Node, ctx *Context) value.Value {
	switch node := node.(type) {
	case *ast.Block:
		return e.block(node, ctx)
	case *ast.UseDir:
		return e.useDir(node, ctx)
	case *ast.ContinueDir:
		return CONTINUE
	case *ast.BreakDir:
		return BREAK
	case *ast.ReserveDir:
		return e.reserveDir(node, ctx)
	case *ast.BreakifDir:
		return e.breakifDir(node, ctx)
	case *ast.ComponentDir:
		return e.compDir(node, ctx)
	case *ast.ContinueifDir:
		return e.continueifDir(node, ctx)
	case *ast.SlotDir:
		return e.slotDir(node, ctx)
	case *ast.SlotifDir:
		return e.slotifDir(node, ctx)
	case *ast.DumpDir:
		return e.dumpDir(node, ctx)
	case *ast.InsertDir:
		return e.insertDir(node, ctx)
	}
	return e.newError(node, ctx, fail.ErrUnknownType, node)
}

func (e *Evaluator) program(prog *ast.Program, ctx *Context) value.Literal {
	var stmts strings.Builder
	stmts.Grow(len(prog.Chunks))

	for i := range prog.Chunks {
		stmt := e.evalIntoLit(prog.Chunks[i], ctx)
		if isError(stmt) {
			return stmt
		}

		stmts.WriteString(stmt.String())
	}

	return &value.Str{Val: stmts.String()}
}

func (e *Evaluator) embedded(embedded *ast.Embedded, ctx *Context) value.Literal {
	var out strings.Builder
	out.Grow(len(embedded.Segments))

	for _, segment := range embedded.Segments {
		val := e.evalIntoLit(segment, ctx)
		if isError(val) {
			return val
		}
		out.WriteString(val.String())
	}

	return &value.Str{Val: out.String()}
}

func (e *Evaluator) ifDir(ifStmt *ast.IfDir, ctx *Context) value.Literal {
	cond := e.evalIntoLit(ifStmt.Cond, ctx)
	if isError(cond) {
		return cond
	}

	ifCtx := NewContext(ctx.scope.Child(), ctx.absPath)
	if isTruthy(cond) {
		return e.evalIntoLit(ifStmt.IfBlock, ifCtx)
	}

	for i := range ifStmt.ElseifDirs {
		elseifNode := ifStmt.ElseifDirs[i]

		cond = e.evalIntoLit(elseifNode.Cond, ifCtx)
		if isError(cond) {
			return cond
		}

		if isTruthy(cond) {
			return e.evalIntoLit(elseifNode.Block, ifCtx)
		}
	}

	if ifStmt.ElseBlock != nil {
		return e.evalIntoLit(ifStmt.ElseBlock, ifCtx)
	}

	return NIL
}

func (e *Evaluator) block(blockStmt *ast.Block, ctx *Context) value.Value {
	if blockStmt == nil {
		return NIL
	}

	chunks := make([]value.Value, 0, len(blockStmt.Chunks))

	for i := range blockStmt.Chunks {
		stmt := e.evalIntoLit(blockStmt.Chunks[i], ctx)
		if isError(stmt) {
			return stmt
		}

		chunks = append(chunks, stmt)
		if hasBreakStmt(stmt) || hasContinueStmt(stmt) {
			break
		}
	}

	return &value.Block{Elements: chunks}
}

func (e *Evaluator) assignStmt(assignStmt *ast.AssignStmt, ctx *Context) value.Literal {
	right := e.evalIntoLit(assignStmt.Right, ctx)
	if isError(right) {
		return right
	}

	return e.assignTo(assignStmt.Left, right, ctx)
}

func (e *Evaluator) assignTo(
	left ast.Expression,
	right value.Literal,
	ctx *Context,
) value.Literal {
	switch l := left.(type) {
	case *ast.IdentExpr:
		return e.assignIdent(l, right, ctx)
	case *ast.IndexExpr:
		return e.assignIndexExp(l, right, ctx)
	case *ast.DotExpr:
		return e.assignDotExp(l, right, ctx)
	}

	return e.newError(
		left,
		ctx,
		fail.ErrNotSupportedAssign,
		value.FromTokenToValueType(left.Tok().Type),
	)
}

func (e *Evaluator) assignIdent(
	ident *ast.IdentExpr,
	val value.Literal,
	ctx *Context,
) value.Literal {
	if err := ctx.scope.Set(ident.Name, val); err != nil {
		return e.newError(ident, ctx, "%s", err.Error())
	}
	return NIL
}

func (e *Evaluator) assignIndexExp(
	indexExp *ast.IndexExpr,
	val value.Literal,
	ctx *Context,
) value.Literal {
	left := e.evalIntoLit(indexExp.Left, ctx)
	if isError(left) {
		return left
	}

	idx := e.evalIntoLit(indexExp.Index, ctx)
	if isError(idx) {
		return idx
	}

	if !left.Is(value.ARR_VAL) {
		return e.newError(indexExp, ctx, fail.ErrIndexNotSupported, left.Type())
	}

	// Index must be integer
	if !idx.Is(value.INT_VAL) {
		return e.newError(indexExp, ctx, fail.ErrArrIndexInt, idx.Type())
	}

	arr := left.(*value.Arr)
	index := idx.(*value.Int).Val

	if index < 0 || index >= int64(len(arr.Elements)) {
		return e.newError(indexExp, ctx, fail.ErrArrIndexOutOfBound, index, len(arr.Elements))
	}

	arr.Elements[index] = val

	return NIL
}

func (e *Evaluator) assignDotExp(
	dotExp *ast.DotExpr,
	val value.Literal,
	ctx *Context,
) value.Literal {
	// Evaluate the left side to get the value
	left := e.evalIntoLit(dotExp.Left, ctx)
	if isError(left) {
		return left
	}

	// Get the key (property name)
	key := dotExp.Key.(*ast.IdentExpr).Name

	// Type assert that left is a value
	obj, ok := left.(*value.Obj)
	if !ok {
		return e.newError(dotExp, ctx, fail.ErrKeyOnNonObj, left.Type(), key)
	}

	obj.Pairs[key] = val

	return NIL
}

func (e *Evaluator) useDir(useStmt *ast.UseDir, ctx *Context) value.Value {
	if useStmt.LayoutProg == nil {
		if e.usingTemplates {
			return e.newError(useStmt, ctx, fail.ErrUseDirMissingLayout, useStmt.Name.Val)
		}
		return e.newError(useStmt, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	e.usingUseStmt = true

	// Make sure that layout is missing @use
	if useStmt.LayoutProg.IsLayout && useStmt.LayoutProg.HasUseDir() {
		return e.newError(useStmt, ctx, fail.ErrDirStmtNotAllowed)
	}

	// Create new layout context and pass inserts to it
	layoutCtx := NewContext(value.NewScope(), useStmt.LayoutProg.AbsPath)

	// Evaluate @inserts and map them into new context for layout
	for name, insertStmt := range useStmt.Inserts {
		insert := e.insertDir(insertStmt, ctx)
		if isError(insert) {
			return insert
		}
		layoutCtx.inserts[name] = insert
	}

	// Evaluate layout program with new context
	layout := e.evalIntoLit(useStmt.LayoutProg, layoutCtx)
	if isError(layout) {
		return layout
	}

	return &value.Use{
		Path:   useStmt.Name.Val,
		Layout: layout,
	}
}

func (e *Evaluator) reserveDir(reserveStmt *ast.ReserveDir, ctx *Context) value.Value {
	if !e.usingTemplates {
		return e.newError(reserveStmt, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	name := reserveStmt.Name.Val
	insert, ok := ctx.inserts[name]
	if !ok {
		// Inserts are optional, NIL when not provided or fallback argument
		if reserveStmt.Fallback == nil {
			return NIL
		}
		return e.evalIntoLit(reserveStmt.Fallback, ctx)
	}

	// delete reserve after it's been used by reserve
	defer delete(ctx.inserts, name)

	return &value.Reserve{
		Name:   name,
		Insert: insert,
	}
}

func (e *Evaluator) compDir(compStmt *ast.ComponentDir, ctx *Context) value.Value {
	if !e.usingTemplates {
		return e.newError(compStmt, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	name := compStmt.Name.Val
	if compStmt.CompProg == nil {
		return e.newError(compStmt, ctx, fail.ErrUndefinedComponent, name)
	}

	compCtx := NewContext(value.NewScope(), compStmt.CompProg.AbsPath)

	// Evaluate local slots and add them to component context
	for _, slotStmt := range compStmt.Slots {
		slot := e.evalIntoLit(slotStmt, ctx)
		if isError(slot) {
			return slot
		}

		if compCtx.slots[name] == nil {
			compCtx.slots[name] = map[string]value.Value{}
		}

		compCtx.slots[name][slotStmt.Name().Val] = slot
	}

	if compStmt.Argument != nil {
		for key, arg := range compStmt.Argument.Pairs {
			obj := e.evalIntoLit(arg, ctx)
			if isError(obj) {
				return obj
			}

			if err := compCtx.scope.Set(key, obj); err != nil {
				return e.newError(compStmt, ctx, "%s", err.Error())
			}
		}
	}

	content := e.evalIntoLit(compStmt.CompProg, compCtx)
	if isError(content) {
		return content
	}

	return &value.Component{
		Name:    name,
		Content: content,
	}
}

func (e *Evaluator) forDir(forStmt *ast.ForDir, ctx *Context) value.Literal {
	forCtx := NewContext(ctx.scope, ctx.absPath)

	var init value.Literal
	if forStmt.Init != nil {
		if init = e.evalIntoLit(forStmt.Init, forCtx); isError(init) {
			return init
		}
	}

	// Evaluate ElseBlock block if user's condition is false
	if forStmt.Cond != nil {
		cond := e.evalIntoLit(forStmt.Cond, forCtx)
		if isError(cond) {
			return cond
		}

		if !isTruthy(cond) && forStmt.ElseBlock != nil {
			return e.evalIntoLit(forStmt.ElseBlock, forCtx)
		}
	}

	var blocks strings.Builder

	// Loop through the block until the user's condition is false
	for {
		if forStmt.Cond != nil {
			cond := e.evalIntoLit(forStmt.Cond, forCtx)
			if isError(cond) {
				return cond
			}

			if !isTruthy(cond) {
				break
			}
		}

		block := e.evalIntoLit(forStmt.Block, forCtx)
		if isError(block) {
			return block
		}

		blocks.WriteString(block.String())

		if forStmt.Post != nil {
			post := e.evalIntoLit(forStmt.Post, forCtx)
			if isError(post) {
				return post
			}
		}

		if hasBreakStmt(block) {
			break
		}

		if hasContinueStmt(block) {
			continue
		}
	}

	return &value.Str{Val: blocks.String()}
}

func (e *Evaluator) eachDir(eachStmt *ast.EachDir, ctx *Context) value.Literal {
	eachCtx := NewContext(ctx.scope.Child(), ctx.absPath)
	varName := eachStmt.Var.Name

	arrObj := e.evalIntoLit(eachStmt.Arr, eachCtx)
	if isError(arrObj) {
		return arrObj
	}

	arr, ok := arrObj.(*value.Arr)
	if !ok {
		return e.newError(eachStmt, eachCtx, fail.ErrEachDirWithNonArrArg, arrObj.Type())
	}

	arrElems := arr.Elements

	// Evaluate ElseBlock when array is empty
	if len(arrElems) == 0 && eachStmt.ElseBlock != nil {
		return e.evalIntoLit(eachStmt.ElseBlock, eachCtx)
	}

	var blocks strings.Builder
	blocks.Grow(len(arrElems))

	for i := range arrElems {
		if err := eachCtx.scope.Set(varName, arrElems[i]); err != nil {
			return e.newError(eachStmt, eachCtx, "%s", err.Error())
		}

		eachCtx.scope.SetLoopVar(map[string]value.Literal{
			"index": &value.Int{Val: int64(i)},
			"first": nativeBoolToBoolObj(i == 0),
			"last":  nativeBoolToBoolObj(i == len(arrElems)-1),
			"iter":  &value.Int{Val: int64(i + 1)},
		})

		block := e.evalIntoLit(eachStmt.Block, eachCtx)
		if isError(block) {
			return block
		}

		blocks.WriteString(block.String())
		if hasBreakStmt(block) {
			break
		}

		if hasContinueStmt(block) {
			continue
		}
	}

	return &value.Str{Val: blocks.String()}
}

func (e *Evaluator) breakifDir(breakifStmt *ast.BreakifDir, ctx *Context) value.Value {
	cond := e.evalIntoLit(breakifStmt.Cond, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return BREAK
	}

	return NIL
}

func (e *Evaluator) continueifDir(contifStmt *ast.ContinueifDir, ctx *Context) value.Value {
	cond := e.evalIntoLit(contifStmt.Cond, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return CONTINUE
	}

	return NIL
}

func (e *Evaluator) slotDir(slotStmt *ast.SlotDir, ctx *Context) value.Value {
	if slotStmt.IsLocal {
		return e.localSlotStmt(slotStmt, ctx)
	}
	return e.externalSlotStmt(slotStmt, ctx)
}

func (e *Evaluator) slotifDir(slotifStmt *ast.SlotifDir, ctx *Context) value.Value {
	cond := e.evalIntoLit(slotifStmt.Cond, ctx)
	if isError(cond) {
		return cond
	}

	if !isTruthy(cond) {
		return NIL
	}

	block := e.evalIntoLit(slotifStmt.Block(), ctx)
	if isError(block) {
		return block
	}

	return &value.Slot{
		Name:    slotifStmt.Name().Val,
		Content: block,
	}
}

func (e *Evaluator) externalSlotStmt(slotStmt *ast.SlotDir, ctx *Context) value.Value {
	name := slotStmt.Name().Val
	compName := slotStmt.CompName

	// Get slot's content from the context
	content, ok := ctx.slots[compName][name]
	if !ok {
		// Slots are optional in component files since v3.1.0
		return NIL
	}

	// delete slot after it's been used by external component
	defer delete(ctx.slots[compName], name)

	return &value.Slot{Name: name, Content: content}
}

func (e *Evaluator) localSlotStmt(slotStmt *ast.SlotDir, ctx *Context) value.Value {
	var block value.Literal = NIL

	if slotStmt.Block() != nil {
		block = e.evalIntoLit(slotStmt.Block(), ctx)
		if isError(block) {
			return block
		}
	}

	return &value.Slot{
		Name:    slotStmt.Name().Val,
		Content: block,
	}
}

func (e *Evaluator) insertDir(insertStmt *ast.InsertDir, ctx *Context) value.Value {
	if !e.usingTemplates {
		return e.newError(insertStmt, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	name := insertStmt.Name.Val
	if !e.usingUseStmt {
		return e.newError(insertStmt, ctx, fail.ErrInsertRequiresUse, name)
	}

	block := e.combineInsertContent(insertStmt, ctx)
	if isError(block) {
		return block
	}

	return &value.Insert{
		Name:  name,
		Block: block,
	}
}

// combineInsertContent combines insert Argument and Block (depending what user has)
// into a single value that we can work with.
func (e *Evaluator) combineInsertContent(insertStmt *ast.InsertDir, ctx *Context) value.Literal {
	if insertStmt.Argument != nil {
		arg := e.evalIntoLit(insertStmt.Argument, ctx)
		if isError(arg) {
			return arg
		}
		return arg
	}

	if insertStmt.Block == nil {
		return e.newError(insertStmt, ctx, fail.ErrInsertMustHaveContent)
	}

	return e.evalIntoLit(insertStmt.Block, ctx)
}

func (e *Evaluator) dumpDir(dumpStmt *ast.DumpDir, ctx *Context) value.Value {
	values := make([]string, 0, len(dumpStmt.Args))

	for i := range dumpStmt.Args {
		evaluated := e.evalIntoLit(dumpStmt.Args[i], ctx)
		values = append(values, evaluated.Dump(0))
	}

	return &value.Dump{Vals: values}
}

func (e *Evaluator) identExpr(ident *ast.IdentExpr, ctx *Context) value.Literal {
	varName := ident.Name
	if varName == "global" && e.config != nil && e.config.GlobalData != nil {
		return value.NativeToValue(e.config.GlobalData)
	}

	if val, ok := ctx.scope.Get(varName); ok {
		return val
	}

	return e.newError(ident, ctx, fail.ErrVariableIsUndefined, ident.Name)
}

func (e *Evaluator) indexExpr(indexExp *ast.IndexExpr, ctx *Context) value.Literal {
	left := e.evalIntoLit(indexExp.Left, ctx)
	if isError(left) {
		return left
	}

	idx := e.evalIntoLit(indexExp.Index, ctx)
	if isError(idx) {
		return idx
	}

	switch {
	case left.Is(value.ARR_VAL) && idx.Is(value.INT_VAL):
		return e.arrIndexExp(left, idx)
	case left.Is(value.OBJ_VAL) && idx.Is(value.STR_VAL):
		return e.objKeyExp(left.(*value.Obj), idx.(*value.Str).Val)
	}

	return e.newError(indexExp, ctx, fail.ErrIndexNotSupported, left.Type())
}

func (e *Evaluator) arrIndexExp(arrObj, idx value.Literal) value.Literal {
	arr := arrObj.(*value.Arr)
	index := idx.(*value.Int).Val
	max := int64(len(arr.Elements) - 1)

	if index < 0 || index > max {
		return NIL
	}

	return arr.Elements[index]
}

func (e *Evaluator) objKeyExp(obj *value.Obj, key string) value.Literal {
	if pair, ok := obj.Pairs[key]; ok {
		return pair
	}

	// Capitalize the first letter of the key and try again to support
	// case insensitive key access for the first key character.
	if pair, ok := obj.Pairs[capitalizeFirst(key)]; ok {
		return pair
	}

	return NIL
}

func (e *Evaluator) dotExpr(dotExp *ast.DotExpr, ctx *Context) value.Literal {
	left := e.evalIntoLit(dotExp.Left, ctx)
	if isError(left) {
		return left
	}

	key := dotExp.Key.(*ast.IdentExpr)
	obj, ok := left.(*value.Obj)
	if !ok {
		return e.newError(dotExp, ctx, fail.ErrKeyOnNonObj, left.Type(), key)
	}

	return e.objKeyExp(obj, key.Name)
}

func (e *Evaluator) strExpr(strLit *ast.StrExpr) value.Literal {
	str := html.EscapeString(strLit.Val)

	// unescape single and double quotes
	str = strings.ReplaceAll(str, "&#34;", `"`)
	str = strings.ReplaceAll(str, "&#39;", `'`)

	return &value.Str{Val: str}
}

func (e *Evaluator) prefixExpr(prefixExp *ast.PrefixExpr, ctx *Context) value.Literal {
	right := e.evalIntoLit(prefixExp.Right, ctx)
	if isError(right) {
		return right
	}

	switch prefixExp.Op {
	case "-":
		return e.minusPrefixOpExp(right, prefixExp, ctx)
	case "!":
		return e.bangOpExp(right, prefixExp, ctx)
	}

	return e.newError(prefixExp, ctx, fail.ErrUnknownOp, prefixExp.Op, right.Type())
}

func (e *Evaluator) ternaryExpr(ternExp *ast.TernaryExpr, ctx *Context) value.Literal {
	cond := e.evalIntoLit(ternExp.Cond, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return e.evalIntoLit(ternExp.IfExpr, ctx)
	}

	return e.evalIntoLit(ternExp.ElseExpr, ctx)
}

func (e *Evaluator) arrExpr(arrLit *ast.ArrExpr, ctx *Context) value.Literal {
	elems := e.evalExpressions(arrLit.Elements, ctx)
	if len(elems) == 1 && isError(elems[0]) {
		return elems[0]
	}

	return &value.Arr{Elements: elems}
}

func (e *Evaluator) objExpr(objLit *ast.ObjExpr, ctx *Context) value.Literal {
	pairs := make(map[string]value.Literal, len(objLit.Pairs))

	for key, val := range objLit.Pairs {
		valObj := e.evalIntoLit(val, ctx)
		if isError(valObj) {
			return valObj
		}

		pairs[key] = valObj
	}

	return value.NewObj(pairs)
}

func (e *Evaluator) evalExpressions(exps []ast.Expression, ctx *Context) []value.Literal {
	evaluatedObjs := make([]value.Literal, 0, len(exps))

	for i := range exps {
		evaluated := e.evalIntoLit(exps[i], ctx)
		if isError(evaluated) {
			return []value.Literal{evaluated}
		}

		evaluatedObjs = append(evaluatedObjs, evaluated)
	}

	return evaluatedObjs
}

func (e *Evaluator) infixExpr(
	op string,
	leftNode,
	rightNode ast.Expression,
	ctx *Context,
) value.Literal {
	left := e.evalIntoLit(leftNode, ctx)
	if isError(left) {
		return left
	}

	if obj, ok := e.shortCircuit(left, op); ok {
		return obj
	}

	right := e.evalIntoLit(rightNode, ctx)
	if isError(right) {
		return right
	}

	if op == "&&" || op == "||" {
		return e.logicalExp(op, right, left, leftNode, ctx)
	}

	return e.operatorExp(op, left, right, leftNode, ctx)
}

// Short-circuit evaluation for logical operators to prevent
// checking conditions if the left side is false.
func (e *Evaluator) shortCircuit(left value.Literal, op string) (value.Literal, bool) {
	if op == "&&" && !isTruthy(left) {
		return FALSE, true
	}
	if op == "||" && isTruthy(left) {
		return TRUE, true
	}

	return nil, false
}

func (e *Evaluator) incStmt(incStmt *ast.IncStmt, ctx *Context) value.Literal {
	left := e.evalIntoLit(incStmt.Left, ctx)
	if isError(left) {
		return left
	}

	switch l := left.(type) {
	case *value.Int:
		return e.assignTo(incStmt.Left, &value.Int{Val: l.Val + 1}, ctx)
	case *value.Float:
		return e.assignTo(incStmt.Left, &value.Float{Val: l.Val + 1}, ctx)
	}

	return e.newError(incStmt, ctx, fail.ErrIllegalTypeForInc, left.Type())
}

func (e *Evaluator) decStmt(decInc *ast.DecStmt, ctx *Context) value.Literal {
	left := e.evalIntoLit(decInc.Left, ctx)
	if isError(left) {
		return left
	}

	switch l := left.(type) {
	case *value.Int:
		return e.assignTo(decInc.Left, &value.Int{Val: l.Val - 1}, ctx)
	case *value.Float:
		err := l.SubtractFromFloat(1)
		if err != nil {
			return e.newError(decInc, ctx, fail.ErrCannotDecFromFloat, left, err)
		}
		return e.assignTo(decInc.Left, l, ctx)
	}

	return e.newError(decInc, ctx, fail.ErrIllegalTypeForDec, left.Type())
}

func (e *Evaluator) callExpr(callExp *ast.CallExpr, ctx *Context) value.Literal {
	receiver := e.evalIntoLit(callExp.Receiver, ctx)
	funcName := callExp.Function.Name
	if isError(receiver) {
		return receiver
	}

	receiverType := receiver.Type()
	typeFuncs, ok := functions[receiverType]
	if !ok {
		return e.newError(callExp, ctx, fail.ErrFuncNotDefined, receiverType, funcName)
	}

	args := e.evalExpressions(callExp.Arguments, ctx)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	buitin, ok := typeFuncs[callExp.Function.Name]
	if ok {
		result, err := buitin.Fn(receiver, args...)
		if err != nil {
			return e.newError(callExp, ctx, "%s", err.Error())
		}
		return result
	}

	if hasCustomFunc(e.customFunc, receiverType, funcName) {
		nativeArgs := e.valuesToNativeType(args)

		switch r := receiver.(type) {
		case *value.Str:
			fun := e.customFunc.Str[funcName]
			res := fun(r.String(), nativeArgs...)
			return value.NativeToValue(res)
		case *value.Arr:
			fun := e.customFunc.Arr[funcName]
			nativeElems := e.valuesToNativeType(r.Elements)
			res := fun(nativeElems, nativeArgs...)
			return value.NativeToValue(res)
		case *value.Int:
			fun := e.customFunc.Int[funcName]
			res := fun(int(r.Val), nativeArgs...)
			return value.NativeToValue(res)
		case *value.Float:
			fun := e.customFunc.Float[funcName]
			res := fun(r.Val, nativeArgs...)
			return value.NativeToValue(res)
		case *value.Bool:
			fun := e.customFunc.Bool[funcName]
			res := fun(r.Val, nativeArgs...)
			return value.NativeToValue(res)
		case *value.Obj:
			fun := e.customFunc.Obj[funcName]
			firstArg := r.Native()
			res := fun(firstArg.(map[string]any), nativeArgs...)
			return value.NativeToValue(res)
		}
	}

	return e.newError(callExp, ctx, fail.ErrFuncNotDefined, receiver.Type(), callExp.Function.Name)
}

func (e *Evaluator) globalCallExpr(globalCallExp *ast.GlobalCallExpr, ctx *Context) value.Literal {
	switch globalCallExp.Function.Name {
	case "defined":
		return e.globalFuncDefined(globalCallExp, ctx)
	case "hasValue":
		return e.globalFuncHasValue(globalCallExp, ctx)
	default:
		return e.newError(
			globalCallExp,
			ctx,
			fail.ErrGlobalFuncMissing,
			globalCallExp.Function.Name,
		)
	}
}

func (e *Evaluator) globalFuncDefined(
	globalCallExp *ast.GlobalCallExpr,
	ctx *Context,
) value.Literal {
	for i := range globalCallExp.Arguments {
		evaluated := e.evalIntoLit(globalCallExp.Arguments[i], ctx)
		if isUndefinedError(evaluated) {
			return FALSE
		}

		if isError(evaluated) {
			return evaluated
		}
	}

	return TRUE
}

func (e *Evaluator) globalFuncHasValue(
	globalCallExp *ast.GlobalCallExpr,
	ctx *Context,
) value.Literal {
	// Check if all arguments are defined first
	areDefined := e.globalFuncDefined(globalCallExp, ctx)
	if isError(areDefined) {
		return areDefined
	}

	if areDefined == FALSE {
		return FALSE
	}

	// At this point, all variables are defined.
	// Checking if they have nullable values.
	for _, exp := range globalCallExp.Arguments {
		arg := e.evalIntoLit(exp, ctx)
		if isError(arg) {
			return arg
		}

		if !isTruthy(arg) {
			return FALSE
		}
	}

	return TRUE
}

func (e *Evaluator) valuesToNativeType(args []value.Literal) []any {
	vals := make([]any, len(args))
	for i := range args {
		vals[i] = args[i].Native()
	}

	return vals
}

func (e *Evaluator) operatorExp(
	op string,
	left,
	right value.Literal,
	leftNode ast.Node,
	ctx *Context,
) value.Literal {
	if op == "==" || op == "!=" {
		return e.comparrisonInfixExp(op, right, left, leftNode, ctx)
	}

	switch l := left.(type) {
	case *value.Int:
		return e.intInfixExp(op, right, l, leftNode, ctx)
	case *value.Float:
		return e.floatInfixExp(op, right, l, leftNode, ctx)
	case *value.Str:
		return e.stringInfixExp(op, right, l, leftNode, ctx)
	}

	return e.newError(leftNode, ctx, fail.ErrCannotUseOperator, op, left.Type(), op, right.Type())
}

func (e *Evaluator) logicalExp(
	op string,
	right,
	left value.Literal,
	leftNode ast.Node,
	ctx *Context,
) value.Literal {
	switch op {
	case "&&":
		return &value.Bool{Val: isTruthy(left) && isTruthy(right)}
	case "||":
		return &value.Bool{Val: isTruthy(left) || isTruthy(right)}
	}

	return e.newError(leftNode, ctx, fail.ErrCannotUseOperator, op, left.Type(), op, right.Type())
}

func (e *Evaluator) intInfixExp(
	op string,
	right value.Literal,
	l *value.Int,
	leftNode ast.Node,
	ctx *Context,
) value.Literal {
	r, ok := right.(*value.Int)
	if !ok {
		return e.newError(leftNode, ctx, fail.ErrCannotUseOperator, op, l.Type(), op, right.Type())
	}

	switch op {
	case "+":
		return &value.Int{Val: l.Val + r.Val}
	case "-":
		return &value.Int{Val: l.Val - r.Val}
	case "*":
		return &value.Int{Val: l.Val * r.Val}
	case "/":
		if r.Val == 0 {
			return e.newError(leftNode, ctx, fail.ErrDivisionByZero)
		}
		return &value.Int{Val: l.Val / r.Val}
	case "%":
		return &value.Int{Val: l.Val % r.Val}
	case ">":
		return nativeBoolToBoolObj(l.Val > r.Val)
	case "<":
		return nativeBoolToBoolObj(l.Val < r.Val)
	case ">=":
		return nativeBoolToBoolObj(l.Val >= r.Val)
	case "<=":
		return nativeBoolToBoolObj(l.Val <= r.Val)
	}

	return e.newError(leftNode, ctx, fail.ErrCannotUseOperator, op, l.Type(), op, right.Type())
}

func (e *Evaluator) comparrisonInfixExp(
	op string, // == or !=
	right,
	left value.Literal,
	leftNode ast.Node,
	ctx *Context,
) value.Literal {
	if left.Type() != right.Type() {
		return nativeBoolToBoolObj(op == "!=")
	}
	var areEqual bool

	switch l := left.(type) {
	case *value.Int:
		areEqual = l.Val == right.(*value.Int).Val
	case *value.Float:
		areEqual = l.Val == right.(*value.Float).Val
	case *value.Str:
		areEqual = l.Val == right.(*value.Str).Val
	case *value.Bool:
		areEqual = l.Val == right.(*value.Bool).Val
	case *value.Arr, *value.Obj:
		areEqual = reflect.DeepEqual(left, right)
	case *value.Nil:
		areEqual = true
	default:
		e.newError(leftNode, ctx, fail.ErrCannotUseOperator, op, l.Type(), op, right.Type())
	}

	if op == "!=" {
		areEqual = !areEqual
	}

	return nativeBoolToBoolObj(areEqual)
}

func (e *Evaluator) stringInfixExp(
	op string,
	right value.Literal,
	l *value.Str,
	leftNode ast.Node,
	ctx *Context,
) value.Literal {
	r, ok := right.(*value.Str)
	if !ok {
		return e.newError(leftNode, ctx, fail.ErrCannotUseOperator, op, l.Type(), op, right.Type())
	}

	if op == "+" {
		return &value.Str{Val: l.Val + r.Val}
	}

	return e.newError(leftNode, ctx, fail.ErrCannotUseOperator, op, l.Type(), op, right.Type())
}

func (e *Evaluator) floatInfixExp(
	op string,
	right value.Literal,
	l *value.Float,
	leftNode ast.Node,
	ctx *Context,
) value.Literal {
	r, ok := right.(*value.Float)
	if !ok {
		return e.newError(leftNode, ctx, fail.ErrCannotUseOperator, op, l.Type(), op, right.Type())
	}

	switch op {
	case "+":
		return &value.Float{Val: l.Val + r.Val}
	case "-":
		return &value.Float{Val: l.Val - r.Val}
	case "*":
		return &value.Float{Val: l.Val * r.Val}
	case "/":
		return &value.Float{Val: l.Val / r.Val}
	case ">":
		return nativeBoolToBoolObj(l.Val > r.Val)
	case "<":
		return nativeBoolToBoolObj(l.Val < r.Val)
	case ">=":
		return nativeBoolToBoolObj(l.Val >= r.Val)
	case "<=":
		return nativeBoolToBoolObj(l.Val <= r.Val)
	}

	return e.newError(leftNode, ctx, fail.ErrCannotUseOperator, op, l.Type(), op, right.Type())
}

func (e *Evaluator) minusPrefixOpExp(
	right value.Literal,
	node ast.Node,
	ctx *Context,
) value.Literal {
	switch r := right.(type) {
	case *value.Int:
		return &value.Int{Val: -r.Val}
	case *value.Float:
		return &value.Float{Val: -r.Val}
	}

	return e.newError(node, ctx, fail.ErrPrefixOpIsWrong, "-", right.Type())
}

func (e *Evaluator) bangOpExp(right value.Literal, node ast.Node, ctx *Context) value.Literal {
	switch right {
	case FALSE:
		return TRUE
	case TRUE:
		return FALSE
	case NIL:
		return TRUE
	}

	return e.newError(node, ctx, fail.ErrPrefixOpIsWrong, "!", right.Type())
}

func (e *Evaluator) newError(node ast.Node, ctx *Context, format string, a ...any) *value.Error {
	return &value.Error{
		Err:     fail.New(node.Line(), ctx.absPath, "evaluator", format, a...),
		ErrorID: format,
	}
}
