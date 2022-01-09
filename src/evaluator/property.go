package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateProperty(node *ast.Property, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)

	if isError(left) {
		return left
	}

	switch left.(type) {
	case *object.LibraryModule:
		property := node.Property.(*ast.Identifier)
		module := left.(*object.LibraryModule)

		if function, ok := module.Properties[property.Value]; ok {
			return unwrapCall(node.Token, function, nil, scope)
		}
	}

	return newError("%d:%d: runtime error: unknown property: %s.%s", node.Token.Line, node.Token.Column, left.String(), node.Property.(*ast.Identifier).Value)
}
