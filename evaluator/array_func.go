package evaluator

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
)

// arrayLenFunc returns the length of the given array
func arrayLenFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	elems := receiver.(*object.Array).Elements
	length := len(elems)
	return &object.Int{Value: int64(length)}, nil
}

// arrayJoinFunc joins the elements of the given array with the given separator
func arrayJoinFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	var separator string

	if len(args) == 0 {
		separator = ","
	} else {
		str, ok := args[0].(*object.Str)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, "join", object.ARR_OBJ)
			return nil, errors.New(msg)
		}

		separator = str.Value
	}

	elems := receiver.(*object.Array).Elements

	var result string

	for i, el := range elems {
		if i > 0 {
			result += separator
		}
		result += el.String()
	}

	return &object.Str{Value: result}, nil
}

// arrayRandFunc returns a random element from the given array
func arrayRandFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	elems := receiver.(*object.Array).Elements
	if len(elems) == 0 {
		return &object.Nil{}, nil
	}

	return elems[0], nil
}

// arrayReverseFunc reverses the elements of the given array
func arrayReverseFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	elems := receiver.(*object.Array).Elements
	if len(elems) == 0 {
		return receiver, nil
	}

	reversed := make([]object.Object, len(elems))

	for i, el := range elems {
		reversed[len(elems)-i-1] = el
	}

	return &object.Array{Elements: reversed}, nil
}

// arraySliceFunc returns a slice of the given array
func arraySliceFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	elems := receiver.(*object.Array).Elements

	argsLen := len(args)
	elemsLen := len(elems)
	if argsLen == 0 {
		msg := fmt.Sprintf(fail.ErrFuncRequiresOneArg, "slice", object.ARR_OBJ)
		return nil, errors.New(msg)
	}

	startFrom, ok := args[0].(*object.Int)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, "slice", object.ARR_OBJ)
		return nil, errors.New(msg)
	}

	start := max(int(startFrom.Value), 0)
	start = min(start, elemsLen)

	if argsLen == 1 {
		return &object.Array{Elements: elems[start:]}, nil
	}

	endAt, ok := args[1].(*object.Int)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncSecondArgInt, "slice", object.ARR_OBJ)
		return nil, errors.New(msg)
	}

	end := int(endAt.Value)
	if end < 0 || end > elemsLen {
		end = elemsLen
	}

	return &object.Array{Elements: elems[start:end]}, nil
}

// arrayShuffleFunc shuffles the elements of the given array
func arrayShuffleFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	elems := receiver.(*object.Array).Elements

	length := len(elems)
	if length == 0 {
		return receiver, nil
	}

	// Create a copy of the elements to shuffle
	shuffled := make([]object.Object, length)
	copy(shuffled, elems)

	// Seed the random number generator to ensure different results
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Perform Fisher-Yates shuffle
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)                               // Pick a random index
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i] // Swap elements
	}

	return &object.Array{Elements: shuffled}, nil
}

// arrayContainsFunc checks if the given array contains the given element
func arrayContainsFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncRequiresOneArg, "contains", object.ARR_OBJ)
		return nil, errors.New(msg)
	}

	elems := receiver.(*object.Array).Elements
	if len(elems) == 0 {
		return FALSE, nil
	}

	target := args[0]

	for _, el := range elems {
		isObj := el.Type() == object.OBJ_OBJ && target.Type() == object.OBJ_OBJ
		isArr := el.Type() == object.ARR_OBJ && target.Type() == object.ARR_OBJ

		if isObj || isArr {
			if reflect.DeepEqual(el, target) {
				return TRUE, nil
			}

			continue
		}

		if el.Val() == target.Val() {
			return TRUE, nil
		}
	}

	return FALSE, nil
}

// arrayAppendFunc appends the given elements to the given array
func arrayAppendFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncRequiresOneArg, "append", object.ARR_OBJ)
		return nil, errors.New(msg)
	}

	elems := receiver.(*object.Array).Elements
	newElems := make([]object.Object, len(elems)+len(args))

	copy(newElems, elems)

	for i, arg := range args {
		newElems[len(elems)+i] = arg
	}

	return &object.Array{Elements: newElems}, nil
}

// arrayPrependFunc prepends the given elements to the given array
func arrayPrependFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	argsLen := len(args)
	if argsLen == 0 {
		msg := fmt.Sprintf(fail.ErrFuncRequiresOneArg, "prepend", object.ARR_OBJ)
		return nil, errors.New(msg)
	}

	elems := receiver.(*object.Array).Elements
	newElems := make([]object.Object, len(elems)+argsLen)

	copy(newElems, args)

	for i, el := range elems {
		newElems[argsLen+i] = el
	}

	return &object.Array{Elements: newElems}, nil
}
