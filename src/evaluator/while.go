package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateWhile(node *ast.While, scope *object.Scope) object.Object {
	for {
		condition := Evaluate(node.Condition, scope)

		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			evaluated := Evaluate(node.Consequence, scope)

			if isTerminator(evaluated) {
				switch val := evaluated.(type) {
				case *object.Error:
					return val
				case *object.Continue:
					//
				case *object.Break:
					return value.NULL
				}
			}
		} else {
			break
		}
	}

	return nil
}
