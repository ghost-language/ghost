package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateExpressions(expressions []ast.ExpressionNode, scope *object.Scope) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		evaluated := Evaluate(expression, scope)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}
