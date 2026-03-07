package object

import "fmt"

type String struct {
	Val string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) String() string {
	return s.Val
}

func (s *String) Dump(ident int) string {
	return fmt.Sprintf(`<span style="%s">%q</span>`, DUMP_STR, s.Val)
}

func (s *String) JSON() (string, error) {
	return fmt.Sprintf("%q", s.Val), nil
}

func (s *String) Native() any {
	return s.Val
}

func (s *String) Is(t ObjectType) bool {
	return t == s.Type()
}
