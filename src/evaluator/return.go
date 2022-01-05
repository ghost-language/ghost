package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateReturn(node *ast.Return, env *object.Environment) object.Object {
	value := Evaluate(node.Value, env)

	return &object.Return{Value: value}
}
