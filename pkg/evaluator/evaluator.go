package evaluator

import (
	"html"
	"strings"

	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/pkg/ast"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/object"
)

var (
	NIL      = &object.Nil{}
	TRUE     = &object.Bool{Value: true}
	FALSE    = &object.Bool{Value: false}
	BREAK    = &object.Break{}
	CONTINUE = &object.Continue{}
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

func (e *Evaluator) Eval(node ast.Node, ctx *Context) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return e.program(node, ctx)
	case *ast.HTMLStmt:
		return &object.HTML{Value: node.String()}
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
	case *ast.BreakIfStmt:
		return e.breakIf(node, ctx)
	case *ast.ComponentStmt:
		return e.component(node, ctx)
	case *ast.ContinueIfStmt:
		return e.continueIf(node, ctx)
	case *ast.SlotStmt:
		return e.slot(node, ctx)
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
		return &object.Int{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return e.stringLit(node)
	case *ast.BooleanLiteral:
		return nativeBoolToBoolObj(node.Value)
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

func (e *Evaluator) program(prog *ast.Program, ctx *Context) object.Object {
	var stmts strings.Builder
	stmts.Grow(len(prog.Statements))

	for i := range prog.Statements {
		stmt := e.Eval(prog.Statements[i], ctx)
		if isError(stmt) {
			return stmt
		}

		stmts.WriteString(stmt.String())
	}

	return &object.HTML{Value: stmts.String()}
}

func (e *Evaluator) _if(ifStmt *ast.IfStmt, ctx *Context) object.Object {
	cond := e.Eval(ifStmt.Condition, ctx)
	if isError(cond) {
		return cond
	}

	ifCtx := NewContext(ctx.scope.Child(), ctx.absPath)
	if isTruthy(cond) {
		return e.Eval(ifStmt.IfBlock, ifCtx)
	}

	for i := range ifStmt.ElseIfStmts {
		elseIfNode, ok := ifStmt.ElseIfStmts[i].(*ast.ElseIfStmt)
		if !ok {
			continue
		}

		cond = e.Eval(elseIfNode.Condition, ifCtx)
		if isError(cond) {
			return cond
		}

		if isTruthy(cond) {
			return e.Eval(elseIfNode.Block, ifCtx)
		}
	}

	if ifStmt.ElseBlock != nil {
		return e.Eval(ifStmt.ElseBlock, ifCtx)
	}

	return NIL
}

func (e *Evaluator) block(blockStmt *ast.BlockStmt, ctx *Context) object.Object {
	stmts := make([]object.Object, 0, len(blockStmt.Statements))

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

	return &object.Block{Elements: stmts}
}

func (e *Evaluator) assign(assignStmt *ast.AssignStmt, ctx *Context) object.Object {
	right := e.Eval(assignStmt.Right, ctx)
	if isError(right) {
		return right
	}

	if err := ctx.scope.Set(assignStmt.Left.Name, right); err != nil {
		return e.newError(assignStmt, ctx, "%s", err.Error())
	}

	return NIL
}

func (e *Evaluator) use(useStmt *ast.UseStmt, ctx *Context) object.Object {
	if useStmt.LayoutProg == nil {
		if e.usingTemplates {
			return e.newError(useStmt, ctx, fail.ErrUseStmtMissingLayout, useStmt.Name.Value)
		}
		return e.newError(useStmt, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	e.usingUseStmt = true

	// Make sure that layout is missing @use
	if useStmt.LayoutProg.IsLayout && useStmt.LayoutProg.HasUseStmt() {
		return e.newError(useStmt, ctx, fail.ErrUseStmtNotAllowed)
	}

	// Create new layout context and pass inserts to it
	layoutCtx := NewContext(object.NewScope(), useStmt.LayoutProg.AbsPath)

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

	return &object.Use{
		Path:   useStmt.Name.Value,
		Layout: layout,
	}
}

func (e *Evaluator) reserve(reserveStmt *ast.ReserveStmt, ctx *Context) object.Object {
	if !e.usingTemplates {
		return e.newError(reserveStmt, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	name := reserveStmt.Name.Value
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

	return &object.Reserve{
		Name:   name,
		Insert: insert,
	}
}

func (e *Evaluator) component(compStmt *ast.ComponentStmt, ctx *Context) object.Object {
	if !e.usingTemplates {
		return e.newError(compStmt, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	name := compStmt.Name.Value
	if compStmt.CompProg == nil {
		return e.newError(compStmt, ctx, fail.ErrUndefinedComponent, name)
	}

	compCtx := NewContext(object.NewScope(), compStmt.CompProg.AbsPath)

	// Evaluate local slots and add them to component context
	for _, slotStmt := range compStmt.Slots {
		slot := e.Eval(slotStmt, ctx)
		if isError(slot) {
			return slot
		}

		if compCtx.slots[name] == nil {
			compCtx.slots[name] = map[string]object.Object{}
		}

		compCtx.slots[name][slotStmt.Name.Value] = slot
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

	return &object.Component{
		Name:    name,
		Content: content,
	}
}

func (e *Evaluator) _for(forStmt *ast.ForStmt, ctx *Context) object.Object {
	forCtx := NewContext(ctx.scope.Child(), ctx.absPath)

	var init object.Object
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
		cond := e.Eval(forStmt.Condition, forCtx)
		if isError(cond) {
			return cond
		}

		if !isTruthy(cond) {
			break
		}

		block := e.Eval(forStmt.Block, forCtx)
		if isError(block) {
			return block
		}

		blocks.WriteString(block.String())
		post := e.Eval(forStmt.Post, forCtx)
		if isError(post) {
			return post
		}

		if forStmt.Init == nil || forStmt.Post == nil {
			continue
		}

		varName := forStmt.Init.(*ast.AssignStmt).Left.Name
		if err := forCtx.scope.Set(varName, post); err != nil {
			return e.newError(forStmt, forCtx, "%s", err.Error())
		}

		if hasBreakStmt(block) {
			break
		}

		if hasContinueStmt(block) {
			continue
		}
	}

	return &object.HTML{Value: blocks.String()}
}

func (e *Evaluator) each(eachStmt *ast.EachStmt, ctx *Context) object.Object {
	eachCtx := NewContext(ctx.scope.Child(), ctx.absPath)
	varName := eachStmt.Var.Name

	arrObj := e.Eval(eachStmt.Array, eachCtx)
	if isError(arrObj) {
		return arrObj
	}

	arr, ok := arrObj.(*object.Array)
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

		eachCtx.scope.SetLoopVar(map[string]object.Object{
			"index": &object.Int{Value: int64(i)},
			"first": nativeBoolToBoolObj(i == 0),
			"last":  nativeBoolToBoolObj(i == len(arrElems)-1),
			"iter":  &object.Int{Value: int64(i + 1)},
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

	return &object.HTML{Value: blocks.String()}
}

func (e *Evaluator) breakIf(breakIfStmt *ast.BreakIfStmt, ctx *Context) object.Object {
	cond := e.Eval(breakIfStmt.Condition, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return BREAK
	}

	return NIL
}

func (e *Evaluator) continueIf(contIfStmt *ast.ContinueIfStmt, ctx *Context) object.Object {
	cond := e.Eval(contIfStmt.Condition, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return CONTINUE
	}

	return NIL
}

func (e *Evaluator) slot(slotStmt *ast.SlotStmt, ctx *Context) object.Object {
	if slotStmt.IsLocal {
		return e.localSlotStmt(slotStmt, ctx)
	}
	return e.externalSlotStmt(slotStmt, ctx)
}

func (e *Evaluator) externalSlotStmt(slotStmt *ast.SlotStmt, ctx *Context) object.Object {
	name := slotStmt.Name.Value
	compName := slotStmt.CompName

	// Get slot's content from the context
	content, ok := ctx.slots[compName][name]
	if !ok {
		// Slots are optional in component files since v3.1.0
		return NIL
	}

	// delete slot after it's been used by external component
	defer delete(ctx.slots[compName], name)

	return &object.Slot{Name: name, Content: content}
}

func (e *Evaluator) localSlotStmt(slotStmt *ast.SlotStmt, ctx *Context) object.Object {
	var block object.Object = NIL

	if slotStmt.Block != nil {
		block = e.Eval(slotStmt.Block, ctx)
		if isError(block) {
			return block
		}
	}

	return &object.Slot{
		Name:    slotStmt.Name.Value,
		Content: block,
	}
}

func (e *Evaluator) insert(insertStmt *ast.InsertStmt, ctx *Context) object.Object {
	if !e.usingTemplates {
		return e.newError(insertStmt, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	name := insertStmt.Name.Value
	if !e.usingUseStmt {
		return e.newError(insertStmt, ctx, fail.ErrInsertRequiresUse, name)
	}

	block := e.combineInsertContent(insertStmt, ctx)
	if isError(block) {
		return block
	}

	return &object.Insert{
		Name:  name,
		Block: block,
	}
}

// combineInsertContent combines insert Argument and Block (depending what user has)
// into a single object that we can work with.
func (e *Evaluator) combineInsertContent(insertStmt *ast.InsertStmt, ctx *Context) object.Object {
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

func (e *Evaluator) dump(dumpStmt *ast.DumpStmt, ctx *Context) object.Object {
	values := make([]string, 0, len(dumpStmt.Arguments))

	for i := range dumpStmt.Arguments {
		evaluated := e.Eval(dumpStmt.Arguments[i], ctx)
		values = append(values, evaluated.Dump(0))
	}

	return &object.Dump{Values: values}
}

func (e *Evaluator) ident(ident *ast.Identifier, ctx *Context) object.Object {
	varName := ident.Name
	if varName == "global" && e.config != nil && e.config.GlobalData != nil {
		return object.NativeToObject(e.config.GlobalData)
	}

	if val, ok := ctx.scope.Get(varName); ok {
		return val
	}

	return e.newError(ident, ctx, fail.ErrVariableIsUndefined, ident.Name)
}

func (e *Evaluator) indexExp(indexExp *ast.IndexExp, ctx *Context) object.Object {
	left := e.Eval(indexExp.Left, ctx)
	if isError(left) {
		return left
	}

	idx := e.Eval(indexExp.Index, ctx)
	if isError(idx) {
		return idx
	}

	switch {
	case left.Is(object.ARR_OBJ) && idx.Is(object.INT_OBJ):
		return e.arrayIndexExp(left, idx)
	case left.Is(object.OBJ_OBJ) && idx.Is(object.STR_OBJ):
		return e.objectKeyExp(left.(*object.Obj), idx.(*object.Str).Value, indexExp.Index, ctx)
	}

	return e.newError(indexExp, ctx, fail.ErrIndexNotSupported, left.Type())
}

func (e *Evaluator) arrayIndexExp(arrObj, idx object.Object) object.Object {
	arr := arrObj.(*object.Array)
	index := idx.(*object.Int).Value
	max := int64(len(arr.Elements) - 1)

	if index < 0 || index > max {
		return NIL
	}

	return arr.Elements[index]
}

func (e *Evaluator) objectKeyExp(
	obj *object.Obj,
	key string,
	node ast.Node,
	ctx *Context,
) object.Object {
	// First, try to get key as it is with regular case.
	if pair, ok := obj.Pairs[key]; ok {
		return pair
	}

	// Makes the first letter uppercase.
	upperFirstCharKey := strings.ToUpper(key[:1]) + key[1:]

	if pair, ok := obj.Pairs[upperFirstCharKey]; ok {
		return pair
	}

	return NIL // Undefined props result in nil
}

func (e *Evaluator) dotExp(dotExp *ast.DotExp, ctx *Context) object.Object {
	left := e.Eval(dotExp.Left, ctx)
	if isError(left) {
		return left
	}

	key := dotExp.Key.(*ast.Identifier)
	obj, ok := left.(*object.Obj)
	if !ok {
		return e.newError(dotExp, ctx, fail.ErrPropertyOnNonObject, left.Type(), key)
	}

	return e.objectKeyExp(obj, key.Name, dotExp, ctx)
}

func (e *Evaluator) stringLit(strLit *ast.StringLiteral) object.Object {
	str := html.EscapeString(strLit.Value)

	// unescape single and double quotes
	str = strings.ReplaceAll(str, "&#34;", `"`)
	str = strings.ReplaceAll(str, "&#39;", `'`)

	return &object.Str{Value: str}
}

func (e *Evaluator) prefixExp(prefixExp *ast.PrefixExp, ctx *Context) object.Object {
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

func (e *Evaluator) ternaryExp(ternExp *ast.TernaryExp, ctx *Context) object.Object {
	cond := e.Eval(ternExp.Condition, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return e.Eval(ternExp.IfBlock, ctx)
	}

	return e.Eval(ternExp.ElseBlock, ctx)
}

func (e *Evaluator) arrayLit(arrLit *ast.ArrayLiteral, ctx *Context) object.Object {
	elems := e.evalExpressions(arrLit.Elements, ctx)
	if len(elems) == 1 && isError(elems[0]) {
		return elems[0]
	}

	return &object.Array{Elements: elems}
}

func (e *Evaluator) objectLit(objLit *ast.ObjectLiteral, ctx *Context) object.Object {
	pairs := make(map[string]object.Object, len(objLit.Pairs))

	for key, val := range objLit.Pairs {
		valObj := e.Eval(val, ctx)
		if isError(valObj) {
			return valObj
		}

		pairs[key] = valObj
	}

	return object.NewObj(pairs)
}

func (e *Evaluator) evalExpressions(exps []ast.Expression, ctx *Context) []object.Object {
	result := make([]object.Object, 0, len(exps))

	for i := range exps {
		evaluated := e.Eval(exps[i], ctx)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) infixExp(
	op string,
	leftNode, rightNode ast.Expression,
	ctx *Context,
) object.Object {
	left := e.Eval(leftNode, ctx)
	if isError(left) {
		return left
	}

	right := e.Eval(rightNode, ctx)
	if isError(right) {
		return right
	}

	if op == "&&" || op == "||" {
		return e.logicalExp(op, right, left, leftNode, ctx)
	}

	return e.infixOpExp(op, left, right, leftNode, ctx)
}

func (e *Evaluator) postfixExp(postfixExp *ast.PostfixExp, ctx *Context) object.Object {
	leftObj := e.Eval(postfixExp.Left, ctx)
	if isError(leftObj) {
		return leftObj
	}

	return e.postfixOpExp(leftObj, postfixExp.Op, postfixExp, ctx)
}

func (e *Evaluator) callExp(callExp *ast.CallExp, ctx *Context) object.Object {
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

		switch receiverType {
		case object.STR_OBJ:
			fun := e.customFunc.Str[funcName]
			res := fun(receiver.String(), nativeArgs...)
			return object.NativeToObject(res)
		case object.ARR_OBJ:
			fun := e.customFunc.Arr[funcName]
			nativeElems := e.objectsToNativeType(receiver.(*object.Array).Elements)
			res := fun(nativeElems, nativeArgs...)
			return object.NativeToObject(res)
		case object.INT_OBJ:
			fun := e.customFunc.Int[funcName]
			res := fun(int(receiver.(*object.Int).Value), nativeArgs...)
			return object.NativeToObject(res)
		case object.FLOAT_OBJ:
			fun := e.customFunc.Float[funcName]
			res := fun(receiver.(*object.Float).Value, nativeArgs...)
			return object.NativeToObject(res)
		case object.BOOL_OBJ:
			fun := e.customFunc.Bool[funcName]
			res := fun(receiver.(*object.Bool).Value, nativeArgs...)
			return object.NativeToObject(res)
		case object.OBJ_OBJ:
			fun := e.customFunc.Obj[funcName]
			firstArg := receiver.(*object.Obj).Val()
			res := fun(firstArg.(map[string]any), nativeArgs...)
			return object.NativeToObject(res)
		}
	}

	return e.newError(callExp, ctx, fail.ErrFuncNotDefined, receiver.Type(), callExp.Function.Name)
}

func (e *Evaluator) globalCallExp(globalCallExp *ast.GlobalCallExp, ctx *Context) object.Object {
	switch globalCallExp.Function.Name {
	case "defined":
		return e.globalFuncDefined(globalCallExp, ctx)
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
) object.Object {
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

func (e *Evaluator) objectsToNativeType(args []object.Object) []any {
	result := make([]any, 0, len(args))
	for i := range args {
		result = append(result, args[i].Val())
	}

	return result
}

func (e *Evaluator) postfixOpExp(
	left object.Object,
	op string,
	node ast.Node,
	ctx *Context,
) object.Object {
	if op == "++" {
		if left.Is(object.INT_OBJ) {
			val := left.(*object.Int).Value + 1
			return &object.Int{Value: val}
		}

		if left.Is(object.FLOAT_OBJ) {
			val := left.(*object.Float).Value + 1
			return &object.Float{Value: val}
		}
	}

	if op == "--" {
		if left.Is(object.INT_OBJ) {
			val := left.(*object.Int).Value - 1
			return &object.Int{Value: val}
		}

		if left.Is(object.FLOAT_OBJ) {
			val := left.(*object.Float).Value
			float := &object.Float{Value: val}

			if err := float.SubtractFromFloat(1); err != nil {
				return e.newError(node, ctx, fail.ErrCannotSubFromFloat, float, err)
			}

			return float
		}
	}

	return e.newError(node, ctx, fail.ErrUnknownOp, left.Type(), op)
}

func (e *Evaluator) infixOpExp(
	op string,
	left,
	right object.Object,
	leftNode ast.Node,
	ctx *Context,
) object.Object {
	if left.Type() != right.Type() {
		return e.newError(leftNode, ctx, fail.ErrTypeMismatch, left.Type(), op, right.Type())
	}

	switch left.Type() {
	case object.INT_OBJ:
		return e.intInfixExp(op, right, left, leftNode, ctx)
	case object.FLOAT_OBJ:
		return e.floatInfixExp(op, right, left, leftNode, ctx)
	case object.STR_OBJ:
		return e.stringInfixExp(op, right, left, leftNode, ctx)
	}

	return e.newError(leftNode, ctx, fail.ErrUnknownTypeForOp, left.Type(), op)
}

func (e *Evaluator) logicalExp(
	op string,
	right,
	left object.Object,
	leftNode ast.Node,
	ctx *Context,
) object.Object {
	switch op {
	case "&&":
		return &object.Bool{Value: isTruthy(left) && isTruthy(right)}
	case "||":
		return &object.Bool{Value: isTruthy(left) || isTruthy(right)}
	}

	return e.newError(leftNode, ctx, fail.ErrUnknownTypeForOp, left.Type(), op)
}

func (e *Evaluator) intInfixExp(
	op string,
	right,
	left object.Object,
	leftNode ast.Node,
	ctx *Context,
) object.Object {
	leftVal := left.(*object.Int).Value
	rightVal := right.(*object.Int).Value

	switch op {
	case "+":
		return &object.Int{Value: leftVal + rightVal}
	case "-":
		return &object.Int{Value: leftVal - rightVal}
	case "*":
		return &object.Int{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return e.newError(leftNode, ctx, fail.ErrDivisionByZero)
		}
		return &object.Int{Value: leftVal / rightVal}
	case "%":
		return &object.Int{Value: leftVal % rightVal}
	case "==":
		return nativeBoolToBoolObj(leftVal == rightVal)
	case "!=":
		return nativeBoolToBoolObj(leftVal != rightVal)
	case ">":
		return nativeBoolToBoolObj(leftVal > rightVal)
	case "<":
		return nativeBoolToBoolObj(leftVal < rightVal)
	case ">=":
		return nativeBoolToBoolObj(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBoolObj(leftVal <= rightVal)
	}

	return e.newError(leftNode, ctx, fail.ErrUnknownTypeForOp, left.Type(), op)
}

func (e *Evaluator) stringInfixExp(
	op string,
	right,
	left object.Object,
	leftNode ast.Node,
	ctx *Context,
) object.Object {
	leftVal := left.(*object.Str).Value
	rightVal := right.(*object.Str).Value

	switch op {
	case "==":
		return nativeBoolToBoolObj(leftVal == rightVal)
	case "!=":
		return nativeBoolToBoolObj(leftVal != rightVal)
	case "+":
		return &object.Str{Value: leftVal + rightVal}
	}

	return e.newError(leftNode, ctx, fail.ErrUnknownTypeForOp, left.Type(), op)
}

func (e *Evaluator) floatInfixExp(
	op string,
	right,
	left object.Object,
	leftNode ast.Node,
	ctx *Context,
) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch op {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "==":
		return nativeBoolToBoolObj(leftVal == rightVal)
	case "!=":
		return nativeBoolToBoolObj(leftVal != rightVal)
	case ">":
		return nativeBoolToBoolObj(leftVal > rightVal)
	case "<":
		return nativeBoolToBoolObj(leftVal < rightVal)
	case ">=":
		return nativeBoolToBoolObj(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBoolObj(leftVal <= rightVal)
	}

	return e.newError(leftNode, ctx, fail.ErrUnknownTypeForOp, left.Type(), op)
}

func (e *Evaluator) minusPrefixOpExp(
	right object.Object,
	node ast.Node,
	ctx *Context,
) object.Object {
	switch right.Type() {
	case object.INT_OBJ:
		val := right.(*object.Int).Value
		return &object.Int{Value: -val}
	case object.FLOAT_OBJ:
		val := right.(*object.Float).Value
		return &object.Float{Value: -val}
	}

	return e.newError(node, ctx, fail.ErrPrefixOpIsWrong, "-", right.Type())
}

func (e *Evaluator) bangOpExp(right object.Object, node ast.Node, ctx *Context) object.Object {
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

func (e *Evaluator) newError(node ast.Node, ctx *Context, format string, a ...any) *object.Error {
	return &object.Error{
		Err:     fail.New(node.Line(), ctx.absPath, "evaluator", format, a...),
		ErrorID: format,
	}
}
