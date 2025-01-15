package object

import (
	"bytes"
	"strings"
)

type Obj struct {
	Pairs map[string]Object
}

func (o *Obj) Type() ObjectType {
	return OBJ_OBJ
}

func (o *Obj) String() string {
	var out bytes.Buffer

	out.WriteString("{")

	idx := 0
	last := len(o.Pairs) - 1

	for key, pair := range o.Pairs {
		out.WriteString(key + ": " + pair.String())

		if idx != last {
			out.WriteString(", ")
		}

		idx++
	}

	out.WriteString("}")

	return out.String()
}

func (o *Obj) Dump(ident int) string {
	spaces := strings.Repeat("  ", ident)
	ident += 1

	var out bytes.Buffer

	out.WriteString("\n" + spaces)
	out.WriteString("<span class='textwire-brace'>{</span>\n")

	ident += 1
	insideSpaces := strings.Repeat("  ", ident)

	for key, pair := range o.Pairs {
		out.WriteString(insideSpaces)
		out.WriteString(`<span class="textwire-prop">"` + key + `"</span>`)
		out.WriteString(": ")
		out.WriteString(pair.Dump(ident))
		out.WriteString(",\n")
	}

	out.WriteString(spaces + "<span class='textwire-brace'>}</span>")

	return out.String()
}

func (o *Obj) Val() interface{} {
	result := make(map[string]interface{})

	for k, v := range o.Pairs {
		result[k] = v.Val()
	}

	return result
}

func (o *Obj) Is(t ObjectType) bool {
	return t == o.Type()
}
