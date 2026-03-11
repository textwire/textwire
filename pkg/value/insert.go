package value

type Insert struct {
	Name  string
	Block Value //@insert(name)<Block>@end or @insert(name, <Block>)
}

func (*Insert) Type() ValueType {
	return RESERVE_VAL
}

func (i *Insert) String() string {
	if i.Block == nil {
		panic("Block field on Insert must not be nil when calling String()")
	}

	return i.Block.String()
}

func (i *Insert) Is(t ValueType) bool {
	return t == i.Type()
}
