package builtins

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/utilities"
)

func init() {
	RegisterFunction("math.abs", mathAbsFunction)
}

func mathAbsFunction(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return utilities.NewError("wrong number of arguments. got=%d, expected=1",
			len(args))
	}

	if args[0].Type() != object.NUMBER_OBJ {
		return utilities.NewError("argument to `math.abs` must be NUMBER, got %s", args[0].Type())
	}

	number := args[0].(*object.Number)

	return &object.Number{Value: number.Value.Abs()}
}
