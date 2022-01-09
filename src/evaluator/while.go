package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateWhile(node *ast.While, scope *object.Scope) object.Object {
	for {
		condition := Evaluate(node.Condition, scope)

		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			Evaluate(node.Consequence, scope)
		} else {
			break
		}
	}

	return nil
}
