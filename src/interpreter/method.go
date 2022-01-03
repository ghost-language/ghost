package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateMethod(node *ast.Method, env *object.Environment) object.Object {
	left := Evaluate(node.Left, env)

	if isError(left) {
		return left
	}

	arguments := evaluateExpressions(node.Arguments, env)

	if len(arguments) == 1 && isError(arguments[0]) {
		return arguments[0]
	}

	if result, ok := left.Method(node.Method.(*ast.Identifier).Value, arguments); ok {
		return result
	}

	switch left.(type) {
	case *object.LibraryModule:
		method := node.Method.(*ast.Identifier)
		module := left.(*object.LibraryModule)

		if function, ok := module.Methods[method.Value]; ok {
			return unwrapCall(node.Token, function, arguments, env)
		}
	}

	return newError("unknown method: %s.%s", left.String(), node.Method.(*ast.Identifier).Value)
}