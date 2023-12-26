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
	case *ast.ReturnStatement:
		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.NilLiteral:
		return NIL
	}

	return nil
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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}
