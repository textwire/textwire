package object

type BuiltinFunction func(receiver Object, args ...Object) (Object, error)

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
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

func (b *Builtin) Val() any {
	return b.Fn
}
