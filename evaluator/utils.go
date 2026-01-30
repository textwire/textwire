package evaluator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/textwire/textwire/v2/config"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
	"github.com/textwire/textwire/v2/utils"
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
func hasCustomFunc(customFunc *config.Func, t object.ObjectType, funcName string) bool {
	if customFunc == nil {
		return false
	}

	switch t {
	case object.STR_OBJ:
		return customFunc.Str[funcName] != nil
	case object.ARR_OBJ:
		return customFunc.Arr[funcName] != nil
	case object.INT_OBJ:
		return customFunc.Int[funcName] != nil
	case object.FLOAT_OBJ:
		return customFunc.Float[funcName] != nil
	case object.BOOL_OBJ:
		return customFunc.Bool[funcName] != nil
	case object.OBJ_OBJ:
		return customFunc.Obj[funcName] != nil
	default:
		return false
	}
}

func addDecimals(
	receiver object.Object,
	objType object.ObjectType,
	args ...object.Object,
) (object.Object, error) {
	var val string

	switch objType {
	case object.STR_OBJ:
		val = receiver.(*object.Str).Value
	case object.INT_OBJ:
		val = receiver.(*object.Int).String()
	}

	if !utils.StrIsInt(val) {
		return &object.Str{Value: val}, nil
	}

	separator := "."
	decimals := 2

	if len(args) > 2 {
		msg := fmt.Sprintf(fail.ErrFuncMaxArgs, "decimal", objType, 2)
		return nil, errors.New(msg)
	}

	if len(args) >= 1 {
		separatorArg, ok := args[0].(*object.Str)

		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, "decimal", objType)
			return nil, errors.New(msg)
		}

		separator = separatorArg.Value
	}

	if len(args) == 2 {
		decimalArg, ok := args[1].(*object.Int)

		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncSecondArgInt, "decimal", objType)
			return nil, errors.New(msg)
		}

		decimals = int(decimalArg.Value)
	}

	zeros := strings.Repeat("0", decimals)
	if decimals == 0 {
		return &object.Str{Value: val}, nil
	}

	return &object.Str{Value: val + separator + zeros}, nil
}

func isUndefinedVarError(obj object.Object) bool {
	err, isErr := obj.(*object.Error)
	return isErr &&
		(err.ErrorID == fail.ErrIdentifierIsUndefined || err.ErrorID == fail.ErrPropertyNotFound)
}
