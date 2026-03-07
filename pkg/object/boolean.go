package object

import "fmt"

type Bool struct {
	Val bool
}

func (b *Bool) Type() ObjectType {
	return BOOL_OBJ
}

func (b *Bool) String() string {
	if b.Val {
		return "1"
	}
	return "0"
}

func (b *Bool) Dump(ident int) string {
	if b.Val {
		return fmt.Sprintf(`<span style="%s">true</span>`, DUMP_KEYWORD)
	}
	return fmt.Sprintf(`<span style="%s">false</span>`, DUMP_KEYWORD)
}

func (b *Bool) JSON() (string, error) {
	return fmt.Sprintf("%t", b.Val), nil
}

func (b *Bool) Native() any {
	return b.Val
}

func (b *Bool) Is(t ObjectType) bool {
	return t == b.Type()
}
