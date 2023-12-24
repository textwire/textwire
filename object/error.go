package object

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return STRING_OBJ
}

func (e *Error) Inspect() string {
	return e.Message
}
