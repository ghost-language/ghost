package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateMethod(node *ast.Method, env *object.Environment) (object.Object, bool) {
	left, ok := Evaluate(node.Left, env)

	if !ok {
		return nil, false
	}

	arguments, ok := evaluateExpressions(node.Arguments, env)

	if !ok {
		return nil, false
	}

	switch left.(type) {
	case *object.LibraryModule:
		method := node.Method.(*ast.Identifier)
		module := left.(*object.LibraryModule)

		if !ok {
			return nil, false
		}

		if function, ok := module.Methods[method.Value]; ok {
			return unwrapCall(node.Token, function, arguments, env)
		}
	}

	return nil, false
}
