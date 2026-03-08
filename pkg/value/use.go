package value

type Use struct {
	Path   string
	Layout Value
}

func (u *Use) Type() ValueType {
	return USE_OBJ
}

func (u *Use) String() string {
	if u.Layout == nil {
		panic("Layout field on Use object must not be nil when calling String()")
	}
	return u.Layout.String()
}

func (u *Use) Dump(ident int) string {
	return ""
}

func (u *Use) JSON() (string, error) {
	return "", nil
}

func (u *Use) Native() any {
	if u.Layout == nil {
		panic("Layout field on Use object must not be nil when calling Native()")
	}
	return u.Layout.Native()
}

func (u *Use) Is(t ValueType) bool {
	return t == u.Type()
}
