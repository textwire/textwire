package object

type Break struct{}

func (b *Break) Type() ObjectType {
	return BREAK_OBJ
}

func (b *Break) String() string {
	return ""
}

func (b *Break) Is(t ObjectType) bool {
	return t == b.Type()
}
