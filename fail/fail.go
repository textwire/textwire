package fail

import "fmt"

const (
	// Parser errors
	ERR_EMPTY_BRACKETS            = "bracket statement must contain an expression '{{ <expression> }}'"
	ERR_WRONG_NEXT_TOKEN          = "expected next token to be '%s', got '%s' instead"
	ERR_EXPECTED_EXPRESSION       = "expected expression, got '}}'"
	ERR_COULD_NOT_PARSE_AS        = "could not parse '%s' as '%s'"
	ERR_NO_PREFIX_PARSE_FUNC      = "no prefix parse function for '%s'"
	ERR_ILLEGAL_TOKEN             = "illegal token '%s' found"
	ERR_ELSEIF_CANNOT_FOLLOW_ELSE = "ELSEIF statement cannot follow ELSE statement"
)

type Error struct {
	message string
	line    uint
	origin  string // "lexer", "parser", "interpreter"
}

func New(line uint, origin string, msg string, args ...interface{}) *Error {
	return &Error{
		line:    line,
		origin:  origin,
		message: fmt.Sprintf(msg, args...),
	}
}

func (e *Error) String() string {
	return fmt.Sprintf("[Textwire error in %s on line %d]: %s", e.origin, e.line, e.message)
}
