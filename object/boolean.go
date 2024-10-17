package object

type Bool struct {
	Value bool
}

func (b *Bool) Type() ObjectType {
	return BOOL_OBJ
}

func (b *Bool) String() string {
	if b.Value {
		return "1"
	}

	return "0"
}

func (b *Bool) Val() interface{} {
	return b.Value
}

func (b *Bool) Is(t ObjectType) bool {
	return t == b.Type()
}
