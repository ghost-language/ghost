package builtins

import (
	"ghostlang.org/ghost/object"
)

func First(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)

	return arr.Elements[0]
}
