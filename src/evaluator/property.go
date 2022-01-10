package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateProperty(node *ast.Property, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)

	if isError(left) {
		return left
	}

	switch left.(type) {
	case *object.Instance:
		property := node.Property.(*ast.Identifier)
		instance := left.(*object.Instance)

		if value, ok := instance.Class.Environment.Get(property.Value); ok {
			return value
		}
	case *object.LibraryModule:
		property := node.Property.(*ast.Identifier)
		module := left.(*object.LibraryModule)

		if function, ok := module.Properties[property.Value]; ok {
			return unwrapCall(node.Token, function, nil, scope)
		}
	case *object.Map:
		property := &object.String{Value: node.Property.(*ast.Identifier).Value}
		mapObj := left.(*object.Map)

		pair, ok := mapObj.Pairs[property.MapKey()]

		if !ok {
			return value.NULL
		}

		return pair.Value
	}

	return newError("%d:%d: runtime error: unknown property: %s.%s", node.Token.Line, node.Token.Column, left.String(), node.Property.(*ast.Identifier).Value)
}
