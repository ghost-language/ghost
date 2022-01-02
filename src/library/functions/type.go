package functions

import (
	"strings"

	"ghostlang.org/x/ghost/object"
)

func Type(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		// TODO: error
		return nil
	}

	objectType := string(args[0].Type())

	return &object.String{Value: strings.ToLower(objectType)}
}
