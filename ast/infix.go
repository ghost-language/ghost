package ast

import "ghostlang.org/x/ghost/token"

type Infix struct {
	ExpressionNode
	Token    token.Token
	Left     ExpressionNode
	Operator token.Type
	Right    ExpressionNode
}
