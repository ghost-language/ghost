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

	arguments := evaluateExpressions(node.Arguments, scope)

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
			return unwrapCall(node.Token, function, arguments, scope)
		}
	}

	return newError("%d:%d: runtime error: unknown method: %s.%s", node.Token.Line, node.Token.Column, left.String(), node.Method.(*ast.Identifier).Value)
}
