package object

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) String() string {
	return e.Message
}

func (e *Error) Is(t ObjectType) bool {
	return t == e.Type()
}
