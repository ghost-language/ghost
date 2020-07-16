package builtins

import (
	"ghostlang.org/ghost/object"
)

func Last(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.LIST_OBJ {
		return newError("argument to `last` must be LIST, got %s", args[0].Type())
	}

	list := args[0].(*object.List)
	length := len(list.Elements)

	return list.Elements[length-1]
}
