package object

import "bytes"

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (a *Array) String() string {
	var result bytes.Buffer

	for _, elem := range a.Elements {
		result.WriteString(elem.String() + ", ")
	}

	if result.Len() > 1 {
		result.Truncate(result.Len() - 2)
	}

	return result.String()
}
