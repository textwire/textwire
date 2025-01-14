package object

type Use struct {
	Path    string
	Content Object
}

func (u *Use) Type() ObjectType {
	return USE_OBJ
}

func (u *Use) String() string {
	return u.Content.String()
}

func (u *Use) Dump(ident int) string {
	return "use"
}

func (u *Use) Val() interface{} {
	return u.Content.Val()
}

func (u *Use) Is(t ObjectType) bool {
	return t == u.Type()
}
