package evaluator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/object"
)

// objCamelFunc converts object keys to camel case recursively
func objCamelFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	obj := receiver.(*object.Obj)
	return &object.Obj{Pairs: obj.ToCamel()}, nil
}

func objGetFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	obj := receiver.(*object.Obj)

	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, object.OBJ_OBJ, "get")
		return nil, errors.New(msg)
	}

	if len(args) > 1 {
		msg := fmt.Sprintf(fail.ErrFuncMaxArgs, object.OBJ_OBJ, "get", 1)
		return nil, errors.New(msg)
	}

	pattern, ok := args[0].(*object.Str)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, object.OBJ_OBJ, "get")
		return nil, errors.New(msg)
	}

	patternStr := pattern.String()

	if result, ok := obj.Pairs[patternStr]; ok {
		return result, nil
	}

	props := strings.Split(patternStr, ".")
	return findObjectKey(props, obj.Pairs), nil
}

func findObjectKey(props []string, pairs map[string]object.Object) object.Object {
	current := pairs

	for i := range props {
		result, ok := current[props[i]]
		if !ok {
			return NIL
		}

		// If this is the last key, return the value
		if i == len(props)-1 {
			return result
		}

		// If not last key, value must be an object to continue
		if result.Type() != object.OBJ_OBJ {
			return NIL
		}

		current = result.(*object.Obj).Pairs
	}

	return NIL
}
