package object

import (
	"fmt"
	"sort"
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

	var out strings.Builder
	out.Grow(2)

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

	var out strings.Builder
	out.Grow(4)

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
	res := map[string]any{}
	for k, v := range o.Pairs {
		res[k] = v.Val()
	}

	return res
}

func (o *Obj) Is(t ObjectType) bool {
	return t == o.Type()
}

func (o *Obj) sortedKeys() []string {
	keys := make([]string, 0, len(o.Pairs))
	for k := range o.Pairs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
