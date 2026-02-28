package object

type HTML struct {
	Value string
}

func (h *HTML) Type() ObjectType {
	return HTML_OBJ
}

func (h *HTML) String() string {
	return h.Value
}

func (h *HTML) Dump(ident int) string {
	return ""
}

func (h *HTML) JSON() (string, error) {
	return "", nil
}

func (h *HTML) Val() any {
	return h.Value
}

func (h *HTML) Is(t ObjectType) bool {
	return t == h.Type()
}
