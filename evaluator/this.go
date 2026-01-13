package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateThis(node *ast.This, scope *object.Scope) object.Object {
	switch scope.Self.(type) {
	case *object.Instance, *object.Class:
		return scope.Self
	}

	return newError("%d:%d:%s: runtime error: 'this' used outside of class context", node.Token.Line, node.Token.Column, node.Token.File)
}
