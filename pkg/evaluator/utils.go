package evaluator

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"slices"
	"strings"
	"unicode"

	"github.com/textwire/textwire/v6/pkg/ast"
	"github.com/textwire/textwire/v6/pkg/fail"
	"github.com/textwire/textwire/v6/pkg/value"
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

// isMissingNil reports whether value is a nil produced by a missing object
// property or array index.
func isMissingNil(obj value.Value) bool {
	nilVal, ok := obj.(*value.Nil)
	if !ok {
		return false
	}

	switch nilVal.Origin() {
	case value.NilOriginMissingProp, value.NilOriginMissingIndex:
		return true
	}

	return false
}

// literalsEqual reports whether two literals represent the same value. Nil
// origins are ignored because they describe where nil came from, not its value.
// The raw mark on strings is ignored too, it only changes how they're printed.
func literalsEqual(left, right value.Literal) bool {
	switch left := left.(type) {
	case *value.Nil:
		_, ok := right.(*value.Nil)
		return ok
	case *value.Str:
		right, ok := right.(*value.Str)
		return ok && left.Val == right.Val
	case *value.Arr:
		right, ok := right.(*value.Arr)
		return ok && arrsEqual(left, right)
	case *value.Obj:
		right, ok := right.(*value.Obj)
		return ok && objsEqual(left, right)
	}

	return reflect.DeepEqual(left, right)
}

func arrsEqual(left, right *value.Arr) bool {
	if len(left.Elements) != len(right.Elements) {
		return false
	}

	for i := range left.Elements {
		if !literalsEqual(left.Elements[i], right.Elements[i]) {
			return false
		}
	}

	return true
}

func objsEqual(left, right *value.Obj) bool {
	if len(left.Pairs) != len(right.Pairs) {
		return false
	}

	for key, leftVal := range left.Pairs {
		rightVal, ok := right.Pairs[key]
		if !ok || !literalsEqual(leftVal, rightVal) {
			return false
		}
	}

	return true
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

// stripControlDirs removes break/continue values from a block so a loop
// doesn't pass control statements it has already handled to its parent loop.
func stripControlDirs(obj value.Value) value.Value {
	block, ok := obj.(*value.Block)
	if !ok {
		return obj
	}

	if !block.HasBreak() && !block.HasContinue() {
		return obj
	}

	stripped := value.NewBlock(len(block.Chunks()))
	for _, chunk := range block.Chunks() {
		if value.IsControlDir(chunk) {
			continue
		}

		stripped.Append(stripControlDirs(chunk))
	}

	return stripped
}

// getDecimalConfig extracts the separator and decimal count from args.
// Both are optional and default to "." and 2. A negative decimal count
// is clamped to 0. It returns an error if there are too many arguments
// or if an argument has the wrong type.
func getDecimalConfig(
	objType value.ValueType,
	args ...value.Literal,
) (string, int, error) {
	if len(args) > 2 {
		return "", 0, fmt.Errorf(fail.ErrFuncMaxArgs, objType, "decimal", 2)
	}

	separator := "."
	decimals := 2

	if len(args) > 0 {
		arg, ok := args[0].(*value.Str)
		if !ok {
			return "", 0, fmt.Errorf(fail.ErrFuncFirstArgStr, objType, "decimal")
		}
		separator = arg.Val
	}

	if len(args) > 1 {
		arg, ok := args[1].(*value.Int)
		if !ok {
			return "", 0, fmt.Errorf(fail.ErrFuncSecondArgInt, objType, "decimal")
		}
		decimals = max(int(arg.Val), 0)
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
// - Changes separator to the one you pass.
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

const zeroWidthJoiner = '\u200D'

// isRegionalIndicator reports whether r is one of the regional indicator
// symbols used to build flag emoji.
func isRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

// isClusterExtender reports whether r continues the cluster of the rune before it.
func isClusterExtender(r rune) bool {
	skinToneMods := r >= 0x1F3FB && r <= 0x1F3FF
	return unicode.IsMark(r) || skinToneMods || r == zeroWidthJoiner
}

// splitClusters splits s into approximate user-perceived characters: each
// base rune together with any combining marks, skin tone modifiers or ZWJ
// sequence parts that follow it. Flags built from two regional indicators
// stay together too. The returned strings are substrings of s, so no
// per-cluster copying occurs.
func splitClusters(s string) []string {
	var clusters []string

	start := 0           // byte offset where the current cluster begins
	afterZWJ := false    // previous rune was a zero width joiner
	hasOpenFlag := false // current cluster ends in a regional indicator still waiting for its pair

	for i, r := range s {
		isRI := isRegionalIndicator(r)
		isExtender := isClusterExtender(r)

		continuesCluster := afterZWJ || isExtender || (isRI && hasOpenFlag)

		// Save the finished cluster when the new one begins
		if i > 0 && !continuesCluster {
			clusters = append(clusters, s[start:i])
			start = i
		}

		switch {
		case isRI:
			// First RI opens a flag, the second one completes it.
			hasOpenFlag = !hasOpenFlag
		case isExtender:
			// Extenders don't change whether a flag is pending.
		default:
			hasOpenFlag = false
		}

		afterZWJ = r == zeroWidthJoiner
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

// floatToInt wraps a float as a value.Int, failing when it is outside the int64 range
func floatToInt(val float64, method string) (value.Literal, error) {
	if !floatInInt64Range(val) {
		msg := fmt.Sprintf(fail.ErrFuncOutOfInt64Range, value.FLOAT_VAL, method)
		return nil, errors.New(msg)
	}

	return &value.Int{Val: int64(val)}, nil
}

// floatInInt64Range reports whether val can be converted to int64 without overflowing
func floatInInt64Range(val float64) bool {
	return val >= math.MinInt64 && val < -math.MinInt64
}

// elementsMatch compares two literals, using deep equality for
// composite types (objects and arrays) and native equality otherwise.
func elementsMatch(a, b value.Literal) bool {
	if isComposite(a, b) {
		return literalsEqual(a, b)
	}

	return a.Native() == b.Native()
}

// isComposite reports whether both literals are objects or both are arrays.
func isComposite(a, b value.Literal) bool {
	if a.Type() != b.Type() {
		return false
	}

	return a.Type() == value.OBJ_VAL || a.Type() == value.ARR_VAL
}
