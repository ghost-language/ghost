package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateClass(node *ast.Class, env *object.Environment) (object.Object, bool) {
	env.Set(node.Name.Lexeme, nil)

	// methods := make(map[string]*object.UserFunction)

	// for _, method := range node.Methods {
	// 	function := &object.UserFunction{
	// 		Parameters: method.Parameters,
	// 		Body: method.Body,
	// 		Defaults: method.Defaults,
	// 		Env: env,
	// 	}

	// 	methods[method.Name] = function
	// }

	class := &object.Class{
		Name: node.Name.Lexeme,
		// Methods: methods,
	}

	env.Set(node.Name.Lexeme, class)

	return value.NULL, true
}