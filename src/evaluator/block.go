package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateBlock(node *ast.Block, scope *object.Scope) object.Object {
	var result object.Object

	for _, statement := range node.Statements {
		result = Evaluate(statement, scope)

		if result != nil {
			switch statement := result.(type) {
			case *object.Error:
			case *object.Return:
				return statement
			}
		}
	}

	return result
}
