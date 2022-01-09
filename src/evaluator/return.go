package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateReturn(node *ast.Return, scope *object.Scope) object.Object {
	value := Evaluate(node.Value, scope)

	return &object.Return{Value: value}
}
