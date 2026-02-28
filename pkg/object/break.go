package object

type Break struct{}

func (b *Break) Type() ObjectType {
	return BREAK_OBJ
}

func (b *Break) String() string {
	return ""
}

func (b *Break) Dump(ident int) string {
	return ""
}

func (b *Break) JSON() (string, error) {
	return "", nil
}

func (b *Break) Val() any {
	return nil
}

func (b *Break) Is(t ObjectType) bool {
	return t == b.Type()
}
