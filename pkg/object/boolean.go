package object

import "fmt"

type Boolean struct {
	Val bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) String() string {
	if b.Val {
		return "1"
	}
	return "0"
}

func (b *Boolean) Dump(ident int) string {
	if b.Val {
		return fmt.Sprintf(`<span style="%s">true</span>`, DUMP_KEYWORD)
	}
	return fmt.Sprintf(`<span style="%s">false</span>`, DUMP_KEYWORD)
}

func (b *Boolean) JSON() (string, error) {
	return fmt.Sprintf("%t", b.Val), nil
}

func (b *Boolean) Native() any {
	return b.Val
}

func (b *Boolean) Is(t ObjectType) bool {
	return t == b.Type()
}
