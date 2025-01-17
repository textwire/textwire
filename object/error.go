package object

import (
	"bytes"
	"fmt"

	"github.com/textwire/textwire/v2/fail"
)

type Error struct {
	Err *fail.Error
}

func (e *Error) Type() ObjectType {
	return ERR_OBJ
}

func (e *Error) String() string {
	return e.Err.String()
}

func (e *Error) Dump(ident int) string {
	var out bytes.Buffer

	out.WriteString("<span class='textwire-meta'>error >>></span>\n")
	out.WriteString(fmt.Sprintf("<span class='textwire-key'>%s</span>\n\n", e.Err.Meta()))
	out.WriteString(fmt.Sprintf("<span class='textwire-str'>%s</span>\n", e.Err.Message()))
	out.WriteString("<span class='textwire-meta'><<<</span>")

	return out.String()
}

func (e *Error) Val() interface{} {
	return e.Err.String()
}

func (e *Error) Is(t ObjectType) bool {
	return t == e.Type()
}
