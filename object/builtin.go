package object

import "github.com/textwire/textwire/v2/ctx"

type BuiltinFunction func(ctx *ctx.EvalCtx, receiver Object, args ...Object) (Object, error)

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) String() string {
	return "builtin function"
}

func (b *Builtin) Dump(ident int) string {
	return "builtin"
}

func (b *Builtin) Val() interface{} {
	return b.Fn
}
