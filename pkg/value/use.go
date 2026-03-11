package value

type Use struct {
	Path   string
	Layout Value
}

func (*Use) Type() ValueType {
	return USE_VAL
}

func (u *Use) String() string {
	if u.Layout == nil {
		panic("Layout field on Use must not be nil when calling String()")
	}
	return u.Layout.String()
}

func (u *Use) Is(t ValueType) bool {
	return t == u.Type()
}
