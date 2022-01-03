package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateIf(node *ast.If, env *object.Environment) object.Object {
	condition := Evaluate(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Evaluate(node.Consequence, env)
	} else if node.Alternative != nil {
		return Evaluate(node.Alternative, env)
	}

	return value.NULL
}
