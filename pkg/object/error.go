package object

import (
	"bytes"
	"fmt"

	"github.com/textwire/textwire/v3/pkg/fail"
)

type Error struct {
	Err *fail.Error

	// ErrorID is a raw error message with %s characters instead of values
	ErrorID string
}

func (e *Error) Type() ObjectType {
	return ERR_OBJ
}

func (e *Error) String() string {
	if e.Err == nil {
		panic("Err field on Error object must not be nil when calling String()")
	}
	return e.Err.String()
}

func (e *Error) Dump(ident int) string {
	var out bytes.Buffer
	out.Grow(4)

	fmt.Fprintf(&out, `<span style="%s">error"""</span>`+"\n", DUMP_META)
	fmt.Fprintf(&out, `<span style="%s">%s</span>`+"\n\n", DUMP_KEY, e.Err.Meta())
	fmt.Fprintf(&out, `<span style="%s">%s</span>`+"\n", DUMP_STR, e.Err.Message())
	fmt.Fprintf(&out, `<span style="%s">"""</span>`, DUMP_META)

	return out.String()
}

func (e *Error) Val() any {
	return e.String()
}

func (e *Error) Is(t ObjectType) bool {
	return t == e.Type()
}
