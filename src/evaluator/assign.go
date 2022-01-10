package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateAssign(node *ast.Assign, scope *object.Scope) object.Object {
	value := Evaluate(node.Value, scope)

	if isError(value) {
		return value
	}

	switch assignment := node.Name.(type) {
	case *ast.Identifier:
		return evaluateIdentifierAssignment(assignment, value, scope)
	case *ast.Index:
		return evaluateIndexAssignment(assignment, value, scope)
	case *ast.Property:
		return evaluatePropertyAssignment(assignment, value, scope)
	}

	object.NewError("%d:%d: runtime error: cannot assign variable to a %T", node.Token.Line, node.Token.Column, node.Name)

	return nil
}

func evaluateIdentifierAssignment(node *ast.Identifier, value object.Object, scope *object.Scope) object.Object {
	scope.Environment.Set(node.Value, value)

	return nil
}

func evaluateIndexAssignment(node *ast.Index, assignmentValue object.Object, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)
	index := Evaluate(node.Index, scope)

	switch obj := left.(type) {
	case *object.List:
		idx := int(index.(*object.Number).Value.IntPart())
		elements := obj.Elements

		if idx < 0 {
			return object.NewError("%d:%d: runtime error: index out of range: %d", node.Token.Line, node.Token.Column, idx)
		}

		if idx >= len(elements) {
			for i := len(elements); i <= idx; i++ {
				elements = append(elements, value.NULL)
			}

			obj.Elements = elements
		}

		elements[idx] = assignmentValue
	case *object.Map:
		key, ok := index.(object.Mappable)

		if !ok {
			return object.NewError("%d:%d: runtime error: unusable as a map key: %s", node.Token.Line, node.Token.Column, index.Type())
		}

		hashed := key.MapKey()
		pair := object.MapPair{Key: index, Value: assignmentValue}
		obj.Pairs[hashed] = pair
	}

	return nil
}

func evaluatePropertyAssignment(node *ast.Property, assignmentValue object.Object, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)

	switch obj := left.(type) {
	case *object.Map:
		key := &object.String{Value: node.Property.(*ast.Identifier).Value}
		hashed := key.MapKey()
		pair := object.MapPair{Key: key, Value: assignmentValue}
		obj.Pairs[hashed] = pair

		return nil
	}

	return object.NewError("%d:%d: runtime error: can only assign properties to maps, got %s", node.Token.Line, node.Token.Column, left.Type())
}
