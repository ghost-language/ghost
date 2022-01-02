package modules

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/version"
)

var Ghost = map[string]*object.LibraryFunction{}

func init() {
	RegisterMethod(Ghost, "version", ghostVersion)
}

func ghostVersion(env *object.Environment, args ...object.Object) object.Object {
	return &object.String{Value: version.Version}
}
