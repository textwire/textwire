package fail

import (
	"errors"
	"fmt"
	"log"
)

// Each string constant here is also an ErrorID on Error object.
// It helps to identify error by checking if ErrorID == fail.ErrEmptyBraces.
const (
	// Parser errors
	ErrEmptyBraces               = "empty expression {{}} - must contain valid code like {{ variable }} or {{ 1 + 2 }}"
	ErrWrongNextToken            = "syntax error: expected '%s' but found '%s'"
	ErrExpectedExpression        = "expected expression before '}}'"
	ErrCouldNotParseAs           = "cannot parse '%s' as %s"
	ErrNoPrefixParseFunc         = "unexpected token '%s' at start of expression"
	ErrIllegalToken              = "illegal token '%s'"
	ErrElseifCannotFollowElse    = "'@elseif' cannot come after '@else'"
	ErrExceptedComponentStmt     = "expected component statement, got %T"
	ErrComponentMustHaveBlock    = "component '%s' missing required block"
	ErrExpectedObjectLiteral     = "expected object literal, got '%s'"
	ErrSlotNotDefined            = "slot '%s' not defined in component '%s'"
	ErrDefaultSlotNotDefined     = "default slot not defined in component '%s'"
	ErrDuplicateSlotUsage        = "slot '%s' used %d times in component '%s'"
	ErrDuplicateDefaultSlotUsage = "default slot used %d times in component '%s'"
	ErrExpectedComponentName     = "component name cannot be empty"
	ErrUndefinedInsert           = "insert '%s' not found in layout - add matching @reserve"
	ErrDuplicateInserts          = "duplicate insert '%s' found"
	ErrUseStmtFirstArgStr        = "first argument of @use must be a string, got '%s'"

	// Evaluator (interpreter) errors
	ErrUnknownNodeType         = "unsupported expression type '%T'"
	ErrInsertMustHaveContent   = "insert must have content or text argument"
	ErrIndexNotSupported       = "type '%s' does not support indexing"
	ErrUnknownOperator         = "unknown operator '%s%s'"
	ErrCannotSubFromFloat      = "cannot decrement from float '%s' due to error: %s"
	ErrTypeMismatch            = "type mismatch: cannot %s %s %s"
	ErrUnknownTypeForOperator  = "operator '%s' not supported for type '%s'"
	ErrPrefixOperatorIsWrong   = "cannot apply prefix '%s' to type '%s'"
	ErrUseStmtMissingLayout    = "@use statement missing layout file"
	ErrIdentifierIsUndefined   = "variable '%s' is not defined"
	ErrReservedIdentifiers     = "'loop' and 'global' are reserved variable names"
	ErrIdentifierTypeMismatch  = "cannot assign identifier '%s' of type '%s' to type '%s'"
	ErrDotOperatorNotSupported = "type '%s' does not support property access"
	ErrPropertyNotFound        = "property '%s' not found on type '%s'"
	ErrDivisionByZero          = "division by zero - divisor cannot be zero"
	ErrSomeDirsOnlyInTemplates = "@use, @insert, @reserve, @component only allowed in templates"
	ErrGlobalFuncMissing       = "global function '%s' not found"
	ErrPropertyOnNonObject     = "cannot get property '%s' from a non-object type '%s'"
	ErrEachDirWithNonArrArg    = "cannot use @each statement with non-array type '%s' after 'in' keyword"

	// Functions
	ErrNoFuncForThisType  = "function '%s' not available for type '%s'"
	ErrFuncRequiresOneArg = "function '%s' on type '%s' requires at least one argument"
	ErrFuncFirstArgInt    = "first argument for function '%s' on type '%s' must be an INTEGER"
	ErrFuncFirstArgStr    = "first argument for function '%s' on type '%s' must be a STRING"
	ErrFuncSecondArgInt   = "second argument for function '%s' on type '%s' must be an INTEGER"
	ErrFuncSecondArgStr   = "second argument for function '%s' on type '%s' must be a STRING"
	ErrFuncMaxArgs        = "function '%s' on type '%s' accepts a maximum of %d arguments"

	// Template errors
	ErrUnsupportedType   = "unsupported value type '%T'"
	ErrUseStmtNotAllowed = "@use not allowed in layout files - causes infinite recursion"
	ErrTemplateNotFound  = "template file '%s' not found"

	// API errors
	ErrFuncAlreadyDefined        = "custom function '%s' already defined for type '%s'"
	ErrCannotOverrideBuiltInFunc = "cannot override built-in function '%s' for '%s'"
	ErrProgramNotFound           = "parsed AST program not found for '%s' file"
	ErrUndefinedComponent        = "component '%s' is not defined"

	NoErrorsFound = "no Textwire errors found"
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

func FromError(err error, line uint, absPath, origin string, args ...any) *Error {
	if err == nil {
		panic("err should never be nil in fail.FromError() function")
	}

	return New(line, absPath, origin, err.Error(), args...)
}
