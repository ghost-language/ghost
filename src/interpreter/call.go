package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
)

func evaluateCall(node *ast.Call, env *object.Environment) (object.Object, bool) {
	callee, ok := Evaluate(node.Callee, env)

	if !ok {
		return nil, false
	}

	arguments, ok := evaluateExpressions(node.Arguments, env)

	if !ok {
		return nil, false
	}

	return unwrapCall(node.Token, callee, arguments, env)
}

func unwrapCall(tok token.Token, callee object.Object, arguments []object.Object, env *object.Environment) (object.Object, bool) {
	switch callee := callee.(type) {
	case *object.LibraryFunction:
		if result := callee.Function(arguments...); result != nil {
			return result, true
		}

		return nil, true
	case *object.Function:
		functionEnvironment := createFunctionEnvironment(callee, arguments)

		Evaluate(callee.Body, functionEnvironment)

		// to do, parse and return return values
		return value.NULL, true
	default:
		log.Debug("found uncallable object: %T", callee.Type)
		return nil, false
	}
}
