package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
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

	return object.NewError(fault.Syntax, node.Token, "cannot assign to this expression").
		WithHelp("only a variable, a list or map entry, or a property can be assigned to")
}

// constructorFieldError rejects `constructor = ...` in a class body. A field by
// that name would never be run as a constructor, so silently accepting it hides
// the mistake.
func constructorFieldError(node *ast.Assign) object.Object {
	return object.NewError(fault.Syntax, node.Token, "`%s` has to be declared as a method, not a field", constructorName).
		WithHelp("write `%s(...) { ... }` rather than `%s = ...`", constructorName, constructorName)
}

func evaluateIdentifierAssignment(node *ast.Identifier, value object.Object, scope *object.Scope) object.Object {
	scope.Environment.Set(node.Value, value)

	return nil
}

func evaluateIndexAssignment(node *ast.Index, assignmentValue object.Object, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)

	if isError(left) {
		return left
	}

	index := Evaluate(node.Index, scope)

	if isError(index) {
		return index
	}

	if left == nil || index == nil {
		return object.NewError(fault.Type, node.Token, "cannot assign into a null value")
	}

	switch obj := left.(type) {
	case *object.List:
		numIdx, ok := index.(*object.Number)
		if !ok {
			return object.NewError(fault.Type, node.Token, "a list index has to be a number, got %s", object.TypeName(index))
		}
		idx := int(numIdx.Int64())
		elements := obj.Elements

		if idx < 0 {
			return object.NewError(fault.Index, node.Token, "cannot assign to index %d, which is before the start of the list", idx)
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
			return object.NewError(fault.Type, node.Token, "%s cannot be used as a map key", object.TypeName(index)).
				WithHelp("a map key has to be a string, a number, or a boolean")
		}

		hashed := key.MapKey()
		pair := object.MapPair{Key: index, Value: assignmentValue}
		obj.Pairs[hashed] = pair
	default:
		return object.NewError(fault.Type, node.Token, "cannot assign into %s", object.TypeName(left)).
			WithHelp("only lists and maps can be assigned into by index")
	}

	return nil
}

func evaluatePropertyAssignment(node *ast.Property, assignmentValue object.Object, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)

	if isError(left) {
		return left
	}

	property, ok := node.Property.(*ast.Identifier)

	if !ok {
		return object.NewError(fault.Syntax, node.Token, "a property has to be named")
	}

	if left == nil {
		return object.NewError(fault.Type, node.Token, "cannot set property `%s` on a null value", property.Value)
	}

	switch obj := left.(type) {
	case *object.Instance:
		obj.Environment.Set(property.Value, assignmentValue)

		return nil
	case *object.Map:
		key := &object.String{Value: property.Value}
		hashed := key.MapKey()
		pair := object.MapPair{Key: key, Value: assignmentValue}
		obj.Pairs[hashed] = pair

		return nil
	}

	return object.NewError(fault.Type, node.Token, "cannot set a property on %s", object.TypeName(left)).
		WithHelp("properties can be set on class instances and on maps")
}
