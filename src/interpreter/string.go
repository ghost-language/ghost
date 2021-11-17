package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateString(node *ast.String, env *object.Environment) (object.Object, bool) {
	return &object.String{Value: node.Value}, true
}
