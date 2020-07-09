package builtins

import (
	"unicode/utf8"

	"ghostlang.org/ghost/object"
)

func Len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, expected=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
	default:
		return newError("argument to `len` not supported, got %s", args[0].Type())
	}
}
