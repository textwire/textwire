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
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return e.program(node, ctx)
	case *ast.HTMLStmt:
		return &value.HTML{Val: node.String()}
	case *ast.ExpressionStmt:
		return e.Eval(node.Expression, ctx)
	case *ast.IfStmt:
		return e._if(node, ctx)
	case *ast.BlockStmt:
		return e.block(node, ctx)
	case *ast.AssignStmt:
		return e.assign(node, ctx)
	case *ast.UseStmt:
		return e.use(node, ctx)
	case *ast.ReserveStmt:
		return e.reserve(node, ctx)
	case *ast.ForStmt:
		return e._for(node, ctx)
	case *ast.EachStmt:
		return e.each(node, ctx)
	case *ast.BreakifStmt:
		return e.breakif(node, ctx)
	case *ast.ComponentStmt:
		return e.component(node, ctx)
	case *ast.ContinueifStmt:
		return e.continueif(node, ctx)
	case *ast.SlotStmt:
		return e.slot(node, ctx)
	case *ast.SlotifStmt:
		return e.slotif(node, ctx)
	case *ast.DumpStmt:
		return e.dump(node, ctx)
	case *ast.InsertStmt:
		return e.insert(node, ctx)
	case *ast.ContinueStmt:
		return CONTINUE
	case *ast.BreakStmt:
		return BREAK
	case *ast.IllegalNode:
		return NIL

	// Expressions
	case *ast.Identifier:
		return e.ident(node, ctx)
	case *ast.IndexExp:
		return e.indexExp(node, ctx)
	case *ast.DotExp:
		return e.dotExp(node, ctx)
	case *ast.IntegerLiteral:
		return &value.Int{Val: node.Val}
	case *ast.FloatLiteral:
		return &value.Float{Val: node.Val}
	case *ast.StringLiteral:
		return e.stringLit(node)
	case *ast.BooleanLiteral:
		return nativeBoolToBoolObj(node.Val)
	case *ast.ObjectLiteral:
		return e.objectLit(node, ctx)
	case *ast.ArrayLiteral:
		return e.arrayLit(node, ctx)
	case *ast.PrefixExp:
		return e.prefixExp(node, ctx)
	case *ast.TernaryExp:
		return e.ternaryExp(node, ctx)
	case *ast.InfixExp:
		return e.infixExp(node.Op, node.Left, node.Right, ctx)
	case *ast.PostfixExp:
		return e.postfixExp(node, ctx)
	case *ast.CallExp:
		return e.callExp(node, ctx)
	case *ast.GlobalCallExp:
		return e.globalCallExp(node, ctx)
	case *ast.NilLiteral:
		return NIL
	}

	return e.newError(node, ctx, fail.ErrUnknownNodeType, node)
}

func (e *Evaluator) program(prog *ast.Program, ctx *Context) value.Value {
	var stmts strings.Builder
	stmts.Grow(len(prog.Statements))

	for i := range prog.Statements {
		stmt := e.Eval(prog.Statements[i], ctx)
		if isError(stmt) {
			return stmt
		}

		stmts.WriteString(stmt.String())
	}

	return &value.HTML{Val: stmts.String()}
}

func (e *Evaluator) _if(ifStmt *ast.IfStmt, ctx *Context) value.Value {
	cond := e.Eval(ifStmt.Condition, ctx)
	if isError(cond) {
		return cond
	}

	ifCtx := NewContext(ctx.scope.Child(), ctx.absPath)
	if isTruthy(cond) {
		return e.Eval(ifStmt.IfBlock, ifCtx)
	}

	for i := range ifStmt.ElseifStmts {
		elseifNode, ok := ifStmt.ElseifStmts[i].(*ast.ElseIfStmt)
		if !ok {
			continue
		}

		cond = e.Eval(elseifNode.Condition, ifCtx)
		if isError(cond) {
			return cond
		}

		if isTruthy(cond) {
			return e.Eval(elseifNode.Block, ifCtx)
		}
	}

	if ifStmt.ElseBlock != nil {
		return e.Eval(ifStmt.ElseBlock, ifCtx)
	}

	return NIL
}

func (e *Evaluator) block(blockStmt *ast.BlockStmt, ctx *Context) value.Value {
	if blockStmt == nil {
		return NIL
	}

	stmts := make([]value.Value, 0, len(blockStmt.Statements))

	for i := range blockStmt.Statements {
		stmt := e.Eval(blockStmt.Statements[i], ctx)
		if isError(stmt) {
			return stmt
		}

		stmts = append(stmts, stmt)
		if hasBreakStmt(stmt) || hasContinueStmt(stmt) {
			break
		}
	}

	return &value.Block{Elements: stmts}
}

func (e *Evaluator) assign(assignStmt *ast.AssignStmt, ctx *Context) value.Value {
	right := e.Eval(assignStmt.Right, ctx)
	if isError(right) {
		return right
	}

	switch left := assignStmt.Left.(type) {
	case *ast.Identifier:
		return e.assignIdentifier(left, right, ctx)
	case *ast.IndexExp:
		return e.assignIndexExp(left, right, ctx)
	case *ast.DotExp:
		return e.assignDotExp(left, right, ctx)
	default:
		return e.newError(
			assignStmt,
			ctx,
			fail.ErrNotSupportedAssign,
			value.FromTokenToObjectType(left.Tok().Type),
		)
	}
}

func (e *Evaluator) assignIdentifier(
	ident *ast.Identifier,
	val value.Value,
	ctx *Context,
) value.Value {
	if err := ctx.scope.Set(ident.Name, val); err != nil {
		return e.newError(ident, ctx, "%s", err.Error())
	}
	return NIL
}

func (e *Evaluator) assignIndexExp(
	indexExp *ast.IndexExp,
	val value.Value,
	ctx *Context,
) value.Value {
	left := e.Eval(indexExp.Left, ctx)
	if isError(left) {
		return left
	}

	idx := e.Eval(indexExp.Index, ctx)
	if isError(idx) {
		return idx
	}

	if !left.Is(value.ARR_OBJ) {
		return e.newError(indexExp, ctx, fail.ErrIndexNotSupported, left.Type())
	}

	// Index must be integer
	if !idx.Is(value.INT_OBJ) {
		return e.newError(indexExp, ctx, fail.ErrArrayIndexInteger, idx.Type())
	}

	arr := left.(*value.Array)
	index := idx.(*value.Int).Val

	if index < 0 || index >= int64(len(arr.Elements)) {
		return e.newError(indexExp, ctx, fail.ErrArrayIndexOutOfBound, index, len(arr.Elements))
	}

	arr.Elements[index] = val

	return NIL
}

func (e *Evaluator) assignDotExp(
	dotExp *ast.DotExp,
	val value.Value,
	ctx *Context,
) value.Value {
	// Evaluate the left side to get the object
	left := e.Eval(dotExp.Left, ctx)
	if isError(left) {
		return left
	}

	// Get the key (property name)
	key := dotExp.Key.(*ast.Identifier).Name

	// Type assert that left is an object
	obj, ok := left.(*value.Obj)
	if !ok {
		return e.newError(dotExp, ctx, fail.ErrKeyOnNonObject, left.Type(), key)
	}

	// Set the value on the object
	obj.Pairs[key] = val

	return NIL
}

func (e *Evaluator) use(useStmt *ast.UseStmt, ctx *Context) value.Value {
	if useStmt.LayoutProg == nil {
		if e.usingTemplates {
			return e.newError(useStmt, ctx, fail.ErrUseStmtMissingLayout, useStmt.Name.Val)
		}
		return e.newError(useStmt, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	e.usingUseStmt = true

	// Make sure that layout is missing @use
	if useStmt.LayoutProg.IsLayout && useStmt.LayoutProg.HasUseStmt() {
		return e.newError(useStmt, ctx, fail.ErrUseStmtNotAllowed)
	}

	// Create new layout context and pass inserts to it
	layoutCtx := NewContext(value.NewScope(), useStmt.LayoutProg.AbsPath)

	// Evaluate @inserts and map them into new context for layout
	for name, insertStmt := range useStmt.Inserts {
		insert := e.insert(insertStmt, ctx)
		if isError(insert) {
			return insert
		}
		layoutCtx.inserts[name] = insert
	}

	// Evaluate layout program with new context
	layout := e.Eval(useStmt.LayoutProg, layoutCtx)
	if isError(layout) {
		return layout
	}

	return &value.Use{
		Path:   useStmt.Name.Val,
		Layout: layout,
	}
}

func (e *Evaluator) reserve(reserveStmt *ast.ReserveStmt, ctx *Context) value.Value {
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
		return e.Eval(reserveStmt.Fallback, ctx)
	}

	// delete reserve after it's been used by reserve
	defer delete(ctx.inserts, name)

	return &value.Reserve{
		Name:   name,
		Insert: insert,
	}
}

func (e *Evaluator) component(compStmt *ast.ComponentStmt, ctx *Context) value.Value {
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
		slot := e.Eval(slotStmt, ctx)
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
			obj := e.Eval(arg, ctx)
			if isError(obj) {
				return obj
			}

			if err := compCtx.scope.Set(key, obj); err != nil {
				return e.newError(compStmt, ctx, "%s", err.Error())
			}
		}
	}

	content := e.Eval(compStmt.CompProg, compCtx)
	if isError(content) {
		return content
	}

	return &value.Component{
		Name:    name,
		Content: content,
	}
}

func (e *Evaluator) _for(forStmt *ast.ForStmt, ctx *Context) value.Value {
	forCtx := NewContext(ctx.scope, ctx.absPath)

	var init value.Value
	if forStmt.Init != nil {
		if init = e.Eval(forStmt.Init, forCtx); isError(init) {
			return init
		}
	}

	// Evaluate ElseBlock block if user's condition is false
	if forStmt.Condition != nil {
		cond := e.Eval(forStmt.Condition, forCtx)
		if isError(cond) {
			return cond
		}

		if !isTruthy(cond) && forStmt.ElseBlock != nil {
			return e.Eval(forStmt.ElseBlock, forCtx)
		}
	}

	var blocks strings.Builder

	// Loop through the block until the user's condition is false
	for {
		if forStmt.Condition != nil {
			cond := e.Eval(forStmt.Condition, forCtx)
			if isError(cond) {
				return cond
			}

			if !isTruthy(cond) {
				break
			}
		}

		block := e.Eval(forStmt.Block, forCtx)
		if isError(block) {
			return block
		}

		blocks.WriteString(block.String())

		var post value.Value
		if forStmt.Post != nil {
			post = e.Eval(forStmt.Post, forCtx)
			if isError(post) {
				return post
			}
		}

		if post != nil {
			varName := forStmt.Init.(*ast.AssignStmt).Left.(*ast.Identifier).Name
			if err := forCtx.scope.Set(varName, post); err != nil {
				return e.newError(forStmt, forCtx, "%s", err.Error())
			}
		}

		if hasBreakStmt(block) {
			break
		}

		if hasContinueStmt(block) {
			continue
		}
	}

	return &value.HTML{Val: blocks.String()}
}

func (e *Evaluator) each(eachStmt *ast.EachStmt, ctx *Context) value.Value {
	eachCtx := NewContext(ctx.scope.Child(), ctx.absPath)
	varName := eachStmt.Var.Name

	arrObj := e.Eval(eachStmt.Array, eachCtx)
	if isError(arrObj) {
		return arrObj
	}

	arr, ok := arrObj.(*value.Array)
	if !ok {
		return e.newError(eachStmt, eachCtx, fail.ErrEachDirWithNonArrArg, arrObj.Type())
	}

	arrElems := arr.Elements

	// Evaluate ElseBlock when array is empty
	if len(arrElems) == 0 && eachStmt.ElseBlock != nil {
		return e.Eval(eachStmt.ElseBlock, eachCtx)
	}

	var blocks strings.Builder
	blocks.Grow(len(arrElems))

	for i := range arrElems {
		if err := eachCtx.scope.Set(varName, arrElems[i]); err != nil {
			return e.newError(eachStmt, eachCtx, "%s", err.Error())
		}

		eachCtx.scope.SetLoopVar(map[string]value.Value{
			"index": &value.Int{Val: int64(i)},
			"first": nativeBoolToBoolObj(i == 0),
			"last":  nativeBoolToBoolObj(i == len(arrElems)-1),
			"iter":  &value.Int{Val: int64(i + 1)},
		})

		block := e.Eval(eachStmt.Block, eachCtx)
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

	return &value.HTML{Val: blocks.String()}
}

func (e *Evaluator) breakif(breakifStmt *ast.BreakifStmt, ctx *Context) value.Value {
	cond := e.Eval(breakifStmt.Condition, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return BREAK
	}

	return NIL
}

func (e *Evaluator) continueif(contifStmt *ast.ContinueifStmt, ctx *Context) value.Value {
	cond := e.Eval(contifStmt.Condition, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return CONTINUE
	}

	return NIL
}

func (e *Evaluator) slot(slotStmt *ast.SlotStmt, ctx *Context) value.Value {
	if slotStmt.IsLocal {
		return e.localSlotStmt(slotStmt, ctx)
	}
	return e.externalSlotStmt(slotStmt, ctx)
}

func (e *Evaluator) slotif(slotifStmt *ast.SlotifStmt, ctx *Context) value.Value {
	cond := e.Eval(slotifStmt.Condition, ctx)
	if isError(cond) {
		return cond
	}

	if !isTruthy(cond) {
		return NIL
	}

	block := e.Eval(slotifStmt.Block(), ctx)
	if isError(block) {
		return block
	}

	return &value.Slot{
		Name:    slotifStmt.Name().Val,
		Content: block,
	}
}

func (e *Evaluator) externalSlotStmt(slotStmt *ast.SlotStmt, ctx *Context) value.Value {
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

func (e *Evaluator) localSlotStmt(slotStmt *ast.SlotStmt, ctx *Context) value.Value {
	var block value.Value = NIL

	if slotStmt.Block() != nil {
		block = e.Eval(slotStmt.Block(), ctx)
		if isError(block) {
			return block
		}
	}

	return &value.Slot{
		Name:    slotStmt.Name().Val,
		Content: block,
	}
}

func (e *Evaluator) insert(insertStmt *ast.InsertStmt, ctx *Context) value.Value {
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
// into a single object that we can work with.
func (e *Evaluator) combineInsertContent(insertStmt *ast.InsertStmt, ctx *Context) value.Value {
	if insertStmt.Argument != nil {
		arg := e.Eval(insertStmt.Argument, ctx)
		if isError(arg) {
			return arg
		}
		return arg
	}

	if insertStmt.Block == nil {
		return e.newError(insertStmt, ctx, fail.ErrInsertMustHaveContent)
	}

	return e.Eval(insertStmt.Block, ctx)
}

func (e *Evaluator) dump(dumpStmt *ast.DumpStmt, ctx *Context) value.Value {
	values := make([]string, 0, len(dumpStmt.Arguments))

	for i := range dumpStmt.Arguments {
		evaluated := e.Eval(dumpStmt.Arguments[i], ctx)
		values = append(values, evaluated.Dump(0))
	}

	return &value.Dump{Vals: values}
}

func (e *Evaluator) ident(ident *ast.Identifier, ctx *Context) value.Value {
	varName := ident.Name
	if varName == "global" && e.config != nil && e.config.GlobalData != nil {
		return value.NativeToObject(e.config.GlobalData)
	}

	if val, ok := ctx.scope.Get(varName); ok {
		return val
	}

	return e.newError(ident, ctx, fail.ErrVariableIsUndefined, ident.Name)
}

func (e *Evaluator) indexExp(indexExp *ast.IndexExp, ctx *Context) value.Value {
	left := e.Eval(indexExp.Left, ctx)
	if isError(left) {
		return left
	}

	idx := e.Eval(indexExp.Index, ctx)
	if isError(idx) {
		return idx
	}

	switch {
	case left.Is(value.ARR_OBJ) && idx.Is(value.INT_OBJ):
		return e.arrayIndexExp(left, idx)
	case left.Is(value.OBJ_OBJ) && idx.Is(value.STR_OBJ):
		return e.objectKeyExp(left.(*value.Obj), idx.(*value.Str).Val)
	}

	return e.newError(indexExp, ctx, fail.ErrIndexNotSupported, left.Type())
}

func (e *Evaluator) arrayIndexExp(arrObj, idx value.Value) value.Value {
	arr := arrObj.(*value.Array)
	index := idx.(*value.Int).Val
	max := int64(len(arr.Elements) - 1)

	if index < 0 || index > max {
		return NIL
	}

	return arr.Elements[index]
}

func (e *Evaluator) objectKeyExp(obj *value.Obj, key string) value.Value {
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

func (e *Evaluator) dotExp(dotExp *ast.DotExp, ctx *Context) value.Value {
	left := e.Eval(dotExp.Left, ctx)
	if isError(left) {
		return left
	}

	key := dotExp.Key.(*ast.Identifier)
	obj, ok := left.(*value.Obj)
	if !ok {
		return e.newError(dotExp, ctx, fail.ErrKeyOnNonObject, left.Type(), key)
	}

	return e.objectKeyExp(obj, key.Name)
}

func (e *Evaluator) stringLit(strLit *ast.StringLiteral) value.Value {
	str := html.EscapeString(strLit.Val)

	// unescape single and double quotes
	str = strings.ReplaceAll(str, "&#34;", `"`)
	str = strings.ReplaceAll(str, "&#39;", `'`)

	return &value.Str{Val: str}
}

func (e *Evaluator) prefixExp(prefixExp *ast.PrefixExp, ctx *Context) value.Value {
	right := e.Eval(prefixExp.Right, ctx)
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

func (e *Evaluator) ternaryExp(ternExp *ast.TernaryExp, ctx *Context) value.Value {
	cond := e.Eval(ternExp.Condition, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return e.Eval(ternExp.IfBlock, ctx)
	}

	return e.Eval(ternExp.ElseBlock, ctx)
}

func (e *Evaluator) arrayLit(arrLit *ast.ArrayLiteral, ctx *Context) value.Value {
	elems := e.evalExpressions(arrLit.Elements, ctx)
	if len(elems) == 1 && isError(elems[0]) {
		return elems[0]
	}

	return &value.Array{Elements: elems}
}

func (e *Evaluator) objectLit(objLit *ast.ObjectLiteral, ctx *Context) value.Value {
	pairs := make(map[string]value.Value, len(objLit.Pairs))

	for key, val := range objLit.Pairs {
		valObj := e.Eval(val, ctx)
		if isError(valObj) {
			return valObj
		}

		pairs[key] = valObj
	}

	return value.NewObj(pairs)
}

func (e *Evaluator) evalExpressions(exps []ast.Expression, ctx *Context) []value.Value {
	evaluatedObjs := make([]value.Value, 0, len(exps))

	for i := range exps {
		evaluated := e.Eval(exps[i], ctx)
		if isError(evaluated) {
			return []value.Value{evaluated}
		}

		evaluatedObjs = append(evaluatedObjs, evaluated)
	}

	return evaluatedObjs
}

func (e *Evaluator) infixExp(
	op string,
	leftNode, rightNode ast.Expression,
	ctx *Context,
) value.Value {
	left := e.Eval(leftNode, ctx)
	if isError(left) {
		return left
	}

	if obj, ok := e.shortCircuit(left, op); ok {
		return obj
	}

	right := e.Eval(rightNode, ctx)
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
func (e *Evaluator) shortCircuit(left value.Value, op string) (value.Value, bool) {
	if op == "&&" && !isTruthy(left) {
		return FALSE, true
	}
	if op == "||" && isTruthy(left) {
		return TRUE, true
	}

	return nil, false
}

func (e *Evaluator) postfixExp(postfixExp *ast.PostfixExp, ctx *Context) value.Value {
	leftObj := e.Eval(postfixExp.Left, ctx)
	if isError(leftObj) {
		return leftObj
	}

	return e.postfixOpExp(leftObj, postfixExp.Op, postfixExp, ctx)
}

func (e *Evaluator) callExp(callExp *ast.CallExp, ctx *Context) value.Value {
	receiver := e.Eval(callExp.Receiver, ctx)
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
		res, err := buitin.Fn(receiver, args...)
		if err != nil {
			return e.newError(callExp, ctx, "%s", err.Error())
		}

		return res
	}

	if hasCustomFunc(e.customFunc, receiverType, funcName) {
		nativeArgs := e.objectsToNativeType(args)

		switch r := receiver.(type) {
		case *value.Str:
			fun := e.customFunc.Str[funcName]
			res := fun(r.String(), nativeArgs...)
			return value.NativeToObject(res)
		case *value.Array:
			fun := e.customFunc.Arr[funcName]
			nativeElems := e.objectsToNativeType(r.Elements)
			res := fun(nativeElems, nativeArgs...)
			return value.NativeToObject(res)
		case *value.Int:
			fun := e.customFunc.Int[funcName]
			res := fun(int(r.Val), nativeArgs...)
			return value.NativeToObject(res)
		case *value.Float:
			fun := e.customFunc.Float[funcName]
			res := fun(r.Val, nativeArgs...)
			return value.NativeToObject(res)
		case *value.Bool:
			fun := e.customFunc.Bool[funcName]
			res := fun(r.Val, nativeArgs...)
			return value.NativeToObject(res)
		case *value.Obj:
			fun := e.customFunc.Obj[funcName]
			firstArg := r.Native()
			res := fun(firstArg.(map[string]any), nativeArgs...)
			return value.NativeToObject(res)
		}
	}

	return e.newError(callExp, ctx, fail.ErrFuncNotDefined, receiver.Type(), callExp.Function.Name)
}

func (e *Evaluator) globalCallExp(globalCallExp *ast.GlobalCallExp, ctx *Context) value.Value {
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
	globalCallExp *ast.GlobalCallExp,
	ctx *Context,
) value.Value {
	for i := range globalCallExp.Arguments {
		evaluated := e.Eval(globalCallExp.Arguments[i], ctx)
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
	globalCallExp *ast.GlobalCallExp,
	ctx *Context,
) value.Value {
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
		arg := e.Eval(exp, ctx)
		if isError(arg) {
			return arg
		}

		if !isTruthy(arg) {
			return FALSE
		}
	}

	return TRUE
}

func (e *Evaluator) objectsToNativeType(args []value.Value) []any {
	vals := make([]any, len(args))
	for i := range args {
		vals[i] = args[i].Native()
	}

	return vals
}

func (e *Evaluator) postfixOpExp(
	left value.Value,
	op string,
	node ast.Node,
	ctx *Context,
) value.Value {
	switch op {
	case "++":
		if int, ok := left.(*value.Int); ok {
			return &value.Int{Val: int.Val + 1}
		}

		if fl, ok := left.(*value.Float); ok {
			return &value.Float{Val: fl.Val + 1}
		}
	case "--":
		if int, ok := left.(*value.Int); ok {
			return &value.Int{Val: int.Val - 1}
		}

		if fl, ok := left.(*value.Float); ok {
			float := &value.Float{Val: fl.Val}

			if err := float.SubtractFromFloat(1); err != nil {
				return e.newError(node, ctx, fail.ErrCannotSubFromFloat, float, err)
			}

			return float
		}
	}

	return e.newError(node, ctx, fail.ErrUnknownOp, left.Type(), op)
}

func (e *Evaluator) operatorExp(
	op string,
	left,
	right value.Value,
	leftNode ast.Node,
	ctx *Context,
) value.Value {
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
	left value.Value,
	leftNode ast.Node,
	ctx *Context,
) value.Value {
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
	right value.Value,
	l *value.Int,
	leftNode ast.Node,
	ctx *Context,
) value.Value {
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
	left value.Value,
	leftNode ast.Node,
	ctx *Context,
) value.Value {
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
	case *value.Array, *value.Obj:
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
	right value.Value,
	l *value.Str,
	leftNode ast.Node,
	ctx *Context,
) value.Value {
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
	right value.Value,
	l *value.Float,
	leftNode ast.Node,
	ctx *Context,
) value.Value {
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
	right value.Value,
	node ast.Node,
	ctx *Context,
) value.Value {
	switch r := right.(type) {
	case *value.Int:
		return &value.Int{Val: -r.Val}
	case *value.Float:
		return &value.Float{Val: -r.Val}
	}

	return e.newError(node, ctx, fail.ErrPrefixOpIsWrong, "-", right.Type())
}

func (e *Evaluator) bangOpExp(right value.Value, node ast.Node, ctx *Context) value.Value {
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
