package object

type ObjectType string

const (
	NIL_OBJ = "NIL"
	ERR_OBJ = "ERROR"

	INT_OBJ     = "INTEGER"
	FLOAT_OBJ   = "FLOAT"
	BOOL_OBJ    = "BOOLEAN"
	STR_OBJ     = "STRING"
	ARR_OBJ     = "SLICE"
	HTML_OBJ    = "HTML"
	USE_OBJ     = "LAYOUT"
	RESERVE_OBJ = "RESERVE"
	INSERT_OBJ  = "INSERT"
	BLOCK_OBJ   = "BLOCK"
	BUILTIN_OBJ = "FUNCTION"
)

type Object interface {
	Type() ObjectType
	String() string
	Is(ObjectType) bool
}
