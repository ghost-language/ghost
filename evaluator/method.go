package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateMethod(node *ast.Method, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)

	if isError(left) {
		return left
	}

	name := node.Method.(*ast.Identifier)

	if left == nil {
		return newError("%d:%d:%s: runtime error: cannot call method %s on a null value", node.Token.Line, node.Token.Column, node.Token.File, name.Value)
	}

	arguments := evaluateExpressions(node.Arguments, scope)

	if len(arguments) == 1 && isError(arguments[0]) {
		return arguments[0]
	}

	result, handled := left.Method(name.Value, arguments)

	if isError(result) {
		return result
	}

	switch receiver := left.(type) {
	case *object.Map:
		property := &object.String{Value: name.Value}

		if function, ok := receiver.Pairs[property.MapKey()]; ok {
			return unwrapCall(node.Token, function.Value, arguments, scope)
		}

		if handled {
			return result
		}

		return newError("%d:%d:%s: runtime error: unknown method: %s.%s", node.Token.Line, node.Token.Column, node.Token.File, receiver.Type(), name.Value)
	case *object.Instance:
		return unwrapReturn(callInstanceMethod(node, receiver, receiver.Class, name.Value, arguments))
	case *object.Super:
		return unwrapReturn(callInstanceMethod(node, receiver.Instance, receiver.Class, name.Value, arguments))
	case *object.Class:
		return newError("%d:%d:%s: runtime error: unknown method %s on class %s; construct instances with `new %s()`", node.Token.Line, node.Token.Column, node.Token.File, name.Value, receiver.Name.Value, receiver.Name.Value)
	case *object.LibraryModule:
		if function, ok := receiver.Methods[name.Value]; ok {
			return unwrapCall(node.Token, function, arguments, scope)
		}
	}

	if !handled || result == nil {
		return newError("%d:%d:%s: runtime error: unknown method: %s.%s", node.Token.Line, node.Token.Column, node.Token.File, left.Type(), name.Value)
	}

	return result
}

// callInstanceMethod resolves a method starting at the given class and invokes
// it against the receiver. `super` calls pass the superclass of the declaring
// class as the starting point; ordinary calls pass the receiver's own class.
func callInstanceMethod(node *ast.Method, receiver *object.Instance, start *object.Class, name string, arguments []object.Object) object.Object {
	member, declaringClass, ok := object.LookupMember(start, name)

	if !ok {
		return newError("%d:%d:%s: runtime error: undefined method %s for class %s", node.Token.Line, node.Token.Column, node.Token.File, name, receiver.Class.Name.Value)
	}

	method, ok := member.(*object.Function)

	if !ok {
		return newError("%d:%d:%s: runtime error: invalid type %s in class %s", node.Token.Line, node.Token.Column, node.Token.File, member.Type(), receiver.Class.Name.Value)
	}

	return invokeMethod(method, receiver, declaringClass, arguments)
}
