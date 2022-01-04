package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateBlock(node *ast.Block, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range node.Statements {
		result = Evaluate(statement, env)

		switch result.(type) {
		case *object.Error:
			return result
		}
	}

	return result
}
