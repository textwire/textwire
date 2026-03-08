package value

import (
	"fmt"
)

type Int struct {
	Val int64
}

func (i *Int) Type() ValueType {
	return INT_VAL
}

func (i *Int) String() string {
	return fmt.Sprintf("%d", i.Val)
}

func (i *Int) Dump(ident int) string {
	return fmt.Sprintf(`<span style="%s">%d</span>`, DUMP_NUM, i.Val)
}

func (i *Int) JSON() (string, error) {
	return fmt.Sprintf("%d", i.Val), nil
}

func (i *Int) Native() any {
	return i.Val
}

func (i *Int) Is(t ValueType) bool {
	return t == i.Type()
}
