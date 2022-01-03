package interpreter

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
	case *object.LibraryFunction:
		if result := callee.Function(env, arguments...); result != nil {
			return result
		}

		return nil
	case *object.Function:
		functionEnvironment := createFunctionEnvironment(callee, arguments)
		Evaluate(callee.Body, functionEnvironment)

		// to do, parse and return return values
		return value.NULL
	default:
		return newError("uncallable object: %s", callee.Type())
	}
}
