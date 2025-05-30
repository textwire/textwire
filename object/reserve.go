package object

type Reserve struct {
	Name     string
	Content  Object
	Argument Object
}

func (r *Reserve) Type() ObjectType {
	return RESERVE_OBJ
}

func (r *Reserve) String() string {
	if r.Argument != nil {
		return r.Argument.String()
	}

	return r.Content.String()
}

func (r *Reserve) Dump(ident int) string {
	return "reserve stmt"
}

func (r *Reserve) Val() any {
	return r.Content.Val()
}

func (r *Reserve) Is(t ObjectType) bool {
	return t == r.Type()
}
