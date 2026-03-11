package value

type Reserve struct {
	Name   string
	Insert Value
}

func (*Reserve) Type() ValueType {
	return RESERVE_VAL
}

func (r *Reserve) String() string {
	if r.Insert == nil {
		panic("Insert field on Reseve must not be nil when calling String()")
	}
	return r.Insert.String()
}

func (r *Reserve) Is(t ValueType) bool {
	return t == r.Type()
}
