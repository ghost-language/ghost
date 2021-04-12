package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateBlock(node *ast.Block, env *object.Environment) (object.Object, bool) {
	blockEnv := object.ExtendEnvironment(env)

	for _, statement := range node.Statements {
		_, err := Evaluate(statement, blockEnv)

		if err != true {
			return nil, err
		}
	}

	return nil, true
}
