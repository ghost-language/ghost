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
			evaluated := Evaluate(node.Consequence, scope)

			switch value := evaluated.(type) {
			case *object.Error:
				return value
			case *object.Return:
				return value.Value
			}
		} else {
			break
		}
	}

	return nil
}
