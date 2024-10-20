package evaluator

import (
	"github.com/textwire/textwire/v2/object"
)

var functions = map[object.ObjectType]map[string]*object.Builtin{
	object.STR_OBJ: {
		"len":   {Fn: strLenFunc},
		"split": {Fn: strSplitFunc},
		"raw":   {Fn: strRawFunc},
		"trim":  {Fn: strTrimFunc},
		"upper": {Fn: strUpperFunc},
		"lower": {Fn: strLowerFunc},
	},
	object.ARR_OBJ: {
		"len":     {Fn: arrayLenFunc},
		"join":    {Fn: arrayJoinFunc},
		"rand":    {Fn: arrayRandFunc},
		"reverse": {Fn: arrayReverseFunc},
		"slice":   {Fn: arraySliceFunc},
	},
	object.FLOAT_OBJ: {
		"int": {Fn: floatIntFunc},
	},
	object.INT_OBJ: {
		"float": {Fn: intFloatFunc},
	},
	object.BOOL_OBJ: {},
}
