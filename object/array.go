package object

import (
	"bytes"
	"strings"
)

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARR_OBJ
}

func (a *Array) String() string {
	var out bytes.Buffer

	for _, elem := range a.Elements {
		out.WriteString(elem.String() + ", ")
	}

	if out.Len() > 1 {
		out.Truncate(out.Len() - 2)
	}

	return out.String()
}

func (a *Array) Dump(ident int) string {
	spaces := strings.Repeat(" ", ident)
	ident += 1

	var out bytes.Buffer

	out.WriteString("<span class='textwire-brace'>[</span>\n")

	idx := 0
	last := len(a.Elements) - 1

	for _, elem := range a.Elements {
		out.WriteString(spaces + elem.Dump(ident))

		if idx != last {
			out.WriteString(",\n")
		}

		idx++
	}

	out.WriteString("<span class='textwire-brace'>]</span>\n")

	return out.String()
}

func (a *Array) Val() interface{} {
	var result []interface{}

	for _, elem := range a.Elements {
		result = append(result, elem.Val())
	}

	return result
}

func (a *Array) Is(t ObjectType) bool {
	return t == a.Type()
}
