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
	if len(a.Elements) == 0 {
		return ""
	}

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

	fmt.Fprintf(&out, `<span style="%s">array:%d </span>`, DUMP_META, len(a.Elements))
	fmt.Fprintf(&out, `<span style="%s">[</span>`, DUMP_BRACE)

	if len(a.Elements) == 0 {
		spaces = ""
	} else {
		out.WriteByte('\n')
	}

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

	return res + spaces + fmt.Sprintf(`<span style="%s">]</span>`, DUMP_BRACE)
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
