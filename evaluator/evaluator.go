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

func (e *Evaluator) Eval(node ast.Node, scope *object.Scope, path string) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return e.evalProgram(node, scope)
	case *ast.HTMLStmt:
		return &object.HTML{Value: node.String()}
	case *ast.ExpressionStmt:
		return e.Eval(node.Expression, scope, path)
	case *ast.IfStmt:
		return e.evalIfStmt(node, scope, path)
	case *ast.BlockStmt:
		return e.evalBlockStmt(node, scope, path)
	case *ast.AssignStmt:
		return e.evalAssignStmt(node, scope, path)
	case *ast.UseStmt:
		return e.evalUseStmt(node, scope, path)
	case *ast.ReserveStmt:
		return e.evalReserveStmt(node, scope, path)
	case *ast.ForStmt:
		return e.evalForStmt(node, scope, path)
	case *ast.EachStmt:
		return e.evalEachStmt(node, scope, path)
	case *ast.BreakIfStmt:
		return e.evalBreakIfStmt(node, scope, path)
	case *ast.ComponentStmt:
		return e.evalComponentStmt(node, scope, path)
	case *ast.ContinueIfStmt:
		return e.evalContinueIfStmt(node, scope, path)
	case *ast.SlotStmt:
		return e.evalSlotStmt(node, scope, path)
	case *ast.DumpStmt:
		return e.evalDumpStmt(node, scope, path)
	case *ast.InsertStmt:
		return e.evalInsertStmt(node, path)
	case *ast.ContinueStmt:
		return CONTINUE
	case *ast.BreakStmt:
		return BREAK
	case *ast.IllegalNode:
		return NIL

	// Expressions
	case *ast.Identifier:
		return e.evalIdentifier(node, scope, path)
	case *ast.IndexExp:
		return e.evalIndexExp(node, scope, path)
	case *ast.DotExp:
		return e.evalDotExp(node, scope, path)
	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return e.evalString(node, scope)
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ObjectLiteral:
		return e.evalObjectLiteral(node, scope, path)
	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(node, scope, path)
	case *ast.PrefixExp:
		return e.evalPrefixExp(node, scope, path)
	case *ast.TernaryExp:
		return e.evalTernaryExp(node, scope, path)
	case *ast.InfixExp:
		return e.evalInfixExp(node.Operator, node.Left, node.Right, scope, path)
	case *ast.PostfixExp:
		return e.evalPostfixExp(node, scope, path)
	case *ast.CallExp:
		return e.evalCallExp(node, scope, path)
	case *ast.GlobalCallExp:
		return e.evalGlobalCallExp(node, scope, path)
	case *ast.NilLiteral:
		return NIL
	}

	return e.newError(node, path, fail.ErrUnknownNodeType, node)
}

func (e *Evaluator) evalProgram(prog *ast.Program, scope *object.Scope) object.Object {
	var stmts strings.Builder
	stmts.Grow(len(prog.Statements))

	for i := range prog.Statements {
		stmt := e.Eval(prog.Statements[i], scope, prog.AbsPath)
		if isError(stmt) {
			return stmt
		}

		stmts.WriteString(stmt.String())
	}

	return &object.HTML{Value: stmts.String()}
}

func (e *Evaluator) evalIfStmt(node *ast.IfStmt, scope *object.Scope, path string) object.Object {
	cond := e.Eval(node.Condition, scope, path)
	if isError(cond) {
		return cond
	}

	childScope := scope.Child()
	if isTruthy(cond) {
		return e.Eval(node.IfBlock, childScope, path)
	}

	for i := range node.ElseIfStmts {
		elseIfNode, ok := node.ElseIfStmts[i].(*ast.ElseIfStmt)
		if !ok {
			continue
		}

		cond = e.Eval(elseIfNode.Condition, scope, path)
		if isError(cond) {
			return cond
		}

		if isTruthy(cond) {
			return e.Eval(elseIfNode.Block, childScope, path)
		}
	}

	if node.ElseBlock != nil {
		return e.Eval(node.ElseBlock, childScope, path)
	}

	return NIL
}

func (e *Evaluator) evalBlockStmt(
	node *ast.BlockStmt,
	scope *object.Scope,
	path string,
) object.Object {
	stmts := make([]object.Object, 0, len(node.Statements))

	for i := range node.Statements {
		stmt := e.Eval(node.Statements[i], scope, path)
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

func (e *Evaluator) evalAssignStmt(
	node *ast.AssignStmt,
	scope *object.Scope,
	path string,
) object.Object {
	right := e.Eval(node.Right, scope, path)
	if isError(right) {
		return right
	}

	if err := scope.Set(node.Left.Name, right); err != nil {
		return e.newError(node, path, "%s", err.Error())
	}

	return NIL
}

func (e *Evaluator) evalUseStmt(node *ast.UseStmt, scope *object.Scope, path string) object.Object {
	if node.Attachment == nil {
		if e.UsingTemplates {
			return e.newError(node, path, fail.ErrUseStmtMissingLayout, node.Name.Value)
		}
		return e.newError(node, path, fail.ErrSomeDirsOnlyInTemplates)
	}

	if node.Attachment.IsLayout && node.Attachment.HasUseStmt() {
		return e.newError(node, path, fail.ErrUseStmtNotAllowed)
	}

	layout := e.Eval(node.Attachment, scope, node.Attachment.AbsPath)
	if isError(layout) {
		return layout
	}

	return &object.Use{
		Path:    node.Name.Value,
		Content: layout,
	}
}

func (e *Evaluator) evalReserveStmt(
	node *ast.ReserveStmt,
	scope *object.Scope,
	path string,
) object.Object {
	if !e.UsingTemplates {
		return e.newError(node, path, fail.ErrSomeDirsOnlyInTemplates)
	}

	reserve := &object.Reserve{Name: node.Name.Value}

	// Inserts are optional statements. If not provided, reserve should be empty.
	if node.Insert == nil {
		return NIL
	}

	if node.Insert.Block != nil {
		block := e.Eval(node.Insert.Block, scope, node.Insert.AbsPath)
		if isError(block) {
			return block
		}

		reserve.Content = block

		return reserve
	}

	if node.Insert.Argument == nil {
		return e.newError(node.Insert, path, fail.ErrInsertMustHaveContent)
	}

	firstArg := e.Eval(node.Insert.Argument, scope, path)
	if isError(firstArg) {
		return firstArg
	}

	reserve.Argument = firstArg

	return reserve
}

func (e *Evaluator) evalComponentStmt(
	node *ast.ComponentStmt,
	scope *object.Scope,
	path string,
) object.Object {
	if !e.UsingTemplates {
		return e.newError(node, path, fail.ErrSomeDirsOnlyInTemplates)
	}

	compName := e.Eval(node.Name, scope, path)
	if isError(compName) {
		return compName
	}

	if node.Attachment == nil {
		return e.newError(node, path, fail.ErrComponentMustHaveBlock, compName)
	}

	comp := &object.Component{Name: compName.String()}
	childScope := object.NewScope()

	if node.Argument != nil {
		for key, arg := range node.Argument.Pairs {
			obj := e.Eval(arg, scope, path)
			if isError(obj) {
				return obj
			}

			err := childScope.Set(key, obj)
			if err != nil {
				return e.newError(node, path, "%s", err.Error())
			}
		}
	}

	blockObj := e.Eval(node.Attachment, childScope, node.Attachment.AbsPath)
	if isError(blockObj) {
		return blockObj
	}

	comp.Content = blockObj

	return comp
}

func (e *Evaluator) evalForStmt(node *ast.ForStmt, scope *object.Scope, path string) object.Object {
	childScope := scope.Child()

	var init object.Object
	if node.Init != nil {
		if init = e.Eval(node.Init, childScope, path); isError(init) {
			return init
		}
	}

	// Evaluate ElseBlock block if user's condition is false
	if node.Condition != nil {
		cond := e.Eval(node.Condition, childScope, path)
		if isError(cond) {
			return cond
		}

		if !isTruthy(cond) && node.ElseBlock != nil {
			return e.Eval(node.ElseBlock, childScope, path)
		}
	}

	var blocks strings.Builder

	// Loop through the block until the user's condition is false
	for {
		cond := e.Eval(node.Condition, childScope, path)
		if isError(cond) {
			return cond
		}

		if !isTruthy(cond) {
			break
		}

		block := e.Eval(node.Block, childScope, path)
		if isError(block) {
			return block
		}

		blocks.WriteString(block.String())
		post := e.Eval(node.Post, childScope, path)
		if isError(post) {
			return post
		}

		if node.Init == nil || node.Post == nil {
			continue
		}

		varName := node.Init.(*ast.AssignStmt).Left.Name
		if err := childScope.Set(varName, post); err != nil {
			return e.newError(node, path, "%s", err.Error())
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

func (e *Evaluator) evalEachStmt(
	node *ast.EachStmt,
	scope *object.Scope,
	path string,
) object.Object {
	childScope := scope.Child()
	varName := node.Var.Name

	arrObj := e.Eval(node.Array, childScope, path)
	if isError(arrObj) {
		return arrObj
	}

	arr, ok := arrObj.(*object.Array)
	if !ok {
		return e.newError(node, path, fail.ErrEachDirWithNonArrArg, arrObj.Type())
	}

	arrElems := arr.Elements

	// Evaluate ElseBlock when array is empty
	if len(arrElems) == 0 && node.ElseBlock != nil {
		return e.Eval(node.ElseBlock, childScope, path)
	}

	var blocks strings.Builder
	blocks.Grow(len(arrElems))

	for i := range arrElems {
		if err := childScope.Set(varName, arrElems[i]); err != nil {
			return e.newError(node, path, "%s", err.Error())
		}

		childScope.SetLoopVar(map[string]object.Object{
			"index": &object.Int{Value: int64(i)},
			"first": nativeBoolToBooleanObject(i == 0),
			"last":  nativeBoolToBooleanObject(i == len(arrElems)-1),
			"iter":  &object.Int{Value: int64(i + 1)},
		})

		block := e.Eval(node.Block, childScope, path)
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

func (e *Evaluator) evalBreakIfStmt(
	node *ast.BreakIfStmt,
	scope *object.Scope,
	path string,
) object.Object {
	cond := e.Eval(node.Condition, scope, path)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return BREAK
	}

	return NIL
}

func (e *Evaluator) evalContinueIfStmt(
	node *ast.ContinueIfStmt,
	scope *object.Scope,
	path string,
) object.Object {
	cond := e.Eval(node.Condition, scope, path)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return CONTINUE
	}

	return NIL
}

func (e *Evaluator) evalSlotStmt(
	node *ast.SlotStmt,
	scope *object.Scope,
	path string,
) object.Object {
	var body object.Object

	if node.Body != nil {
		body = e.Eval(node.Body, scope, path)
		if isError(body) {
			return body
		}
	} else {
		body = NIL
	}

	return &object.Slot{Name: node.Name.Value, Content: body}
}

func (e *Evaluator) evalInsertStmt(node *ast.InsertStmt, path string) object.Object {
	if !e.UsingTemplates {
		return e.newError(node, path, fail.ErrSomeDirsOnlyInTemplates)
	}

	// we do not evaluate inserts, they are getting attached
	// to @reserve directive.
	return NIL
}

func (e *Evaluator) evalDumpStmt(node *ast.DumpStmt, scope *object.Scope, path string) object.Object {
	values := make([]string, 0, len(node.Arguments))

	for i := range node.Arguments {
		evaluated := e.Eval(node.Arguments[i], scope, path)
		values = append(values, evaluated.Dump(0))
	}

	return &object.Dump{Values: values}
}

func (e *Evaluator) evalIdentifier(
	node *ast.Identifier,
	scope *object.Scope,
	path string,
) object.Object {
	varName := node.Name
	if varName == "global" && e.Config != nil && e.Config.GlobalData != nil {
		return object.NativeToObject(e.Config.GlobalData)
	}

	if val, ok := scope.Get(varName); ok {
		return val
	}

	return e.newError(node, path, fail.ErrIdentifierIsUndefined, node.Name)
}

func (e *Evaluator) evalIndexExp(
	node *ast.IndexExp,
	scope *object.Scope,
	path string,
) object.Object {
	left := e.Eval(node.Left, scope, path)
	if isError(left) {
		return left
	}

	idx := e.Eval(node.Index, scope, path)
	if isError(idx) {
		return idx
	}

	switch {
	case left.Is(object.ARR_OBJ) && idx.Is(object.INT_OBJ):
		return e.evalArrayIndexExp(left, idx)
	case left.Is(object.OBJ_OBJ) && idx.Is(object.STR_OBJ):
		return e.evalObjectKeyExp(left.(*object.Obj), idx.(*object.Str).Value, node.Index, path)
	}

	return e.newError(node, path, fail.ErrIndexNotSupported, left.Type())
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
	path string,
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

	return e.newError(node, path, fail.ErrPropertyNotFound, key, object.OBJ_OBJ)
}

func (e *Evaluator) evalDotExp(node *ast.DotExp, scope *object.Scope, path string) object.Object {
	left := e.Eval(node.Left, scope, path)
	if isError(left) {
		return left
	}

	key := node.Key.(*ast.Identifier)
	obj, ok := left.(*object.Obj)
	if !ok {
		return e.newError(node, path, fail.ErrPropertyOnNonObject, key, left.Type())
	}

	return e.evalObjectKeyExp(obj, key.Name, node, path)
}

func (e *Evaluator) evalString(node *ast.StringLiteral, _ *object.Scope) object.Object {
	str := html.EscapeString(node.Value)

	// unescape single and double quotes
	str = strings.ReplaceAll(str, "&#34;", `"`)
	str = strings.ReplaceAll(str, "&#39;", `'`)

	return &object.Str{Value: str}
}

func (e *Evaluator) evalPrefixExp(node *ast.PrefixExp, scope *object.Scope, path string) object.Object {
	right := e.Eval(node.Right, scope, path)
	if isError(right) {
		return right
	}

	switch node.Operator {
	case "-":
		return e.evalMinusPrefixOperatorExp(right, node, path)
	case "!":
		return e.evalBangOperatorExp(right, node, path)
	}

	return e.newError(node, path, fail.ErrUnknownOperator, node.Operator, right.Type())
}

func (e *Evaluator) evalTernaryExp(
	node *ast.TernaryExp,
	scope *object.Scope,
	path string,
) object.Object {
	cond := e.Eval(node.Condition, scope, path)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return e.Eval(node.IfBlock, scope, path)
	}

	return e.Eval(node.ElseBlock, scope, path)
}

func (e *Evaluator) evalArrayLiteral(
	node *ast.ArrayLiteral,
	scope *object.Scope,
	path string,
) object.Object {
	elems := e.evalExpressions(node.Elements, scope, path)
	if len(elems) == 1 && isError(elems[0]) {
		return elems[0]
	}

	return &object.Array{Elements: elems}
}

func (e *Evaluator) evalObjectLiteral(
	node *ast.ObjectLiteral,
	scope *object.Scope,
	path string,
) object.Object {
	pairs := make(map[string]object.Object, len(node.Pairs))

	for key, val := range node.Pairs {
		valObj := e.Eval(val, scope, path)
		if isError(valObj) {
			return valObj
		}

		pairs[key] = valObj
	}

	return object.NewObj(pairs)
}

func (e *Evaluator) evalExpressions(
	exps []ast.Expression,
	scope *object.Scope,
	path string,
) []object.Object {
	res := make([]object.Object, 0, len(exps))

	for i := range exps {
		evaluated := e.Eval(exps[i], scope, path)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		res = append(res, evaluated)
	}

	return res
}

func (e *Evaluator) evalInfixExp(
	operator string,
	leftNode,
	rightNode ast.Expression,
	scope *object.Scope,
	path string,
) object.Object {
	left := e.Eval(leftNode, scope, path)
	if isError(left) {
		return left
	}

	right := e.Eval(rightNode, scope, path)
	if isError(right) {
		return right
	}

	return e.evalInfixOperatorExp(operator, left, right, leftNode, path)
}

func (e *Evaluator) evalPostfixExp(
	node *ast.PostfixExp,
	scope *object.Scope,
	path string,
) object.Object {
	leftObj := e.Eval(node.Left, scope, path)
	if isError(leftObj) {
		return leftObj
	}

	return e.evalPostfixOperatorExp(leftObj, node.Operator, node, path)
}

func (e *Evaluator) evalCallExp(
	node *ast.CallExp,
	scope *object.Scope,
	path string,
) object.Object {
	receiver := e.Eval(node.Receiver, scope, path)
	funcName := node.Function.Name
	if isError(receiver) {
		return receiver
	}

	receiverType := receiver.Type()
	typeFuncs, ok := functions[receiverType]
	if !ok {
		return e.newError(node, path, fail.ErrNoFuncForThisType, funcName, receiverType)
	}

	args := e.evalExpressions(node.Arguments, scope, path)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	buitin, ok := typeFuncs[node.Function.Name]
	if ok {
		res, err := buitin.Fn(receiver, args...)
		if err != nil {
			return e.newError(node, path, "%s", err.Error())
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

	return e.newError(node, path, fail.ErrNoFuncForThisType, node.Function.Name, receiver.Type())
}

func (e *Evaluator) evalGlobalCallExp(
	node *ast.GlobalCallExp,
	scope *object.Scope,
	path string,
) object.Object {
	switch node.Function.Name {
	case "defined":
		return e.evalGlobalFuncDefined(node, scope, path)
	default:
		return e.newError(node, path, fail.ErrGlobalFuncMissing, node.Function.Name)
	}
}

func (e *Evaluator) evalGlobalFuncDefined(
	node *ast.GlobalCallExp,
	scope *object.Scope,
	path string,
) object.Object {
	definedVars := make([]bool, 0, len(node.Arguments))
	for i := range node.Arguments {
		evaluated := e.Eval(node.Arguments[i], scope, path)
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

func (e *Evaluator) evalPostfixOperatorExp(
	left object.Object,
	op string,
	node ast.Node,
	path string,
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
				return e.newError(node, path, fail.ErrCannotSubFromFloat, float, err)
			}

			return float
		}
	}

	return e.newError(node, path, fail.ErrUnknownOperator, left.Type(), op)
}

func (e *Evaluator) evalInfixOperatorExp(
	op string,
	left,
	right object.Object,
	leftNode ast.Node,
	path string,
) object.Object {
	if left.Type() != right.Type() {
		return e.newError(leftNode, path, fail.ErrTypeMismatch, left.Type(), op, right.Type())
	}

	switch left.Type() {
	case object.INT_OBJ:
		return e.evalIntegerInfixExp(op, right, left, leftNode, path)
	case object.FLOAT_OBJ:
		return e.evalFloatInfixExp(op, right, left, leftNode, path)
	case object.BOOL_OBJ:
		return e.evalBooleanInfixExp(op, right, left, leftNode, path)
	case object.STR_OBJ:
		return e.evalStringInfixExp(op, right, left, leftNode, path)
	}

	return e.newError(leftNode, path, fail.ErrUnknownTypeForOperator, left.Type(), op)
}

func (e *Evaluator) evalBooleanInfixExp(
	operator string,
	right,
	left object.Object,
	leftNode ast.Node,
	path string,
) object.Object {
	leftVal := left.(*object.Bool).Value
	rightVal := right.(*object.Bool).Value

	switch operator {
	case "&&":
		return &object.Bool{Value: leftVal && rightVal}
	case "||":
		return &object.Bool{Value: leftVal || rightVal}
	}

	return e.newError(leftNode, path, fail.ErrUnknownTypeForOperator, left.Type(), operator)
}

func (e *Evaluator) evalIntegerInfixExp(
	operator string,
	right,
	left object.Object,
	leftNode ast.Node,
	path string,
) object.Object {
	leftVal := left.(*object.Int).Value
	rightVal := right.(*object.Int).Value

	switch operator {
	case "+":
		return &object.Int{Value: leftVal + rightVal}
	case "-":
		return &object.Int{Value: leftVal - rightVal}
	case "*":
		return &object.Int{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return e.newError(leftNode, path, fail.ErrDivisionByZero)
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

	return e.newError(leftNode, path, fail.ErrUnknownTypeForOperator, left.Type(), operator)
}

func (e *Evaluator) evalStringInfixExp(
	operator string,
	right,
	left object.Object,
	leftNode ast.Node,
	path string,
) object.Object {
	leftVal := left.(*object.Str).Value
	rightVal := right.(*object.Str).Value

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "+":
		return &object.Str{Value: leftVal + rightVal}
	}

	return e.newError(leftNode, path, fail.ErrUnknownTypeForOperator, left.Type(), operator)
}

func (e *Evaluator) evalFloatInfixExp(
	operator string,
	right,
	left object.Object,
	leftNode ast.Node,
	path string,
) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
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

	return e.newError(leftNode, path, fail.ErrUnknownTypeForOperator, left.Type(), operator)
}

func (e *Evaluator) evalMinusPrefixOperatorExp(
	right object.Object,
	node ast.Node,
	path string,
) object.Object {
	switch right.Type() {
	case object.INT_OBJ:
		val := right.(*object.Int).Value
		return &object.Int{Value: -val}
	case object.FLOAT_OBJ:
		val := right.(*object.Float).Value
		return &object.Float{Value: -val}
	}

	return e.newError(node, path, fail.ErrPrefixOperatorIsWrong, "-", right.Type())
}

func (e *Evaluator) evalBangOperatorExp(
	right object.Object,
	node ast.Node,
	path string,
) object.Object {
	switch right {
	case FALSE:
		return TRUE
	case TRUE:
		return FALSE
	case NIL:
		return TRUE
	}

	return e.newError(node, path, fail.ErrPrefixOperatorIsWrong, "!", right.Type())
}

func (e *Evaluator) newError(node ast.Node, path, format string, a ...any) *object.Error {
	return &object.Error{
		Err:     fail.New(node.Line(), path, "evaluator", format, a...),
		ErrorID: format,
	}
}
