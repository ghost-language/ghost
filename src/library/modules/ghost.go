package modules

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"ghostlang.org/x/ghost/version"
)

var Ghost = map[string]*object.LibraryFunction{}

func init() {
	RegisterMethod(Ghost, "abort", ghostAbort)
	RegisterMethod(Ghost, "execute", ghostExecute)
	RegisterMethod(Ghost, "identifiers", ghostIdentifiers)
	RegisterMethod(Ghost, "version", ghostVersion)
}

func ghostAbort(env *object.Environment, args ...object.Object) object.Object {
	switch obj := args[0].(type) {
	case *object.Null:
		return nil
	case *object.String:
		return object.NewError(obj.Value)
	}

	return object.NewError("ghost.abort must contain either a null or string value")
}

func ghostExecute(env *object.Environment, args ...object.Object) object.Object {
	scanner := scanner.New(args[0].(*object.String).Value)
	parser := parser.New(scanner)
	program := parser.Parse()

	return evaluate(program, env)
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
