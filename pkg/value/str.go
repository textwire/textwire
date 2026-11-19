package value

import (
	"fmt"
	"html"
)

type Str struct {
	Val string

	// raw is true when the string is marked with the raw() function.
	// A raw string is printed without HTML escaping.
	raw bool
}

func (*Str) Type() ValueType {
	return STR_VAL
}

func (s *Str) String() string {
	return s.Val
}

func (s *Str) Dump(ident int) string {
	// Dump output is HTML, so a raw string is escaped here too
	str := fmt.Sprintf(`<span style="%s">%q</span>`, DUMP_STR, html.EscapeString(s.Val))
	if s.raw {
		return fmt.Sprintf(`<span style="%s">raw </span>`, DUMP_META) + str
	}
	return str
}

func (s *Str) JSON() (string, error) {
	return marshalJSONString(s.Val)
}

func (s *Str) Native() any {
	return s.Val
}

func (s *Str) Is(t ValueType) bool {
	return t == s.Type()
}
