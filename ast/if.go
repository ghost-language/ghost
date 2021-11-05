package ast

import "ghostlang.org/x/ghost/token"

type If struct {
	ExpressionNode
	Token       token.Token
	Condition   ExpressionNode
	Consequence *Block
	Alternative *Block
}
