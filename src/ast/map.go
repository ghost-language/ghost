package ast

import "ghostlang.org/x/ghost/token"

type Map struct {
	ExpressionNode
	Token token.Token
	Pairs map[ExpressionNode]ExpressionNode
}
