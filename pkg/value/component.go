package value

type Component struct {
	Name    string
	Content Value
}

func (*Component) Type() ValueType {
	return COMPONENT_VAL
}

func (c *Component) String() string {
	return c.Content.String()
}

func (c *Component) Is(t ValueType) bool {
	return t == c.Type()
}
