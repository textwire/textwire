package object

import (
	"bytes"
	"fmt"
	"strings"
)

type Obj struct {
	Pairs map[string]Object
}

func NewObj(pairs map[string]Object) *Obj {
	if pairs == nil {
		pairs = map[string]Object{}
	}
	return &Obj{Pairs: pairs}
}

func (o *Obj) Type() ObjectType {
	return OBJ_OBJ
}

func (o *Obj) String() string {
	if o.Pairs == nil {
		return "{}"
	}

	keys := o.sortedKeys()

	// Estimate size: { + keys + values + commas + quotes
	estimatedSize := 2 + len(keys)*10
	for _, k := range keys {
		estimatedSize += len(k) + len(o.Pairs[k].String())
	}

	var out strings.Builder

	out.Grow(estimatedSize)
	out.WriteString("{")

	for i, k := range keys {
		pair := o.Pairs[k]
		if i > 0 {
			out.WriteString(", ")
		}

		if _, isStr := pair.(*Str); isStr {
			out.WriteString(k + `: "` + pair.String() + `"`)
		} else {
			out.WriteString(k + ": " + pair.String())
		}
	}

	out.WriteString("}")
	return out.String()
}

func (o *Obj) Dump(ident int) string {
	spaces := strings.Repeat("  ", ident)
	ident += 1

	var out bytes.Buffer

	fmt.Fprintf(&out, "<span class='textwire-meta'>object:%d </span>", len(o.Pairs))
	out.WriteString("<span class='textwire-brace'>{</span>\n")

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

func (o *Obj) Val() any {
	result := map[string]any{}

	for k, v := range o.Pairs {
		result[k] = v.Val()
	}

	return result
}

func (o *Obj) Is(t ObjectType) bool {
	return t == o.Type()
}
