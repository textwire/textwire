package value

// circularReference is rendered by String and Dump when a value contains a
// cycle. JSON returns an error instead.
const circularReference = "[circular reference]"

const (
	unvisited visitState = iota
	visiting
	done
)

type visitState uint8

// References reports whether target is reachable from root through arrays and
// objects. It walks root's arrays and objects depth first and returns true
// if it reaches target. It has no depth limit.
func References(root, target Literal) bool {
	if root == nil || target == nil {
		return false
	}

	seen := map[Literal]bool{}
	stack := []Literal{root}

	for len(stack) > 0 {
		// The last item is the top of the stack
		current := stack[len(stack)-1]

		// Removes current item by shortening the slice by one
		stack = stack[:len(stack)-1]

		if current == target {
			return true
		}

		kids, ok := children(current)
		if !ok || seen[current] {
			continue
		}

		seen[current] = true
		stack = append(stack, kids...)
	}

	return false
}

// HasCycle reports whether root contains itself, directly or through nested
// arrays and objects. For example, arr[0] = arr makes arr contain itself.
func HasCycle(root Literal) bool {
	return cycleFrom(root, map[Literal]visitState{})
}

// cycleFrom walks depth first. An array or object is marked visiting while its
// elements are walked, so reaching it again means a cycle. Once marked done, it
// can't be part of a cycle, so it's skipped.
func cycleFrom(v Literal, state map[Literal]visitState) bool {
	kids, ok := children(v)
	if !ok {
		return false
	}

	switch state[v] {
	case visiting:
		return true
	case done:
		return false
	}

	state[v] = visiting

	for _, kid := range kids {
		if cycleFrom(kid, state) {
			return true
		}
	}

	state[v] = done

	return false
}

// children returns the elements of an array or the values of an object. Object
// keys are plain strings, so they aren't included. The boolean is false when v
// isn't a non-nil array or object.
func children(v Literal) ([]Literal, bool) {
	switch v := v.(type) {
	case *Arr:
		if v == nil {
			return nil, false
		}
		return v.Elements, true

	case *Obj:
		if v == nil {
			return nil, false
		}
		return objectValues(v), true

	default:
		return nil, false
	}
}

func objectValues(o *Obj) []Literal {
	values := make([]Literal, 0, len(o.Pairs))
	for _, value := range o.Pairs {
		values = append(values, value)
	}
	return values
}
