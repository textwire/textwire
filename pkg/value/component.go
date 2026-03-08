package value

type Component struct {
	Name    string
	Content Value
}

func (c *Component) Type() ValueType {
	return COMPONENT_VAL
}

func (c *Component) String() string {
	return c.Content.String()
}

func (c *Component) Dump(ident int) string {
	return ""
}

func (c *Component) JSON() (string, error) {
	return "", nil
}

func (c *Component) Native() any {
	if c.Content == nil {
		panic("Content field on Component object must not be nil when calling Native()")
	}

	return c.Content.Native()
}

func (c *Component) Is(t ValueType) bool {
	return t == c.Type()
}
