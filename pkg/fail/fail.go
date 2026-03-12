package fail

import (
	"errors"
	"fmt"
	"log"

	"github.com/textwire/textwire/v3/pkg/position"
)

// Each string constant here is also an ErrorID on Error object.
// It helps to identify error by checking if ErrorID == fail.ErrEmptyBraces.
const (
	// Parser errors
	ErrEmptyBraces            = "empty expression {{}} - must contain valid code like {{ variable }} or {{ 1 + 2 }}"
	ErrWrongNextToken         = "syntax error: expected '%s' but found '%s'"
	ErrExpectedExpression     = "expected expression before '}}'"
	ErrCouldNotParseAs        = "cannot parse '%s' as %s"
	ErrIllegalToken           = "illegal token '%s'"
	ErrElseifCannotFollowElse = "'@elseif' cannot come after '@else'"
	ErrExpectedObjLit         = "expected object literal, got '%s'"
	ErrObjKeyUseGet           = "to access a key that starts with a number, use object.get('key') function"
	ErrSlotNotDefined         = "@component('%s') references @slot('%s') which doesn't exist in the component file"
	ErrDuplicateReserves      = "found duplicate @reserve('%s') inside of a layout file %s"
	ErrDuplicateSlot          = "@slot('%s') used %d times in @component('%s')"
	ErrDuplicateDefaultSlot   = "default @slot used %d times in @component('%s')"
	ErrExpectedComponentName  = "@component('') cannot have empty name"
	ErrExpectedUseName        = "@use('') cannot have empty name"
	ErrUnusedInsertDetected   = "@insert('%s') needs to have a matching @reserve('%s') in layout file"
	ErrDuplicateInserts       = "duplicate @insert('%s') found"
	ErrUseDirFirstArgStr      = "argument 1 of @use(str) must be a string, got @use('%s')"
	ErrOnlyOneUseDir          = "@use() directive can only be used once per template"
	ErrForLoopExpectStmt      = "@for() expects statement as post conditional, got expression '%s', like 'i++', 'i = i + 2', etc"

	// Evaluator (interpreter) errors
	ErrUnknownType             = "unsupported type '%T'"
	ErrInsertMustHaveContent   = "@insert() requires either a block (body) or a second argument"
	ErrIndexNotSupported       = "type '%s' does not support indexing"
	ErrUnknownOp               = "unknown operator '%s%s'"
	ErrCannotUseOperator       = "operator '%s' is not supported for the combination '%s' %s '%s'"
	ErrCannotDecFromFloat      = "cannot decrement from float '%s' due to error: %s"
	ErrPrefixOpIsWrong         = "cannot apply prefix '%s' to type '%s'"
	ErrVariableIsUndefined     = "variable '%s' is not defined"
	ErrReservedIdentifiers     = "'loop' and 'global' are reserved variable names"
	ErrIdentifierTypeMismatch  = "cannot assign identifier '%s' of type '%s' to type '%s'"
	ErrNotSupportedAssign      = "left side of an assign statement must be an identifier, index expression, or object key access, got '%s'"
	ErrDivisionByZero          = "division by zero - divisor cannot be zero"
	ErrEachDirWithNonArrArg    = "cannot use @each(item in array) with non-array type '%s' after 'in' keyword"
	ErrArrIndexInt             = "array index must be an integer, got '%s'"
	ErrArrIndexOutOfBound      = "index %d out of bounds for array of length %d"
	ErrSomeDirsOnlyInTemplates = "@use, @insert, @reserve, @component only allowed in templates"
	ErrInsertRequiresUse       = "@insert('%s') cannot be used without @use()"
	ErrUseDirMissingLayout     = "@use('%s') missing layout file"
	ErrGlobalFuncMissing       = "global function %s() not found"
	ErrKeyOnNonObj             = "'%s' type does not support attribute '%s' access"
	ErrIllegalTypeForInc       = "cannot increment '%s', only integer and float are allowed"
	ErrIllegalTypeForDec       = "cannot decrement '%s', only integer and float are allowed"

	// Functions
	ErrFuncNotDefined   = "%s.%s() is not defined"
	ErrFuncMissingArg   = "%s.%s(arg) missing required argument"
	ErrFuncFirstArgInt  = "argument 1 on %s.%s() must be 'integer'"
	ErrFuncFirstArgStr  = "argument 1 on %s.%s() must be 'string'"
	ErrFuncSecondArgInt = "argument 2 on %s.%s() must be 'integer'"
	ErrFuncSecondArgStr = "argument 2 on %s.%s() must be 'string'"
	ErrFuncMaxArgs      = "%s.%s() takes at most %d arguments"

	// Template errors
	ErrUnsupportedType       = "unsupported value type '%T'"
	ErrDirStmtNotAllowed     = "@use() not allowed in layout files - causes infinite recursion"
	ErrTemplateNotFound      = "template file '%s' not found"
	ErrDefaultSlotNotDefined = "default @slot not defined in @component('%s')"

	// API errors
	ErrFuncAlreadyDefined = "custom function '%s' already defined for type '%s'"
	ErrUndefinedComponent = "@component('%s') missing required component file"
)

// Error is the main error type for Textwire that contains all the necessary
// information about the error like the line number, file path, etc.
type Error struct {
	pos      *position.Pos
	origin   string // "parser" | "evaluator" | "template" | "API"
	filepath string
	message  string
}

// New creates a new Error instance of Error
func New(pos *position.Pos, filepath, origin, msg string, args ...any) *Error {
	if pos == nil {
		pos = &position.Pos{}
	}

	return &Error{
		pos:      pos,
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

func (e *Error) Pos() *position.Pos {
	return e.pos
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
	return fmt.Sprintf("Textwire ERROR%s:%d", path, e.pos.Line())
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

func FromError(err error, pos *position.Pos, absPath, origin string, args ...any) *Error {
	if err == nil {
		panic("err should never be nil in fail.FromError() function")
	}

	if pos == nil {
		pos = &position.Pos{}
	}

	return New(pos, absPath, origin, err.Error(), args...)
}
