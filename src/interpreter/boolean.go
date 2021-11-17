package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateBoolean(node *ast.Boolean, env *object.Environment) (object.Object, bool) {
	return toBooleanValue(node.Value), true
}
