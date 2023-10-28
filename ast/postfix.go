package ast

import "ghostlang.org/x/ghost/token"

type Postfix struct {
	ExpressionNode
	Token    token.Token
	Operator string
}
