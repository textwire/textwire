package object

type Obj struct {
	Pairs map[string]Object
}

func (o *Obj) Type() ObjectType {
	return OBJ_OBJ
}

func (a *Obj) String() string {
	return ""
}

func (o *Obj) Is(t ObjectType) bool {
	return t == o.Type()
}
