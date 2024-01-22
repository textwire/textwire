package object

type Bool struct {
	Value bool
}

func (b *Bool) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Bool) String() string {
	if b.Value {
		return "1"
	}

	return "0"
}
