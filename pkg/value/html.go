package value

type HTML struct {
	Val string
}

func (h *HTML) Type() ValueType {
	return HTML_OBJ
}

func (h *HTML) String() string {
	return h.Val
}

func (h *HTML) Dump(ident int) string {
	return ""
}

func (h *HTML) JSON() (string, error) {
	return "", nil
}

func (h *HTML) Native() any {
	return h.Val
}

func (h *HTML) Is(t ValueType) bool {
	return t == h.Type()
}
