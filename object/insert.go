package object

type Insert struct {
	Name    string
	Content Object
}

func (i *Insert) Type() ObjectType {
	return INSERT_OBJ
}

func (i *Insert) String() string {
	return i.Content.String()
}
