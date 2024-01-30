package object

type Html struct {
	Value string
}

func (h *Html) Type() ObjectType {
	return HTML_OBJ
}

func (h *Html) String() string {
	return h.Value
}

func (h *Html) Is(t ObjectType) bool {
	return t == h.Type()
}
