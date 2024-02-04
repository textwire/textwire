package evaluator

import (
	"bytes"
	"html"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/fail"
	"github.com/textwire/textwire/object"
)

var (
	NIL   = &object.Nil{}
	TRUE  = &object.Bool{Value: true}
	FALSE = &object.Bool{Value: false}
)

type Evaluator struct {
	ctx *EvalContext
}

func New(ctx *EvalContext) *Evaluator {
	return &Evaluator{ctx: ctx}
}

func (e *Evaluator) Eval(node ast.Node, env *object.Env) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return e.evalProgram(node, env)
	case *ast.HTMLStatement:
		return &object.HTML{Value: node.String()}
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression, env)
	case *ast.IfStatement:
		return e.evalIfStatement(node, env)
	case *ast.BlockStatement:
		return e.evalBlockStatement(node, env)
	case *ast.DefineStatement:
		return e.evalDeclStatement(node, env)
	case *ast.UseStatement:
		return e.evalUseStatement(node, env)
	case *ast.InsertStatement:
		return NIL
	case *ast.ReserveStatement:
		return e.evalReserveStatement(node, env)

	// Expressions
	case *ast.Identifier:
		return e.evalIdentifier(node, env)
	case *ast.IndexExpression:
		return e.evalIndexExpression(node, env)
	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return &object.Str{Value: html.EscapeString(node.Value)}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(node, env)
	case *ast.PrefixExpression:
		return e.evalPrefixExpression(node, env)
	case *ast.TernaryExpression:
		return e.evalTernaryExpression(node, env)
	case *ast.InfixExpression:
		return e.evalInfixExpression(node.Operator, node.Left, node.Right, env)
	case *ast.PostfixExpression:
		return e.evalPostfixExpression(node, env)
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

func (e *Evaluator) evalIfStatement(node *ast.IfStatement, env *object.Env) object.Object {
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

func (e *Evaluator) evalBlockStatement(block *ast.BlockStatement, env *object.Env) object.Object {
	var elems []object.Object

	for _, statement := range block.Statements {
		stmtObj := e.Eval(statement, env)

		if isError(stmtObj) {
			return stmtObj
		}

		elems = append(elems, stmtObj)
	}

	return &object.Block{Elements: elems}
}

func (e *Evaluator) evalDeclStatement(node *ast.DefineStatement, env *object.Env) object.Object {
	val := e.Eval(node.Value, env)

	if isError(val) {
		return val
	}

	env.Set(node.Name.Value, val)

	return NIL
}

func (e *Evaluator) evalUseStatement(node *ast.UseStatement, env *object.Env) object.Object {
	if node.Program == nil {
		return e.newError(node, fail.ErrUseStmtMustHaveProgram)
	}

	if node.Program.IsLayout {
		if hasUseStmt, _ := node.Program.HasUseStmt(); hasUseStmt {
			return e.newError(node, fail.ErrUseStmtNotAllowed)
		}
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

func (e *Evaluator) evalReserveStatement(node *ast.ReserveStatement, env *object.Env) object.Object {
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

func (e *Evaluator) evalIdentifier(node *ast.Identifier, env *object.Env) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	return e.newError(node, fail.ErrIdentifierNotFound, node.Value)
}

func (e *Evaluator) evalIndexExpression(node *ast.IndexExpression, env *object.Env) object.Object {
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
		return e.evalArrayIndexExpression(left, idx)
	}

	return e.newError(node, fail.ErrIndexNotSupported, left.Type())
}

func (e *Evaluator) evalArrayIndexExpression(arr, idx object.Object) object.Object {
	arrObj := arr.(*object.Array)
	index := idx.(*object.Int).Value
	max := int64(len(arrObj.Elements) - 1)

	if index < 0 || index > max {
		return NIL
	}

	return arrObj.Elements[index]
}

func (e *Evaluator) evalPrefixExpression(node *ast.PrefixExpression, env *object.Env) object.Object {
	right := e.Eval(node.Right, env)

	if isError(right) {
		return right
	}

	switch node.Operator {
	case "-":
		return e.evalMinusPrefixOperatorExpression(right, node)
	case "!":
		return e.evalBangOperatorExpression(right, node)
	}

	return e.newError(node, fail.ErrUnknownOperator,
		node.Operator, right.Type())
}

func (e *Evaluator) evalTernaryExpression(node *ast.TernaryExpression, env *object.Env) object.Object {
	condition := e.Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.Eval(node.Consequence, env)
	}

	return e.Eval(node.Alternative, env)
}

func (e *Evaluator) evalArrayLiteral(node *ast.ArrayLiteral, env *object.Env) object.Object {
	elems := e.evalExpressions(node.Elements, env)

	if len(elems) == 1 && isError(elems[0]) {
		return elems[0]
	}

	return &object.Array{Elements: elems}
}

func (e *Evaluator) evalExpressions(exps []ast.Expression, env *object.Env) []object.Object {
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

func (e *Evaluator) evalInfixExpression(operator string, left, right ast.Expression, env *object.Env) object.Object {
	leftObj := e.Eval(left, env)

	if isError(leftObj) {
		return leftObj
	}

	rightObj := e.Eval(right, env)

	if isError(rightObj) {
		return rightObj
	}

	return e.evalInfixOperatorExpression(operator, leftObj, rightObj, left)
}

func (e *Evaluator) evalPostfixExpression(node *ast.PostfixExpression, env *object.Env) object.Object {
	leftObj := e.Eval(node.Left, env)

	if isError(leftObj) {
		return leftObj
	}

	return e.evalPostfixOperatorExpression(leftObj, node.Operator, node)
}

func (e *Evaluator) evalPostfixOperatorExpression(left object.Object, operator string, node ast.Node) object.Object {
	if operator == "++" {
		if left.Is(object.INT_OBJ) {
			value := left.(*object.Int).Value
			return &object.Int{Value: value + 1}
		}

		if left.Is(object.FLOAT_OBJ) {
			value := left.(*object.Float).Value
			return &object.Float{Value: value + 1}
		}
	}

	if operator == "--" {
		if left.Is(object.INT_OBJ) {
			value := left.(*object.Int).Value
			return &object.Int{Value: value - 1}
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

func (e *Evaluator) evalInfixOperatorExpression(operator string, left, right object.Object, leftNode ast.Node) object.Object {
	if left.Type() != right.Type() {
		return e.newError(leftNode, fail.ErrTypeMismatch,
			left.Type(), operator, right.Type())
	}

	if operator == "+" && left.Is(object.STR_OBJ) {
		return e.evalStringInfixExpression(operator, right, left)
	}

	switch left.Type() {
	case object.INT_OBJ:
		return e.evalIntegerInfixExpression(operator, right, left, leftNode)
	case object.FLOAT_OBJ:
		return e.evalFloatInfixExpression(operator, right, left, leftNode)
	}

	return e.newError(leftNode, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
}

func (e *Evaluator) evalStringInfixExpression(operator string, right, left object.Object) object.Object {
	leftVal := left.(*object.Str).Value
	rightVal := right.(*object.Str).Value

	return &object.Str{Value: leftVal + rightVal}
}

func (e *Evaluator) evalIntegerInfixExpression(operator string, right, left object.Object, leftNode ast.Node) object.Object {
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

func (e *Evaluator) evalFloatInfixExpression(operator string, right, left object.Object, leftNode ast.Node) object.Object {
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

func (e *Evaluator) evalMinusPrefixOperatorExpression(right object.Object, node ast.Node) object.Object {
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

func (e *Evaluator) evalBangOperatorExpression(right object.Object, node ast.Node) object.Object {
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

func (e *Evaluator) newError(node ast.Node, format string, a ...interface{}) *object.Error {
	err := fail.New(node.Line(), e.ctx.absPath, "evaluator", format, a...)
	return &object.Error{Err: err}
}
