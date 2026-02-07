package object

import (
	"bytes"
	"fmt"
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
	spaces := strings.Repeat("  ", ident)
	ident += 1

	var out bytes.Buffer

	fmt.Fprintf(&out, "<span class='textwire-meta'>array:%d </span>", len(a.Elements))
	out.WriteString("<span class='textwire-brace'>[</span>\n")

	insideSpaces := strings.Repeat("  ", ident)

	for _, elem := range a.Elements {
		out.WriteString(insideSpaces)
		out.WriteString(elem.Dump(ident))
		out.WriteString(",\n")
	}

	res := out.String()

	// take last 8 characters of the string
	lastChar := res[len(res)-8:]

	if lastChar == "}</span>" {
		res += "\n"
	}

	return res + spaces + "<span class='textwire-brace'>]</span>"
}

func (a *Array) Val() any {
	var result []any

	for _, elem := range a.Elements {
		result = append(result, elem.Val())
	}

	return result
}

func (a *Array) Is(t ObjectType) bool {
	return t == a.Type()
}
