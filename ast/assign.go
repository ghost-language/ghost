package ast

import "ghostlang.org/x/ghost/token"

type Assign struct {
	ExpressionNode
	Name  token.Token
	Value ExpressionNode
}
