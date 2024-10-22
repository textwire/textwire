package evaluator

import (
	"github.com/textwire/textwire/v2/config"
	"github.com/textwire/textwire/v2/object"
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
	case nil:
		return false
	}

	return true
}

func isError(obj object.Object) bool {
	return obj.Is(object.ERR_OBJ)
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}

	return FALSE
}

func hasBreakStmt(obj object.Object) bool {
	return hasControlStmt(obj, object.BREAK_OBJ)
}

func hasContinueStmt(obj object.Object) bool {
	return hasControlStmt(obj, object.CONTINUE_OBJ)
}

func hasControlStmt(obj object.Object, controlType object.ObjectType) bool {
	block, isBlock := obj.(*object.Block)

	if !isBlock {
		return obj.Is(controlType)
	}

	// also check recursively for nested blocks
	for _, elem := range block.Elements {
		if hasControlStmt(elem, controlType) {
			return true
		}
	}

	return false
}

// hasCustomFunc checks if the object has a custom function
func hasCustomFunc(customFunc *config.Func, t object.ObjectType) bool {
	if customFunc == nil {
		return false
	}

	switch t {
	case object.STR_OBJ:
		return len(customFunc.Str) > 0
	case object.ARR_OBJ:
		return len(customFunc.Arr) > 0
	case object.INT_OBJ:
		return len(customFunc.Int) > 0
	case object.FLOAT_OBJ:
		return len(customFunc.Float) > 0
	case object.BOOL_OBJ:
		return len(customFunc.Bool) > 0
	default:
		return false
	}
}
