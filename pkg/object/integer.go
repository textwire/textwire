package object

import (
	"fmt"
)

type Integer struct {
	Val int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) String() string {
	return fmt.Sprintf("%d", i.Val)
}

func (i *Integer) Dump(ident int) string {
	return fmt.Sprintf(`<span style="%s">%d</span>`, DUMP_NUM, i.Val)
}

func (i *Integer) JSON() (string, error) {
	return fmt.Sprintf("%d", i.Val), nil
}

func (i *Integer) Native() any {
	return i.Val
}

func (i *Integer) Is(t ObjectType) bool {
	return t == i.Type()
}
