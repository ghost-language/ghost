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
	case *ast.ListPattern:
		return evaluateListPatternAssignment(assignment, value, scope)
	case *ast.MapPattern:
		return evaluateMapPatternAssignment(assignment, value, scope)
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

// evaluateListPatternAssignment binds each target to the value at its
// position, or to null past the end of a shorter list - the same leniency
// `list[i]` itself already has for an out-of-range read (§13.6).
func evaluateListPatternAssignment(node *ast.ListPattern, assignmentValue object.Object, scope *object.Scope) object.Object {
	list, ok := assignmentValue.(*object.List)

	if !ok {
		return object.NewError(fault.Type, node.Token, "cannot destructure %s as a list", object.TypeName(assignmentValue)).
			WithHelp("a list pattern needs a list on the right of `=`")
	}

	for index, target := range node.Targets {
		if index < len(list.Elements) {
			scope.Environment.Set(target.Value, list.Elements[index])
		} else {
			scope.Environment.Set(target.Value, value.NULL)
		}
	}

	return nil
}

// evaluateMapPatternAssignment binds each target to the map's value under
// its source key, or to null for a missing key - the same leniency a map
// read (`map.key`, `map["key"]`) already has.
func evaluateMapPatternAssignment(node *ast.MapPattern, assignmentValue object.Object, scope *object.Scope) object.Object {
	mapValue, ok := assignmentValue.(*object.Map)

	if !ok {
		return object.NewError(fault.Type, node.Token, "cannot destructure %s as a map", object.TypeName(assignmentValue)).
			WithHelp("a map pattern needs a map on the right of `=`")
	}

	for _, pair := range node.Pairs {
		key := &object.String{Value: pair.Source.Value}

		if entry, ok := mapValue.Pairs[key.MapKey()]; ok {
			scope.Environment.Set(pair.Target.Value, entry.Value)
		} else {
			scope.Environment.Set(pair.Target.Value, value.NULL)
		}
	}

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
		obj.SetPair(hashed, object.MapPair{Key: index, Value: assignmentValue})
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
		obj.SetPair(key.MapKey(), object.MapPair{Key: key, Value: assignmentValue})

		return nil
	}

	return object.NewError(fault.Type, node.Token, "cannot set a property on %s", object.TypeName(left)).
		WithHelp("properties can be set on class instances and on maps")
}
