package value

import "fmt"

// NilOrigin describes where a Nil value came from. It is used to tell a
// missing object property or array index apart from a value that exists and
// happens to be nil.
type NilOrigin string

const (
	NilOriginUnset        NilOrigin = ""
	NilOriginMissingProp  NilOrigin = "missing property"
	NilOriginMissingIndex NilOrigin = "missing index"
)

type Nil struct {
	origin NilOrigin
}

func NewNil(origin NilOrigin) *Nil {
	return &Nil{origin: origin}
}

func (n *Nil) Origin() NilOrigin {
	return n.origin
}

func (*Nil) Type() ValueType {
	return NIL_VAL
}

func (n *Nil) String() string {
	return ""
}

func (n *Nil) Dump(ident int) string {
	return fmt.Sprintf(`<span style="%s">nil</span>`, DUMP_KEYWORD)
}

func (n *Nil) JSON() (string, error) {
	return "null", nil
}

func (n *Nil) Native() any {
	return nil
}

func (n *Nil) Is(t ValueType) bool {
	return t == n.Type()
}
