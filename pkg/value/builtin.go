package value

type BuiltinFunc func(receiver Literal, args ...Literal) (Literal, error)

type Builtin struct {
	Fn BuiltinFunc
}

func (*Builtin) Type() ValueType {
	return BUILTIN_VAL
}

func (*Builtin) String() string {
	return ""
}
