package object

type BuiltinFunction func(receiver Object, args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) Val() interface{} {
	return b.Fn
}

func (b *Builtin) String() string {
	return "builtin function"
}
