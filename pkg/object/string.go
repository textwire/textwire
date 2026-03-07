package object

import "fmt"

type Str struct {
	Val string
}

func (s *Str) Type() ObjectType {
	return STR_OBJ
}

func (s *Str) String() string {
	return s.Val
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

func (s *Str) Is(t ObjectType) bool {
	return t == s.Type()
}
