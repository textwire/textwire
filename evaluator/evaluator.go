package evaluator

import (
	"html"
	"strings"

	"github.com/textwire/textwire/v3/ast"
	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/object"
)

var (
	NIL      = &object.Nil{}
	TRUE     = &object.Bool{Value: true}
	FALSE    = &object.Bool{Value: false}
	BREAK    = &object.Break{}
	CONTINUE = &object.Continue{}
)

type Evaluator struct {
	CustomFunc     *config.Func
	UsingTemplates bool

	// Config can be nil when Textwire is used for simple string and
	// file evaluation. If config is not nil, it means we use templates.
	Config *config.Config
}

func New(customFunc *config.Func, conf *config.Config) *Evaluator {
	return &Evaluator{
		CustomFunc:     customFunc,
		Config:         conf,
		UsingTemplates: conf != nil,
	}
}

func (e *Evaluator) Eval(node ast.Node, ctx *Context) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return e.evalProgram(node, ctx)
	case *ast.HTMLStmt:
		return &object.HTML{Value: node.String()}
	case *ast.ExpressionStmt:
		return e.Eval(node.Expression, ctx)
	case *ast.IfStmt:
		return e.evalIfStmt(node, ctx)
	case *ast.BlockStmt:
		return e.evalBlockStmt(node, ctx)
	case *ast.AssignStmt:
		return e.evalAssignStmt(node, ctx)
	case *ast.UseStmt:
		return e.evalUseStmt(node, ctx)
	case *ast.ReserveStmt:
		return e.evalReserveStmt(node, ctx)
	case *ast.ForStmt:
		return e.evalForStmt(node, ctx)
	case *ast.EachStmt:
		return e.evalEachStmt(node, ctx)
	case *ast.BreakIfStmt:
		return e.evalBreakIfStmt(node, ctx)
	case *ast.ComponentStmt:
		return e.evalComponentStmt(node, ctx)
	case *ast.ContinueIfStmt:
		return e.evalContinueIfStmt(node, ctx)
	case *ast.SlotStmt:
		return e.evalSlotStmt(node, ctx)
	case *ast.DumpStmt:
		return e.evalDumpStmt(node, ctx)
	case *ast.InsertStmt:
		return e.evalInsertStmt(node, ctx)
	case *ast.ContinueStmt:
		return CONTINUE
	case *ast.BreakStmt:
		return BREAK
	case *ast.IllegalNode:
		return NIL

	// Expressions
	case *ast.Identifier:
		return e.evalIdentifier(node, ctx)
	case *ast.IndexExp:
		return e.evalIndexExp(node, ctx)
	case *ast.DotExp:
		return e.evalDotExp(node, ctx)
	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return e.evalString(node)
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ObjectLiteral:
		return e.evalObjectLiteral(node, ctx)
	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(node, ctx)
	case *ast.PrefixExp:
		return e.evalPrefixExp(node, ctx)
	case *ast.TernaryExp:
		return e.evalTernaryExp(node, ctx)
	case *ast.InfixExp:
		return e.evalInfixExp(node.Op, node.Left, node.Right, ctx)
	case *ast.PostfixExp:
		return e.evalPostfixExp(node, ctx)
	case *ast.CallExp:
		return e.evalCallExp(node, ctx)
	case *ast.GlobalCallExp:
		return e.evalGlobalCallExp(node, ctx)
	case *ast.NilLiteral:
		return NIL
	}

	return e.newError(node, ctx, fail.ErrUnknownNodeType, node)
}

func (e *Evaluator) evalProgram(prog *ast.Program, ctx *Context) object.Object {
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

func (e *Evaluator) evalIfStmt(node *ast.IfStmt, ctx *Context) object.Object {
	cond := e.Eval(node.Condition, ctx)
	if isError(cond) {
		return cond
	}

	ifCtx := NewContext(ctx.scope.Child(), ctx.absPath)
	if isTruthy(cond) {
		return e.Eval(node.IfBlock, ifCtx)
	}

	for i := range node.ElseIfStmts {
		elseIfNode, ok := node.ElseIfStmts[i].(*ast.ElseIfStmt)
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

	if node.ElseBlock != nil {
		return e.Eval(node.ElseBlock, ifCtx)
	}

	return NIL
}

func (e *Evaluator) evalBlockStmt(node *ast.BlockStmt, ctx *Context) object.Object {
	stmts := make([]object.Object, 0, len(node.Statements))

	for i := range node.Statements {
		stmt := e.Eval(node.Statements[i], ctx)
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

func (e *Evaluator) evalAssignStmt(node *ast.AssignStmt, ctx *Context) object.Object {
	right := e.Eval(node.Right, ctx)
	if isError(right) {
		return right
	}

	if err := ctx.scope.Set(node.Left.Name, right); err != nil {
		return e.newError(node, ctx, "%s", err.Error())
	}

	return NIL
}

func (e *Evaluator) evalUseStmt(node *ast.UseStmt, ctx *Context) object.Object {
	if node.LayoutProg == nil {
		if e.UsingTemplates {
			return e.newError(node, ctx, fail.ErrUseStmtMissingLayout, node.Name.Value)
		}
		return e.newError(node, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	if node.LayoutProg.IsLayout && node.LayoutProg.HasUseStmt() {
		return e.newError(node, ctx, fail.ErrUseStmtNotAllowed)
	}

	useStmtCtx := NewContext(ctx.scope, node.LayoutProg.AbsPath)
	layout := e.Eval(node.LayoutProg, useStmtCtx)
	if isError(layout) {
		return layout
	}

	return &object.Use{
		Path:    node.Name.Value,
		Content: layout,
	}
}

func (e *Evaluator) evalReserveStmt(node *ast.ReserveStmt, ctx *Context) object.Object {
	if !e.UsingTemplates {
		return e.newError(node, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	reserve := &object.Reserve{Name: node.Name.Value}

	// Inserts are optional statements. If not provided, reserve should be empty.
	if node.Insert == nil {
		return NIL
	}

	if node.Insert.Block != nil {
		reserveCtx := NewContext(ctx.scope, node.Insert.AbsPath)
		block := e.Eval(node.Insert.Block, reserveCtx)
		if isError(block) {
			return block
		}

		reserve.Content = block

		return reserve
	}

	if node.Insert.Argument == nil {
		return e.newError(node.Insert, ctx, fail.ErrInsertMustHaveContent)
	}

	firstArg := e.Eval(node.Insert.Argument, ctx)
	if isError(firstArg) {
		return firstArg
	}

	reserve.Argument = firstArg

	return reserve
}

func (e *Evaluator) evalComponentStmt(node *ast.ComponentStmt, ctx *Context) object.Object {
	if !e.UsingTemplates {
		return e.newError(node, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	compName := e.Eval(node.Name, ctx)
	if isError(compName) {
		return compName
	}

	if node.CompProg == nil {
		return e.newError(node, ctx, fail.ErrComponentMustHaveBlock, compName)
	}

	comp := &object.Component{Name: compName.String()}
	compCtx := NewContext(object.NewScope(), node.CompProg.AbsPath)

	if node.Argument != nil {
		for key, arg := range node.Argument.Pairs {
			obj := e.Eval(arg, ctx)
			if isError(obj) {
				return obj
			}

			if err := compCtx.scope.Set(key, obj); err != nil {
				return e.newError(node, ctx, "%s", err.Error())
			}
		}
	}

	blockObj := e.Eval(node.CompProg, compCtx)
	if isError(blockObj) {
		return blockObj
	}

	comp.Content = blockObj

	return comp
}

func (e *Evaluator) evalForStmt(node *ast.ForStmt, ctx *Context) object.Object {
	forCtx := NewContext(ctx.scope.Child(), ctx.absPath)

	var init object.Object
	if node.Init != nil {
		if init = e.Eval(node.Init, forCtx); isError(init) {
			return init
		}
	}

	// Evaluate ElseBlock block if user's condition is false
	if node.Condition != nil {
		cond := e.Eval(node.Condition, forCtx)
		if isError(cond) {
			return cond
		}

		if !isTruthy(cond) && node.ElseBlock != nil {
			return e.Eval(node.ElseBlock, forCtx)
		}
	}

	var blocks strings.Builder

	// Loop through the block until the user's condition is false
	for {
		cond := e.Eval(node.Condition, forCtx)
		if isError(cond) {
			return cond
		}

		if !isTruthy(cond) {
			break
		}

		block := e.Eval(node.Block, forCtx)
		if isError(block) {
			return block
		}

		blocks.WriteString(block.String())
		post := e.Eval(node.Post, forCtx)
		if isError(post) {
			return post
		}

		if node.Init == nil || node.Post == nil {
			continue
		}

		varName := node.Init.(*ast.AssignStmt).Left.Name
		if err := forCtx.scope.Set(varName, post); err != nil {
			return e.newError(node, forCtx, "%s", err.Error())
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

func (e *Evaluator) evalEachStmt(node *ast.EachStmt, ctx *Context) object.Object {
	eachCtx := NewContext(ctx.scope.Child(), ctx.absPath)
	varName := node.Var.Name

	arrObj := e.Eval(node.Array, eachCtx)
	if isError(arrObj) {
		return arrObj
	}

	arr, ok := arrObj.(*object.Array)
	if !ok {
		return e.newError(node, eachCtx, fail.ErrEachDirWithNonArrArg, arrObj.Type())
	}

	arrElems := arr.Elements

	// Evaluate ElseBlock when array is empty
	if len(arrElems) == 0 && node.ElseBlock != nil {
		return e.Eval(node.ElseBlock, eachCtx)
	}

	var blocks strings.Builder
	blocks.Grow(len(arrElems))

	for i := range arrElems {
		if err := eachCtx.scope.Set(varName, arrElems[i]); err != nil {
			return e.newError(node, eachCtx, "%s", err.Error())
		}

		eachCtx.scope.SetLoopVar(map[string]object.Object{
			"index": &object.Int{Value: int64(i)},
			"first": nativeBoolToBooleanObject(i == 0),
			"last":  nativeBoolToBooleanObject(i == len(arrElems)-1),
			"iter":  &object.Int{Value: int64(i + 1)},
		})

		block := e.Eval(node.Block, eachCtx)
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

func (e *Evaluator) evalBreakIfStmt(node *ast.BreakIfStmt, ctx *Context) object.Object {
	cond := e.Eval(node.Condition, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return BREAK
	}

	return NIL
}

func (e *Evaluator) evalContinueIfStmt(node *ast.ContinueIfStmt, ctx *Context) object.Object {
	cond := e.Eval(node.Condition, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return CONTINUE
	}

	return NIL
}

func (e *Evaluator) evalSlotStmt(node *ast.SlotStmt, ctx *Context) object.Object {
	var body object.Object

	if node.Body != nil {
		body = e.Eval(node.Body, ctx)
		if isError(body) {
			return body
		}
	} else {
		body = NIL
	}

	return &object.Slot{Name: node.Name.Value, Content: body}
}

func (e *Evaluator) evalInsertStmt(node *ast.InsertStmt, ctx *Context) object.Object {
	if !e.UsingTemplates {
		return e.newError(node, ctx, fail.ErrSomeDirsOnlyInTemplates)
	}

	// we do not evaluate inserts, they are linked to ast.ReserveStmt as
	// AST programs by this time.
	return NIL
}

func (e *Evaluator) evalDumpStmt(node *ast.DumpStmt, ctx *Context) object.Object {
	values := make([]string, 0, len(node.Arguments))

	for i := range node.Arguments {
		evaluated := e.Eval(node.Arguments[i], ctx)
		values = append(values, evaluated.Dump(0))
	}

	return &object.Dump{Values: values}
}

func (e *Evaluator) evalIdentifier(node *ast.Identifier, ctx *Context) object.Object {
	varName := node.Name
	if varName == "global" && e.Config != nil && e.Config.GlobalData != nil {
		return object.NativeToObject(e.Config.GlobalData)
	}

	if val, ok := ctx.scope.Get(varName); ok {
		return val
	}

	return e.newError(node, ctx, fail.ErrIdentifierIsUndefined, node.Name)
}

func (e *Evaluator) evalIndexExp(node *ast.IndexExp, ctx *Context) object.Object {
	left := e.Eval(node.Left, ctx)
	if isError(left) {
		return left
	}

	idx := e.Eval(node.Index, ctx)
	if isError(idx) {
		return idx
	}

	switch {
	case left.Is(object.ARR_OBJ) && idx.Is(object.INT_OBJ):
		return e.evalArrayIndexExp(left, idx)
	case left.Is(object.OBJ_OBJ) && idx.Is(object.STR_OBJ):
		return e.evalObjectKeyExp(left.(*object.Obj), idx.(*object.Str).Value, node.Index, ctx)
	}

	return e.newError(node, ctx, fail.ErrIndexNotSupported, left.Type())
}

func (e *Evaluator) evalArrayIndexExp(
	arrObj,
	idx object.Object,
) object.Object {
	arr := arrObj.(*object.Array)
	index := idx.(*object.Int).Value
	max := int64(len(arr.Elements) - 1)

	if index < 0 || index > max {
		return NIL
	}

	return arr.Elements[index]
}

func (e *Evaluator) evalObjectKeyExp(
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

	return e.newError(node, ctx, fail.ErrPropertyNotFound, key, object.OBJ_OBJ)
}

func (e *Evaluator) evalDotExp(node *ast.DotExp, ctx *Context) object.Object {
	left := e.Eval(node.Left, ctx)
	if isError(left) {
		return left
	}

	key := node.Key.(*ast.Identifier)
	obj, ok := left.(*object.Obj)
	if !ok {
		return e.newError(node, ctx, fail.ErrPropertyOnNonObject, left.Type(), key)
	}

	return e.evalObjectKeyExp(obj, key.Name, node, ctx)
}

func (e *Evaluator) evalString(node *ast.StringLiteral) object.Object {
	str := html.EscapeString(node.Value)

	// unescape single and double quotes
	str = strings.ReplaceAll(str, "&#34;", `"`)
	str = strings.ReplaceAll(str, "&#39;", `'`)

	return &object.Str{Value: str}
}

func (e *Evaluator) evalPrefixExp(node *ast.PrefixExp, ctx *Context) object.Object {
	right := e.Eval(node.Right, ctx)
	if isError(right) {
		return right
	}

	switch node.Op {
	case "-":
		return e.evalMinusPrefixOpExp(right, node, ctx)
	case "!":
		return e.evalBangOpExp(right, node, ctx)
	}

	return e.newError(node, ctx, fail.ErrUnknownOp, node.Op, right.Type())
}

func (e *Evaluator) evalTernaryExp(node *ast.TernaryExp, ctx *Context) object.Object {
	cond := e.Eval(node.Condition, ctx)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return e.Eval(node.IfBlock, ctx)
	}

	return e.Eval(node.ElseBlock, ctx)
}

func (e *Evaluator) evalArrayLiteral(node *ast.ArrayLiteral, ctx *Context) object.Object {
	elems := e.evalExpressions(node.Elements, ctx)
	if len(elems) == 1 && isError(elems[0]) {
		return elems[0]
	}

	return &object.Array{Elements: elems}
}

func (e *Evaluator) evalObjectLiteral(node *ast.ObjectLiteral, ctx *Context) object.Object {
	pairs := make(map[string]object.Object, len(node.Pairs))

	for key, val := range node.Pairs {
		valObj := e.Eval(val, ctx)
		if isError(valObj) {
			return valObj
		}

		pairs[key] = valObj
	}

	return object.NewObj(pairs)
}

func (e *Evaluator) evalExpressions(exps []ast.Expression, ctx *Context) []object.Object {
	res := make([]object.Object, 0, len(exps))

	for i := range exps {
		evaluated := e.Eval(exps[i], ctx)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		res = append(res, evaluated)
	}

	return res
}

func (e *Evaluator) evalInfixExp(
	op string,
	leftNode,
	rightNode ast.Expression,
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

	return e.evalInfixOpExp(op, left, right, leftNode, ctx)
}

func (e *Evaluator) evalPostfixExp(node *ast.PostfixExp, ctx *Context) object.Object {
	leftObj := e.Eval(node.Left, ctx)
	if isError(leftObj) {
		return leftObj
	}

	return e.evalPostfixOpExp(leftObj, node.Op, node, ctx)
}

func (e *Evaluator) evalCallExp(node *ast.CallExp, ctx *Context) object.Object {
	receiver := e.Eval(node.Receiver, ctx)
	funcName := node.Function.Name
	if isError(receiver) {
		return receiver
	}

	receiverType := receiver.Type()
	typeFuncs, ok := functions[receiverType]
	if !ok {
		return e.newError(node, ctx, fail.ErrFuncNotDefined, receiverType, funcName)
	}

	args := e.evalExpressions(node.Arguments, ctx)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	buitin, ok := typeFuncs[node.Function.Name]
	if ok {
		res, err := buitin.Fn(receiver, args...)
		if err != nil {
			return e.newError(node, ctx, "%s", err.Error())
		}

		return res
	}

	if hasCustomFunc(e.CustomFunc, receiverType, funcName) {
		nativeArgs := e.objectsToNativeType(args)

		switch receiverType {
		case object.STR_OBJ:
			fun := e.CustomFunc.Str[funcName]
			res := fun(receiver.String(), nativeArgs...)
			return object.NativeToObject(res)
		case object.ARR_OBJ:
			fun := e.CustomFunc.Arr[funcName]
			nativeElems := e.objectsToNativeType(receiver.(*object.Array).Elements)
			res := fun(nativeElems, nativeArgs...)
			return object.NativeToObject(res)
		case object.INT_OBJ:
			fun := e.CustomFunc.Int[funcName]
			res := fun(int(receiver.(*object.Int).Value), nativeArgs...)
			return object.NativeToObject(res)
		case object.FLOAT_OBJ:
			fun := e.CustomFunc.Float[funcName]
			res := fun(receiver.(*object.Float).Value, nativeArgs...)
			return object.NativeToObject(res)
		case object.BOOL_OBJ:
			fun := e.CustomFunc.Bool[funcName]
			res := fun(receiver.(*object.Bool).Value, nativeArgs...)
			return object.NativeToObject(res)
		case object.OBJ_OBJ:
			fun := e.CustomFunc.Obj[funcName]
			firstArg := receiver.(*object.Obj).Val()
			res := fun(firstArg.(map[string]any), nativeArgs...)
			return object.NativeToObject(res)
		}
	}

	return e.newError(node, ctx, fail.ErrFuncNotDefined, receiver.Type(), node.Function.Name)
}

func (e *Evaluator) evalGlobalCallExp(node *ast.GlobalCallExp, ctx *Context) object.Object {
	switch node.Function.Name {
	case "defined":
		return e.evalGlobalFuncDefined(node, ctx)
	default:
		return e.newError(node, ctx, fail.ErrGlobalFuncMissing, node.Function.Name)
	}
}

func (e *Evaluator) evalGlobalFuncDefined(node *ast.GlobalCallExp, ctx *Context) object.Object {
	definedVars := make([]bool, 0, len(node.Arguments))
	for i := range node.Arguments {
		evaluated := e.Eval(node.Arguments[i], ctx)
		definedVars = append(definedVars, !isUndefinedVarError(evaluated))
	}

	for _, defined := range definedVars {
		if !defined {
			return FALSE
		}
	}

	return TRUE
}

func (e *Evaluator) objectsToNativeType(args []object.Object) []any {
	res := make([]any, 0, len(args))
	for i := range args {
		res = append(res, args[i].Val())
	}

	return res
}

func (e *Evaluator) evalPostfixOpExp(
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

func (e *Evaluator) evalInfixOpExp(
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
		return e.evalIntegerInfixExp(op, right, left, leftNode, ctx)
	case object.FLOAT_OBJ:
		return e.evalFloatInfixExp(op, right, left, leftNode, ctx)
	case object.BOOL_OBJ:
		return e.evalBooleanInfixExp(op, right, left, leftNode, ctx)
	case object.STR_OBJ:
		return e.evalStringInfixExp(op, right, left, leftNode, ctx)
	}

	return e.newError(leftNode, ctx, fail.ErrUnknownTypeForOp, left.Type(), op)
}

func (e *Evaluator) evalBooleanInfixExp(
	op string,
	right,
	left object.Object,
	leftNode ast.Node,
	ctx *Context,
) object.Object {
	leftVal := left.(*object.Bool).Value
	rightVal := right.(*object.Bool).Value

	switch op {
	case "&&":
		return &object.Bool{Value: leftVal && rightVal}
	case "||":
		return &object.Bool{Value: leftVal || rightVal}
	}

	return e.newError(leftNode, ctx, fail.ErrUnknownTypeForOp, left.Type(), op)
}

func (e *Evaluator) evalIntegerInfixExp(
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
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	}

	return e.newError(leftNode, ctx, fail.ErrUnknownTypeForOp, left.Type(), op)
}

func (e *Evaluator) evalStringInfixExp(
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
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "+":
		return &object.Str{Value: leftVal + rightVal}
	}

	return e.newError(leftNode, ctx, fail.ErrUnknownTypeForOp, left.Type(), op)
}

func (e *Evaluator) evalFloatInfixExp(
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
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	}

	return e.newError(leftNode, ctx, fail.ErrUnknownTypeForOp, left.Type(), op)
}

func (e *Evaluator) evalMinusPrefixOpExp(
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

func (e *Evaluator) evalBangOpExp(
	right object.Object,
	node ast.Node,
	ctx *Context,
) object.Object {
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
