package evaluator

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/textwire/textwire/v2/ast"
	"github.com/textwire/textwire/v2/ctx"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
)

var (
	NIL      = &object.Nil{}
	TRUE     = &object.Bool{Value: true}
	FALSE    = &object.Bool{Value: false}
	BREAK    = &object.Break{}
	CONTINUE = &object.Continue{}
)

type Evaluator struct {
	ctx *ctx.EvalCtx
}

func New(ctx *ctx.EvalCtx) *Evaluator {
	return &Evaluator{ctx: ctx}
}

func (e *Evaluator) Eval(node ast.Node, env *object.Env) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return e.evalProgram(node, env)
	case *ast.HTMLStmt:
		return &object.HTML{Value: node.String()}
	case *ast.ExpressionStmt:
		return e.Eval(node.Expression, env)
	case *ast.IfStmt:
		return e.evalIfStmt(node, env)
	case *ast.BlockStmt:
		return e.evalBlockStmt(node, env)
	case *ast.AssignStmt:
		return e.evalAssignStmt(node, env)
	case *ast.UseStmt:
		return e.evalUseStmt(node, env)
	case *ast.ReserveStmt:
		return e.evalReserveStmt(node, env)
	case *ast.ForStmt:
		return e.evalForStmt(node, env)
	case *ast.EachStmt:
		return e.evalEachStmt(node, env)
	case *ast.BreakIfStmt:
		return e.evalBreakIfStmt(node, env)
	case *ast.ComponentStmt:
		return e.evalComponentStmt(node, env)
	case *ast.ContinueIfStmt:
		return e.evalContinueIfStmt(node, env)
	case *ast.SlotStmt:
		return e.evalSlotStmt(node, env)
	case *ast.DumpStmt:
		return e.evalDumpStmt(node, env)
	case *ast.BreakStmt:
		return BREAK
	case *ast.ContinueStmt:
		return CONTINUE
	case *ast.InsertStmt:
		return NIL

	// Expressions
	case *ast.Identifier:
		return e.evalIdentifier(node, env)
	case *ast.IndexExp:
		return e.evalIndexExp(node, env)
	case *ast.DotExp:
		return e.evalDotExp(node, env)
	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return e.evalString(node, env)
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ObjectLiteral:
		return e.evalObjectLiteral(node, env)
	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(node, env)
	case *ast.PrefixExp:
		return e.evalPrefixExp(node, env)
	case *ast.TernaryExp:
		return e.evalTernaryExp(node, env)
	case *ast.InfixExp:
		return e.evalInfixExp(node.Operator, node.Left, node.Right, env)
	case *ast.PostfixExp:
		return e.evalPostfixExp(node, env)
	case *ast.CallExp:
		return e.evalCallExp(node, env)
	case *ast.NilLiteral:
		return NIL
	}

	return e.newError(node, fail.ErrUnknownNodeType, node)
}

func (e *Evaluator) evalProgram(prog *ast.Program, env *object.Env) object.Object {
	var out bytes.Buffer

	for _, statement := range prog.Statements {
		stmtObj := e.Eval(statement, env)

		if isError(stmtObj) {
			return stmtObj
		}

		out.WriteString(stmtObj.String())
	}

	return &object.HTML{Value: out.String()}
}

func (e *Evaluator) evalIfStmt(node *ast.IfStmt, env *object.Env) object.Object {
	condition := e.Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	newEnv := object.NewEnclosedEnv(env)

	if isTruthy(condition) {
		return e.Eval(node.Consequence, newEnv)
	}

	for _, alt := range node.Alternatives {
		condition = e.Eval(alt.Condition, env)

		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return e.Eval(alt.Consequence, newEnv)
		}
	}

	if node.Alternative != nil {
		return e.Eval(node.Alternative, newEnv)
	}

	return NIL
}

func (e *Evaluator) evalBlockStmt(block *ast.BlockStmt, env *object.Env) object.Object {
	var elems []object.Object

	for _, stmt := range block.Statements {
		obj := e.Eval(stmt, env)

		if isError(obj) {
			return obj
		}

		elems = append(elems, obj)

		if hasBreakStmt(obj) || hasContinueStmt(obj) {
			break
		}
	}

	return &object.Block{Elements: elems}
}

func (e *Evaluator) evalAssignStmt(node *ast.AssignStmt, env *object.Env) object.Object {
	val := e.Eval(node.Value, env)

	if isError(val) {
		return val
	}

	err := env.Set(node.Name.Value, val)

	if err != nil {
		return e.newError(node, "%s", err.Error())
	}

	return NIL
}

func (e *Evaluator) evalUseStmt(node *ast.UseStmt, env *object.Env) object.Object {
	if node.Program == nil {
		return e.newError(node, fail.ErrUseStmtMustHaveProgram)
	}

	if node.Program.IsLayout && node.Program.HasUseStmt() {
		return e.newError(node, fail.ErrUseStmtNotAllowed)
	}

	layoutContent := e.Eval(node.Program, env)

	if isError(layoutContent) {
		return layoutContent
	}

	return &object.Use{
		Path:    node.Name.Value,
		Content: layoutContent,
	}
}

func (e *Evaluator) evalReserveStmt(node *ast.ReserveStmt, env *object.Env) object.Object {
	stmt := &object.Reserve{Name: node.Name.Value}

	if node.Insert.Block != nil {
		result := e.Eval(node.Insert.Block, env)

		if isError(result) {
			return result
		}

		stmt.Content = result

		return stmt
	}

	if node.Insert.Argument == nil {
		return e.newError(node.Insert, fail.ErrInsertMustHaveContent)
	}

	firstArg := e.Eval(node.Insert.Argument, env)

	if isError(firstArg) {
		return firstArg
	}

	stmt.Argument = firstArg

	return stmt
}

func (e *Evaluator) evalComponentStmt(node *ast.ComponentStmt, env *object.Env) object.Object {
	name := e.Eval(node.Name, env)

	if isError(name) {
		return name
	}

	if node.Block == nil {
		return e.newError(node, fail.ErrComponentMustHaveBlock, name.String())
	}

	stmt := &object.Component{Name: name.String()}

	newEnv := object.NewEnclosedEnv(env)

	if node.Argument != nil {
		for key, arg := range node.Argument.Pairs {
			val := e.Eval(arg, env)

			if isError(val) {
				return val
			}

			newEnv.Set(key, val)
		}
	}

	content := e.Eval(node.Block, newEnv)

	if isError(content) {
		return content
	}

	stmt.Content = content

	return stmt
}

func (e *Evaluator) evalForStmt(node *ast.ForStmt, env *object.Env) object.Object {
	newEnv := object.NewEnclosedEnv(env)

	var init object.Object
	var blocks bytes.Buffer

	if node.Init != nil {
		if init = e.Eval(node.Init, newEnv); isError(init) {
			return init
		}
	}

	// evaluate alternative block if condition is false
	if node.Condition != nil {
		cond := e.Eval(node.Condition, newEnv)

		if isError(cond) {
			return cond
		}

		if !isTruthy(cond) && node.Alternative != nil {
			return e.Eval(node.Alternative, newEnv)
		}
	}

	for {
		cond := e.Eval(node.Condition, newEnv)

		if isError(cond) {
			return cond
		}

		if !isTruthy(cond) {
			break
		}

		block := e.Eval(node.Block, newEnv)

		if isError(block) {
			return block
		}

		blocks.WriteString(block.String())

		post := e.Eval(node.Post, newEnv)

		if isError(post) {
			return post
		}

		if node.Init == nil || node.Post == nil {
			continue
		}

		varName := node.Init.(*ast.AssignStmt).Name.Value

		err := newEnv.Set(varName, post)

		if err != nil {
			return e.newError(node, "%s", err.Error())
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

func (e *Evaluator) evalEachStmt(node *ast.EachStmt, env *object.Env) object.Object {
	newEnv := object.NewEnclosedEnv(env)

	var blocks bytes.Buffer

	varName := node.Var.Value
	arrObj := e.Eval(node.Array, newEnv)

	if isError(arrObj) {
		return arrObj
	}

	elems := arrObj.(*object.Array).Elements
	elemsLen := len(elems)

	// evaluate alternative block if array is empty
	if elemsLen == 0 && node.Alternative != nil {
		return e.Eval(node.Alternative, newEnv)
	}

	for i, elem := range elems {
		err := newEnv.Set(varName, elem)

		if err != nil {
			return e.newError(node, "%s", err.Error())
		}

		newEnv.SetLoopVar(map[string]object.Object{
			"index": &object.Int{Value: int64(i)},
			"first": nativeBoolToBooleanObject(i == 0),
			"last":  nativeBoolToBooleanObject(i == elemsLen-1),
			"iter":  &object.Int{Value: int64(i + 1)},
		})

		block := e.Eval(node.Block, newEnv)

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
	env *object.Env,
) object.Object {
	condition := e.Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return BREAK
	}

	return NIL
}

func (e *Evaluator) evalContinueIfStmt(
	node *ast.ContinueIfStmt,
	env *object.Env,
) object.Object {
	condition := e.Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return CONTINUE
	}

	return NIL
}

func (e *Evaluator) evalSlotStmt(
	node *ast.SlotStmt,
	env *object.Env,
) object.Object {
	var body object.Object

	if node.Body != nil {
		body = e.Eval(node.Body, env)

		if isError(body) {
			return body
		}
	} else {
		body = NIL
	}

	return &object.Slot{Name: node.Name.Value, Content: body}
}

func (e *Evaluator) evalDumpStmt(node *ast.DumpStmt, env *object.Env) object.Object {
	var values []string

	for _, arg := range node.Arguments {
		val := e.Eval(arg, env)
		values = append(values, fmt.Sprintf(object.DumpHTML, val.String()))
	}

	return &object.Dump{Values: values}
}

func (e *Evaluator) evalIdentifier(
	node *ast.Identifier,
	env *object.Env,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	return e.newError(node, fail.ErrIdentifierNotFound, node.Value)
}

func (e *Evaluator) evalIndexExp(
	node *ast.IndexExp,
	env *object.Env,
) object.Object {
	left := e.Eval(node.Left, env)

	if isError(left) {
		return left
	}

	idx := e.Eval(node.Index, env)

	if isError(idx) {
		return idx
	}

	switch {
	case left.Is(object.ARR_OBJ) && idx.Is(object.INT_OBJ):
		return e.evalArrayIndexExp(left, idx)
	case left.Is(object.OBJ_OBJ) && idx.Is(object.STR_OBJ):
		return e.evalObjectIndexExp(left.(*object.Obj), idx.(*object.Str).Value, node.Index)
	}

	return e.newError(node, fail.ErrIndexNotSupported, left.Type())
}

func (e *Evaluator) evalArrayIndexExp(
	arr,
	idx object.Object,
) object.Object {
	arrObj := arr.(*object.Array)
	index := idx.(*object.Int).Value
	max := int64(len(arrObj.Elements) - 1)

	if index < 0 || index > max {
		return NIL
	}

	return arrObj.Elements[index]
}

func (e *Evaluator) evalObjectIndexExp(
	obj object.Object,
	idx string,
	node ast.Node,
) object.Object {
	objObj := obj.(*object.Obj)
	pair, ok := objObj.Pairs[idx]

	if ok {
		return pair
	}

	// make first letter lowercase on idx
	idxUpper := strings.ToUpper(idx[:1]) + idx[1:]

	if pair, ok = objObj.Pairs[idxUpper]; !ok {
		return e.newError(node, fail.ErrPropertyNotFound, idx, object.OBJ_OBJ)
	}

	return pair
}

func (e *Evaluator) evalDotExp(node *ast.DotExp, env *object.Env) object.Object {
	left := e.Eval(node.Left, env)

	if isError(left) {
		return left
	}

	key := node.Key.(*ast.Identifier)

	return e.evalObjectIndexExp(left.(*object.Obj), key.Value, node)
}

func (e *Evaluator) evalString(node *ast.StringLiteral, _ *object.Env) object.Object {
	str := html.EscapeString(node.Value)

	// unescape single and double quotes
	str = strings.ReplaceAll(str, "&#34;", `"`)
	str = strings.ReplaceAll(str, "&#39;", `'`)

	return &object.Str{Value: str}
}

func (e *Evaluator) evalPrefixExp(node *ast.PrefixExp, env *object.Env) object.Object {
	right := e.Eval(node.Right, env)

	if isError(right) {
		return right
	}

	switch node.Operator {
	case "-":
		return e.evalMinusPrefixOperatorExp(right, node)
	case "!":
		return e.evalBangOperatorExp(right, node)
	}

	return e.newError(node, fail.ErrUnknownOperator,
		node.Operator, right.Type())
}

func (e *Evaluator) evalTernaryExp(
	node *ast.TernaryExp,
	env *object.Env,
) object.Object {
	condition := e.Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.Eval(node.Consequence, env)
	}

	return e.Eval(node.Alternative, env)
}

func (e *Evaluator) evalArrayLiteral(
	node *ast.ArrayLiteral,
	env *object.Env,
) object.Object {
	elems := e.evalExpressions(node.Elements, env)

	if len(elems) == 1 && isError(elems[0]) {
		return elems[0]
	}

	return &object.Array{Elements: elems}
}

func (e *Evaluator) evalObjectLiteral(node *ast.ObjectLiteral, env *object.Env) object.Object {
	pairs := make(map[string]object.Object)

	for key, value := range node.Pairs {
		valueObj := e.Eval(value, env)

		if isError(valueObj) {
			return valueObj
		}

		pairs[key] = valueObj
	}

	return &object.Obj{Pairs: pairs}
}

func (e *Evaluator) evalExpressions(
	exps []ast.Expression,
	env *object.Env,
) []object.Object {
	var result []object.Object

	for _, expr := range exps {
		evaluated := e.Eval(expr, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) evalInfixExp(
	operator string,
	left,
	right ast.Expression,
	env *object.Env,
) object.Object {
	leftObj := e.Eval(left, env)

	if isError(leftObj) {
		return leftObj
	}

	rightObj := e.Eval(right, env)

	if isError(rightObj) {
		return rightObj
	}

	return e.evalInfixOperatorExp(operator, leftObj, rightObj, left)
}

func (e *Evaluator) evalPostfixExp(
	node *ast.PostfixExp,
	env *object.Env,
) object.Object {
	leftObj := e.Eval(node.Left, env)

	if isError(leftObj) {
		return leftObj
	}

	return e.evalPostfixOperatorExp(leftObj, node.Operator, node)
}

func (e *Evaluator) evalCallExp(
	node *ast.CallExp,
	env *object.Env,
) object.Object {
	receiverObj := e.Eval(node.Receiver, env)

	if isError(receiverObj) {
		return receiverObj
	}

	receiverType := receiverObj.Type()
	funcName := node.Function.Value

	typeFuncs, ok := functions[receiverType]

	if !ok {
		return e.newError(node, fail.ErrNoFuncForThisType, funcName, receiverType)
	}

	args := e.evalExpressions(node.Arguments, env)

	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	buitin, ok := typeFuncs[node.Function.Value]

	if ok {
		res, err := buitin.Fn(e.ctx, receiverObj, args...)

		if err != nil {
			return e.newError(node, "%s", err.Error())
		}

		return res
	}

	if hasCustomFunc(e.ctx.CustomFunc, receiverType, funcName) {
		nativeArgs := e.objectsToNativeType(args)

		switch receiverType {
		case object.STR_OBJ:
			fun := e.ctx.CustomFunc.Str[funcName]
			res := fun(receiverObj.String(), nativeArgs...)
			return object.NativeToObject(res)
		case object.ARR_OBJ:
			fun := e.ctx.CustomFunc.Arr[funcName]
			nativeElems := e.objectsToNativeType(receiverObj.(*object.Array).Elements)
			res := fun(nativeElems, nativeArgs...)
			return object.NativeToObject(res)
		case object.INT_OBJ:
			fun := e.ctx.CustomFunc.Int[funcName]
			res := fun(int(receiverObj.(*object.Int).Value), nativeArgs...)
			return object.NativeToObject(res)
		case object.FLOAT_OBJ:
			fun := e.ctx.CustomFunc.Float[funcName]
			res := fun(receiverObj.(*object.Float).Value, nativeArgs...)
			return object.NativeToObject(res)
		case object.BOOL_OBJ:
			fun := e.ctx.CustomFunc.Bool[funcName]
			res := fun(receiverObj.(*object.Bool).Value, nativeArgs...)
			return object.NativeToObject(res)
		}
	}

	return e.newError(node, fail.ErrNoFuncForThisType, node.Function.Value, receiverObj.Type())
}

func (e *Evaluator) objectsToNativeType(args []object.Object) []interface{} {
	var result []interface{}

	for _, arg := range args {
		result = append(result, arg.Val())
	}

	return result
}

func (e *Evaluator) evalPostfixOperatorExp(
	left object.Object,
	operator string,
	node ast.Node,
) object.Object {
	if operator == "++" {
		if left.Is(object.INT_OBJ) {
			value := left.(*object.Int).Value + 1
			return &object.Int{Value: value}
		}

		if left.Is(object.FLOAT_OBJ) {
			value := left.(*object.Float).Value + 1
			return &object.Float{Value: value}
		}
	}

	if operator == "--" {
		if left.Is(object.INT_OBJ) {
			value := left.(*object.Int).Value - 1
			return &object.Int{Value: value}
		}

		if left.Is(object.FLOAT_OBJ) {
			value := left.(*object.Float).Value
			float := &object.Float{Value: value}
			float.SubtractFromFloat(1)
			return float
		}
	}

	return e.newError(node, fail.ErrUnknownOperator,
		left.Type(), operator)
}

func (e *Evaluator) evalInfixOperatorExp(
	operator string,
	left,
	right object.Object,
	leftNode ast.Node,
) object.Object {
	if left.Type() != right.Type() {
		return e.newError(leftNode, fail.ErrTypeMismatch,
			left.Type(), operator, right.Type())
	}

	switch left.Type() {
	case object.INT_OBJ:
		return e.evalIntegerInfixExp(operator, right, left, leftNode)
	case object.FLOAT_OBJ:
		return e.evalFloatInfixExp(operator, right, left, leftNode)
	case object.STR_OBJ:
		return e.evalStringInfixExp(operator, right, left, leftNode)
	}

	return e.newError(leftNode, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
}

func (e *Evaluator) evalIntegerInfixExp(
	operator string,
	right,
	left object.Object,
	leftNode ast.Node,
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
			return e.newError(leftNode, fail.ErrDivisionByZero)
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

	return e.newError(leftNode, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
}

func (e *Evaluator) evalStringInfixExp(
	operator string,
	right,
	left object.Object,
	leftNode ast.Node,
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

	return e.newError(leftNode, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
}

func (e *Evaluator) evalFloatInfixExp(
	operator string,
	right,
	left object.Object,
	leftNode ast.Node,
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

	return e.newError(leftNode, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
}

func (e *Evaluator) evalMinusPrefixOperatorExp(
	right object.Object,
	node ast.Node,
) object.Object {
	switch right.Type() {
	case object.INT_OBJ:
		value := right.(*object.Int).Value
		return &object.Int{Value: -value}
	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}

	return e.newError(node, fail.ErrPrefixOperatorIsWrong,
		"-", right.Type())
}

func (e *Evaluator) evalBangOperatorExp(
	right object.Object,
	node ast.Node,
) object.Object {
	switch right {
	case FALSE:
		return TRUE
	case TRUE:
		return FALSE
	case NIL:
		return TRUE
	}

	return e.newError(node, fail.ErrPrefixOperatorIsWrong,
		"!", right.Type())
}

func (e *Evaluator) newError(
	node ast.Node,
	format string,
	a ...interface{},
) *object.Error {
	err := fail.New(node.Line(), e.ctx.AbsPath, "evaluator", format, a...)
	return &object.Error{Err: err}
}
