package object

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type ObjectType string

const (
	ERR_OBJ = "ERROR"

	// Literals
	INTEGER_OBJ ObjectType = "INTEGER"
	FLOAT_OBJ   ObjectType = "FLOAT"
	BOOLEAN_OBJ ObjectType = "BOOLEAN"
	STRING_OBJ  ObjectType = "STRING"
	ARRARY_OBJ  ObjectType = "ARRAY"
	MAP_OBJ     ObjectType = "OBJECT"
	NIL_OBJ                = "NIL"

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

	DUMP_PROP    = "color: #f8f8f2 !important"
	DUMP_STR     = "color: #c3e88d !important"
	DUMP_NUM     = "color: #76a8ff !important"
	DUMP_KEYWORD = "color: #c792ea !important"
	DUMP_BRACE   = "color: #e99f33 !important"
	DUMP_META    = "color: #2c8ed0 !important"
	DUMP_KEY     = "color: #ffcb8b !important"
)

type Object interface {
	Type() ObjectType
	String() string
	Is(ObjectType) bool
	Native() any
	// move Dump and JSON to Literal object in v4.0.0
	Dump(ident int) string
	JSON() (string, error)
}

func FromTokenToObjectType(astType token.TokenType) ObjectType {
	switch astType {
	case token.INT:
		return INTEGER_OBJ
	case token.LBRACE:
		return MAP_OBJ
	case token.LBRACKET:
		return ARRARY_OBJ
	case token.FLOAT:
		return FLOAT_OBJ
	case token.TRUE, token.FALSE:
		return BOOLEAN_OBJ
	case token.STR:
		return STRING_OBJ
	case token.NIL:
		return NIL_OBJ
	default:
		return ERR_OBJ
	}
}
