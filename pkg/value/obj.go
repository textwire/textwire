package value

import (
	"fmt"
	"html"
	"sort"
	"strings"

	"codeberg.org/textwire/textwire/v6/pkg/fail"
	"codeberg.org/textwire/textwire/v6/pkg/utils"
)

type Obj struct {
	Pairs map[string]Literal
}

func NewObj(pairs map[string]Literal) *Obj {
	if pairs == nil {
		pairs = map[string]Literal{}
	}
	return &Obj{Pairs: pairs}
}

func (*Obj) Type() ValueType {
	return OBJ_VAL
}

func (o *Obj) String() string {
	if HasCycle(o) {
		return circularReference
	}
	return o.stringNoCycleCheck()
}

// stringNoCycleCheck is the String() without the HasCycle check.
// It builds the output from the elements, and for each element it skips
// the cycle check for child arrays and objects.
func (o *Obj) stringNoCycleCheck() string {
	return o.joinPairs(writeStringPair)
}

// htmlNoCycleCheck is stringNoCycleCheck with HTML escaping, see escapedString.
func (o *Obj) htmlNoCycleCheck() string {
	return o.joinPairs(writeHTMLPair)
}

// joinPairs writes each key and value with writePair, sorted by key.
func (o *Obj) joinPairs(writePair func(out *strings.Builder, key string, val Literal)) string {
	if o.Pairs == nil {
		return "{}"
	}

	keys := o.sortedKeys()
	var out strings.Builder
	out.Grow(2)

	out.WriteByte('{')

	for i, key := range keys {
		if i > 0 {
			out.WriteString(", ")
		}
		writePair(&out, key, o.Pairs[key])
	}

	out.WriteByte('}')

	return out.String()
}

// writeStringPair writes `key: "string"` for strings and `key: value` for
// other values.
func writeStringPair(out *strings.Builder, key string, val Literal) {
	out.WriteString(key)
	out.WriteString(": ")

	if str, isStr := val.(*Str); isStr {
		out.WriteString(`"`)
		out.WriteString(str.Val)
		out.WriteString(`"`)
		return
	}

	out.WriteString(stringOf(val))
}

// writeHTMLPair is writeStringPair with HTML escaping. The key is always
// escaped. A raw string is written with its quotes as is, other strings get
// their quotes escaped as &#34;, like html.EscapeString does.
func writeHTMLPair(out *strings.Builder, key string, val Literal) {
	out.WriteString(html.EscapeString(key))
	out.WriteString(": ")

	str, isStr := val.(*Str)

	switch {
	case isStr && str.raw:
		out.WriteString(`"`)
		out.WriteString(str.Val)
		out.WriteString(`"`)
	case isStr:
		out.WriteString("&#34;")
		out.WriteString(html.EscapeString(str.Val))
		out.WriteString("&#34;")
	default:
		out.WriteString(htmlOf(val))
	}
}

func (o *Obj) JSON() (string, error) {
	if HasCycle(o) {
		return "", fmt.Errorf(fail.ErrCircularValue, o.Type())
	}
	return o.jsonNoCycleCheck()
}

func (o *Obj) jsonNoCycleCheck() (string, error) {
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

		encodedKey, err := marshalJSONString(k)
		if err != nil {
			return "", err
		}

		jsonVal, err := jsonOf(pair)
		if err != nil {
			return "", err
		}

		out.WriteString(encodedKey)
		out.WriteByte(':')
		out.WriteString(jsonVal)
	}

	out.WriteByte('}')

	return out.String(), nil
}

func (o *Obj) Dump(ident int) string {
	if HasCycle(o) {
		return circularReference
	}
	return o.dumpNoCycleCheck(ident)
}

func (o *Obj) dumpNoCycleCheck(ident int) string {
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
		out.WriteString(html.EscapeString(key))
		fmt.Fprintf(&out, `"</span>: `)
		out.WriteString(dumpOf(pair, ident))
		out.WriteString(",\n")
	}

	out.WriteString(spaces)
	fmt.Fprintf(&out, `<span style="%s">}</span>`, DUMP_BRACE)

	return out.String()
}

func (o *Obj) Native() any {
	if HasCycle(o) {
		return nil
	}
	return o.nativeNoCycleCheck()
}

func (o *Obj) nativeNoCycleCheck() any {
	res := map[string]any{}
	for k, v := range o.Pairs {
		res[k] = nativeOf(v)
	}
	return res
}

func (o *Obj) Is(t ValueType) bool {
	return t == o.Type()
}

// ToCamel converts each key in a pair to camel case and returns it
// without mutating it.
func (o Obj) ToCamel() map[string]Literal {
	res := make(map[string]Literal, len(o.Pairs))
	for k, v := range o.Pairs {
		key := utils.ToCamel(k)
		if v.Is(OBJ_VAL) {
			v = &Obj{Pairs: v.(*Obj).ToCamel()}
		}
		res[key] = v
	}

	return res
}

func (o *Obj) sortedKeys() []string {
	keys := make([]string, 0, len(o.Pairs))
	for k := range o.Pairs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
