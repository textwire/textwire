package object

type Nil struct{}

func (n *Nil) Type() ObjectType {
	return NIL_OBJ
}

func (n *Nil) String() string {
	return ""
}

func (n *Nil) Dump(ident int) string {
	return "<span class='textwire-keyword'>nil</span>"
}

func (n *Nil) Val() any {
	return nil
}

func (n *Nil) Is(t ObjectType) bool {
	return t == n.Type()
}
