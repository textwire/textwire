package value

import "fmt"

type Nil struct{}

func (n *Nil) Type() ValueType {
	return NIL_OBJ
}

func (n *Nil) String() string {
	return ""
}

func (n *Nil) Dump(ident int) string {
	return fmt.Sprintf(`<span style="%s">nil</span>`, DUMP_KEYWORD)
}

func (n *Nil) JSON() (string, error) {
	return "null", nil
}

func (n *Nil) Native() any {
	return nil
}

func (n *Nil) Is(t ValueType) bool {
	return t == n.Type()
}
