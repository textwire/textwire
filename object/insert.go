package object

import "fmt"

type Insert struct {
	Name  string
	Block Object //@insert(name)<Block>@end or @insert(name, <Block>)
}

func (i *Insert) Type() ObjectType {
	return RESERVE_OBJ
}

func (i *Insert) String() string {
	if i.Block == nil {
		panic("Block field on Insert object must not be nil when calling String()")
	}

	return i.Block.String()
}

func (r *Insert) Dump(ident int) string {
	return fmt.Sprintf("@insert(%q)", r.Name)
}

func (r *Insert) Val() any {
	if r.Block == nil {
		panic("Block field on Insert object must not be nil when calling Val()")
	}

	return r.Block.Val()
}

func (i *Insert) Is(t ObjectType) bool {
	return t == i.Type()
}
