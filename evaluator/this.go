package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
)

func evaluateThis(node *ast.This, scope *object.Scope) object.Object {
	switch scope.Self.(type) {
	case *object.Instance, *object.Class:
		return scope.Self
	}

	return object.NewError(fault.Name, node.Token, "`this` can only be used inside a class")
}
