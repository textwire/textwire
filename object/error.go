package object

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return STRING_OBJ
}

func (e *Error) String() string {
	return e.Message
}
