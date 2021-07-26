package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateReturn(node *ast.Return, env *object.Environment) (object.Object, bool) {
	val, ok := Evaluate(node.Value, env)

	if !ok {
		return nil, ok
	}

	return &object.Return{Value: val}, ok
}