package fail

import (
	"fmt"
	"log"
)

const (
	// Parser errors
	ErrEmptyBrackets          = "Bracket statement must contain an expression '{{ <expression> }}'"
	ErrWrongNextToken         = "Expected next token to be '%s', got '%s' instead"
	ErrExpectedExpression     = "Expected expression, got '}}'"
	ErrCouldNotParseAs        = "Could not parse '%s' as '%s'"
	ErrNoPrefixParseFunc      = "No prefix parse function for '%s'"
	ErrIllegalToken           = "Illegal token '%s' found"
	ErrElseifCannotFollowElse = "'@elseif' directive cannot follow '@else'"
	ErrExceptedReserveStmt    = "Expected *ReserveStatement, got %T"
	ErrInsertStmtNotDefined   = "The insert statement named '%s' is not defined"
	ErrExpectedIdentifier     = "Expected identifier, got '%s' instead"

	// Interpreter (evaluator) errors
	ErrUnknownNodeType        = "Unknown node type '%T'"
	ErrInsertMustHaveContent  = "The INSERT statement must have a content or a text argument"
	ErrIdentifierNotFound     = "Identifier '%s' not found"
	ErrIndexNotSupported      = "The index operator '%s' is not supported"
	ErrUnknownOperator        = "Unknown operator '%s%s'"
	ErrTypeMismatch           = "Type mismatch '%s %s %s'"
	ErrUnknownTypeForOperator = "Unknown type '%s' for '%s' operator"
	ErrPrefixOperatorIsWrong  = "Prefix operator '%s' cannot be applied to '%s'"
	ErrUseStmtMustHaveProgram = "The 'use' statement must have a program attached"
	ErrFuncDoNotExist         = "Function '%s' doesn't exist. Try to update the Textwire to the latest version or checkout the documentation https://textwire.github.io"
	ErrNoFuncForThisType      = "Function '%s' doesn't exist for type '%s'"

	// Template errors
	ErrUnsupportedType   = "Unsupported type '%T'"
	ErrTemplateNotFound  = "Template not found"
	ErrUseStmtNotAllowed = "The 'use' statement is not allowed in a layout file. It will cause infinite recursion"

	NoErrorsFound = "There are no Textwire errors"
)

// Error is the main error type for Textwire that contains all the necessary
// information about the error like the line number, file path, etc.
type Error struct {
	message  string
	line     uint
	filepath string
	origin   string // "parser" | "evaluator" | "template"
}

// New creates a new Error instance of Error
func New(line uint, filepath, origin, msg string, args ...interface{}) *Error {
	return &Error{
		line:     line,
		origin:   origin,
		filepath: filepath,
		message:  fmt.Sprintf(msg, args...),
	}
}

// String returns the full error message with all the details
func (e *Error) String() string {
	path := ""

	if e.filepath != "" {
		path = fmt.Sprintf(" in %s", e.filepath)
	}

	return fmt.Sprintf("[Textwire ERROR%s:%d]: %s",
		path, e.line, e.message)
}

// FatalOnError calls log.Fatal if the error message is not empty
func (e *Error) FatalOnError() {
	if e.message == "" {
		return
	}

	log.Fatal(e.String())
}

// PanicOnError panics if the error message is not empty
func (e *Error) PanicOnError() {
	if e.message == "" {
		return
	}

	panic(e.String())
}

// PrintOnError prints the error message to the standard output
// when the error message is not empty
func (e *Error) PrintOnError() {
	if e.message == "" {
		return
	}

	log.Println(e.String())
}

func (e *Error) ToSlice() []*Error {
	return []*Error{e}
}

func FromError(err error, line uint, absPath, origin string, args ...interface{}) *Error {
	if err == nil {
		return nil
	}

	return New(line, absPath, origin, err.Error(), args...)
}
