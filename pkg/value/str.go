package value

import (
	"fmt"
	"html"
)

type Str struct {
	Val   string
	IsRaw bool
}

func (*Str) Type() ValueType {
	return STR_VAL
}

func (s *Str) String() string {
	if s.IsRaw {
		return s.Val
	}
	return html.EscapeString(s.Val)
}

func (s *Str) Dump(ident int) string {
	return fmt.Sprintf(`<span style="%s">%q</span>`, DUMP_STR, s.Val)
}

func (s *Str) JSON() (string, error) {
	return fmt.Sprintf("%q", s.Val), nil
}

func (s *Str) Native() any {
	return s.Val
}

func (s *Str) Is(t ValueType) bool {
	return t == s.Type()
}
