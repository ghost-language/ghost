package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateIf(node *ast.If, env *environment.Environment) (object.Object, bool) {
	condition, ok := Evaluate(node.Condition, env)

	if !ok {
		return nil, false
	}

	conditionIsTrue := isTruthy(condition)

	if conditionIsTrue {
		return Evaluate(node.Consequence, env)
	} else if node.Alternative != nil && !conditionIsTrue {
		return Evaluate(node.Alternative, env)
	}

	return value.NULL, true
}
