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
	ErrEmptyBraces            = "empty expression {{}} - must contain valid code like {{ variable }} or {{ 1 + 2 }}"
	ErrWrongNextToken         = "syntax error: expected '%s' but found '%s'"
	ErrExpectedExpression     = "expected expression before '}}'"
	ErrCouldNotParseAs        = "cannot parse '%s' as %s"
	ErrNoPrefixParseFunc      = "unexpected token '%s' at start of expression"
	ErrIllegalToken           = "illegal token '%s'"
	ErrElseifCannotFollowElse = "'@elseif' cannot come after '@else'"
	ErrExpectedObjectLiteral  = "expected object literal, got '%s'"
	ErrSlotNotDefined         = "@component('%s') references @slot('%s') which doesn't exist in the component file"
	ErrDuplicateReserves      = "found duplicate @reserve('%s') inside of a layout file %s"
	ErrDuplicateSlot          = "@slot('%s') used %d times in @component('%s')"
	ErrDuplicateDefaultSlot   = "default @slot used %d times in @component('%s')"
	ErrExpectedComponentName  = "@component('') cannot have empty name"
	ErrUnusedInsertDetected   = "@insert('%s') needs to have a matching @reserve('%s') in layout file"
	ErrDuplicateInserts       = "duplicate @insert('%s') found"
	ErrUseStmtFirstArgStr     = "argument 1 of @use(STR) must be a string, got @use('%s')"
	ErrOnlyOneUseDir          = "@use() directive can only be used once per template"
	ErrSlotifPosition         = "@slotif() directive can only be used inside @component directive"

	// Evaluator (interpreter) errors
	ErrUnknownNodeType         = "unsupported expression type '%T'"
	ErrInsertMustHaveContent   = "@insert() requires either a block (body) or a second argument"
	ErrIndexNotSupported       = "type '%s' does not support indexing"
	ErrUnknownOp               = "unknown operator '%s%s'"
	ErrCannotSubFromFloat      = "cannot decrement from float '%s' due to error: %s"
	ErrTypeMismatch            = "type mismatch: cannot %s %s %s"
	ErrUnknownTypeForOp        = "operator '%s' not supported for type '%s'"
	ErrPrefixOpIsWrong         = "cannot apply prefix '%s' to type '%s'"
	ErrVariableIsUndefined     = "variable '%s' is not defined"
	ErrReservedIdentifiers     = "'loop' and 'global' are reserved variable names"
	ErrIdentifierTypeMismatch  = "cannot assign identifier '%s' of type '%s' to type '%s'"
	ErrDivisionByZero          = "division by zero - divisor cannot be zero"
	ErrEachDirWithNonArrArg    = "cannot use @each(item in ARRAY) with non-array type '%s' after 'in' keyword"
	ErrSomeDirsOnlyInTemplates = "@use, @insert, @reserve, @component only allowed in templates"
	ErrInsertRequiresUse       = "@insert('%s') cannot be used without @use()"
	ErrUseStmtMissingLayout    = "@use('%s') missing layout file"
	ErrGlobalFuncMissing       = "global function %s() not found"
	ErrPropertyOnNonObject     = "'%s' type does not support attribute '%s' access"

	// Functions
	ErrFuncNotDefined   = "%s.%s() is not defined"
	ErrFuncMissingArg   = "%s.%s(arg) missing required argument"
	ErrFuncFirstArgInt  = "argument 1 on %s.%s() must be INTEGER"
	ErrFuncFirstArgStr  = "argument 1 on %s.%s() must be STRING"
	ErrFuncSecondArgInt = "argument 2 on %s.%s() must be INTEGER"
	ErrFuncSecondArgStr = "argument 2 on %s.%s() must be STRING"
	ErrFuncMaxArgs      = "%s.%s() takes at most %d arguments"

	// Template errors
	ErrUnsupportedType       = "unsupported value type '%T'"
	ErrUseStmtNotAllowed     = "@use() not allowed in layout files - causes infinite recursion"
	ErrTemplateNotFound      = "template file '%s' not found"
	ErrDefaultSlotNotDefined = "default @slot not defined in @component('%s')"

	// API errors
	ErrFuncAlreadyDefined = "custom function '%s' already defined for type '%s'"
	ErrUndefinedComponent = "@component('%s') missing required component file"
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

func (e *Error) Origin() string {
	return e.origin
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
