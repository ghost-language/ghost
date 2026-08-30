package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateTernary(node *ast.Ternary, scope *object.Scope) object.Object {
	condition := Evaluate(node.Condition, scope)

	if isError(condition) {
		return condition
	}

	if object.IsTrue(condition) {
		return Evaluate(node.IfTrue, scope)
	}

	return Evaluate(node.IfFalse, scope)
}
