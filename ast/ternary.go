package ast

import "ghostlang.org/x/ghost/token"

type Ternary struct {
	ExpressionNode
	Token     token.Token
	Condition ExpressionNode
	IfTrue    ExpressionNode
	IfFalse   ExpressionNode
}
