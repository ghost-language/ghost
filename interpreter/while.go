package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/helper"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateWhile(node *ast.While, env *object.Environment) (object.Object, bool) {
	for {
		condition, ok := Evaluate(node.Condition, env)

		if !ok {
			return nil, ok
		}

		if helper.IsTruthy(condition) {
			Evaluate(node.Body, env)
		} else {
			break
		}
	}

	return value.NULL, true
}
