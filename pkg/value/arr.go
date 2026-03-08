package value

import (
	"bytes"
	"fmt"
	"strings"
)

type Arr struct {
	Elements []Value
}

func (a *Arr) Type() ValueType {
	return ARR_VAL
}

func (a *Arr) String() string {
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

func (a *Arr) Dump(ident int) string {
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

func (a *Arr) JSON() (string, error) {
	var out strings.Builder
	out.Grow(len(a.Elements) + 2)

	out.WriteByte('[')

	for i := range a.Elements {
		if i > 0 {
			out.WriteByte(',')
		}

		val, err := a.Elements[i].JSON()
		if err != nil {
			return "", err
		}

		out.WriteString(val)
	}

	out.WriteByte(']')

	return out.String(), nil
}

func (a *Arr) Native() any {
	var vals []any

	for _, elem := range a.Elements {
		vals = append(vals, elem.Native())
	}

	return vals
}

func (a *Arr) Is(t ValueType) bool {
	return t == a.Type()
}
