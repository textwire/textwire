package evaluator

import (
	"github.com/textwire/textwire/v3/pkg/object"
)

var functions = map[object.ObjectType]map[string]*object.Builtin{
	object.STR_OBJ: {
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
	object.ARR_OBJ: {
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
	object.FLOAT_OBJ: {
		"int":   {Fn: floatIntFunc},
		"str":   {Fn: floatStrFunc},
		"abs":   {Fn: floatAbsFunc},
		"ceil":  {Fn: floatCeilFunc},
		"floor": {Fn: floatFloorFunc},
		"round": {Fn: floatRoundFunc},
	},
	object.INT_OBJ: {
		"float":   {Fn: intFloatFunc},
		"abs":     {Fn: intAbsFunc},
		"str":     {Fn: intStrFunc},
		"len":     {Fn: intLenFunc},
		"decimal": {Fn: intDecimalFunc},
	},
	object.BOOL_OBJ: {
		"binary": {Fn: boolBinaryFunc},
		"then":   {Fn: boolThenFunc},
	},
	object.OBJ_OBJ: {
		"json": {Fn: jsonFunc},
	},
}

// jsonFunc convert value to json representation
func jsonFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	json, err := receiver.JSON()
	if err != nil {
		return nil, err
	}
	return &object.Str{Value: json}, nil
}
