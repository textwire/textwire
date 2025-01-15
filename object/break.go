package object

type Break struct{}

func (b *Break) Type() ObjectType {
	return BREAK_OBJ
}

func (b *Break) String() string {
	return ""
}

func (b *Break) Dump(ident int) string {
	return "break stmt"
}

func (b *Break) Val() interface{} {
	return nil
}

func (b *Break) Is(t ObjectType) bool {
	return t == b.Type()
}
