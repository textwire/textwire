package object

type Reserve struct {
	Name    string
	Content Object
}

func (r *Reserve) Type() ObjectType {
	return RESERVE_OBJ
}

func (r *Reserve) String() string {
	return r.Content.String()
}
