package value

import (
	"bytes"
	"fmt"
	"strings"

	"codeberg.org/textwire/textwire/v6/pkg/fail"
)

type Arr struct {
	Elements []Literal
}

func (*Arr) Type() ValueType {
	return ARR_VAL
}

func (a *Arr) String() string {
	if HasCycle(a) {
		return circularReference
	}
	return a.stringNoCycleCheck()
}

// stringNoCycleCheck is the String() without the HasCycle check.
// It builds the output from the elements, and for each element it skips
// the cycle check for child arrays and objects.
func (a *Arr) stringNoCycleCheck() string {
	return a.joinElements(stringOf)
}

// htmlNoCycleCheck is stringNoCycleCheck with HTML escaping, see escapedString.
func (a *Arr) htmlNoCycleCheck() string {
	return a.joinElements(htmlOf)
}

// joinElements converts each element with elemOf and joins them with commas.
func (a *Arr) joinElements(toString func(Literal) string) string {
	if len(a.Elements) == 0 {
		return ""
	}

	var out bytes.Buffer
	for _, elem := range a.Elements {
		out.WriteString(toString(elem))
		out.WriteString(", ")
	}

	if out.Len() > 1 {
		out.Truncate(out.Len() - 2)
	}

	return out.String()
}

func (a *Arr) Dump(indent int) string {
	if HasCycle(a) {
		return circularReference
	}
	return a.dumpNoCycleCheck(indent)
}

func (a *Arr) dumpNoCycleCheck(indent int) string {
	spaces := strings.Repeat("  ", indent)
	indent += 1

	var out bytes.Buffer

	fmt.Fprintf(&out, `<span style="%s">array:%d </span>`, DUMP_META, len(a.Elements))
	fmt.Fprintf(&out, `<span style="%s">[</span>`, DUMP_BRACE)

	if len(a.Elements) == 0 {
		spaces = ""
	} else {
		out.WriteByte('\n')
	}

	insideSpaces := strings.Repeat("  ", indent)

	for _, elem := range a.Elements {
		out.WriteString(insideSpaces)
		out.WriteString(dumpOf(elem, indent))
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
	if HasCycle(a) {
		return "", fmt.Errorf(fail.ErrCircularValue, a.Type())
	}
	return a.jsonNoCycleCheck()
}

func (a *Arr) jsonNoCycleCheck() (string, error) {
	var out strings.Builder
	out.Grow(len(a.Elements) + 2)

	out.WriteByte('[')

	for i := range a.Elements {
		if i > 0 {
			out.WriteByte(',')
		}

		val, err := jsonOf(a.Elements[i])
		if err != nil {
			return "", err
		}

		out.WriteString(val)
	}

	out.WriteByte(']')

	return out.String(), nil
}

func (a *Arr) Native() any {
	if HasCycle(a) {
		return nil
	}
	return a.nativeNoCycleCheck()
}

func (a *Arr) nativeNoCycleCheck() any {
	var vals []any
	for _, elem := range a.Elements {
		vals = append(vals, nativeOf(elem))
	}
	return vals
}

func (a *Arr) Is(t ValueType) bool {
	return t == a.Type()
}
