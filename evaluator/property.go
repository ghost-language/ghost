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
		return evaluateInstanceProperty(left, node)
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

func evaluateInstanceProperty(left object.Object, node *ast.Property) object.Object {
	var val object.Object

	instance := left.(*object.Instance)
	property := node.Property.(*ast.Identifier)

	if instance.Environment.Has(property.Value) {
		val, _ = instance.Environment.Get(property.Value)

		return val
	}

	if instance.Class.Environment.Has(property.Value) {
		val, _ = instance.Class.Environment.Get(property.Value)

		return val
	}

	for _, trait := range instance.Class.Traits {
		if trait.Environment.Has(property.Value) {
			val, _ = trait.Environment.Get(property.Value)

			return val
		}
	}

	instance.Environment.Set(property.Value, value.NULL)

	val, _ = instance.Environment.Get(property.Value)

	return val
}
