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
			resultType := result.Type()

			if resultType == object.ERROR || resultType == object.RETURN {
				return result
			}
		}
	}

	return result
}
