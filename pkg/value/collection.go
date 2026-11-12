package value

import "html"

// Compile-time checks that *Arr and *Obj implement collection. If one of them
// loses a method, the build fails here instead of at runtime.
var (
	_ collection = (*Arr)(nil)
	_ collection = (*Obj)(nil)
)

// collection is implemented by *Arr and *Obj. Its methods skip the cycle
// check, because the outermost String, escapedString, Dump, JSON or Native
// call already checked the whole tree.
type collection interface {
	stringNoCycleCheck() string
	htmlNoCycleCheck() string
	dumpNoCycleCheck(indent int) string
	jsonNoCycleCheck() (string, error)
	nativeNoCycleCheck() any
}

// stringOf, htmlOf, dumpOf, jsonOf and nativeOf convert an
// element of an array or object. Collections skip the cycle check, other
// values call their own method.
func stringOf(v Literal) string {
	if c, ok := v.(collection); ok {
		return c.stringNoCycleCheck()
	}

	return v.String()
}

func htmlOf(v Literal) string {
	if c, ok := v.(collection); ok {
		return c.htmlNoCycleCheck()
	}

	if str, ok := v.(*Str); ok && str.raw {
		return str.Val
	}

	return html.EscapeString(v.String())
}

// escapedString returns v as a string with HTML characters escaped, like
// {{ }} prints it. Raw strings are not escaped, also inside arrays and objects.
func escapedString(v Literal) string {
	if _, ok := v.(collection); ok && HasCycle(v) {
		return circularReference
	}

	return htmlOf(v)
}

func dumpOf(v Literal, indent int) string {
	if c, ok := v.(collection); ok {
		return c.dumpNoCycleCheck(indent)
	}

	return v.Dump(indent)
}

func jsonOf(v Literal) (string, error) {
	if c, ok := v.(collection); ok {
		return c.jsonNoCycleCheck()
	}

	return v.JSON()
}

func nativeOf(v Literal) any {
	if c, ok := v.(collection); ok {
		return c.nativeNoCycleCheck()
	}

	return v.Native()
}
