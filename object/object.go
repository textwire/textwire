package object

type ObjectType string

const (
	NIL_OBJ   = "NIL"
	ERROR_OBJ = "ERROR"

	INTEGER_OBJ          = "INTEGER"
	UNSIGNED_INTEGER_OBJ = "UNSIGNED_INTEGER"
	FLOAT_OBJ            = "FLOAT"
	BOOLEAN_OBJ          = "BOOLEAN"
	STRING_OBJ           = "STRING"
	HTML_OBJ             = "HTML"

	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

type Object interface {
	Type() ObjectType
	String() string
}
