package object

import "github.com/textwire/textwire/v2/fail"

type Error struct {
	Err *fail.Error
}

func (e *Error) Type() ObjectType {
	return ERR_OBJ
}

func (e *Error) String() string {
	return e.Err.String()
}

func (e *Error) Is(t ObjectType) bool {
	return t == e.Type()
}
