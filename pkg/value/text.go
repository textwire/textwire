package value

type Text struct {
	Val string
}

func (h *Text) Type() ValueType {
	return TEXT_VAL
}

func (h *Text) String() string {
	return h.Val
}

func (h *Text) Dump(ident int) string {
	return ""
}

func (h *Text) JSON() (string, error) {
	return "", nil
}

func (h *Text) Native() any {
	return h.Val
}

func (h *Text) Is(t ValueType) bool {
	return t == h.Type()
}
