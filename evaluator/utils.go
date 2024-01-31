package evaluator

import (
	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/fail"
	"github.com/textwire/textwire/object"
)

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Bool:
		return obj.Value
	case *object.Int:
		return obj.Value != 0
	case *object.Float:
		return obj.Value != 0.0
	case *object.Str:
		return obj.Value != ""
	case *object.Nil:
		return false
	}

	return true
}

func newError(node ast.Node, format string, a ...interface{}) *object.Error {
	err := fail.New(node.LineNum(), "interpreter", format, a...)
	return &object.Error{Err: err}
}

func isError(obj object.Object) bool {
	return obj.Is(object.ERROR_OBJ)
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}

	return FALSE
}
