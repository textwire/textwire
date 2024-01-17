package object

type ObjectType string

const (
	NIL_OBJ   = "NIL"
	ERROR_OBJ = "ERROR"

	INT_OBJ     = "INT"
	FLOAT_OBJ   = "FLOAT"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"
	HTML_OBJ    = "HTML"
	LAYOUT_OBJ  = "LAYOUT"
	BLOCK_OBJ   = "BLOCK"

	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

type Object interface {
	Type() ObjectType
	String() string
}
