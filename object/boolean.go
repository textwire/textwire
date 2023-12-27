package object

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) String() string {
	if b.Value {
		return "1"
	}

	return "0"
}
