package builtins

import (
	"ghostlang.org/ghost/object"
)

func Tail(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, expected=1", len(args))
	}

	if args[0].Type() != object.LIST_OBJ {
		return newError("argument to `tail` must be LIST, got %s", args[0].Type())
	}

	list := args[0].(*object.List)
	length := len(list.Elements)

	if length > 0 {
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, list.Elements[1:length])

		return &object.List{Elements: newElements}
	}

	return NULL
}
