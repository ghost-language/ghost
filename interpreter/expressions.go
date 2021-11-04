package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateExpressions(expressions []ast.ExpressionNode) ([]object.Object, bool) {
	var result []object.Object

	for _, expression := range expressions {
		evaluated, ok := Evaluate(expression)

		if !ok {
			return nil, false
		}

		result = append(result, evaluated)
	}

	return result, true
}
