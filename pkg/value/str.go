package value

import (
	"fmt"
)

type Str struct {
	Val string
}

func (*Str) Type() ValueType {
	return STR_VAL
}

func (s *Str) String() string {
	return s.Val
}

func (s *Str) Dump(ident int) string {
	return fmt.Sprintf(`<span style="%s">%q</span>`, DUMP_STR, s.Val)
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
