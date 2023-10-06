package ast

import "ghostlang.org/x/ghost/token"

type While struct {
	ExpressionNode
	Token       token.Token
	Condition   ExpressionNode
	Consequence *Block
}
