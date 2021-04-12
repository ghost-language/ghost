package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateDeclaration(node *ast.Declaration, env *object.Environment) (object.Object, bool) {
	if node.Initializer != nil {
		val, success := Evaluate(node.Initializer, env)

		if success != true {
			return nil, false
		}

		env.Declare(node.Name.Lexeme, val)
	} else {
		env.Declare(node.Name.Lexeme, value.NULL)
	}

	return nil, true
}
