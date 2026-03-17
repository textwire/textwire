package evaluator

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/textwire/textwire/v4/config"
	"github.com/textwire/textwire/v4/pkg/fail"
	"github.com/textwire/textwire/v4/pkg/value"
)

func isTruthy(obj value.Value) bool {
	switch obj := obj.(type) {
	case *value.Bool:
		return obj.Val
	case *value.Int:
		return obj.Val != 0
	case *value.Float:
		return obj.Val != 0.0
	case *value.Str:
		return obj.Val != ""
	case *value.Nil:
		return false
	case *value.Obj:
		return len(obj.Pairs) > 0
	case *value.Arr:
		return len(obj.Elements) != 0
	case nil:
		return false
	}

	return true
}

func isError(obj value.Value) bool {
	return obj.Is(value.ERR_VAL)
}

func isUndefinedError(obj value.Value) bool {
	undefinedErrors := []string{
		fail.ErrVariableIsUndefined,
		fail.ErrKeyOnNonObj,
	}

	err, isErr := obj.(*value.Error)
	return isErr && slices.Contains(undefinedErrors, err.ErrorID)
}

func nativeBoolToBoolObj(input bool) value.Literal {
	if input {
		return TRUE
	}
	return FALSE
}

func hasBreak(obj value.Value) bool {
	return hasControlStmt(obj, value.BREAK_VAL)
}

func hasContinue(obj value.Value) bool {
	return hasControlStmt(obj, value.CONTINUE_VAL)
}

func hasControlStmt(obj value.Value, controlType value.ValueType) bool {
	block, ok := obj.(*value.Block)
	if !ok {
		return obj.Is(controlType)
	}

	// also check recursively for nested blocks
	for _, elem := range block.Chunks {
		if hasControlStmt(elem, controlType) {
			return true
		}
	}

	return false
}

// hasCustomFunc checks if the object has a custom function
func hasCustomFunc(customFunc *config.Func, t value.ValueType, funcName string) bool {
	if customFunc == nil {
		return false
	}

	switch t {
	case value.STR_VAL:
		return customFunc.Str[funcName] != nil
	case value.ARR_VAL:
		return customFunc.Arr[funcName] != nil
	case value.INT_VAL:
		return customFunc.Int[funcName] != nil
	case value.FLOAT_VAL:
		return customFunc.Float[funcName] != nil
	case value.BOOL_VAL:
		return customFunc.Bool[funcName] != nil
	case value.OBJ_VAL:
		return customFunc.Obj[funcName] != nil
	default:
		return false
	}
}

// getDecimalConfig extracts separator and decimal count from arguments.
// Returns error if arguments are invalid.
func getDecimalConfig(
	objType value.ValueType,
	args ...value.Literal,
) (separator string, decimals int, err error) {
	separator = "."
	decimals = 2

	if len(args) > 2 {
		return "", 0, fmt.Errorf(fail.ErrFuncMaxArgs, objType, "decimal", 2)
	}

	if len(args) >= 1 {
		separatorArg, ok := args[0].(*value.Str)
		if !ok {
			return "", 0, fmt.Errorf(fail.ErrFuncFirstArgStr, objType, "decimal")
		}
		separator = separatorArg.Val
	}

	if len(args) == 2 {
		decimalArg, ok := args[1].(*value.Int)
		if !ok {
			return "", 0, fmt.Errorf(fail.ErrFuncSecondArgInt, objType, "decimal")
		}
		decimals = int(decimalArg.Val)
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
	formatted := fmt.Sprintf("%."+strconv.Itoa(decimals)+"f", f)
	if separator != "." {
		formatted = strings.Replace(formatted, ".", separator, 1)
	}
	return formatted
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
