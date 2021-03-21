package standard

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/version"
)

func init() {
	RegisterFunction("version", versionFunction)
}

func versionFunction(args ...[]object.Object) object.Object {
	v := version.Version

	return &object.String{Value: v}
}
