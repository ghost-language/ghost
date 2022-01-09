package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateClass(node *ast.Class, scope *object.Scope) object.Object {
	class := &object.Class{
		Name:  node.Name,
		Scope: scope,
	}

	scope.Environment.Set(node.Name.Value, class)

	return class
}
