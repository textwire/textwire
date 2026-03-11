package value

type Break struct{}

func (*Break) Type() ValueType {
	return BREAK_VAL
}

func (*Break) String() string {
	return ""
}

func (b *Break) Is(t ValueType) bool {
	return t == b.Type()
}
