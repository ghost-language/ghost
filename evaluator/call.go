package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
)

func evaluateCall(node *ast.Call, scope *object.Scope) object.Object {
	callee := Evaluate(node.Callee, scope)

	if isError(callee) {
		return callee
	}

	arguments := evaluateExpressions(node.Arguments, scope)

	if len(arguments) == 1 && isError(arguments[0]) {
		return arguments[0]
	}

	return unwrapCall(node.Token, callee, arguments, scope)
}

func unwrapCall(tok token.Token, callee object.Object, arguments []object.Object, scope *object.Scope) object.Object {
	if callee == nil {
		return newError("%d:%d:%s: runtime error: uncallable object: null", tok.Line, tok.Column, tok.File)
	}

	switch callee := callee.(type) {
	case *object.LibraryFunction:
		if result := callee.Function(scope, tok, arguments...); result != nil {
			return result
		}

		return nil
	case *object.LibraryProperty:
		if result := callee.Property(scope, tok); result != nil {
			return result
		}

		return nil
	case *object.Function:
		functionEnvironment := createFunctionEnvironment(callee, arguments)
		functionScope := &object.Scope{Self: callee, Environment: functionEnvironment}

		// A function declared inside a method keeps that method's receiver, so
		// `this` and `super` still work from within a closure.
		if callee.Scope != nil {
			if instance, ok := callee.Scope.Self.(*object.Instance); ok {
				functionScope.Self = instance
				functionScope.Class = callee.Scope.Class
			}
		}

		evaluated := Evaluate(callee.Body, functionScope)

		return unwrapReturn(evaluated)
	default:
		return newError("%d:%d:%s: runtime error: uncallable object: %s", tok.Line, tok.Column, tok.File, callee.Type())
	}
}

func unwrapReturn(obj object.Object) object.Object {
	switch value := obj.(type) {
	case *object.Error:
		return obj
	case *object.Return:
		return value.Value
	}

	return value.NULL
}
