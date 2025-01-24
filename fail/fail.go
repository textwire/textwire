package fail

import (
	"errors"
	"fmt"
	"log"
)

const (
	// Parser errors
	ErrEmptyBrackets             = "bracket statement must contain an expression '{{ <expression> }}'"
	ErrWrongNextToken            = "expected next token to be '%s', got '%s' instead"
	ErrExpectedExpression        = "expected expression, got '}}'"
	ErrCouldNotParseAs           = "could not parse '%s' as '%s'"
	ErrNoPrefixParseFunc         = "no prefix parse function for '%s'"
	ErrIllegalToken              = "illegal token '%s' found"
	ErrElseifCannotFollowElse    = "'@elseif' directive cannot follow '@else'"
	ErrExpectedIdentifier        = "expected identifier, got '%s' instead"
	ErrExceptedComponentStmt     = "expected *ComponentStmt, got %T"
	ErrComponentMustHaveBlock    = "the component '%s' must have a block"
	ErrExpectedObjectLiteral     = "expected object literal, got '%s' instead"
	ErrSlotNotDefined            = "'%s' slot is not defined in the component '%s'"
	ErrDefaultSlotNotDefined     = "default slot is not defined in the component '%s'"
	ErrDuplicateSlotUsage        = "duplicate slot usage '%s' found %d times in the component '%s'"
	ErrDuplicateDefaultSlotUsage = "duplicate default slot usage found %d times in the component '%s'"
	ErrExpectedComponentName     = "expected component name, got empty string instead"
	ErrUndefinedComponent        = "component '%s' is not defined. Check if component exists"

	// Evaluator (interpreter) errors
	ErrUnknownNodeType         = "unknown node type '%T'"
	ErrInsertMustHaveContent   = "the INSERT statement must have a content or a text argument"
	ErrIdentifierNotFound      = "identifier '%s' not found"
	ErrIndexNotSupported       = "the index operator '%s' is not supported"
	ErrUnknownOperator         = "unknown operator '%s%s'"
	ErrTypeMismatch            = "type mismatch '%s %s %s'"
	ErrUnknownTypeForOperator  = "unknown type '%s' for '%s' operator"
	ErrPrefixOperatorIsWrong   = "prefix operator '%s' cannot be applied to '%s'"
	ErrUseStmtMustHaveProgram  = "the 'use' statement must have a program attached"
	ErrLoopVariableIsReserved  = "the 'loop' variable is reserved. You cannot use it as a variable name"
	ErrVariableTypeMismatch    = "cannot assign variable '%s' of type '%s' to type '%s'"
	ErrDotOperatorNotSupported = "the dot operator is not supported for type '%s'"
	ErrPropertyNotFound        = "property '%s' not found in type '%s'"
	ErrDivisionByZero          = "division by zero error. The right-hand side of the division operator must not be zero"

	// Functions
	ErrNoFuncForThisType  = "function '%s' doesn't exist for type '%s'"
	ErrFuncRequiresOneArg = "function '%s' on type '%s' requires at least one argument"
	ErrFuncFirstArgInt    = "the first argument for function '%s' on type '%s' must be an INTEGER"
	ErrFuncFirstArgStr    = "the first argument for function '%s' on type '%s' must be a STRING"
	ErrFuncSecondArgInt   = "the second argument for function '%s' on type '%s' must be an INTEGER"
	ErrFuncSecondArgStr   = "the second argument for function '%s' on type '%s' must be a STRING"
	ErrFuncMaxArgs        = "function '%s' on type '%s' accepts a maximum of '%d' arguments"

	// Template errors
	ErrUnsupportedType   = "unsupported type '%T'"
	ErrTemplateNotFound  = "template not found"
	ErrUseStmtNotAllowed = "the 'use' statement is not allowed in a layout file. It will cause infinite recursion"

	// API errors
	ErrFuncAlreadyDefined        = "custom function '%s' already defined for '%s'"
	ErrCannotOverrideBuiltInFunc = "cannot override built-in function '%s' for '%s'"

	NoErrorsFound = "there are no Textwire errors"
)

// Error is the main error type for Textwire that contains all the necessary
// information about the error like the line number, file path, etc.
type Error struct {
	message  string
	line     uint
	filepath string
	origin   string // "parser" | "evaluator" | "template" | "API"
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

// Meta returns the error meta information like the file path and line number
func (e *Error) Meta() string {
	var path string

	if e.filepath != "" {
		path = fmt.Sprintf(" in %s", e.filepath)
	}

	return fmt.Sprintf("Textwire %s ERROR%s:%d", e.origin, path, e.line)
}

// Message returns the error message
func (e *Error) Message() string {
	return e.message
}

// String returns the full error message with all the details
func (e *Error) String() string {
	return fmt.Sprintf("[%s]: %s", e.Meta(), e.Message())
}

// FatalOnError calls log.Fatal if the error message is not empty
func (e *Error) FatalOnError() {
	if e == nil {
		return
	}

	log.Fatal(e.String())
}

// PanicOnError panics if the error message is not empty
func (e *Error) PanicOnError() {
	if e == nil {
		return
	}

	panic(e.String())
}

// PrintOnError prints the error message to the standard output
// when the error message is not empty
func (e *Error) PrintOnError() {
	if e == nil {
		return
	}

	log.Println(e.String())
}

func (e *Error) Error() error {
	return errors.New(e.String())
}

func FromError(err error, line uint, absPath, origin string, args ...interface{}) *Error {
	if err == nil {
		return nil
	}

	return New(line, absPath, origin, err.Error(), args...)
}
