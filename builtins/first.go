package builtins

import (
	"ghostlang.org/ghost/object"
)

func First(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.LIST_OBJ {
		return newError("argument to `first` must be LIST, got %s", args[0].Type())
	}

	list := args[0].(*object.List)

	return list.Elements[0]
}
