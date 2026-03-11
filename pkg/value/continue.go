package value

type Continue struct{}

func (c *Continue) Type() ValueType {
	return CONTINUE_VAL
}

func (*Continue) String() string {
	return ""
}

func (c *Continue) Is(t ValueType) bool {
	return t == c.Type()
}
