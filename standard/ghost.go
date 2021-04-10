package standard

import (
	"fmt"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/version"
)

func init() {
	RegisterFunction("version", versionFunction)
	RegisterFunction("dd", ddFunction)
}

func versionFunction(args []object.Object) object.Object {
	v := version.Version

	return &object.String{Value: v}
}

func ddFunction(args []object.Object) object.Object {
	panic(fmt.Sprintf("%+v", args))
}
