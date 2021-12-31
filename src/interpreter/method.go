package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
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

		// libraryModule, ok := library.Modules[module.Name]

		if !ok {
			return nil, false
		}

		log.Debug("module: %s, method: %s, library: %q", module.Name, method.Value, module.Methods)
	}

	return unwrapCall(node.Token, left, arguments, env)
}
