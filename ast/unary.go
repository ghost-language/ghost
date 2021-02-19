package ast

import "ghostlang.org/x/ghost/token"

// Unary structures are for unary expressions.
type Unary struct {
	ExpressionNode
	Operator token.Token
	Right    ExpressionNode
}
