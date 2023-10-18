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

		if !instance.Environment.Has(property.Value) {
			instance.Environment.Set(property.Value, value.NULL)
		}

		val, _ := instance.Environment.Get(property.Value)

		return val
	case *object.LibraryModule:
		property := node.Property.(*ast.Identifier)
		module := left.(*object.LibraryModule)

		if function, ok := module.Properties[property.Value]; ok {
			return unwrapCall(node.Token, function, nil, scope)
		}

		return newError("%d:%d:%s: runtime error: unknown property: %s.%s", node.Token.Line, node.Token.Column, node.Token.File, module.Name, property.Value)
	case *object.Map:
		property := &object.String{Value: node.Property.(*ast.Identifier).Value}
		mapObj := left.(*object.Map)

		pair, ok := mapObj.Pairs[property.MapKey()]

		if !ok {
			return value.NULL
		}

		return pair.Value
	}

	return nil
}
