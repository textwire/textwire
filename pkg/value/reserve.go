package value

type Reserve struct {
	Name   string
	Insert Value
}

func (r *Reserve) Type() ValueType {
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

func (r *Reserve) Native() any {
	if r.Insert == nil {
		panic("Insert field on Reseve object must not be nil when calling Native()")
	}
	return r.Insert.Native()
}

func (r *Reserve) Is(t ValueType) bool {
	return t == r.Type()
}
