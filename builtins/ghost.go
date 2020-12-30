package builtins

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/version"
)

func init() {
	RegisterFunction("Ghost.version", ghostVersionFunction)
}

// ghostVersionFunction returns the current version of Ghost.
func ghostVersionFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	v := version.Version

	return &object.String{Value: v}
}
