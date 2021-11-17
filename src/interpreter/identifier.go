package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/library"
	"ghostlang.org/x/ghost/object"
)

func evaluateIdentifier(node *ast.Identifier, env *object.Environment) (object.Object, bool) {
	if libraryFunction, ok := library.Functions[node.Value]; ok {
		return libraryFunction, true
	}

	value, ok := env.Get(node.Value)

	if !ok {
		return nil, false
	}

	return value, true
}
