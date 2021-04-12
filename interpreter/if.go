package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/helper"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateIf(node *ast.If, env *object.Environment) (object.Object, bool) {
	if conditionValue, success := Evaluate(node.Condition, env); success == true {
		if helper.IsTruthy(conditionValue) {
			return Evaluate(node.Then, env)
		} else if node.Else != nil {
			return Evaluate(node.Else, env)
		}
	}

	return value.NULL, true
}
