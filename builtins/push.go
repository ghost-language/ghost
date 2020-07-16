package builtins

import (
	"ghostlang.org/ghost/decimal"
	"ghostlang.org/ghost/object"
)

func Push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, expected=2", len(args))
	}

	if args[0].Type() != object.LIST_OBJ {
		return newError("argument to `push` must be LIST, got %s", args[0].Type())
	}

	list := args[0].(*object.List)
	length := len(list.Elements)

	newElements := make([]object.Object, length+1, length+1)
	copy(newElements, list.Elements)
	newElements[length] = args[1]

	list.Elements = newElements

	return &object.Number{Value: decimal.NewFromInt(int64(len(list.Elements)))}
}
