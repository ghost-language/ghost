package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
)

func evaluateCall(node *ast.Call, env *environment.Environment) (object.Object, bool) {
	callee, ok := Evaluate(node.Callee, env)

	if !ok {
		return nil, false
	}

	arguments, ok := evaluateExpressions(node.Arguments, env)

	if !ok {
		return nil, false
	}

	return unwrapCall(node.Token, callee, arguments)
}

func unwrapCall(tok token.Token, callee object.Object, arguments []object.Object) (object.Object, bool) {
	switch callee := callee.(type) {
	case *object.LibraryFunction:
		if result := callee.Function(arguments...); result != nil {
			return result, true
		}

		return value.NULL, true
	default:
		return nil, false
	}
}
