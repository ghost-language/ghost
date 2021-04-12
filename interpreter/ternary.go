package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/helper"
	"ghostlang.org/x/ghost/object"
)

func evaluateTernary(node *ast.Ternary, env *object.Environment) (object.Object, bool) {
	condition, _ := Evaluate(node.Condition, env)

	if helper.IsTruthy(condition) {
		return Evaluate(node.Then, env)
	}

	return Evaluate(node.Else, env)
}
