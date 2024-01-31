package object

import "github.com/textwire/textwire/fail"

type Error struct {
	Err *fail.Error
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) String() string {
	return e.Err.String()
}

func (e *Error) Is(t ObjectType) bool {
	return t == e.Type()
}
