package object

type Reserve struct {
	Name string
}

func (r *Reserve) Type() ObjectType {
	return RESERVE_OBJ
}

func (r *Reserve) String() string {
	return ""
}
