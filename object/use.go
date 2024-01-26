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
