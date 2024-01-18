package object

type Layout struct {
	Path    string
	Content Object
}

func (l *Layout) Type() ObjectType {
	return LAYOUT_OBJ
}

func (l *Layout) String() string {
	return l.Content.String()
}
