package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
)

func evaluateInfix(node *ast.Infix, env *environment.Environment) (object.Object, bool) {
	left, ok := Evaluate(node.Left, env)

	if !ok {
		return nil, false
	}

	right, ok := Evaluate(node.Right, env)

	if !ok {
		return nil, false
	}

	switch {
	case left.Type() == object.NUMBER && right.Type() == object.NUMBER:
		return evaluateNumberInfix(node, left, right)
	default:
		return nil, false
	}
}
