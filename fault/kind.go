package fault

// Kind names what sort of thing went wrong. It exists so that a reader can tell
// a typo from a bad argument from a missing file before reading the sentence
// that follows, and so that the same sort of problem is described the same way
// no matter which part of Ghost noticed it.
//
// Every failure Ghost reports is one of these. There is no general-purpose
// "runtime error": if a new failure does not obviously belong to one of these
// kinds, that is a sign the message has not been thought through yet.
type Kind int

const (
	// Syntax is source that could not be read: a stray character, an
	// unterminated string, a token where another was required.
	Syntax Kind = iota

	// Name is a reference to something that was never defined.
	Name

	// Type is an operation applied to the wrong sort of value.
	Type

	// Argument is a call whose arguments do not fit the thing being called.
	Argument

	// Index is a subscript outside the bounds of what it indexes, or a key that
	// cannot be used as one.
	Index

	// Value is a value of the right type that the operation cannot accept —
	// dividing by zero, a negative length, a singular matrix.
	Value

	// Property is a member that the value or class does not have.
	Property

	// Import is a module that could not be found, read, or that does not export
	// what was asked of it.
	Import

	// System is the world outside the program failing: a file that will not
	// open, a socket that will not bind, a plugin that will not load.
	System

	// Internal is a bug in Ghost itself. A reader seeing one of these has found
	// something worth reporting, and the message says so.
	Internal
)

var kindNames = [...]string{
	Syntax:   "syntax error",
	Name:     "name error",
	Type:     "type error",
	Argument: "argument error",
	Index:    "index error",
	Value:    "value error",
	Property: "property error",
	Import:   "import error",
	System:   "system error",
	Internal: "internal error",
}

// String returns the heading this kind is reported under.
func (kind Kind) String() string {
	if int(kind) < 0 || int(kind) >= len(kindNames) {
		return "error"
	}

	return kindNames[kind]
}
