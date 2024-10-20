package object

type ObjectType string

const (
	NIL_OBJ = "NIL"
	ERR_OBJ = "ERROR"

	INT_OBJ       = "INTEGER"
	FLOAT_OBJ     = "FLOAT"
	BOOL_OBJ      = "BOOLEAN"
	STR_OBJ       = "STRING"
	ARR_OBJ       = "ARRAY"
	OBJ_OBJ       = "OBJECT"
	HTML_OBJ      = "HTML"
	USE_OBJ       = "LAYOUT"
	RESERVE_OBJ   = "RESERVE"
	INSERT_OBJ    = "INSERT"
	BLOCK_OBJ     = "BLOCK"
	BUILTIN_OBJ   = "FUNCTION"
	COMPONENT_OBJ = "COMPONENT"
	SLOT_OBJ      = "SLOT"

	BREAK_OBJ       = "BREAK"
	BREAK_IF_OBJ    = "BREAK_IF"
	CONTINUE_OBJ    = "CONTINUE"
	CONTINUE_IF_OBJ = "CONTINUE_IF"
)

type Object interface {
	Type() ObjectType
	String() string
	Is(ObjectType) bool
	Val() interface{}
}
