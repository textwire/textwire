package evaluator

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/object"
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
	case *object.Obj:
		return len(obj.Pairs) > 0
	case *object.Array:
		return len(obj.Elements) != 0
	case nil:
		return false
	}

	return true
}

func isError(obj object.Object) bool {
	return obj.Is(object.ERR_OBJ)
}

func isUndefinedError(obj object.Object) bool {
	undefinedErrors := []string{
		fail.ErrVariableIsUndefined,
		fail.ErrKeyOnNonObject,
	}

	err, isErr := obj.(*object.Error)
	return isErr && slices.Contains(undefinedErrors, err.ErrorID)
}

func nativeBoolToBoolObj(input bool) object.Object {
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
	block, ok := obj.(*object.Block)
	if !ok {
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

// getDecimalConfig extracts separator and decimal count from arguments.
// Returns error if arguments are invalid.
func getDecimalConfig(
	objType object.ObjectType,
	args ...object.Object,
) (separator string, decimals int, err error) {
	separator = "."
	decimals = 2

	if len(args) > 2 {
		return "", 0, fmt.Errorf(fail.ErrFuncMaxArgs, objType, "decimal", 2)
	}

	if len(args) >= 1 {
		separatorArg, ok := args[0].(*object.Str)
		if !ok {
			return "", 0, fmt.Errorf(fail.ErrFuncFirstArgStr, objType, "decimal")
		}
		separator = separatorArg.Value
	}

	if len(args) == 2 {
		decimalArg, ok := args[1].(*object.Int)
		if !ok {
			return "", 0, fmt.Errorf(fail.ErrFuncSecondArgInt, objType, "decimal")
		}
		decimals = int(decimalArg.Value)
	}

	return separator, decimals, nil
}

// formatIntDecimals appends decimal places to an integer string.
func formatIntDecimals(val, separator string, decimals int) string {
	if decimals == 0 {
		return val
	}
	return val + separator + strings.Repeat("0", decimals)
}

// formatFloatDecimals ensures a float string has at least the requested decimal places.
func formatFloatDecimals(val, separator string, decimals int) string {
	parts := strings.Split(val, ".")
	if len(parts) == 2 && len(parts[1]) >= decimals {
		return val
	}

	f, _ := strconv.ParseFloat(val, 64)
	result := fmt.Sprintf("%."+strconv.Itoa(decimals)+"f", f)
	if separator != "." {
		result = strings.Replace(result, ".", separator, 1)
	}
	return result
}

// isValidFloat checks if string represents a valid float.
func isValidFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func strIsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// capitalizeFirst returns the string with its first character uppercased using
// a stack buffer to avoid heap allocations. Returns the original string
// unchanged if the first character is not a lowercase ASCII letter.
func capitalizeFirst(s string) string {
	first := s[0]
	if first < 'a' || first > 'z' {
		return s
	}

	var buf [64]byte
	if len(s) > 64 {
		return string(first-32) + s[1:]
	}

	buf[0] = first - 32
	copy(buf[1:], s[1:])

	return string(buf[:len(s)])
}
