package evaluator

import (
	"bytes"
	"html"

	"github.com/textwire/textwire/ast"
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
		return &object.Html{Value: node.String()}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.DefineStatement:
		return evalDeclStatement(node, env)
	case *ast.LayoutStatement:
		return evalLayoutStatement(node, env)
	case *ast.InsertStatement:
		return NIL
	case *ast.ReserveStatement:
		return evalReserveStatement(node, env)

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(node, env)
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
	case *ast.NilLiteral:
		return NIL
	}

	return newError("Unknown node type: %T", node)
}

func evalProgram(prog *ast.Program, env *object.Env) object.Object {
	var result bytes.Buffer

	for _, statement := range prog.Statements {
		stmtObj := Eval(statement, env)

		if err, ok := stmtObj.(*object.Error); ok {
			return err
		}

		result.WriteString(stmtObj.String())
	}

	return &object.Html{Value: result.String()}
}

func evalIfStatement(node *ast.IfStatement, env *object.Env) object.Object {
	condition := Eval(node.Condition, env)
	newEnv := object.NewEnclosedEnv(env)

	if isTruthy(condition) {
		return Eval(node.Consequence, newEnv)
	}

	for _, alt := range node.Alternatives {
		condition = Eval(alt.Condition, env)

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

		if err, ok := stmtObj.(*object.Error); ok {
			return err
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

func evalLayoutStatement(node *ast.LayoutStatement, env *object.Env) object.Object {
	if node.Program == nil {
		return newError("The layout statement must have a program attached")
	}

	return &object.Layout{
		Path:    node.Path.Value,
		Content: Eval(node.Program, env),
	}
}

func evalReserveStatement(node *ast.ReserveStatement, env *object.Env) object.Object {
	stmt := &object.Reserve{Name: node.Name.Value}

	if node.Insert.Block != nil {
		stmt.Content = Eval(node.Insert.Block, env)
		return stmt
	}

	if node.Insert.Argument == nil {
		return newError("The 'insert' statement must have a content or a text argument")
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

	return newError(`Identifier "` + node.Value + `" not found`)
}

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Env) object.Object {
	right := Eval(node.Right, env)

	if isError(right) {
		return right
	}

	switch node.Operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "!":
		return evalBangOperatorExpression(right)
	}

	return newError("Unknown operator: %s%s", node.Operator, right.Type())
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

	return evalInfixOperatorExpression(operator, leftObj, rightObj)
}

func evalInfixOperatorExpression(operator string, left, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return newError("Type mismatch: %s + %s", left.Type(), right.Type())
	}

	if operator == "+" && left.Type() == object.STRING_OBJ {
		leftVal := left.(*object.Str).Value
		rightVal := right.(*object.Str).Value
		return &object.Str{Value: leftVal + rightVal}
	}

	switch left.Type() {
	case object.INT_OBJ:
		return evalIntegerInfixExpression(operator, right, left)
	case object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, right, left)
	}

	return newError("Unknown type for %s operator: %s", operator, left.Type())
}

func evalIntegerInfixExpression(operator string, right, left object.Object) object.Object {
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
	}

	return newError("Unknown type for %s operator: %s", operator, left.Type())
}

func evalFloatInfixExpression(operator string, right, left object.Object) object.Object {
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
	}

	return newError("Unknown type for %s operator: %s", operator, left.Type())
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.INT_OBJ:
		value := right.(*object.Int).Value
		return &object.Int{Value: -value}
	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}

	return newError("Unknown operator: -%s", right.Type())
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case FALSE:
		return TRUE
	case TRUE:
		return FALSE
	case NIL:
		return TRUE
	}

	return newError("Unknown operator: !%s", right.Type())
}
