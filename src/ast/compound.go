package ast

import "ghostlang.org/x/ghost/token"

type Compound struct {
	ExpressionNode
	Token    token.Token
	Left     ExpressionNode
	Operator string
	Right    ExpressionNode
}
