package object

type Str struct {
	Value string
}

func (s *Str) Type() ObjectType {
	return STR_OBJ
}

func (s *Str) String() string {
	return s.Value
}

func (s *Str) Is(t ObjectType) bool {
	return t == s.Type()
}
