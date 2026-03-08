package value

type Insert struct {
	Name  string
	Block Value //@insert(name)<Block>@end or @insert(name, <Block>)
}

func (i *Insert) Type() ValueType {
	return RESERVE_VAL
}

func (i *Insert) String() string {
	if i.Block == nil {
		panic("Block field on Insert object must not be nil when calling String()")
	}

	return i.Block.String()
}

func (r *Insert) Dump(ident int) string {
	return ""
}

func (i *Insert) JSON() (string, error) {
	return "", nil
}

func (r *Insert) Native() any {
	if r.Block == nil {
		panic("Block field on Insert object must not be nil when calling Native()")
	}

	return r.Block.Native()
}

func (i *Insert) Is(t ValueType) bool {
	return t == i.Type()
}
