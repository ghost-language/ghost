package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
)

func evaluateBlock(node *ast.Block, env *environment.Environment) (object.Object, bool) {
	env = environment.Extend(env)

	for _, statement := range node.Statements {
		_, err := Evaluate(statement, env)

		if err != true {
			return nil, err
		}
	}

	return nil, true
}
