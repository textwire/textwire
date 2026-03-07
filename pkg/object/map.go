package object

import (
	"fmt"
	"sort"
	"strings"

	"github.com/textwire/textwire/v3/pkg/utils"
)

type Map struct {
	Pairs map[string]Object
}

func NewObj(pairs map[string]Object) *Map {
	if pairs == nil {
		pairs = map[string]Object{}
	}
	return &Map{Pairs: pairs}
}

func (o *Map) Type() ObjectType {
	return MAP_OBJ
}

func (o *Map) String() string {
	if o.Pairs == nil {
		return "{}"
	}

	keys := o.sortedKeys()
	var out strings.Builder
	out.Grow(2)

	out.WriteByte('{')

	for i, k := range keys {
		pair := o.Pairs[k]
		if i > 0 {
			out.WriteString(", ")
		}

		if _, isStr := pair.(*String); isStr {
			out.WriteString(k + `: "` + pair.String() + `"`)
		} else {
			out.WriteString(k + ": " + pair.String())
		}
	}

	out.WriteByte('}')

	return out.String()
}

func (o *Map) JSON() (string, error) {
	if o.Pairs == nil {
		return "{}", nil
	}

	keys := o.sortedKeys()
	var out strings.Builder
	out.Grow(2)

	out.WriteByte('{')

	for i, k := range keys {
		pair := o.Pairs[k]
		if i > 0 {
			out.WriteByte(',')
		}

		jsonVal, err := pair.JSON()
		if err != nil {
			return "", err
		}

		fmt.Fprintf(&out, `"%s":%s`, k, jsonVal)
	}

	out.WriteByte('}')

	return out.String(), nil
}

func (o *Map) Dump(ident int) string {
	if o.Pairs == nil {
		return "{}"
	}

	spaces := strings.Repeat("  ", ident)
	ident += 1

	var out strings.Builder
	out.Grow(4)

	fmt.Fprintf(&out, `<span style="%s">object:%d </span>`, DUMP_META, len(o.Pairs))
	fmt.Fprintf(&out, `<span style="%s">{</span>`, DUMP_BRACE)

	if len(o.Pairs) == 0 {
		spaces = ""
	} else {
		out.WriteByte('\n')
	}

	insideSpaces := strings.Repeat("  ", ident)

	for key, pair := range o.Pairs {
		out.WriteString(insideSpaces)
		fmt.Fprintf(&out, `<span style="%s">"`, DUMP_PROP)
		out.WriteString(key)
		fmt.Fprintf(&out, `"</span>: `)
		out.WriteString(pair.Dump(ident))
		out.WriteString(",\n")
	}

	out.WriteString(spaces)
	fmt.Fprintf(&out, `<span style="%s">}</span>`, DUMP_BRACE)

	return out.String()
}

func (o *Map) Native() any {
	res := map[string]any{}
	for k, v := range o.Pairs {
		res[k] = v.Native()
	}

	return res
}

func (o *Map) Is(t ObjectType) bool {
	return t == o.Type()
}

// ToCamel converts each key in a pair to camel case and returns it
// without mutating it.
func (o Map) ToCamel() map[string]Object {
	res := make(map[string]Object, len(o.Pairs))
	for k, v := range o.Pairs {
		key := utils.ToCamel(k)
		if v.Is(MAP_OBJ) {
			v.(*Map).Pairs = v.(*Map).ToCamel()
		}
		res[key] = v
	}

	return res
}

func (o *Map) sortedKeys() []string {
	keys := make([]string, 0, len(o.Pairs))
	for k := range o.Pairs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
