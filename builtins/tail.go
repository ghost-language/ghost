package builtins

import (
	"ghostlang.org/ghost/object"
)

func Tail(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `tail` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	if length > 0 {
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, arr.Elements[1:length])

		return &object.Array{Elements: newElements}
	}

	return NULL
}
