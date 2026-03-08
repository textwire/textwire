package value

type Continue struct{}

func (c *Continue) Type() ValueType {
	return CONTINUE_VAL
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

func (c *Continue) Native() any {
	return nil
}

func (c *Continue) Is(t ValueType) bool {
	return t == c.Type()
}
