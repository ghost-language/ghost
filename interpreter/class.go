package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateClass(node *ast.Class, env *object.Environment) (object.Object, bool) {
	env.Declare(node.Name.Lexeme, nil)

	methods := make(map[string]*object.UserFunction)

	for _, method := range node.Methods {
		function := &object.UserFunction{
			Parameters: method.Parameters,
			Body: method.Body,
			Defaults: method.Defaults,
			Env: env,
		}

		methods[method.Name] = function
	}

	class := &object.Class{
		Name: node.Name.Lexeme,
		Methods: methods,
	}

	env.Assign(node.Name, class)

	return nil, true
}