package value

type Text struct {
	Val string
}

func (*Text) Type() ValueType {
	return TEXT_VAL
}

func (s *Text) String() string {
	return s.Val
}

func (s *Text) Is(t ValueType) bool {
	return t == s.Type()
}
