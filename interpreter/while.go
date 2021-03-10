package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/helper"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateWhile(node *ast.While, env *environment.Environment) (object.Object, bool) {
	for {
		condition, success := Evaluate(node.Condition, env)

		if !success {
			return nil, success
		}

		if helper.IsTruthy(condition) {
			break
		}

		_, success = Evaluate(node.Body, env)

		if !success {
			return nil, success
		}
	}

	return value.NULL, true
}
