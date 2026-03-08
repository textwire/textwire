package value

type BuiltinFunction func(receiver Value, args ...Value) (Value, error)

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ValueType {
	return BUILTIN_OBJ
}

func (b *Builtin) String() string {
	return "func()"
}

func (b *Builtin) Dump(ident int) string {
	return ""
}

func (b *Builtin) JSON() (string, error) {
	return "", nil
}

func (b *Builtin) Native() any {
	return b.Fn
}
