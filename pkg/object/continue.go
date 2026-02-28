package object

type Continue struct{}

func (c *Continue) Type() ObjectType {
	return CONTINUE_OBJ
}

func (c *Continue) String() string {
	return ""
}

func (c *Continue) Dump(ident int) string {
	return ""
}

func (c *Continue) JSON() (string, error) {
	return "", nil
}

func (c *Continue) Val() any {
	return nil
}

func (c *Continue) Is(t ObjectType) bool {
	return t == c.Type()
}
