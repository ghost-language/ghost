package modules

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/version"
)

var Ghost = map[string]*object.LibraryFunction{}

func init() {
	RegisterMethod(Ghost, "identifiers", ghostIdentifiers)
	RegisterMethod(Ghost, "version", ghostVersion)
}

func ghostIdentifiers(env *object.Environment, args ...object.Object) object.Object {
	identifiers := []object.Object{}

	store := env.All()

	for identifier := range store {
		identifiers = append(identifiers, &object.String{Value: identifier})
	}

	return &object.List{Elements: identifiers}
}

func ghostVersion(env *object.Environment, args ...object.Object) object.Object {
	return &object.String{Value: version.Version}
}
