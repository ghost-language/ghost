package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
)

func evaluateCall(node *ast.Call, env *object.Environment) object.Object {
	callee := Evaluate(node.Callee, env)

	if isError(callee) {
		return callee
	}

	arguments := evaluateExpressions(node.Arguments, env)

	if len(arguments) == 1 && isError(arguments[0]) {
		return arguments[0]
	}

	return unwrapCall(node.Token, callee, arguments, env)
}

func unwrapCall(tok token.Token, callee object.Object, arguments []object.Object, env *object.Environment) object.Object {
	switch callee := callee.(type) {
	case *object.Class:
		instance := &object.Instance{Class: callee}

		return instance
	case *object.LibraryFunction:
		if result := callee.Function(env, tok, arguments...); result != nil {
			return result
		}

		return nil
	case *object.LibraryProperty:
		if result := callee.Property(env, tok); result != nil {
			return result
		}

		return nil
	case *object.Function:
		functionEnvironment := createFunctionEnvironment(callee, arguments)
		evaluated := Evaluate(callee.Body, functionEnvironment)

		return unwrapReturn(evaluated)
	default:
		return newError("%d:%d: runtime error: uncallable object: %s", tok.Line, tok.Column, callee.Type())
	}
}

func unwrapReturn(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.Return); ok {
		return returnValue.Value
	}

	return value.NULL
}
