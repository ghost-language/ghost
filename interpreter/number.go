package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateNumber(node *ast.Number) (object.Object, bool) {
	return &object.Number{Value: node.Value}, true
}
