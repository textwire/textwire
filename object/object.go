package object

type ObjectType string

const (
	NIL_OBJ   = "NIL"
	ERROR_OBJ = "ERROR"

	INT_OBJ     = "INT"
	FLOAT_OBJ   = "FLOAT"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"
	ARRAY_OBJ   = "SLICE"
	HTML_OBJ    = "HTML"
	USE_OBJ     = "LAYOUT"
	RESERVE_OBJ = "RESERVE"
	INSERT_OBJ  = "INSERT"
	BLOCK_OBJ   = "BLOCK"

	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

type Object interface {
	Type() ObjectType
	String() string
}
