package value

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type ValueType string

const (
	ERR_OBJ = "ERROR"

	// Literal
	NIL_OBJ   ValueType = "NIL"
	INT_OBJ   ValueType = "INTEGER"
	FLOAT_OBJ ValueType = "FLOAT"
	BOOL_OBJ  ValueType = "BOOLEAN"
	STR_OBJ   ValueType = "STRING"
	ARR_OBJ   ValueType = "ARRAY"
	OBJ_OBJ   ValueType = "OBJECT"

	HTML_OBJ      ValueType = "HTML"
	USE_OBJ       ValueType = "LAYOUT"
	RESERVE_OBJ   ValueType = "RESERVE"
	INSERT_OBJ    ValueType = "INSERT"
	BLOCK_OBJ     ValueType = "BLOCK"
	BUILTIN_OBJ   ValueType = "FUNCTION"
	COMPONENT_OBJ ValueType = "COMPONENT"
	SLOT_OBJ      ValueType = "SLOT"
	DUMP_OBJ      ValueType = "DUMP"

	BREAK_OBJ       ValueType = "BREAK"
	BREAK_IF_OBJ    ValueType = "BREAK_IF"
	CONTINUE_OBJ    ValueType = "CONTINUE"
	CONTINUE_IF_OBJ ValueType = "CONTINUE_IF"

	DUMP_PROP    = "color: #f8f8f2 !important"
	DUMP_STR     = "color: #c3e88d !important"
	DUMP_NUM     = "color: #76a8ff !important"
	DUMP_KEYWORD = "color: #c792ea !important"
	DUMP_BRACE   = "color: #e99f33 !important"
	DUMP_META    = "color: #2c8ed0 !important"
	DUMP_KEY     = "color: #ffcb8b !important"
)

type Value interface {
	Type() ValueType
	String() string
	Is(ValueType) bool
	Native() any
	// move Dump and JSON to Literal object in v4.0.0
	Dump(ident int) string
	JSON() (string, error)
}

func FromTokenToObjectType(astType token.TokenType) ValueType {
	switch astType {
	case token.INT:
		return INT_OBJ
	case token.LBRACE:
		return OBJ_OBJ
	case token.LBRACKET:
		return ARR_OBJ
	case token.FLOAT:
		return FLOAT_OBJ
	case token.TRUE, token.FALSE:
		return BOOL_OBJ
	case token.STR:
		return STR_OBJ
	case token.NIL:
		return NIL_OBJ
	default:
		return ERR_OBJ
	}
}
