package object

type Reserve struct {
	Name *String
}

func (r *Reserve) Type() ObjectType {
	return RESERVE_OBJ
}

func (r *Reserve) String() string {
	return ""
}
