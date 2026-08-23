package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateProperty(node *ast.Property, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)

	if isError(left) {
		return left
	}

	property, ok := node.Property.(*ast.Identifier)

	if !ok {
		return object.NewError(fault.Syntax, node.Token, "a property has to be named")
	}

	if left == nil {
		return object.NewError(fault.Type, node.Token, "cannot read property `%s` of a null value", property.Value)
	}

	switch left := left.(type) {
	case *object.Instance:
		return evaluateInstanceProperty(left, left.Class, property)
	case *object.Super:
		return evaluateInstanceProperty(left.Instance, left.Class, property)
	case *object.LibraryModule:
		if function, ok := left.Properties[property.Value]; ok {
			return unwrapCall(property.Token, function, nil, scope)
		}

		// A class is a module member like any other — reading it off the
		// module (rather than calling it, or `new`-ing it through a further
		// dotted call) just answers the class value itself, the way
		// `audio.Audio` has to for `new audio.Audio(...)` to work.
		if class, ok := left.Classes[property.Value]; ok {
			return class
		}

		return object.NewError(fault.Property, property.Token, "module `%s` has no property `%s`", left.Name, property.Value).
			WithHelp("%s", modulePropertySuggestion(left, property.Value))
	case *object.Class:
		return object.NewError(fault.Property, node.Token, "class `%s` has no property `%s` to read on the class itself", left.Name.Value, property.Value).
			WithHelp("properties are read on instances: `new %s().%s`", left.Name.Value, property.Value)
	case *object.NativeClass:
		return object.NewError(fault.Property, node.Token, "class `%s` has no property `%s` to read on the class itself", left.Name, property.Value).
			WithHelp("properties are read on instances: `new %s().%s`", left.Name, property.Value)
	case *object.Map:
		key := &object.String{Value: property.Value}

		pair, ok := left.Pairs[key.MapKey()]

		if !ok {
			return value.NULL
		}

		return pair.Value
	}

	return object.NewError(fault.Property, node.Token, "cannot read property `%s` of %s", property.Value, object.TypeName(left))
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
