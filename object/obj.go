package object

type Obj struct {
	Pairs map[string]Object
}

func (o *Obj) Type() ObjectType {
	return OBJ_OBJ
}

func (o *Obj) String() string {
	return ""
}

func (o *Obj) Val() interface{} {
	result := make(map[string]interface{})

	for k, v := range o.Pairs {
		result[k] = v.Val()
	}

	return result
}

func (o *Obj) Is(t ObjectType) bool {
	return t == o.Type()
}
