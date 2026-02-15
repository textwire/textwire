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

func (c *Component) Dump(ident int) string {
	if c.Content == nil {
		panic("Content field on Component object must not be nil when calling String()")
	}
	return c.Content.Dump(ident)
}

func (c *Component) Val() any {
	if c.Content == nil {
		panic("Content field on Component object must not be nil when calling Val()")
	}

	return c.Content.Val()
}

func (c *Component) Is(t ObjectType) bool {
	return t == c.Type()
}
