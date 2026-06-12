package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateCompound(node *ast.Compound, scope *object.Scope) object.Object {
	infix := &ast.Infix{
		Token:    node.Token,
		Left:     node.Left,
		Operator: node.Operator[:len(node.Operator)-1],
		Right:    node.Right,
	}

	value := evaluateInfix(infix, scope)

	if isError(value) {
		return value
	}

	switch target := node.Left.(type) {
	case *ast.Identifier:
		scope.Environment.Set(target.Value, value)
	case *ast.Index:
		return evaluateIndexAssignment(target, value, scope)
	case *ast.Property:
		return evaluatePropertyAssignment(target, value, scope)
	default:
		return newError("%d:%d:%s: runtime error: invalid compound assignment target: %T", node.Token.Line, node.Token.Column, node.Token.File, node.Left)
	}

	return nil
}
