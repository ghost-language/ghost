package ast

import "ghostlang.org/x/ghost/token"

type Property struct {
	ExpressionNode
	AssignmentNode
	Token    token.Token
	Left     ExpressionNode
	Property ExpressionNode
}
