package evaluator

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateThis(node *ast.This, scope *object.Scope) object.Object {
	if scope.Self != nil {
		return scope.Self
	}

	pairs := make(map[object.MapKey]object.MapPair)

	fmt.Printf("self: %v\n", scope.Self)

	return &object.Map{Pairs: pairs}

	// return object.NewError("%d:%d:%s: runtime error: cannot call 'this' outside of scope", node.Token.Line, node.Token.Column, node.Token.File)
}
