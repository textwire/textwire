package evaluator

import (
	"github.com/textwire/textwire/v3/pkg/value"
)

var functions = map[value.ValueType]map[string]*value.Builtin{
	value.STR_OBJ: {
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
	value.ARR_OBJ: {
		"len":      {Fn: arrayLenFunc},
		"join":     {Fn: arrayJoinFunc},
		"rand":     {Fn: arrayRandFunc},
		"reverse":  {Fn: arrayReverseFunc},
		"slice":    {Fn: arraySliceFunc},
		"shuffle":  {Fn: arrayShuffleFunc},
		"contains": {Fn: arrayContainsFunc},
		"append":   {Fn: arrayAppendFunc},
		"prepend":  {Fn: arrayPrependFunc},
		"json":     {Fn: jsonFunc},
	},
	value.FLOAT_OBJ: {
		"int":   {Fn: floatIntFunc},
		"str":   {Fn: floatStrFunc},
		"abs":   {Fn: floatAbsFunc},
		"ceil":  {Fn: floatCeilFunc},
		"floor": {Fn: floatFloorFunc},
		"round": {Fn: floatRoundFunc},
	},
	value.INT_OBJ: {
		"float":   {Fn: intFloatFunc},
		"abs":     {Fn: intAbsFunc},
		"str":     {Fn: intStrFunc},
		"len":     {Fn: intLenFunc},
		"decimal": {Fn: intDecimalFunc},
	},
	value.BOOL_OBJ: {
		"binary": {Fn: boolBinaryFunc},
		"then":   {Fn: boolThenFunc},
	},
	value.OBJ_OBJ: {
		"json":  {Fn: jsonFunc},
		"camel": {Fn: objCamelFunc},
		"get":   {Fn: objGetFunc},
	},
}

// jsonFunc convert value to json representation
func jsonFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	json, err := receiver.JSON()
	if err != nil {
		return nil, err
	}
	return &value.Str{Val: json}, nil
}
