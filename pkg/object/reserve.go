package object

type Reserve struct {
	Name   string
	Insert Object
}

func (r *Reserve) Type() ObjectType {
	return RESERVE_OBJ
}

func (r *Reserve) String() string {
	if r.Insert == nil {
		panic("Insert field on Reseve object must not be nil when calling String()")
	}
	return r.Insert.String()
}

func (r *Reserve) Dump(ident int) string {
	return ""
}

func (r *Reserve) JSON() (string, error) {
	return "", nil
}

func (r *Reserve) Val() any {
	if r.Insert == nil {
		panic("Insert field on Reseve object must not be nil when calling Val()")
	}
	return r.Insert.Val()
}

func (r *Reserve) Is(t ObjectType) bool {
	return t == r.Type()
}
