package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/library"
	"ghostlang.org/x/ghost/object"
)

func evaluateIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if libraryModule, ok := library.Modules[node.Value]; ok {
		return libraryModule
	}

	if libraryFunction, ok := library.Functions[node.Value]; ok {
		return libraryFunction
	}

	if identifier, ok := env.Get(node.Value); ok {
		return identifier
	}

	return newError("%d:__: runtime error: unknown identifier: %s", node.Token.Line, node.Value)
}
