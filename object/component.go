package object

type Component struct {
	Name    string
	Content Object
}

func (c *Component) Type() ObjectType {
	return COMPONENT_OBJ
}

func (c *Component) String() string {
	return c.Content.String()
}

func (c *Component) Val() interface{} {
	return c.Content.Val()
}

func (c *Component) Is(t ObjectType) bool {
	return t == c.Type()
}
