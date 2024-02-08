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
}
