package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateBreak(node *ast.Break, scope *object.Scope) object.Object {
	return value.BREAK
}
