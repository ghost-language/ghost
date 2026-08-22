package object

import (
	"ghostlang.org/x/ghost/color"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/token"
)

// Error is a fault travelling back up through the evaluator.
//
// Ghost has no exceptions: an operation that fails returns an error object, and
// every caller that receives one returns it unchanged until it reaches the top.
// The error carries the fault rather than a formatted string so that whatever
// finally prints it — the CLI, the REPL, a Go program embedding Ghost — decides
// how much of it to show and whether to color it.
type Error struct {
	Fault *fault.Fault
}

// String represents the error as a single line, in the file:line:column form
// tooling knows how to follow.
func (err *Error) String() string {
	return err.Fault.String()
}

// Message returns the sentence describing what went wrong, without the position
// or the kind that precede it in a report.
func (err *Error) Message() string {
	if err.Fault == nil {
		return ""
	}

	return err.Fault.Message
}

// Render lays the error out in full, quoting the line it happened on.
func (err *Error) Render(profile color.Profile) string {
	return err.Fault.Render(profile)
}

// Type returns the error object type.
func (err *Error) Type() Type {
	return ERROR
}

// Method defines the set of methods available on error objects.
func (err *Error) Method(method string, tok token.Token, args []Object) (Object, bool) {
	return nil, false
}

// IsError determines if the referenced object is an error.
func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR
	}

	return false
}

// NewError reports a failure at the token the evaluator was looking at. The
// kind says what sort of failure it is and the message says what happened;
// neither should repeat the position, which the report adds for itself.
func NewError(kind fault.Kind, tok token.Token, format string, arguments ...interface{}) *Error {
	return &Error{Fault: fault.At(kind, tok, format, arguments...)}
}

// NewErrorFrom wraps a fault that was built elsewhere — by the parser, or by a
// boundary catching a Go panic — as an error object.
func NewErrorFrom(raised *fault.Fault) *Error {
	return &Error{Fault: raised}
}

// WithHelp attaches a suggested fix to the error, and returns it so that a
// failing branch stays a single expression.
func (err *Error) WithHelp(format string, arguments ...interface{}) *Error {
	err.Fault.WithHelp(format, arguments...)

	return err
}

// WithFrame records a call the error passed through on its way out.
func (err *Error) WithFrame(name string, tok token.Token) *Error {
	err.Fault.WithFrame(name, tok)

	return err
}
