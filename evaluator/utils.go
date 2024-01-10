package evaluator

import "github.com/textwire/textwire/object"

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Int:
		return obj.Value != 0
	case *object.Float:
		return obj.Value != 0.0
	case *object.String:
		return obj.Value != ""
	}

	return true
}
