package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluateString(node *ast.String, scope *object.Scope) object.Object {
	return &object.String{Value: node.Value}
}

func evaluateStringInfix(node *ast.Infix, left object.Object, right object.Object) object.Object {
	leftValue := left.String()
	rightValue := right.String()

	switch node.Operator {
	case token.PLUS:
		return &object.String{Value: leftValue + rightValue}
	case token.LESS:
		return &object.Boolean{Value: leftValue < rightValue}
	case token.LESSEQUAL:
		return &object.Boolean{Value: leftValue <= rightValue}
	case token.GREATER:
		return &object.Boolean{Value: leftValue > rightValue}
	case token.GREATEREQUAL:
		return &object.Boolean{Value: leftValue >= rightValue}
	case token.EQUALEQUAL:
		return &object.Boolean{Value: leftValue == rightValue}
	case token.BANGEQUAL:
		return &object.Boolean{Value: leftValue != rightValue}
	}

	return newError("%d:%d:%s: runtime error: unknown operator: %s %s %s", node.Token.Line, node.Token.Column, node.Token.File, right.Type(), node.Operator, left.Type())
}
