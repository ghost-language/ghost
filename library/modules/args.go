package modules

import (
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// Modules convert Ghost values into Go values constantly, and doing it inline
// turns every method into a wall of type assertions. These helpers keep the
// conversions in one place so that argument errors read the same no matter
// which module reported them.
//
// They are thin wrappers over the checks in the object package rather than a
// second set of rules. A module method, a method on a string, and a call into a
// class all describe a bad argument with the same sentence, because all three
// of them end up in the same function.

// arity checks that a method received an exact number of arguments.
func arity(name string, tok token.Token, args []object.Object, want int) *object.Error {
	return object.Arity(signature(name), tok, args, want)
}

// arityRange checks that a method received a number of arguments within bounds.
func arityRange(name string, tok token.Token, args []object.Object, low, high int) *object.Error {
	return object.ArityRange(signature(name), tok, args, low, high)
}

// arityAtLeast checks that a method received no fewer than a number of
// arguments, which suits the methods that read a whole run of values.
func arityAtLeast(name string, tok token.Token, args []object.Object, low int) *object.Error {
	return object.ArityAtLeast(signature(name), tok, args, low)
}

// numberAt reads a numeric argument.
func numberAt(name string, tok token.Token, args []object.Object, index int) (*object.Number, *object.Error) {
	return object.NumberArgument(signature(name), tok, args, index)
}

// floatAt reads a numeric argument as a float.
func floatAt(name string, tok token.Token, args []object.Object, index int) (float64, *object.Error) {
	number, err := numberAt(name, tok, args, index)

	if err != nil {
		return 0, err
	}

	return number.Float64(), nil
}

// integerAt reads a numeric argument as a whole number.
func integerAt(name string, tok token.Token, args []object.Object, index int) (int64, *object.Error) {
	number, err := numberAt(name, tok, args, index)

	if err != nil {
		return 0, err
	}

	return number.Int64(), nil
}

// booleanAt reads a boolean argument.
func booleanAt(name string, tok token.Token, args []object.Object, index int) (bool, *object.Error) {
	boolean, err := object.BooleanArgument(signature(name), tok, args, index)

	if err != nil {
		return false, err
	}

	return boolean.Value, nil
}

// listAt reads a list argument.
func listAt(name string, tok token.Token, args []object.Object, index int) (*object.List, *object.Error) {
	return object.ListArgument(signature(name), tok, args, index)
}

// stringAt reads a string argument.
func stringAt(name string, tok token.Token, args []object.Object, index int) (string, *object.Error) {
	str, err := object.StringArgument(signature(name), tok, args, index)

	if err != nil {
		return "", err
	}

	return str.Value, nil
}

// systemFailure reports the world outside the program refusing to cooperate: a
// file that will not open, a socket that will not bind. The Go error is quoted
// as-is, because what it says — "no such file or directory", "permission
// denied" — is already the clearest description available.
func systemFailure(name string, tok token.Token, failure error) *object.Error {
	return object.NewError(fault.System, tok, "`%s` failed: %s", signature(name), failure)
}

// functionAt reads a callable argument.
func functionAt(name string, tok token.Token, args []object.Object, index int) (*object.Function, *object.Error) {
	return object.FunctionArgument(signature(name), tok, args, index)
}

// signature renders a module method's name the way it is written in Ghost, so
// that a message quotes something the reader can search their own source for.
func signature(name string) string {
	return name + "()"
}

// gatherNumbers reads a run of arguments as a flat sequence of numbers. Lists
// are flattened to any depth, so a method written against it accepts its values
// spread across the call, collected in a list, or arranged as a matrix, without
// having to spell out the three cases itself.
func gatherNumbers(name string, tok token.Token, args []object.Object) ([]*object.Number, *object.Error) {
	numbers := make([]*object.Number, 0, len(args))

	for _, arg := range args {
		collected, err := appendNumbers(name, tok, numbers, arg)

		if err != nil {
			return nil, err
		}

		numbers = collected
	}

	return numbers, nil
}

// gatherFloats is gatherNumbers for the methods that have no reason to preserve
// whole numbers, such as the statistics whose results are always fractional.
func gatherFloats(name string, tok token.Token, args []object.Object) ([]float64, *object.Error) {
	numbers, err := gatherNumbers(name, tok, args)

	if err != nil {
		return nil, err
	}

	return toFloats(numbers), nil
}

func appendNumbers(name string, tok token.Token, numbers []*object.Number, arg object.Object) ([]*object.Number, *object.Error) {
	switch arg := arg.(type) {
	case *object.Number:
		return append(numbers, arg), nil
	case *object.List:
		for _, element := range arg.Elements {
			collected, err := appendNumbers(name, tok, numbers, element)

			if err != nil {
				return nil, err
			}

			numbers = collected
		}

		return numbers, nil
	}

	return nil, object.NewError(fault.Argument, tok, "`%s` expects numbers or lists of numbers, got %s", signature(name), object.TypeName(arg))
}

// toFloats unwraps numbers for the Go operations that work in float64.
func toFloats(numbers []*object.Number) []float64 {
	values := make([]float64, len(numbers))

	for index, number := range numbers {
		values[index] = number.Float64()
	}

	return values
}

// floatList builds a Ghost list from Go floats.
func floatList(values []float64) *object.List {
	elements := make([]object.Object, len(values))

	for index, value := range values {
		elements[index] = object.NewFloat(value)
	}

	return &object.List{Elements: elements}
}

// integerList builds a Ghost list from whole numbers.
func integerList(values []int64) *object.List {
	elements := make([]object.Object, len(values))

	for index, value := range values {
		elements[index] = object.NewInt(value)
	}

	return &object.List{Elements: elements}
}

// numberList builds a Ghost list from numbers that have already been boxed,
// which keeps whole numbers whole through methods that only reorder values.
func numberList(numbers []*object.Number) *object.List {
	elements := make([]object.Object, len(numbers))

	for index, number := range numbers {
		elements[index] = number
	}

	return &object.List{Elements: elements}
}
