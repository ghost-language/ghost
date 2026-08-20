package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateAssign(node *ast.Assign, scope *object.Scope) object.Object {
	// A bare assignment in a class or trait body declares a field. The
	// initializer is recorded unevaluated and re-evaluated for each instance,
	// so instances do not share one copy of a mutable value.
	if identifier, ok := node.Name.(*ast.Identifier); ok {
		if declaration, ok := scope.Self.(object.FieldDeclarer); ok {
			if identifier.Value == constructorName {
				return constructorFieldError(node)
			}

			declaration.SetField(identifier.Value, node.Value)

			return nil
		}
	}

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

	return object.NewError("%d:%d:%s: runtime error: cannot assign variable to a %T", node.Token.Line, node.Token.Column, node.Token.File, node.Name)
}

// constructorFieldError rejects `constructor = ...` in a class body. A field by
// that name would never be run as a constructor, so silently accepting it hides
// the mistake.
func constructorFieldError(node *ast.Assign) object.Object {
	return object.NewError("%d:%d:%s: runtime error: '%s' must be declared as a method, not a field", node.Token.Line, node.Token.Column, node.Token.File, constructorName)
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
		numIdx, ok := index.(*object.Number)
		if !ok {
			return object.NewError("%d:%d:%s: runtime error: list index must be a number, got %s", node.Token.Line, node.Token.Column, node.Token.File, index.Type())
		}
		idx := int(numIdx.Int64())
		elements := obj.Elements

		if idx < 0 {
			return object.NewError("%d:%d:%s: runtime error: index out of range: %d", node.Token.Line, node.Token.Column, node.Token.File, idx)
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
			return object.NewError("%d:%d:%s: runtime error: unusable as a map key: %s", node.Token.Line, node.Token.Column, node.Token.File, index.Type())
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
	case *object.Instance:
		obj.Environment.Set(node.Property.(*ast.Identifier).Value, assignmentValue)

		return nil
	case *object.Map:
		key := &object.String{Value: node.Property.(*ast.Identifier).Value}
		hashed := key.MapKey()
		pair := object.MapPair{Key: key, Value: assignmentValue}
		obj.Pairs[hashed] = pair

		return nil
	}

	return object.NewError("%d:%d:%s: runtime error: can only assign properties to maps, got %s", node.Token.Line, node.Token.Column, node.Token.File, left.Type())
}
