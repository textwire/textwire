package evaluator

import (
	"github.com/textwire/textwire/v2/object"
)

var functions = map[object.ObjectType]map[string]*object.Builtin{
	object.STR_OBJ: {
		"len":        {Fn: strLenFunc},
		"split":      {Fn: strSplitFunc},
		"raw":        {Fn: strRawFunc},
		"trim":       {Fn: strTrimFunc},
		"upper":      {Fn: strUpperFunc},
		"lower":      {Fn: strLowerFunc},
		"capitalize": {Fn: strCapitalizeFunc},
		"reverse":    {Fn: strReverseFunc},
		"contains":   {Fn: strContainsFunc},
		"truncate":   {Fn: strTruncateFunc},
	},
	object.ARR_OBJ: {
		"len":     {Fn: arrayLenFunc},
		"join":    {Fn: arrayJoinFunc},
		"rand":    {Fn: arrayRandFunc},
		"reverse": {Fn: arrayReverseFunc},
		"slice":   {Fn: arraySliceFunc},
		"shuffle": {Fn: arrayShuffleFunc},
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
		"float": {Fn: intFloatFunc},
		"abs":   {Fn: intAbsFunc},
		"str":   {Fn: intStrFunc},
	},
	object.BOOL_OBJ: {
		"binary": {Fn: boolBinaryFunc},
	},
}
