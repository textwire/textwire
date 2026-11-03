package evaluator

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"codeberg.org/textwire/textwire/v5/config"
	"codeberg.org/textwire/textwire/v5/pkg/ast"
	"codeberg.org/textwire/textwire/v5/pkg/fail"
	"codeberg.org/textwire/textwire/v5/pkg/value"
)

func isTruthy(obj value.Literal) bool {
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

// exprToText formats an expression the way it appears in source code, without
// the decorative parentheses used by ast.Expression.String().
func exprToText(expr ast.Expression) string {
	if expr == nil {
		return ""
	}

	switch expr := expr.(type) {
	case *ast.IdentExpr:
		return expr.Name
	case *ast.IndexExpr:
		return exprToText(expr.Left) + "[" + exprToText(expr.Index) + "]"
	case *ast.DotExpr:
		return exprToText(expr.Left) + "." + exprToText(expr.Key)
	}

	text := expr.String()
	if len(text) >= 2 && text[0] == '(' && text[len(text)-1] == ')' {
		return text[1 : len(text)-1]
	}

	return text
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

// stripControlStmts removes break/continue values from a block so a loop
// doesn't pass control statements it has already handled to its parent loop.
func stripControlStmts(obj value.Value) value.Value {
	block, ok := obj.(*value.Block)
	if !ok {
		return obj
	}

	chunks := make([]value.Value, 0, len(block.Chunks))
	for _, chunk := range block.Chunks {
		if chunk.Is(value.BREAK_VAL) || chunk.Is(value.CONTINUE_VAL) {
			continue
		}

		chunks = append(chunks, stripControlStmts(chunk))
	}

	block.Chunks = chunks

	return block
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
		decimals = max(int(decimalArg.Val), 0)
	}

	return separator, decimals, nil
}

// checkDecimalResultSize returns an error if adding decimals to val would make
// the result longer than maxRepeatLen bytes. Nothing is added when decimals is
// not positive, so that case always passes.
func checkDecimalResultSize(
	objType value.ValueType,
	val, separator string,
	decimals int,
) error {
	if decimals <= 0 {
		return nil
	}

	// Bytes left for decimal digits after the value and separator.
	room := maxRepeatLen - len(val) - len(separator)
	if decimals > room {
		return fmt.Errorf(fail.ErrFuncResultTooLarge, objType, "decimal", maxRepeatLen)
	}

	return nil
}

// formatIntDecimals appends decimal places to an integer string.
func formatIntDecimals(val, separator string, decimals int) string {
	if decimals == 0 {
		return val
	}
	return val + separator + strings.Repeat("0", decimals)
}

// ensureMinFloatDecimals ensures a float string has at least the requested
// decimal places.
// - Only adds decimals, never removes.
// - Only works with floats with dot separator.
// - Changes seperator to the one you pass.
func ensureMinFloatDecimals(val, separator string, decimals int) string {
	intPart, fracPart, found := strings.Cut(val, ".")
	if !found {
		return val
	}

	missing := decimals - len(fracPart)
	if missing > 0 {
		fracPart += strings.Repeat("0", missing)
	}

	return intPart + separator + fracPart
}

// isValidFloat checks if string represents a valid float, e.g. "1.5",
// "-0.5", "1.". Special values like "nan" or "inf" and exponent notation
// are not valid.
func isValidFloat(s string) bool {
	intPart, fracPart, found := strings.Cut(trimSign(s), ".")
	if !found || intPart+fracPart == "" {
		return false
	}

	return isDigits(intPart) && isDigits(fracPart)
}

func strIsInt(s string) bool {
	digits := trimSign(s)
	return digits != "" && isDigits(digits)
}

// trimSign removes one leading '+' or '-' sign from s.
func trimSign(s string) string {
	if s == "" {
		return s
	}

	switch s[0] {
	case '+', '-':
		return s[1:]
	}

	return s
}

// isDigits reports whether s contains only ASCII digits.
// Empty string is valid.
func isDigits(str string) bool {
	for _, char := range str {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// capitalizeFirst returns the string with its first character uppercased using
// a stack buffer to avoid heap allocations. Returns the original string
// unchanged if the first character is not a lowercase ASCII letter.
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}

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

// isClusterExtender reports whether r continues the cluster of the rune before it.
func isClusterExtender(r rune) bool {
	return unicode.IsMark(r) || (r >= 0x1F3FB && r <= 0x1F3FF)
}

// splitClusters splits s into approximate user-perceived characters: each
// base rune together with any combining marks or skin tone modifiers that
// follow it. The returned
// strings are substrings of s, so no per-cluster copying occurs.
func splitClusters(s string) []string {
	var clusters []string
	start := 0

	for i, r := range s {
		// A non-mark rune (other than the first) begins a new cluster,
		// so close out the previous one.
		if i > 0 && !isClusterExtender(r) {
			clusters = append(clusters, s[start:i])
			start = i
		}
	}

	if start < len(s) {
		clusters = append(clusters, s[start:])
	}

	return clusters
}

// getDateTimeLayout takes string date and returns Go's layout format.
// Go's time.Time layouts: https://pkg.go.dev/time#pkg-constants
func getDateTimeLayout(date string) string {
	switch {
	case len(date) == 8 && date[2] == ':' && date[5] == ':':
		return "15:04:05"
	case len(date) == 19 && date[4] == '-' && date[7] == '-' && date[10] == ' ' && date[13] == ':' && date[16] == ':':
		return "2006-01-02 15:04:05"
	case len(date) == 10 && date[4] == '-' && date[7] == '-':
		return "2006-01-02"
	}
	return ""
}
