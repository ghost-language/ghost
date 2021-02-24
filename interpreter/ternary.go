package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/helper"
	"ghostlang.org/x/ghost/object"
)

func evaluateTernary(node *ast.Ternary) (object.Object, bool) {
	condition, _ := Evaluate(node.Condition)

	if helper.IsTruthy(condition) {
		return Evaluate(node.Then)
	}

	return Evaluate(node.Else)
}
