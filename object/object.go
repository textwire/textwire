package object

type ObjectType string

const (
	NIL_OBJ = "NIL"
	ERR_OBJ = "ERROR"

	INT_OBJ       ObjectType = "INTEGER"
	FLOAT_OBJ     ObjectType = "FLOAT"
	BOOL_OBJ      ObjectType = "BOOLEAN"
	STR_OBJ       ObjectType = "STRING"
	ARR_OBJ       ObjectType = "ARRAY"
	OBJ_OBJ       ObjectType = "OBJECT"
	HTML_OBJ      ObjectType = "HTML"
	USE_OBJ       ObjectType = "LAYOUT"
	RESERVE_OBJ   ObjectType = "RESERVE"
	INSERT_OBJ    ObjectType = "INSERT"
	BLOCK_OBJ     ObjectType = "BLOCK"
	BUILTIN_OBJ   ObjectType = "FUNCTION"
	COMPONENT_OBJ ObjectType = "COMPONENT"
	SLOT_OBJ      ObjectType = "SLOT"
	DUMP_OBJ      ObjectType = "DUMP"

	BREAK_OBJ       ObjectType = "BREAK"
	BREAK_IF_OBJ    ObjectType = "BREAK_IF"
	CONTINUE_OBJ    ObjectType = "CONTINUE"
	CONTINUE_IF_OBJ ObjectType = "CONTINUE_IF"
)

type Object interface {
	Type() ObjectType
	String() string
	Dump(ident int) string
	Is(ObjectType) bool
	Val() interface{}
}
