package evaluator

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
)

// arrayLenFunc returns the length of the given array
func arrayLenFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	length := len(receiver.(*object.Array).Elements)
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
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, "join")
			return nil, errors.New(msg)
		}

		separator = str.Value
	}

	elements := receiver.(*object.Array).Elements

	var result string

	for i, el := range elements {
		if i > 0 {
			result += separator
		}
		result += el.String()
	}

	return &object.Str{Value: result}, nil
}

// arrayRandFunc returns a random element from the given array
func arrayRandFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	elements := receiver.(*object.Array).Elements
	length := len(elements)

	if length == 0 {
		return &object.Nil{}, nil
	}

	return elements[0], nil
}

// arrayReverseFunc reverses the elements of the given array
func arrayReverseFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	elements := receiver.(*object.Array).Elements
	length := len(elements)

	if length == 0 {
		return receiver, nil
	}

	reversed := make([]object.Object, length)

	for i, el := range elements {
		reversed[length-i-1] = el
	}

	return &object.Array{Elements: reversed}, nil
}

// arraySliceFunc returns a slice of the given array
func arraySliceFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	arr := receiver.(*object.Array)

	argsLen := len(args)
	elemsLen := len(arr.Elements)

	if argsLen < 1 {
		msg := fmt.Sprintf(fail.ErrFuncRequiresOneArg, "slice")
		return nil, errors.New(msg)
	}

	startFrom, ok := args[0].(*object.Int)

	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, "slice")
		return nil, errors.New(msg)
	}

	start := int(startFrom.Value)

	if start < 0 {
		start = 0
	}

	if start > elemsLen {
		start = elemsLen
	}

	if argsLen == 1 {
		return &object.Array{Elements: arr.Elements[start:]}, nil
	}

	endAt, ok := args[1].(*object.Int)

	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncSecondArgInt, "slice")
		return nil, errors.New(msg)
	}

	end := int(endAt.Value)

	if end < 0 || end > elemsLen {
		end = elemsLen
	}

	return &object.Array{Elements: arr.Elements[start:end]}, nil
}

// arrayShuffleFunc shuffles the elements of the given array
func arrayShuffleFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
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
