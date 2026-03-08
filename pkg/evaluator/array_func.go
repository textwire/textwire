package evaluator

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/value"
)

// arrayLenFunc returns the length of the given array
func arrayLenFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	elems := receiver.(*value.Arr).Elements
	length := len(elems)
	return &value.Int{Val: int64(length)}, nil
}

// arrayJoinFunc joins the elements of the given array with the given separator
func arrayJoinFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	var separator string

	if len(args) == 0 {
		separator = ","
	} else {
		str, ok := args[0].(*value.Str)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, value.ARR_VAL, "join")
			return nil, errors.New(msg)
		}

		separator = str.Val
	}

	elems := receiver.(*value.Arr).Elements

	var out strings.Builder
	out.Grow(len(elems))

	for i := range elems {
		if i > 0 {
			out.WriteString(separator)
		}

		out.WriteString(elems[i].String())
	}

	return &value.Str{Val: out.String()}, nil
}

// arrayRandFunc returns a random element from the given array
func arrayRandFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	elems := receiver.(*value.Arr).Elements
	if len(elems) == 0 {
		return &value.Nil{}, nil
	}

	return elems[0], nil
}

// arrayReverseFunc reverses the elements of the given array
func arrayReverseFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	elems := receiver.(*value.Arr).Elements
	if len(elems) == 0 {
		return receiver, nil
	}

	slices.Reverse(elems)

	return &value.Arr{Elements: elems}, nil
}

// arraySliceFunc returns a slice of the given array
func arraySliceFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	elems := receiver.(*value.Arr).Elements

	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, value.ARR_VAL, "slice")
		return nil, errors.New(msg)
	}

	startFrom, ok := args[0].(*value.Int)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, value.ARR_VAL, "slice")
		return nil, errors.New(msg)
	}

	start := max(int(startFrom.Val), 0)
	start = min(start, len(elems))

	if len(args) == 1 {
		return &value.Arr{Elements: elems[start:]}, nil
	}

	endAt, ok := args[1].(*value.Int)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncSecondArgInt, value.ARR_VAL, "slice")
		return nil, errors.New(msg)
	}

	end := int(endAt.Val)
	if end < 0 || end > len(elems) {
		end = len(elems)
	}

	return &value.Arr{Elements: elems[start:end]}, nil
}

// arrayShuffleFunc shuffles the elements of the given array
func arrayShuffleFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	elems := receiver.(*value.Arr).Elements

	length := len(elems)
	if length == 0 {
		return receiver, nil
	}

	// Create a copy of the elements to shuffle
	shuffled := make([]value.Value, length)
	copy(shuffled, elems)

	// Seed the random number generator to ensure different results
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Perform Fisher-Yates shuffle
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)                               // Pick a random index
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i] // Swap elements
	}

	return &value.Arr{Elements: shuffled}, nil
}

// arrayContainsFunc checks if the given array contains the given element
func arrayContainsFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, value.ARR_VAL, "contains")
		return nil, errors.New(msg)
	}

	elems := receiver.(*value.Arr).Elements
	if len(elems) == 0 {
		return FALSE, nil
	}

	target := args[0]

	for _, el := range elems {
		isObj := el.Type() == value.OBJ_VAL && target.Type() == value.OBJ_VAL
		isArr := el.Type() == value.ARR_VAL && target.Type() == value.ARR_VAL

		if isObj || isArr {
			if reflect.DeepEqual(el, target) {
				return TRUE, nil
			}

			continue
		}

		if el.Native() == target.Native() {
			return TRUE, nil
		}
	}

	return FALSE, nil
}

// arrayAppendFunc appends the given elements to the given array
func arrayAppendFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, value.ARR_VAL, "append")
		return nil, errors.New(msg)
	}

	arr := receiver.(*value.Arr)
	newElems := make(
		[]value.Value,
		len(arr.Elements)+len(args),
	)

	copy(newElems, arr.Elements)

	for i := range args {
		newElems[len(arr.Elements)+i] = args[i]
	}

	return &value.Arr{Elements: newElems}, nil
}

// arrayPrependFunc prepends the given elements to the given array
func arrayPrependFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, value.ARR_VAL, "prepend")
		return nil, errors.New(msg)
	}

	arr := receiver.(*value.Arr)
	newElems := make([]value.Value, len(arr.Elements)+len(args))

	copy(newElems, args)

	for i := range arr.Elements {
		newElems[len(args)+i] = arr.Elements[i]
	}

	return &value.Arr{Elements: newElems}, nil
}
