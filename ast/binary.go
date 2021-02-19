package ast

import "ghostlang.org/x/ghost/token"

// Binary structures are for binary expressions.
type Binary struct {
	ExpressionNode
	Left     ExpressionNode
	Operator token.Token
	Right    ExpressionNode
}
