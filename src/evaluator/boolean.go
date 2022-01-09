package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateBoolean(node *ast.Boolean, scope *object.Scope) object.Object {
	return toBooleanValue(node.Value)
}
