package evaluator

import (
	"bytes"
	"fmt"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/object"
)

var (
	NIL = &object.Nil{}
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

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node.Operator, node.Left, node.Right, env)
	case *ast.NilLiteral:
		return NIL
	}

	return newError("Unknown node type: %T", node)
}

func evalProgram(program *ast.Program, env *object.Env) object.Object {
	var result bytes.Buffer

	for _, statement := range program.Statements {
		stmtObj := Eval(statement, env)

		if err, ok := stmtObj.(*object.Error); ok {
			return err
		}

		result.WriteString(stmtObj.String())
	}

	return &object.Html{Value: result.String()}
}

func evalIdentifier(node *ast.Identifier, env *object.Env) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	return newError("Identifier not found: " + node.Value)
}

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Env) object.Object {
	right := Eval(node.Right, env)

	if isError(right) {
		return right
	}

	switch node.Operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	}

	return newError("Unknown operator: %s%s", node.Operator, right.Type())
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

	switch operator {
	case "+":
		return evalPlusInfixOperatorExpression(leftObj, rightObj)
	}

	return newError("Unknown operator: %s %s %s", leftObj.Type(), operator, rightObj.Type())
}

func evalPlusInfixOperatorExpression(left, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return newError("Type mismatch: %s + %s", left.Type(), right.Type())
	}

	switch left.Type() {
	case object.STRING_OBJ:
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.String{Value: leftVal + rightVal}
	case object.INT_OBJ:
		leftVal := left.(*object.Int).Value
		right := right.(*object.Int)
		return &object.Int{Value: leftVal + right.Value}
	}

	return newError("Unknown type for + operator: %s", left.Type())
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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}
