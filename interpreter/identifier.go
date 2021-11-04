package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/library"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateIdentifier(node *ast.Identifier) (object.Object, bool) {
	if libraryFunction, ok := library.Functions[node.Value]; ok {
		return libraryFunction, true
	}

	return value.TRUE, true
}
