package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/helper"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluateLogical(node *ast.Logical, env *environment.Environment) (object.Object, bool) {
	left, err := Evaluate(node.Left, env)

	if !err {
		return nil, err
	}

	if node.Operator.Type == token.OR {
		if helper.IsTruthy(left) {
			return left, true
		}
	} else if node.Operator.Type == token.AND {
		if !helper.IsTruthy(left) {
			return left, true
		}
	}

	return Evaluate(node.Right, env)
}
