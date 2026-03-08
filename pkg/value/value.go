package value

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type ValueType string

const (
	ERR_VAL = "error"

	// Literal
	NIL_VAL   ValueType = "nil"
	INT_VAL   ValueType = "integer"
	FLOAT_VAL ValueType = "float"
	BOOL_VAL  ValueType = "boolean"
	STR_VAL   ValueType = "string"
	ARR_VAL   ValueType = "array"
	OBJ_VAL   ValueType = "object"

	TEXT_VAL      ValueType = "text"
	USE_VAL       ValueType = "layout"
	RESERVE_VAL   ValueType = "reserve"
	INSERT_VAL    ValueType = "insert"
	BLOCK_VAL     ValueType = "block"
	BUILTIN_VAL   ValueType = "function"
	COMPONENT_VAL ValueType = "component"
	SLOT_VAL      ValueType = "slot"
	DUMP_VAL      ValueType = "dump"

	BREAK_VAL      ValueType = "break"
	BREAKIF_VAL    ValueType = "breakif"
	CONTINUE_VAL   ValueType = "continue"
	CONTINUEIF_VAL ValueType = "continueif"

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
	// move Dump and JSON to LiteralValue in v4.0.0
	Dump(ident int) string
	JSON() (string, error)
}

func FromTokenToValueType(astType token.TokenType) ValueType {
	switch astType {
	case token.INT:
		return INT_VAL
	case token.LBRACE:
		return OBJ_VAL
	case token.LBRACKET:
		return ARR_VAL
	case token.FLOAT:
		return FLOAT_VAL
	case token.TRUE, token.FALSE:
		return BOOL_VAL
	case token.STR:
		return STR_VAL
	case token.NIL:
		return NIL_VAL
	default:
		return ERR_VAL
	}
}
