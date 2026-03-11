package evaluator

import (
	"github.com/textwire/textwire/v3/pkg/value"
)

var functions = map[value.ValueType]map[string]*value.Builtin{
	value.STR_VAL: {
		"len":        {Fn: strLenFunc},
		"split":      {Fn: strSplitFunc},
		"raw":        {Fn: strRawFunc},
		"trim":       {Fn: strTrimFunc},
		"trimRight":  {Fn: strTrimRightFunc},
		"trimLeft":   {Fn: strTrimLeftFunc},
		"upper":      {Fn: strUpperFunc},
		"lower":      {Fn: strLowerFunc},
		"capitalize": {Fn: strCapitalizeFunc},
		"reverse":    {Fn: strReverseFunc},
		"contains":   {Fn: strContainsFunc},
		"truncate":   {Fn: strTruncateFunc},
		"decimal":    {Fn: strDecimalFunc},
		"at":         {Fn: strAtFunc},
		"first":      {Fn: strFirstFunc},
		"last":       {Fn: strLastFunc},
		"repeat":     {Fn: strRepeatFunc},
		"format":     {Fn: strFormatFunc},
	},
	value.ARR_VAL: {
		"len":      {Fn: arrLenFunc},
		"join":     {Fn: arrJoinFunc},
		"rand":     {Fn: arrRandFunc},
		"reverse":  {Fn: arrReverseFunc},
		"slice":    {Fn: arrSliceFunc},
		"shuffle":  {Fn: arrShuffleFunc},
		"contains": {Fn: arrContainsFunc},
		"append":   {Fn: arrAppendFunc},
		"prepend":  {Fn: arrPrependFunc},
		"json":     {Fn: jsonFunc},
	},
	value.FLOAT_VAL: {
		"int":   {Fn: floatIntFunc},
		"str":   {Fn: floatStrFunc},
		"abs":   {Fn: floatAbsFunc},
		"ceil":  {Fn: floatCeilFunc},
		"floor": {Fn: floatFloorFunc},
		"round": {Fn: floatRoundFunc},
	},
	value.INT_VAL: {
		"float":   {Fn: intFloatFunc},
		"abs":     {Fn: intAbsFunc},
		"str":     {Fn: intStrFunc},
		"len":     {Fn: intLenFunc},
		"decimal": {Fn: intDecimalFunc},
	},
	value.BOOL_VAL: {
		"binary": {Fn: boolBinaryFunc},
		"then":   {Fn: boolThenFunc},
	},
	value.OBJ_VAL: {
		"json":  {Fn: jsonFunc},
		"camel": {Fn: objCamelFunc},
		"get":   {Fn: objGetFunc},
	},
}

// jsonFunc convert value to json representation
func jsonFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	json, err := receiver.JSON()
	if err != nil {
		return nil, err
	}
	return &value.Str{Val: json}, nil
}
