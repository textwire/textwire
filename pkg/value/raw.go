package value

// Raw returns v with every string in it marked as raw, so it's printed
// without HTML escaping. Arrays and objects are traversed recursively and
// copied, so v itself is never changed. Other literal types are returned as is.
func Raw(v Literal) Literal {
	return rawCopy(v, map[Literal]Literal{})
}

// rawCopy returns v with its strings marked as raw. The copies map stores
// each array and object that is already copied, with the original as the key
// and its copy as the value.
func rawCopy(v Literal, copies map[Literal]Literal) Literal {
	switch v := v.(type) {
	case *Str:
		return rawCopyStr(v)
	case *Arr:
		return rawCopyArr(v, copies)
	case *Obj:
		return rawCopyObj(v, copies)
	default:
		return v
	}
}

// rawCopyStr returns s if it is already raw, otherwise a raw copy of it.
func rawCopyStr(s *Str) *Str {
	if s.raw {
		return s
	}
	return &Str{Val: s.Val, raw: true}
}

// rawCopyArr returns a deep copy of arr with its strings marked as raw.
// If arr is already in copies, its stored copy is returned. The new copy is
// stored before its elements are copied, so an array that contains itself
// resolves to its copy instead of looping forever.
func rawCopyArr(arr *Arr, copies map[Literal]Literal) *Arr {
	if stored, ok := copies[arr]; ok {
		return stored.(*Arr)
	}

	arrayCopy := &Arr{Elements: make([]Literal, len(arr.Elements))}
	copies[arr] = arrayCopy

	for i, elem := range arr.Elements {
		arrayCopy.Elements[i] = rawCopy(elem, copies)
	}

	return arrayCopy
}

// rawCopyObj is rawCopyArray for objects.
func rawCopyObj(obj *Obj, copies map[Literal]Literal) *Obj {
	if stored, ok := copies[obj]; ok {
		return stored.(*Obj)
	}

	objectCopy := NewObj(make(map[string]Literal, len(obj.Pairs)))
	copies[obj] = objectCopy

	for key, val := range obj.Pairs {
		objectCopy.Pairs[key] = rawCopy(val, copies)
	}

	return objectCopy
}
