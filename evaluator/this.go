package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateThis(node *ast.This, scope *object.Scope) object.Object {
	if scope.Self != nil {
		return scope.Self
	}

	pairs := make(map[object.MapKey]object.MapPair)

	return &object.Map{Pairs: pairs}
}
