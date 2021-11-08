package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateIf(node *ast.If) (object.Object, bool) {
	condition, ok := Evaluate(node.Condition)

	if !ok {
		return nil, false
	}

	conditionIsTrue := isTruthy(condition)

	if conditionIsTrue {
		return Evaluate(node.Consequence)
	} else if node.Alternative != nil && !conditionIsTrue {
		return Evaluate(node.Alternative)
	}

	return value.NULL, true
}
