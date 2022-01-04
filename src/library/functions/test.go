package functions

import (
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func Test(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 2 {
		// TODO: error
		return nil
	}

	description := args[0].(*object.String).Value
	callback := args[1].(*object.Function)

	result := callback.Evaluate(nil, nil)

	if object.IsError(result) {
		log.Error("%s failed.", description)

		return result
	}

	if object.IsFalse(result) {
		return object.NewError("%s failed.", description)
	}

	log.Info("%s passed.", description)

	return value.TRUE
}
