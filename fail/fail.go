package fail

import (
	"errors"
	"fmt"
	"log"
)

const (
	// Parser errors
	ErrEmptyBraces               = "bracket statement must contain an expression '{{ <expression> }}'"
	ErrWrongNextToken            = "expected next token to be '%s', got '%s' instead"
	ErrExpectedExpression        = "expected expression, got '}}'"
	ErrCouldNotParseAs           = "could not parse '%s' as '%s'"
	ErrNoPrefixParseFunc         = "no prefix parse function for '%s'"
	ErrIllegalToken              = "illegal token '%s' found"
	ErrElseifCannotFollowElse    = "'@elseif' directive cannot follow '@else'"
	ErrExceptedComponentStmt     = "expected *ComponentStmt, got %T"
	ErrComponentMustHaveBlock    = "the component '%s' must have a block"
	ErrExpectedObjectLiteral     = "expected object literal, got '%s' instead"
	ErrSlotNotDefined            = "'%s' slot is not defined in the component '%s'"
	ErrDefaultSlotNotDefined     = "default slot is not defined in the component '%s'"
	ErrDuplicateSlotUsage        = "duplicate slot usage '%s' found %d times in the component '%s'"
	ErrDuplicateDefaultSlotUsage = "duplicate default slot usage found %d times in the component '%s'"
	ErrExpectedComponentName     = "expected component name, got empty string instead"
	ErrUndefinedComponent        = "component '%s' is not defined. Check if component exists"
	ErrUndefinedInsert           = "insert with the name '%s' is not defined in layout. Check if you have a matching reserve statement with the same name"
	ErrDuplicateInserts          = "duplicate insert statements with the name '%s' found"

	// Evaluator (interpreter) errors
	ErrUnknownNodeType         = "unknown node type '%T'"
	ErrInsertMustHaveContent   = "insert statement must have a content or a text argument"
	ErrIdentifierNotFound      = "identifier '%s' not found"
	ErrIndexNotSupported       = "index operator '%s' is not supported"
	ErrUnknownOperator         = "unknown operator '%s%s'"
	ErrTypeMismatch            = "type mismatch '%s %s %s'"
	ErrUnknownTypeForOperator  = "unknown type '%s' for '%s' operator"
	ErrPrefixOperatorIsWrong   = "prefix operator '%s' cannot be applied to '%s'"
	ErrUseStmtMustHaveProgram  = "use statement must have a program attached"
	ErrLoopVariableIsReserved  = "loop variable is reserved. You cannot use it as a variable name"
	ErrVariableTypeMismatch    = "cannot assign variable '%s' of type '%s' to type '%s'"
	ErrDotOperatorNotSupported = "the dot operator is not supported for type '%s'"
	ErrPropertyNotFound        = "property '%s' not found in type '%s'"
	ErrDivisionByZero          = "division by zero error. The right-hand side of the division operator must not be zero"

	// Functions
	ErrNoFuncForThisType  = "function '%s' doesn't exist for type '%s'"
	ErrFuncRequiresOneArg = "function '%s' on type '%s' requires at least one argument"
	ErrFuncFirstArgInt    = "first argument for function '%s' on type '%s' must be an INTEGER"
	ErrFuncFirstArgStr    = "first argument for function '%s' on type '%s' must be a STRING"
	ErrFuncSecondArgInt   = "second argument for function '%s' on type '%s' must be an INTEGER"
	ErrFuncSecondArgStr   = "second argument for function '%s' on type '%s' must be a STRING"
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
func New(line uint, filepath, origin, msg string, args ...any) *Error {
	return &Error{
		line:     line,
		origin:   origin,
		filepath: filepath,
		message:  fmt.Sprintf(msg, args...),
	}
}

func (e *Error) Filepath() string {
	return e.filepath
}

func (e *Error) Line() uint {
	return e.line
}

func (e *Error) Message() string {
	return e.message
}

// Meta returns the error meta information like the file path and line number
func (e *Error) Meta() string {
	var path string

	if e.filepath != "" {
		path = fmt.Sprintf(" in %s", e.filepath)
	}

	return fmt.Sprintf("Textwire ERROR%s:%d", path, e.line)
}

// String returns the full error message with all the details
func (e *Error) String() string {
	return fmt.Sprintf("[%s]:\n%s", e.Meta(), e.Message())
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

func FromError(err error, line uint, absPath, origin string, args ...any) *Error {
	if err == nil {
		return nil
	}

	return New(line, absPath, origin, err.Error(), args...)
}
