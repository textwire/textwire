package object

type Use struct {
	Path   string
	Layout Object
}

func (u *Use) Type() ObjectType {
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

func (u *Use) Val() any {
	if u.Layout == nil {
		panic("Layout field on Use object must not be nil when calling Val()")
	}
	return u.Layout.Val()
}

func (u *Use) Is(t ObjectType) bool {
	return t == u.Type()
}
