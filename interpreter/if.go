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

	if isTruthy(condition) {
		return Evaluate(node.Consequence)
	} else if node.Alternative != nil {
		return Evaluate(node.Alternative)
	}

	return value.NULL, true
}
