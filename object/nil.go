package object

type Nil struct{}

func (n *Nil) Type() ObjectType {
	return NIL_OBJ
}

func (n *Nil) String() string {
	return ""
}

func (n *Nil) Is(t ObjectType) bool {
	return t == n.Type()
}
