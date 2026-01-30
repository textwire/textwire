package evaluator

import (
	"bytes"
	"html"
	"strings"

	"github.com/textwire/textwire/v2/ast"
	"github.com/textwire/textwire/v2/config"
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

func (e *Evaluator) Eval(node ast.Node, env *object.Env, path string) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return e.evalProgram(node, env)
	case *ast.HTMLStmt:
		return &object.HTML{Value: node.String()}
	case *ast.ExpressionStmt:
		return e.Eval(node.Expression, env, path)
	case *ast.IfStmt:
		return e.evalIfStmt(node, env, path)
	case *ast.BlockStmt:
		return e.evalBlockStmt(node, env, path)
	case *ast.AssignStmt:
		return e.evalAssignStmt(node, env, path)
	case *ast.UseStmt:
		return e.evalUseStmt(node, env, path)
	case *ast.ReserveStmt:
		return e.evalReserveStmt(node, env, path)
	case *ast.ForStmt:
		return e.evalForStmt(node, env, path)
	case *ast.EachStmt:
		return e.evalEachStmt(node, env, path)
	case *ast.BreakIfStmt:
		return e.evalBreakIfStmt(node, env, path)
	case *ast.ComponentStmt:
		return e.evalComponentStmt(node, env, path)
	case *ast.ContinueIfStmt:
		return e.evalContinueIfStmt(node, env, path)
	case *ast.SlotStmt:
		return e.evalSlotStmt(node, env, path)
	case *ast.DumpStmt:
		return e.evalDumpStmt(node, env, path)
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
		return e.evalIdentifier(node, env, path)
	case *ast.IndexExp:
		return e.evalIndexExp(node, env, path)
	case *ast.DotExp:
		return e.evalDotExp(node, env, path)
	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return e.evalString(node, env)
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ObjectLiteral:
		return e.evalObjectLiteral(node, env, path)
	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(node, env, path)
	case *ast.PrefixExp:
		return e.evalPrefixExp(node, env, path)
	case *ast.TernaryExp:
		return e.evalTernaryExp(node, env, path)
	case *ast.InfixExp:
		return e.evalInfixExp(node.Operator, node.Left, node.Right, env, path)
	case *ast.PostfixExp:
		return e.evalPostfixExp(node, env, path)
	case *ast.CallExp:
		return e.evalCallExp(node, env, path)
	case *ast.GlobalCallExp:
		return e.evalGlobalCallExp(node, env, path)
	case *ast.NilLiteral:
		return NIL
	}

	return e.newError(node, path, fail.ErrUnknownNodeType, node)
}

func (e *Evaluator) evalProgram(prog *ast.Program, env *object.Env) object.Object {
	var out bytes.Buffer

	for _, statement := range prog.Statements {
		stmtObj := e.Eval(statement, env, prog.Filepath)
		if isError(stmtObj) {
			return stmtObj
		}

		out.WriteString(stmtObj.String())
	}

	return &object.HTML{Value: out.String()}
}

func (e *Evaluator) evalIfStmt(node *ast.IfStmt, env *object.Env, path string) object.Object {
	condition := e.Eval(node.Condition, env, path)
	if isError(condition) {
		return condition
	}

	newEnv := object.NewEnclosedEnv(env)
	if isTruthy(condition) {
		return e.Eval(node.Consequence, newEnv, path)
	}

	for _, alt := range node.Alternatives {
		if ifStmt, ok := alt.(*ast.ElseIfStmt); ok {
			condition = e.Eval(ifStmt.Condition, env, path)
			if isError(condition) {
				return condition
			}

			if isTruthy(condition) {
				return e.Eval(ifStmt.Consequence, newEnv, path)
			}
		}
	}

	if node.Alternative != nil {
		return e.Eval(node.Alternative, newEnv, path)
	}

	return NIL
}

func (e *Evaluator) evalBlockStmt(
	block *ast.BlockStmt,
	env *object.Env,
	path string,
) object.Object {
	var elems []object.Object

	for _, stmt := range block.Statements {
		obj := e.Eval(stmt, env, path)
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

func (e *Evaluator) evalAssignStmt(
	node *ast.AssignStmt,
	env *object.Env,
	path string,
) object.Object {
	val := e.Eval(node.Right, env, path)
	if isError(val) {
		return val
	}

	err := env.Set(node.Left.Name, val)
	if err != nil {
		return e.newError(node, path, "%s", err.Error())
	}

	return NIL
}

func (e *Evaluator) evalUseStmt(node *ast.UseStmt, env *object.Env, path string) object.Object {
	if node.Layout == nil {
		if e.UsingTemplates {
			return e.newError(node, path, fail.ErrUseStmtMissingLayout)
		}
		return e.newError(node, path, fail.ErrSomeDirsOnlyInTemplates)
	}

	if node.Layout.IsLayout && node.Layout.HasUseStmt() {
		return e.newError(node, path, fail.ErrUseStmtNotAllowed)
	}

	layoutContent := e.Eval(node.Layout, env, node.Layout.Filepath)
	if isError(layoutContent) {
		return layoutContent
	}

	return &object.Use{
		Path:    node.Name.Value,
		Content: layoutContent,
	}
}

func (e *Evaluator) evalReserveStmt(
	node *ast.ReserveStmt,
	env *object.Env,
	path string,
) object.Object {
	if !e.UsingTemplates {
		return e.newError(node, path, fail.ErrSomeDirsOnlyInTemplates)
	}

	stmt := &object.Reserve{Name: node.Name.Value}

	// Inserts are optional statements.
	// If not provided, reserve should be empty
	if node.Insert == nil {
		return NIL
	}

	if node.Insert.Block != nil {
		result := e.Eval(node.Insert.Block, env, node.Insert.FilePath)
		if isError(result) {
			return result
		}

		stmt.Content = result

		return stmt
	}

	if node.Insert.Argument == nil {
		return e.newError(node.Insert, path, fail.ErrInsertMustHaveContent)
	}

	firstArg := e.Eval(node.Insert.Argument, env, path)
	if isError(firstArg) {
		return firstArg
	}

	stmt.Argument = firstArg

	return stmt
}

func (e *Evaluator) evalComponentStmt(
	node *ast.ComponentStmt,
	env *object.Env,
	path string,
) object.Object {
	if !e.UsingTemplates {
		return e.newError(node, path, fail.ErrSomeDirsOnlyInTemplates)
	}

	name := e.Eval(node.Name, env, path)
	if isError(name) {
		return name
	}

	if node.Block == nil {
		return e.newError(node, path, fail.ErrComponentMustHaveBlock, name.String())
	}

	stmt := &object.Component{Name: name.String()}
	newEnv := object.NewEnclosedEnv(env)

	if node.Argument != nil {
		for key, arg := range node.Argument.Pairs {
			val := e.Eval(arg, env, path)
			if isError(val) {
				return val
			}

			err := newEnv.Set(key, val)
			if err != nil {
				return e.newError(node, path, "%s", err.Error())
			}
		}
	}

	content := e.Eval(node.Block, newEnv, node.Block.Filepath)
	if isError(content) {
		return content
	}

	stmt.Content = content

	return stmt
}

func (e *Evaluator) evalForStmt(node *ast.ForStmt, env *object.Env, path string) object.Object {
	newEnv := object.NewEnclosedEnv(env)

	var init object.Object
	var blocks bytes.Buffer

	if node.Init != nil {
		if init = e.Eval(node.Init, newEnv, path); isError(init) {
			return init
		}
	}

	// evaluate alternative block if user's condition is false
	if node.Condition != nil {
		cond := e.Eval(node.Condition, newEnv, path)
		if isError(cond) {
			return cond
		}

		if !isTruthy(cond) && node.Alternative != nil {
			return e.Eval(node.Alternative, newEnv, path)
		}
	}

	// loop through the block until the user's condition is false
	for {
		cond := e.Eval(node.Condition, newEnv, path)
		if isError(cond) {
			return cond
		}

		if !isTruthy(cond) {
			break
		}

		block := e.Eval(node.Block, newEnv, path)
		if isError(block) {
			return block
		}

		blocks.WriteString(block.String())

		post := e.Eval(node.Post, newEnv, path)
		if isError(post) {
			return post
		}

		if node.Init == nil || node.Post == nil {
			continue
		}

		varName := node.Init.(*ast.AssignStmt).Left.Name
		err := newEnv.Set(varName, post)
		if err != nil {
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
	env *object.Env,
	path string,
) object.Object {
	newEnv := object.NewEnclosedEnv(env)

	var blocks bytes.Buffer

	varName := node.Var.Name
	arrObj := e.Eval(node.Array, newEnv, path)
	if isError(arrObj) {
		return arrObj
	}

	elems := arrObj.(*object.Array).Elements
	elemsLen := len(elems)

	// evaluate alternative block if array is empty
	if elemsLen == 0 && node.Alternative != nil {
		return e.Eval(node.Alternative, newEnv, path)
	}

	for i, elem := range elems {
		err := newEnv.Set(varName, elem)
		if err != nil {
			return e.newError(node, path, "%s", err.Error())
		}

		newEnv.SetLoopVar(map[string]object.Object{
			"index": &object.Int{Value: int64(i)},
			"first": nativeBoolToBooleanObject(i == 0),
			"last":  nativeBoolToBooleanObject(i == elemsLen-1),
			"iter":  &object.Int{Value: int64(i + 1)},
		})

		block := e.Eval(node.Block, newEnv, path)
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
	path string,
) object.Object {
	condition := e.Eval(node.Condition, env, path)

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
	path string,
) object.Object {
	condition := e.Eval(node.Condition, env, path)
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
	path string,
) object.Object {
	var body object.Object

	if node.Body != nil {
		body = e.Eval(node.Body, env, path)
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

func (e *Evaluator) evalDumpStmt(node *ast.DumpStmt, env *object.Env, path string) object.Object {
	var values []string

	for _, arg := range node.Arguments {
		val := e.Eval(arg, env, path)
		values = append(values, val.Dump(0))
	}

	return &object.Dump{Values: values}
}

func (e *Evaluator) evalIdentifier(
	node *ast.Identifier,
	env *object.Env,
	path string,
) object.Object {
	varName := node.Name
	if varName == "global" && e.Config != nil && e.Config.GlobalData != nil {
		return object.NativeToObject(e.Config.GlobalData)
	}

	if val, ok := env.Get(varName); ok {
		return val
	}

	return e.newError(node, path, fail.ErrIdentifierIsUndefined, node.Name)
}

func (e *Evaluator) evalIndexExp(
	node *ast.IndexExp,
	env *object.Env,
	path string,
) object.Object {
	left := e.Eval(node.Left, env, path)
	if isError(left) {
		return left
	}

	idx := e.Eval(node.Index, env, path)
	if isError(idx) {
		return idx
	}

	switch {
	case left.Is(object.ARR_OBJ) && idx.Is(object.INT_OBJ):
		return e.evalArrayIndexExp(left, idx)
	case left.Is(object.OBJ_OBJ) && idx.Is(object.STR_OBJ):
		return e.evalObjectKeyExp(left.(*object.Obj),
			idx.(*object.Str).Value, node.Index, path)
	}

	return e.newError(node, path, fail.ErrIndexNotSupported, left.Type())
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

func (e *Evaluator) evalObjectKeyExp(
	obj *object.Obj,
	key string,
	node ast.Node,
	path string,
) object.Object {
	pair, ok := obj.Pairs[key]
	if ok {
		return pair
	}

	// make first letter lowercase on key
	keyUpper := strings.ToUpper(key[:1]) + key[1:]

	if pair, ok = obj.Pairs[keyUpper]; !ok {
		return e.newError(node, path, fail.ErrPropertyNotFound, key, object.OBJ_OBJ)
	}

	return pair
}

func (e *Evaluator) evalDotExp(node *ast.DotExp, env *object.Env, path string) object.Object {
	left := e.Eval(node.Left, env, path)
	if isError(left) {
		return left
	}

	key := node.Key.(*ast.Identifier)

	return e.evalObjectKeyExp(left.(*object.Obj), key.Name, node, path)
}

func (e *Evaluator) evalString(node *ast.StringLiteral, _ *object.Env) object.Object {
	str := html.EscapeString(node.Value)

	// unescape single and double quotes
	str = strings.ReplaceAll(str, "&#34;", `"`)
	str = strings.ReplaceAll(str, "&#39;", `'`)

	return &object.Str{Value: str}
}

func (e *Evaluator) evalPrefixExp(node *ast.PrefixExp, env *object.Env, path string) object.Object {
	right := e.Eval(node.Right, env, path)
	if isError(right) {
		return right
	}

	switch node.Operator {
	case "-":
		return e.evalMinusPrefixOperatorExp(right, node, path)
	case "!":
		return e.evalBangOperatorExp(right, node, path)
	}

	return e.newError(node, path, fail.ErrUnknownOperator,
		node.Operator, right.Type())
}

func (e *Evaluator) evalTernaryExp(
	node *ast.TernaryExp,
	env *object.Env,
	path string,
) object.Object {
	condition := e.Eval(node.Condition, env, path)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.Eval(node.Consequence, env, path)
	}

	return e.Eval(node.Alternative, env, path)
}

func (e *Evaluator) evalArrayLiteral(
	node *ast.ArrayLiteral,
	env *object.Env,
	path string,
) object.Object {
	elems := e.evalExpressions(node.Elements, env, path)
	if len(elems) == 1 && isError(elems[0]) {
		return elems[0]
	}

	return &object.Array{Elements: elems}
}

func (e *Evaluator) evalObjectLiteral(
	node *ast.ObjectLiteral,
	env *object.Env,
	path string,
) object.Object {
	pairs := map[string]object.Object{}

	for key, value := range node.Pairs {
		valueObj := e.Eval(value, env, path)
		if isError(valueObj) {
			return valueObj
		}

		pairs[key] = valueObj
	}

	return object.NewObj(pairs)
}

func (e *Evaluator) evalExpressions(
	exps []ast.Expression,
	env *object.Env,
	path string,
) []object.Object {
	var result []object.Object

	for _, expr := range exps {
		evaluated := e.Eval(expr, env, path)
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
	path string,
) object.Object {
	leftObj := e.Eval(left, env, path)
	if isError(leftObj) {
		return leftObj
	}

	rightObj := e.Eval(right, env, path)
	if isError(rightObj) {
		return rightObj
	}

	return e.evalInfixOperatorExp(operator, leftObj, rightObj, left, path)
}

func (e *Evaluator) evalPostfixExp(
	node *ast.PostfixExp,
	env *object.Env,
	path string,
) object.Object {
	leftObj := e.Eval(node.Left, env, path)
	if isError(leftObj) {
		return leftObj
	}

	return e.evalPostfixOperatorExp(leftObj, node.Operator, node, path)
}

func (e *Evaluator) evalCallExp(
	node *ast.CallExp,
	env *object.Env,
	path string,
) object.Object {
	receiver := e.Eval(node.Receiver, env, path)
	funcName := node.Function.Name
	if isError(receiver) {
		return receiver
	}

	receiverType := receiver.Type()
	typeFuncs, ok := functions[receiverType]
	if !ok {
		return e.newError(node, path, fail.ErrNoFuncForThisType,
			funcName, receiverType)
	}

	args := e.evalExpressions(node.Arguments, env, path)
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

	return e.newError(node, path, fail.ErrNoFuncForThisType,
		node.Function.Name, receiver.Type())
}

func (e *Evaluator) evalGlobalCallExp(
	node *ast.GlobalCallExp,
	env *object.Env,
	path string,
) object.Object {
	funcName := node.Function.Name
	switch funcName {
	case "defined":
		return e.evalGlobalFuncDefined(node, env, path)
	default:
		return e.newError(node, path, fail.ErrGlobalFuncMissing, funcName)
	}
}

func (e *Evaluator) evalGlobalFuncDefined(
	node *ast.GlobalCallExp,
	env *object.Env,
	path string,
) object.Object {
	var definedVars []bool
	for _, expr := range node.Arguments {
		evaluated := e.Eval(expr, env, path)
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
	var result []any
	for _, arg := range args {
		result = append(result, arg.Val())
	}

	return result
}

func (e *Evaluator) evalPostfixOperatorExp(
	left object.Object,
	operator string,
	node ast.Node,
	path string,
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

			err := float.SubtractFromFloat(1)
			if err != nil {
				return e.newError(node, path, fail.ErrCannotSubFromFloat,
					float.String())
			}

			return float
		}
	}

	return e.newError(node, path, fail.ErrUnknownOperator,
		left.Type(), operator)
}

func (e *Evaluator) evalInfixOperatorExp(
	operator string,
	left,
	right object.Object,
	leftNode ast.Node,
	path string,
) object.Object {
	if left.Type() != right.Type() {
		return e.newError(leftNode, path, fail.ErrTypeMismatch,
			left.Type(), operator, right.Type())
	}

	switch left.Type() {
	case object.INT_OBJ:
		return e.evalIntegerInfixExp(operator, right, left, leftNode, path)
	case object.FLOAT_OBJ:
		return e.evalFloatInfixExp(operator, right, left, leftNode, path)
	case object.BOOL_OBJ:
		return e.evalBooleanInfixExp(operator, right, left, leftNode, path)
	case object.STR_OBJ:
		return e.evalStringInfixExp(operator, right, left, leftNode, path)
	}

	return e.newError(leftNode, path, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
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

	return e.newError(leftNode, path, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
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

	return e.newError(leftNode, path, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
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

	return e.newError(leftNode, path, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
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

	return e.newError(leftNode, path, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
}

func (e *Evaluator) evalMinusPrefixOperatorExp(
	right object.Object,
	node ast.Node,
	path string,
) object.Object {
	switch right.Type() {
	case object.INT_OBJ:
		value := right.(*object.Int).Value
		return &object.Int{Value: -value}
	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}

	return e.newError(node, path, fail.ErrPrefixOperatorIsWrong,
		"-", right.Type())
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
