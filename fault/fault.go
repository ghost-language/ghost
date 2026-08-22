// Package fault is how Ghost says that something went wrong.
//
// A fault is a structured value, not a formatted string: it knows what kind of
// failure it is, where in the source it happened, how wide the offending
// lexeme is, what would fix it, and which calls were in flight at the time.
// Keeping those apart is what lets one renderer produce every message Ghost
// prints, so a type error from the evaluator and an argument error from a
// library module are laid out identically and colored identically.
//
// Nothing constructs a message by hand. Whatever notices the problem describes
// it, and this package decides how it reads.
package fault

import (
	"fmt"

	"ghostlang.org/x/ghost/token"
)

// Position locates a fault in a source file. Length is how many characters the
// offending lexeme covers, so a report can underline the whole of it rather
// than pointing at its first character.
type Position struct {
	File   string
	Line   int
	Column int
	Length int
}

// Known reports whether a position is specific enough to quote source for.
func (position Position) Known() bool {
	return position.File != "" && position.Line > 0 && position.Column > 0
}

// String renders a position the way editors and build tools expect to read one.
func (position Position) String() string {
	if !position.Known() {
		return ""
	}

	return fmt.Sprintf("%s:%d:%d", position.File, position.Line, position.Column)
}

// PositionOf reads the position of a token, including how wide it is.
func PositionOf(tok token.Token) Position {
	length := len([]rune(tok.Lexeme))

	if tok.Length > 0 {
		length = tok.Length
	}

	if length < 1 {
		length = 1
	}

	return Position{File: tok.File, Line: tok.Line, Column: tok.Column, Length: length}
}

// maxFrames bounds how many calls a trace records. A runaway recursion has
// thousands of identical frames and none of them says anything the first few do
// not; the count of what was left out does.
const maxFrames = 8

// Frame is one call that was in progress when a fault was raised. Ghost records
// them only as an error unwinds, so an ordinary call costs nothing.
type Frame struct {
	Name     string
	Position Position
}

// Fault is a single failure, described rather than formatted.
type Fault struct {
	Kind     Kind
	Message  string
	Help     string
	Position Position
	Trace    []Frame

	// Hidden counts the frames left out of the trace once it filled up.
	Hidden int
}

// New builds a fault with no position. Use it only where there genuinely is no
// source to point at — a failure in an embedder's own call into Ghost, say.
func New(kind Kind, format string, arguments ...interface{}) *Fault {
	return &Fault{Kind: kind, Message: fmt.Sprintf(format, arguments...)}
}

// At builds a fault that points at a token. This is the constructor almost
// everything should use: the token is what the scanner, parser, or evaluator
// was looking at when it gave up, and it carries the file, line, column, and
// width needed to quote the source back.
func At(kind Kind, tok token.Token, format string, arguments ...interface{}) *Fault {
	return &Fault{
		Kind:     kind,
		Message:  fmt.Sprintf(format, arguments...),
		Position: PositionOf(tok),
	}
}

// From builds a fault at an already-computed position.
func From(kind Kind, position Position, format string, arguments ...interface{}) *Fault {
	return &Fault{Kind: kind, Message: fmt.Sprintf(format, arguments...), Position: position}
}

// WithHelp attaches a suggested fix. Help is for what the reader should do
// next; it is not a second sentence of explanation, and a fault with nothing
// useful to suggest should not have one.
func (fault *Fault) WithHelp(format string, arguments ...interface{}) *Fault {
	fault.Help = fmt.Sprintf(format, arguments...)

	return fault
}

// WithFrame records a call the fault passed through on its way out. Frames are
// added as the error unwinds, so they arrive innermost first, which is the
// order they are printed in.
func (fault *Fault) WithFrame(name string, tok token.Token) *Fault {
	if len(fault.Trace) >= maxFrames {
		fault.Hidden++

		return fault
	}

	fault.Trace = append(fault.Trace, Frame{Name: name, Position: PositionOf(tok)})

	return fault
}

// Error makes a fault usable as a Go error, so the boundaries where Ghost is
// embedded in a Go program can hand one back without unwrapping it first.
func (fault *Fault) Error() string {
	return fault.String()
}

// String renders the fault on a single line, in the file:line:column form that
// editors and tooling know how to jump to. It is what a fault reduces to when
// there is no room to lay it out properly.
func (fault *Fault) String() string {
	if fault == nil {
		return ""
	}

	if !fault.Position.Known() {
		return fmt.Sprintf("%s: %s", fault.Kind, fault.Message)
	}

	return fmt.Sprintf("%s: %s: %s", fault.Position, fault.Kind, fault.Message)
}
