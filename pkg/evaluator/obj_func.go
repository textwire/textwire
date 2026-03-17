package evaluator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/textwire/textwire/v4/pkg/fail"
	"github.com/textwire/textwire/v4/pkg/value"
)

// objCamelFunc converts object keys to camel case recursively
func objCamelFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	obj := receiver.(*value.Obj)
	return &value.Obj{Pairs: obj.ToCamel()}, nil
}

func objGetFunc(receiver value.Literal, args ...value.Literal) (value.Literal, error) {
	obj := receiver.(*value.Obj)

	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, value.OBJ_VAL, "get")
		return nil, errors.New(msg)
	}

	if len(args) > 1 {
		msg := fmt.Sprintf(fail.ErrFuncMaxArgs, value.OBJ_VAL, "get", 1)
		return nil, errors.New(msg)
	}

	pattern, ok := args[0].(*value.Str)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, value.OBJ_VAL, "get")
		return nil, errors.New(msg)
	}

	patternStr := pattern.String()

	if result, ok := obj.Pairs[patternStr]; ok {
		return result, nil
	}

	props := strings.Split(patternStr, ".")
	return findObjKey(props, obj.Pairs), nil
}

func findObjKey(props []string, pairs map[string]value.Literal) value.Literal {
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
		if result.Type() != value.OBJ_VAL {
			return NIL
		}

		current = result.(*value.Obj).Pairs
	}

	return NIL
}
