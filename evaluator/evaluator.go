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

func Eval(node ast.Node, env *object.Env) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.HTMLStatement:
		return &object.HTML{Value: node.String()}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.DefineStatement:
		return evalDeclStatement(node, env)
	case *ast.UseStatement:
		return evalUseStatement(node, env)
	case *ast.InsertStatement:
		return NIL
	case *ast.ReserveStatement:
		return evalReserveStatement(node, env)

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return &object.Str{Value: html.EscapeString(node.Value)}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.TernaryExpression:
		return evalTernaryExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node.Operator, node.Left, node.Right, env)
	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)
	case *ast.NilLiteral:
		return NIL
	}

	return newError(node, fail.ErrUnknownNodeType, node)
}

func evalProgram(prog *ast.Program, env *object.Env) object.Object {
	var out bytes.Buffer

	for _, statement := range prog.Statements {
		stmtObj := Eval(statement, env)

		if isError(stmtObj) {
			return stmtObj
		}

		out.WriteString(stmtObj.String())
	}

	return &object.HTML{Value: out.String()}
}

func evalIfStatement(node *ast.IfStatement, env *object.Env) object.Object {
	condition := Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	newEnv := object.NewEnclosedEnv(env)

	if isTruthy(condition) {
		return Eval(node.Consequence, newEnv)
	}

	for _, alt := range node.Alternatives {
		condition = Eval(alt.Condition, env)

		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(alt.Consequence, newEnv)
		}
	}

	if node.Alternative != nil {
		return Eval(node.Alternative, newEnv)
	}

	return NIL
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Env) object.Object {
	var elems []object.Object

	for _, statement := range block.Statements {
		stmtObj := Eval(statement, env)

		if isError(stmtObj) {
			return stmtObj
		}

		elems = append(elems, stmtObj)
	}

	return &object.Block{Elements: elems}
}

func evalDeclStatement(node *ast.DefineStatement, env *object.Env) object.Object {
	val := Eval(node.Value, env)

	if isError(val) {
		return val
	}

	env.Set(node.Name.Value, val)

	return NIL
}

func evalUseStatement(node *ast.UseStatement, env *object.Env) object.Object {
	if node.Program == nil {
		return newError(node, "The 'use' statement must have a program attached")
	}

	layoutContent := Eval(node.Program, env)

	if isError(layoutContent) {
		return layoutContent
	}

	return &object.Use{
		Path:    node.Name.Value,
		Content: layoutContent,
	}
}

func evalReserveStatement(node *ast.ReserveStatement, env *object.Env) object.Object {
	stmt := &object.Reserve{Name: node.Name.Value}

	if node.Insert.Block != nil {
		result := Eval(node.Insert.Block, env)

		if isError(result) {
			return result
		}

		stmt.Content = result

		return stmt
	}

	if node.Insert.Argument == nil {
		return newError(node.Insert, fail.ErrInsertMustHaveContent)
	}

	firstArg := Eval(node.Insert.Argument, env)

	if isError(firstArg) {
		return firstArg
	}

	stmt.Argument = firstArg

	return stmt
}

func evalIdentifier(node *ast.Identifier, env *object.Env) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	return newError(node, fail.ErrIdentifierNotFound, node.Value)
}

func evalIndexExpression(node *ast.IndexExpression, env *object.Env) object.Object {
	left := Eval(node.Left, env)

	if isError(left) {
		return left
	}

	idx := Eval(node.Index, env)

	if isError(idx) {
		return idx
	}

	switch {
	case left.Is(object.ARR_OBJ) && idx.Is(object.INT_OBJ):
		return evalArrayIndexExpression(left, idx)
	}

	return newError(node, fail.ErrIndexNotSupported, left.Type())
}

func evalArrayIndexExpression(arr, idx object.Object) object.Object {
	arrObj := arr.(*object.Array)
	index := idx.(*object.Int).Value
	max := int64(len(arrObj.Elements) - 1)

	if index < 0 || index > max {
		return NIL
	}

	return arrObj.Elements[index]
}

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Env) object.Object {
	right := Eval(node.Right, env)

	if isError(right) {
		return right
	}

	switch node.Operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right, node)
	case "!":
		return evalBangOperatorExpression(right, node)
	}

	return newError(node, fail.ErrUnknownOperator,
		node.Operator, right.Type())
}

func evalTernaryExpression(node *ast.TernaryExpression, env *object.Env) object.Object {
	condition := Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	}

	return Eval(node.Alternative, env)
}

func evalArrayLiteral(node *ast.ArrayLiteral, env *object.Env) object.Object {
	elems := evalExpressions(node.Elements, env)

	if len(elems) == 1 && isError(elems[0]) {
		return elems[0]
	}

	return &object.Array{Elements: elems}
}

func evalExpressions(exps []ast.Expression, env *object.Env) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalInfixExpression(operator string, left, right ast.Expression, env *object.Env) object.Object {
	leftObj := Eval(left, env)

	if isError(leftObj) {
		return leftObj
	}

	rightObj := Eval(right, env)

	if isError(rightObj) {
		return rightObj
	}

	return evalInfixOperatorExpression(operator, leftObj, rightObj, left)
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Env) object.Object {
	leftObj := Eval(node.Left, env)

	if isError(leftObj) {
		return leftObj
	}

	return evalPostfixOperatorExpression(leftObj, node.Operator, node)
}

func evalPostfixOperatorExpression(left object.Object, operator string, node ast.Node) object.Object {
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

	return newError(node, fail.ErrUnknownOperator,
		left.Type(), operator)
}

func evalInfixOperatorExpression(operator string, left, right object.Object, leftNode ast.Node) object.Object {
	if left.Type() != right.Type() {
		return newError(leftNode, fail.ErrTypeMismatch,
			left.Type(), operator, right.Type())
	}

	if operator == "+" && left.Is(object.STR_OBJ) {
		return evalStringInfixExpression(operator, right, left)
	}

	switch left.Type() {
	case object.INT_OBJ:
		return evalIntegerInfixExpression(operator, right, left, leftNode)
	case object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, right, left, leftNode)
	}

	return newError(leftNode, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
}

func evalStringInfixExpression(operator string, right, left object.Object) object.Object {
	leftVal := left.(*object.Str).Value
	rightVal := right.(*object.Str).Value

	return &object.Str{Value: leftVal + rightVal}
}

func evalIntegerInfixExpression(operator string, right, left object.Object, leftNode ast.Node) object.Object {
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

	return newError(leftNode, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
}

func evalFloatInfixExpression(operator string, right, left object.Object, leftNode ast.Node) object.Object {
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

	return newError(leftNode, fail.ErrUnknownTypeForOperator,
		left.Type(), operator)
}

func evalMinusPrefixOperatorExpression(right object.Object, node ast.Node) object.Object {
	switch right.Type() {
	case object.INT_OBJ:
		value := right.(*object.Int).Value
		return &object.Int{Value: -value}
	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}

	return newError(node, fail.ErrPrefixOperatorIsWrong,
		"-", right.Type())
}

func evalBangOperatorExpression(right object.Object, node ast.Node) object.Object {
	switch right {
	case FALSE:
		return TRUE
	case TRUE:
		return FALSE
	case NIL:
		return TRUE
	}

	return newError(node, fail.ErrPrefixOperatorIsWrong,
		"!", right.Type())
}
