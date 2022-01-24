package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateContinue(node *ast.Continue, scope *object.Scope) object.Object {
	return value.CONTINUE
}
