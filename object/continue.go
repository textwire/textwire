package object

type Continue struct{}

func (c *Continue) Type() ObjectType {
	return CONTINUE_OBJ
}

func (c *Continue) String() string {
	return ""
}

func (c *Continue) Val() interface{} {
	return nil
}

func (c *Continue) Is(t ObjectType) bool {
	return t == c.Type()
}
