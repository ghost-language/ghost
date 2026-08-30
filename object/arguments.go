package object

import (
	"fmt"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/token"
)

// Ghost has three kinds of callable — methods on values, functions in the
// library, and modules — and every one of them has to answer the same two
// questions about a call: were there enough arguments, and were they the right
// sort of thing. These helpers answer them in one place so the answer reads the
// same wherever it comes from. A reader who has seen
//
//	argument error: `list.join()` expects argument 1 to be a string, got number
//
// once should be able to predict the wording of every other argument error in
// the language.

// Arity checks that a call received an exact number of arguments.
func Arity(name string, tok token.Token, args []Object, want int) *Error {
	if len(args) == want {
		return nil
	}

	return NewError(fault.Argument, tok, "`%s` expects %s, got %d", name, plural(want), len(args))
}

// ArityRange checks that a call received a number of arguments within bounds,
// which is what an optional trailing argument amounts to.
func ArityRange(name string, tok token.Token, args []Object, low int, high int) *Error {
	if len(args) >= low && len(args) <= high {
		return nil
	}

	return NewError(fault.Argument, tok, "`%s` expects between %d and %d arguments, got %d", name, low, high, len(args))
}

// ArityAtLeast checks that a call received no fewer than a number of arguments,
// which suits the calls that read a whole run of values.
func ArityAtLeast(name string, tok token.Token, args []Object, low int) *Error {
	if len(args) >= low {
		return nil
	}

	return NewError(fault.Argument, tok, "`%s` expects at least %s, got %d", name, plural(low), len(args))
}

// NumberArgument reads an argument that has to be a number.
func NumberArgument(name string, tok token.Token, args []Object, index int) (*Number, *Error) {
	value, err := argument(name, tok, args, index, "a number", NUMBER)

	if err != nil {
		return nil, err
	}

	return value.(*Number), nil
}

// StringArgument reads an argument that has to be a string.
func StringArgument(name string, tok token.Token, args []Object, index int) (*String, *Error) {
	value, err := argument(name, tok, args, index, "a string", STRING)

	if err != nil {
		return nil, err
	}

	return value.(*String), nil
}

// BooleanArgument reads an argument that has to be a boolean.
func BooleanArgument(name string, tok token.Token, args []Object, index int) (*Boolean, *Error) {
	value, err := argument(name, tok, args, index, "a boolean", BOOLEAN)

	if err != nil {
		return nil, err
	}

	return value.(*Boolean), nil
}

// ListArgument reads an argument that has to be a list.
func ListArgument(name string, tok token.Token, args []Object, index int) (*List, *Error) {
	value, err := argument(name, tok, args, index, "a list", LIST)

	if err != nil {
		return nil, err
	}

	return value.(*List), nil
}

// MapArgument reads an argument that has to be a map.
func MapArgument(name string, tok token.Token, args []Object, index int) (*Map, *Error) {
	value, err := argument(name, tok, args, index, "a map", MAP)

	if err != nil {
		return nil, err
	}

	return value.(*Map), nil
}

// DateArgument reads an argument that has to be a date.
func DateArgument(name string, tok token.Token, args []Object, index int) (*Date, *Error) {
	value, err := argument(name, tok, args, index, "a date", DATE)

	if err != nil {
		return nil, err
	}

	return value.(*Date), nil
}

// DurationArgument reads an argument that has to be a duration.
func DurationArgument(name string, tok token.Token, args []Object, index int) (*Duration, *Error) {
	value, err := argument(name, tok, args, index, "a duration", DURATION)

	if err != nil {
		return nil, err
	}

	return value.(*Duration), nil
}

// FunctionArgument reads an argument that has to be something callable.
func FunctionArgument(name string, tok token.Token, args []Object, index int) (*Function, *Error) {
	value, err := argument(name, tok, args, index, "a function", FUNCTION)

	if err != nil {
		return nil, err
	}

	return value.(*Function), nil
}

// AnyArgument reads an argument of any type, checking only that it is there.
func AnyArgument(name string, tok token.Token, args []Object, index int) (Object, *Error) {
	if index >= len(args) {
		return nil, missing(name, tok, index)
	}

	return args[index], nil
}

// argument reads and type-checks one argument. Everything above is a thin
// wrapper over it, which is what keeps the two messages it can produce the only
// two messages in the language for a bad argument.
func argument(name string, tok token.Token, args []Object, index int, article string, want Type) (Object, *Error) {
	if index >= len(args) {
		return nil, missing(name, tok, index)
	}

	value := args[index]

	if value == nil || value.Type() != want {
		return nil, NewError(fault.Argument, tok, "`%s` expects argument %d to be %s, got %s", name, index+1, article, TypeName(value))
	}

	return value, nil
}

// missing reports an argument that was never passed. It is kept apart from a
// wrong-typed one because the fix is different: one is a value to change, the
// other is a value to add.
func missing(name string, tok token.Token, index int) *Error {
	return NewError(fault.Argument, tok, "`%s` is missing argument %d", name, index+1)
}

// plural writes an argument count as the phrase it belongs in, so a message
// says "1 argument" rather than "1 arguments".
func plural(count int) string {
	if count == 1 {
		return "1 argument"
	}

	return fmt.Sprintf("%d arguments", count)
}
