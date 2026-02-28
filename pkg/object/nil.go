package object

import "fmt"

type Nil struct{}

func (n *Nil) Type() ObjectType {
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

func (n *Nil) Val() any {
	return nil
}

func (n *Nil) Is(t ObjectType) bool {
	return t == n.Type()
}
