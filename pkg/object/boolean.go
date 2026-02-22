package object

import "fmt"

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
		return fmt.Sprintf(`<span style="%s">true</span>`, DUMP_KEYWORD)
	}
	return fmt.Sprintf(`<span style="%s">false</span>`, DUMP_KEYWORD)
}

func (b *Bool) Val() any {
	return b.Value
}

func (b *Bool) Is(t ObjectType) bool {
	return t == b.Type()
}
