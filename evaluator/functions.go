package evaluator

import (
	"github.com/textwire/textwire/object"
)

var functions = map[object.ObjectType]map[string]*object.Builtin{
	object.STR_OBJ: {
		"len":   {Fn: strLenFunc},
		"split": {Fn: strSplitFunc},
		"raw":   {Fn: strRawFunc},
		"trim":  {Fn: strTrimFunc},
	},
	object.ARR_OBJ: {
		"len":  {Fn: arrayLenFunc},
		"join": {Fn: arrayJoinFunc},
	},
	object.FLOAT_OBJ: {
		"int": {Fn: floatIntFunc},
	},
	object.INT_OBJ: {
		"float": {Fn: intFloatFunc},
	},
}
