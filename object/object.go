package object

type ObjectType string

const (
	NIL_OBJ   = "NIL"
	ERROR_OBJ = "ERROR"

	INT_OBJ   = "INT"
	INT64_OBJ = "INT64"
	INT32_OBJ = "INT32"
	INT16_OBJ = "INT16"
	INT8_OBJ  = "INT8"

	UINT_OBJ   = "UINT"
	UINT64_OBJ = "UINT64"
	UINT32_OBJ = "UINT32"
	UINT16_OBJ = "UINT16"
	UINT8_OBJ  = "UINT8"

	FLOAT64_OBJ = "FLOAT64"
	FLOAT32_OBJ = "FLOAT32"

	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"
	HTML_OBJ    = "HTML"

	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

type Object interface {
	Type() ObjectType
	String() string
}
