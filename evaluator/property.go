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

	property := node.Property.(*ast.Identifier)

	if left == nil {
		return newError("%d:%d:%s: runtime error: cannot read property %s of a null value", node.Token.Line, node.Token.Column, node.Token.File, property.Value)
	}

	switch left := left.(type) {
	case *object.Instance:
		return evaluateInstanceProperty(left, left.Class, property)
	case *object.Super:
		return evaluateInstanceProperty(left.Instance, left.Class, property)
	case *object.LibraryModule:
		if function, ok := left.Properties[property.Value]; ok {
			return unwrapCall(node.Token, function, nil, scope)
		}

		return newError("%d:%d:%s: runtime error: unknown property: %s.%s", node.Token.Line, node.Token.Column, node.Token.File, left.Name, property.Value)
	case *object.Map:
		key := &object.String{Value: property.Value}

		pair, ok := left.Pairs[key.MapKey()]

		if !ok {
			return value.NULL
		}

		return pair.Value
	}

	return newError("%d:%d:%s: runtime error: cannot read property %s of %s", node.Token.Line, node.Token.Column, node.Token.File, property.Value, left.Type())
}

// evaluateInstanceProperty reads a property off an instance: its own fields
// first, then members declared on the class chain starting at `start`. Only the
// instance's own bindings are consulted, never the enclosing lexical scope, so
// a global of the same name cannot pose as a property.
func evaluateInstanceProperty(instance *object.Instance, start *object.Class, property *ast.Identifier) object.Object {
	if start == instance.Class {
		if val, ok := instance.Environment.GetLocal(property.Value); ok {
			return val
		}
	}

	if member, _, ok := object.LookupMember(start, property.Value); ok {
		return member
	}

	return value.NULL
}
