package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateFunction(node *ast.Function, env *object.Environment) (object.Object, bool) {
	return value.NULL, true
}