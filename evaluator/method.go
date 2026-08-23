package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
)

func evaluateMethod(node *ast.Method, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)

	if isError(left) {
		return left
	}

	name, ok := node.Method.(*ast.Identifier)

	if !ok {
		return object.NewError(fault.Syntax, node.Token, "a method call needs a method name")
	}

	// Errors about the method point at the method, not at the dot before it.
	at := name.Token

	if left == nil {
		return object.NewError(fault.Type, at, "cannot call `%s` on a null value", name.Value)
	}

	arguments := evaluateExpressions(node.Arguments, scope)

	if len(arguments) == 1 && isError(arguments[0]) {
		return arguments[0]
	}

	result, handled := left.Method(name.Value, at, arguments)

	if isError(result) {
		return result
	}

	switch receiver := left.(type) {
	case *object.Map:
		property := &object.String{Value: name.Value}

		if function, ok := receiver.Pairs[property.MapKey()]; ok {
			return unwrapCall(at, function.Value, arguments, scope)
		}

		if handled {
			return result
		}

		return object.NewError(fault.Property, at, "a map has no method `%s`", name.Value)
	case *object.Instance:
		return unwrapReturn(callInstanceMethod(node, receiver, receiver.Class, name.Value, arguments, scope))
	case *object.Super:
		return unwrapReturn(callInstanceMethod(node, receiver.Instance, receiver.Class, name.Value, arguments, scope))
	case *object.Class:
		return object.NewError(fault.Property, at, "class `%s` has no method `%s` to call on the class itself", receiver.Name.Value, name.Value).
			WithHelp("methods are called on instances: `new %s().%s()`", receiver.Name.Value, name.Value)
	case *object.NativeClass:
		return object.NewError(fault.Property, at, "class `%s` has no method `%s` to call on the class itself", receiver.Name, name.Value).
			WithHelp("methods are called on instances: `new %s().%s()`", receiver.Name, name.Value)
	case *object.LibraryModule:
		if function, ok := receiver.Methods[name.Value]; ok {
			return unwrapCall(at, function, arguments, scope)
		}

		return object.NewError(fault.Property, at, "module `%s` has no method `%s`", receiver.Name, name.Value).
			WithHelp("%s", moduleSuggestion(receiver, name.Value))
	}

	if !handled || result == nil {
		return object.NewError(fault.Property, at, "%s has no method `%s`", object.TypeName(left), name.Value)
	}

	return result
}

// callInstanceMethod resolves a method starting at the given class and invokes
// it against the receiver. `super` calls pass the superclass of the declaring
// class as the starting point; ordinary calls pass the receiver's own class.
func callInstanceMethod(node *ast.Method, receiver *object.Instance, start *object.Class, name string, arguments []object.Object, scope *object.Scope) object.Object {
	at := node.Token

	if identifier, ok := node.Method.(*ast.Identifier); ok {
		at = identifier.Token
	}

	member, declaringClass, ok := object.LookupMember(start, name)

	if !ok {
		return unknownMethod(at, receiver.Class, name)
	}

	method, ok := member.(*object.Function)

	if !ok {
		return object.NewError(fault.Type, at, "`%s.%s` is a %s, which cannot be called", receiver.Class.Name.Value, name, object.TypeName(member))
	}

	result := invokeMethod(method, receiver, declaringClass, arguments, scope, at)

	if failed, ok := result.(*object.Error); ok {
		return failed.WithFrame(receiver.Class.Name.Value+"."+name+"()", at)
	}

	return result
}
