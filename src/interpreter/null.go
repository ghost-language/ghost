package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateNull(node *ast.Null, env *object.Environment) (object.Object, bool) {
	return value.NULL, true
}
