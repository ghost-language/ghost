package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/library"
	"ghostlang.org/x/ghost/object"
)

func evaluateIdentifier(node *ast.Identifier, env *object.Environment) (object.Object, bool) {
	if libraryModule, ok := library.Modules[node.Value]; ok {
		return libraryModule, true
	}

	if libraryFunction, ok := library.Functions[node.Value]; ok {
		return libraryFunction, true
	}

	value, ok := env.Get(node.Value)

	if !ok {
		err := newError("unkown identifier: %s", node.Value)

		return err, false
	}

	return value, true
}
