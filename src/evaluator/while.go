package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateWhile(node *ast.While, env *object.Environment) object.Object {
	for {
		condition := Evaluate(node.Condition, env)

		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			Evaluate(node.Consequence, env)
		} else {
			break
		}
	}

	return nil
}
