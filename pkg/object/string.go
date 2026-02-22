package object

import "fmt"

type Str struct {
	Value string
}

func (s *Str) Type() ObjectType {
	return STR_OBJ
}

func (s *Str) String() string {
	return s.Value
}

func (s *Str) Dump(ident int) string {
	return fmt.Sprintf(`<span style="%s">%q</span>`, DUMP_STR, s.Value)
}

func (s *Str) Val() any {
	return s.Value
}

func (s *Str) Is(t ObjectType) bool {
	return t == s.Type()
}
