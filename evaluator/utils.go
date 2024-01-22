package evaluator

import (
	"fmt"

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

func newError(format string, a ...interface{}) *object.Error {
	message := fmt.Sprintf("TEXTWIRE ERROR: "+format, a...)
	return &object.Error{Message: message}
}

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}

	return FALSE
}
