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

func (b *Bool) Dump(ident int) string {
	if b.Value {
		return "<span class='textwire-keyword'>true</span>"
	}

	return "<span class='textwire-keyword'>false</span>"
}

func (b *Bool) Val() interface{} {
	return b.Value
}

func (b *Bool) Is(t ObjectType) bool {
	return t == b.Type()
}
