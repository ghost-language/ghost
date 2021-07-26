package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateBlock(node *ast.Block, env *object.Environment) (object.Object, bool) {
	for _, statement := range node.Statements {
		_, ok := Evaluate(statement, env)

		if !ok {
			return nil, ok
		}
	}

	return nil, true
}
