package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateExpressions(expressions []ast.ExpressionNode, env *object.Environment) ([]object.Object, bool) {
	var result []object.Object

	for _, expression := range expressions {
		evaluated, ok := Evaluate(expression, env)

		if !ok {
			return nil, false
		}

		result = append(result, evaluated)
	}

	return result, true
}
