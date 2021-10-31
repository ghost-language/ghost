package ast

import "ghostlang.org/x/ghost/token"

type Prefix struct {
	ExpressionNode
	Token    token.Token
	Operator string
	Right    ExpressionNode
}
