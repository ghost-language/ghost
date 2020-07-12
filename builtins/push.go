package builtins

import (
	"ghostlang.org/ghost/decimal"
	"ghostlang.org/ghost/object"
)

func Push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, expected=2", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	newElements := make([]object.Object, length+1, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	arr.Elements = newElements

	return &object.Number{Value: decimal.NewFromInt(int64(len(arr.Elements)))}
}
