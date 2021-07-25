package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateFunction(node *ast.Function, env *object.Environment) (object.Object, bool) {
	name := node.Name.Lexeme

	function := &object.UserFunction{Env: env, Body: node.Body}

	env.Set(name, function)

	return value.NULL, true
}