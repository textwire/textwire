package fail

import "fmt"

const (
	// Parser errors
	ErrEmptyBrackets          = "bracket statement must contain an expression '{{ <expression> }}'"
	ErrWrongNextToken         = "expected next token to be '%s', got '%s' instead"
	ErrExpectedExpression     = "expected expression, got '}}'"
	ErrCouldNotParseAs        = "could not parse '%s' as '%s'"
	ErrNoPrefixParseFunc      = "no prefix parse function for '%s'"
	ErrIllegalToken           = "illegal token '%s' found"
	ErrElseifCannotFollowElse = "'@elseif' directive cannot follow '@else'"

	// Interpreter (evaluator) errors
	ErrUnknownNodeType        = "unknown node type '%T'"
	ErrInsertMustHaveContent  = "the INSERT statement must have a content or a text argument"
	ErrIdentifierNotFound     = "identifier '%s' not found"
	ErrIndexNotSupported      = "the index operator '%s' is not supported"
	ErrUnknownOperator        = "unknown operator '%s%s'"
	ErrTypeMismatch           = "type mismatch '%s %s %s'"
	ErrUnknownTypeForOperator = "unknown type '%s' for '%s' operator"
	ErrPrefixOperatorIsWrong  = "prefix operator '%s' cannot be applied to '%s'"
)

type Error struct {
	message  string
	line     uint
	filepath string
	origin   string // "lexer", "parser", "interpreter"
}

func New(line uint, origin string, msg string, args ...interface{}) *Error {
	return &Error{
		line:    line,
		origin:  origin,
		message: fmt.Sprintf(msg, args...),
	}
}

func (e *Error) String() string {
	suffix := ""

	if e.filepath != "" {
		suffix = fmt.Sprintf(" in %s", e.filepath)
	}

	return fmt.Sprintf("[Textwire error in %s on line %d]: %s%s",
		e.origin, e.line, e.message, suffix)
}
